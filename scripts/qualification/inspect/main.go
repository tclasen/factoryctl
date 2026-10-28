// inspect reads foreign binaries without executing them.
package main

import (
	"bytes"
	"crypto/sha256"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func inspect(file string) (map[string]any, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(data)
	result := map[string]any{"path": file, "bytes": len(data), "sha256": hex.EncodeToString(sum[:]), "execution": "not executed"}
	target, err := readInfo(file)
	if err != nil {
		return nil, err
	}
	result["build_settings"] = target
	if target["CGO_ENABLED"] != "0" {
		return nil, fmt.Errorf("target does not prove CGO_ENABLED=0")
	}
	// These unique embedded UI/default strings are checked in every target, not
	// merely in the source tree. This is packaging evidence, not browser execution.
	for _, token := range []string{"<!DOCTYPE html>", "Factoryctl", "Record Definition"} {
		if !bytes.Contains(bytes.ToLower(data), bytes.ToLower([]byte(token))) {
			return nil, fmt.Errorf("embedded marker absent: %s", token)
		}
	}
	result["embedded_ui_and_defaults_markers"] = true
	switch target["GOOS"] {
	case "linux":
		f, e := elf.Open(file)
		if e != nil {
			return nil, e
		}
		defer f.Close()
		libs, e := f.ImportedLibraries()
		if e != nil {
			return nil, e
		}
		result["dynamic_libraries"] = libs
		for _, p := range f.Progs {
			if p.Type == elf.PT_INTERP {
				return nil, fmt.Errorf("Linux binary has dynamic loader")
			}
		}
		if len(libs) != 0 {
			return nil, fmt.Errorf("Linux binary has dynamic libraries")
		}
		result["linkage"] = "no ELF interpreter or dynamic library dependencies"
	case "darwin":
		f, e := macho.Open(file)
		if e != nil {
			return nil, e
		}
		defer f.Close()
		libs, e := f.ImportedLibraries()
		if e != nil {
			return nil, e
		}
		result["dynamic_libraries"] = libs
		result["linkage"] = "CGO-disabled self-contained application; Darwin system interfaces retained, not fully static"
	case "windows":
		f, e := pe.Open(file)
		if e != nil {
			return nil, e
		}
		defer f.Close()
		// debug/pe.ImportedLibraries is currently a TODO returning nil.
		// ImportedSymbols actually decodes the import table (name:DLL).
		symbols, e := f.ImportedSymbols()
		if e != nil {
			return nil, e
		}
		libs := []string{}
		seen := map[string]bool{}
		for _, symbol := range symbols {
			parts := strings.Split(symbol, ":")
			if len(parts) != 2 || parts[1] == "" {
				return nil, fmt.Errorf("malformed PE import: %s", symbol)
			}
			lib := strings.ToLower(parts[1])
			if !seen[lib] {
				libs = append(libs, lib)
				seen[lib] = true
			}
		}
		sort.Strings(libs)
		result["dynamic_libraries"] = libs
		result["imported_symbol_count"] = len(symbols)
		result["import_scope"] = "PE startup imports only; OS libraries loaded dynamically by the Go runtime are not enumerated"
		if len(libs) == 0 {
			return nil, fmt.Errorf("Windows startup import table unexpectedly empty")
		}
		for _, lib := range libs {
			if lib != "kernel32.dll" {
				return nil, fmt.Errorf("unqualified Windows DLL dependency: %s", lib)
			}
		}
		result["linkage"] = "CGO-disabled self-contained application; Windows platform DLL interface retained"
	default:
		return nil, fmt.Errorf("unexpected target OS")
	}
	return result, nil
}
func main() {
	results := []map[string]any{}
	for _, file := range os.Args[1:] {
		r, e := inspect(file)
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		results = append(results, r)
	}
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "supply target binaries")
		os.Exit(1)
	}
	_ = json.NewEncoder(os.Stdout).Encode(results)
}

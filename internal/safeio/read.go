// Package safeio opens only regular files and bounds reads of untrusted paths.
package safeio

import (
	"fmt"
	"io"
	"os"
)

// OpenRegular checks before and after open. Unix opens are nonblocking and
// no-follow, so a racing replacement with a FIFO cannot block the process.
func OpenRegular(name string, flags int) (*os.File, error) {
	st, e := os.Lstat(name)
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file: %s", name)
	}
	f, e := os.OpenFile(name, flags|readFlags, 0)
	return checked(f, e)
}
func OpenRegularRoot(root *os.Root, name string, flags int) (*os.File, error) {
	st, e := root.Lstat(name)
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file: %s", name)
	}
	f, e := root.OpenFile(name, flags|readFlags, 0)
	return checked(f, e)
}
func checked(f *os.File, e error) (*os.File, error) {
	if e != nil {
		return nil, e
	}
	st, e := f.Stat()
	if e != nil {
		f.Close()
		return nil, e
	}
	if !st.Mode().IsRegular() {
		f.Close()
		return nil, fmt.Errorf("not a regular file")
	}
	return f, nil
}
func ReadFile(name string, max int64) ([]byte, error) {
	f, e := OpenRegular(name, os.O_RDONLY)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return Read(f, max)
}
func ReadRoot(root *os.Root, name string, max int64) ([]byte, error) {
	f, e := OpenRegularRoot(root, name, os.O_RDONLY)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return Read(f, max)
}

// Read returns max+1 bytes for caller-specific over-limit diagnostics.
func Read(f *os.File, max int64) ([]byte, error) { return io.ReadAll(io.LimitReader(f, max+1)) }

// OpenDirectoryRoot is used only for best-effort directory flushes. The same
// nonblocking flags prevent a racing special-file substitution from hanging.
func OpenDirectoryRoot(root *os.Root, name string) (*os.File, error) {
	st, e := root.Lstat(name)
	if e != nil {
		return nil, e
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("not a directory")
	}
	f, e := root.OpenFile(name, os.O_RDONLY|readFlags, 0)
	if e != nil {
		return nil, e
	}
	st, e = f.Stat()
	if e != nil {
		f.Close()
		return nil, e
	}
	if !st.IsDir() {
		f.Close()
		return nil, fmt.Errorf("not a directory")
	}
	return f, nil
}

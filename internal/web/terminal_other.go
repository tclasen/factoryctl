//go:build !darwin && !linux && !windows

package web

import "os"

func isTerminal(file *os.File) bool { return false }

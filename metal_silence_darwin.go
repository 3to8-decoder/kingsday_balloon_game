//go:build darwin

package main

import (
	"os"
	"syscall"
)

func suppressMetalStartupNoise() func() {
	saved, err := syscall.Dup(syscall.Stderr)
	if err != nil {
		return func() {}
	}
	devNull, err := syscall.Open(os.DevNull, syscall.O_WRONLY, 0)
	if err != nil {
		syscall.Close(saved)
		return func() {}
	}
	syscall.Dup2(devNull, syscall.Stderr)
	syscall.Close(devNull)
	return func() {
		syscall.Dup2(saved, syscall.Stderr)
		syscall.Close(saved)
	}
}

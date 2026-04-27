//go:build !darwin

package main

func suppressMetalStartupNoise() func() {
	return func() {}
}

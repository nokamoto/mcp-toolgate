//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = CI

var g0 = sh.RunCmd("go")

// CI runs the continuous integration tasks.
func CI() {
	mg.Deps(Format, Tidy, Test)
}

// Format formats the code using goimports and gofumpt.
func Format() error {
	chain := [][]string{
		{"go", "install", "golang.org/x/tools/cmd/goimports@latest"},
		{"go", "install", "mvdan.cc/gofumpt@latest"},
		{"go", "install", "github.com/google/yamlfmt/cmd/yamlfmt@latest"},
		{"goimports", "-l", "-w", "."},
		{"gofumpt", "-l", "-w", "."},
		{"yamlfmt", "."},
	}
	for _, c := range chain {
		if err := sh.RunV(c[0], c[1:]...); err != nil {
			return fmt.Errorf("failed to run command %v: %w", c, err)
		}
	}
	return nil
}

// Test runs the tests in the project.
func Test() error {
	if err := g0("test", "./..."); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}
	return nil
}

// Tidy tidies the go.mod file.
func Tidy() error {
	if err := g0("mod", "tidy"); err != nil {
		return fmt.Errorf("tidy failed: %w", err)
	}
	return nil
}

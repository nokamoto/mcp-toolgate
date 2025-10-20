//go:build mage
// +build mage

package main

import (
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
func Format() {
	g0("install", "golang.org/x/tools/cmd/goimports@latest")
	sh.Run("goimports", "-l", "-w", ".")

	g0("install", "mvdan.cc/gofumpt@latest")
	sh.Run("gofumpt", "-l", "-w", ".")

	g0("install", "github.com/google/yamlfmt/cmd/yamlfmt@latest")
	sh.Run("yamlfmt", ".")
}

// Test runs the tests in the project.
func Test() {
	g0("test", "./...")
}

// Tidy tidies the go.mod file.
func Tidy() {
	g0("mod", "tidy")
}

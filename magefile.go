//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Lint() {
	mg.Deps(BufFormat, BufLint)
}

func BufLint() error {
	mg.Deps(InstallBuf)
	return sh.RunV(buf(), "lint")
}

func BufFormat() error {
	mg.Deps(InstallBuf)
	return sh.RunV(buf(), "format", "-w")
}

// Generate generates the golang files for all the protobuf apis.
func Generate() error {
	mg.Deps(InstallBuf)

	err := sh.Rm("gen")
	if err != nil {
		return fmt.Errorf("failed to clean up generated files: %w", err)
	}

	err = sh.Run(buf(), "dep", "update")
	if err != nil {
		return fmt.Errorf("failed to update protobuf dependencies: %w", err)
	}

	err = sh.Run(buf(), "generate")
	if err != nil {
		return fmt.Errorf("failed to generate golang files: %w", err)
	}

	return nil
}

var localBin string

func LocalBin() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir: %w", err)
	}

	localBin = wd + "/bin"

	err = os.MkdirAll(localBin, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create local bin dir: %w", err)
	}

	return nil
}

// Bin installs binaries needed for generation and validation.
func Bin() {
	fmt.Println("Installing Deps...")

	mg.Deps(LocalBin, InstallBuf)
}

func goInstall(repo, version string) error {
	return sh.RunWith(map[string]string{"GOBIN": localBin}, "go", "install", repo+version)
}

const bufVersion = "1.57.0"

func buf() string {
	return localBin + "/buf-" + bufVersion
}

// InstallBuf installs the Buf CLI into the local bin.
func InstallBuf() error {
	mg.Deps(LocalBin)

	err := goInstall("github.com/bufbuild/buf/cmd/buf@v", bufVersion)
	if err != nil {
		return fmt.Errorf("failed to install buf: %s", err)
	}

	err = os.Rename(localBin+"/buf", buf())
	if err != nil {
		return fmt.Errorf("failed to rename buf binary: %w", err)
	}

	return nil
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func installBinary(path string) {
	// check for sudo
	if os.Geteuid() != 0 {
		fmt.Println("Error: pax requires root privileges to install. Please run with sudo.")
		os.Exit(1)
	}

	// move the binary to /usr/bin
	destPath := "/usr/bin/" + filepath.Base(path)
	if err := os.Rename(path, destPath); err != nil {
		fmt.Printf("Failed to move binary: %v\n", err)
		os.Exit(1)
	}
	if err := os.Chmod(destPath, 0755); err != nil {
		fmt.Printf("Failed to set permissions: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Installed to %s\n", destPath)
}

func uninstallBinary(name string) {
	// check for sudo
	if os.Geteuid() != 0 {
		fmt.Println("Error: pax requires root privileges to uninstall. Please run with sudo.")
		os.Exit(1)
	}

	// remove the binary from /usr/bin
	destPath := "/usr/bin/" + name
	if err := os.Remove(destPath); err != nil {
		fmt.Printf("Failed to remove binary: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Uninstalled %s\n", destPath)
}

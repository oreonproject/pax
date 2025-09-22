package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Name                string   `yaml:"name"`
	Description         string   `yaml:"description"`
	Version             string   `yaml:"version"`
	Origin              string   `yaml:"origin"`
	Dependencies        []string `yaml:"dependencies"`
	RuntimeDependencies []string `yaml:"runtime_dependencies"`
	Build               string   `yaml:"build"`
	Binary              string   `yaml:"binary"`
	Install             string   `yaml:"install"`
	Uninstall           string   `yaml:"uninstall"`
}

func installBinary(path string) {
	pkg := getPackageFromYaml(path)

	fmt.Printf("Installing package: %s (%s)\n", pkg.Name, " "+pkg.Version)

	// git clone the origin
	exec.Command("git", "clone", pkg.Origin).Run()

	// get another package struct for each dependency yaml file
	downloadTemporaryDependencies(pkg.Dependencies)

	var wg sync.WaitGroup
	for i, dep := range pkg.Dependencies {
		wg.Add(1)
		go func(dep string, i int) {
			defer wg.Done()
			installBinary(dep)
			fmt.Println("done", i)
		}(dep, i)
	}
	wg.Wait()

	// build the package
	exec.Command(pkg.Build).Run()

	// move the produced binary to /usr/bin
	if err := os.Rename(pkg.Binary, "/usr/bin/"+pkg.Name); err != nil {
		fmt.Printf("Failed to move binary: %v\n", err)
	}
}

func downloadTemporaryDependencies(paths []string) {
	for _, p := range paths {
		pkg := getPackageFromYaml(p)
		exec.Command("git", "clone", pkg.Origin).Run()
		os.Mkdir("tmp/pax/"+pkg.Name, 0755)

		// build the package
		exec.Command(pkg.Build).Run()
		// get the produced binary
		src, err := os.Open("tmp/pax/" + pkg.Name)
		if err != nil {
			fmt.Printf("Failed to open source binary: %v\n", err)
			continue
		}
		defer src.Close()
		dst, err := os.Create("/tmp/pax/" + pkg.Name)
		if err != nil {
			fmt.Printf("Failed to create destination binary: %v\n", err)
			continue
		}
		defer dst.Close()
		if _, err := io.Copy(dst, src); err != nil {
			fmt.Printf("Failed to copy binary: %v\n", err)
		}

		// update path
		if err := exec.Command("export", "PATH=$PATH:tmp/pax/").Run(); err != nil {
			fmt.Printf("Failed to update PATH: %v\n", err)
		}
	}
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

func getPackageFromYaml(path string) Package {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	var pkg Package
	if err := yaml.Unmarshal(data, &pkg); err != nil {
		fmt.Printf("Failed to unmarshal yaml: %v\n", err)
		os.Exit(1)
	}

	return pkg
}

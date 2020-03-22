//+build build

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	strip  bool
	goos   string
	outDir string

	help bool
)

type cmd func()

var commands = map[string]cmd{
	"build": build,
	"clean": clean,
	"test":  test,
}

var descriptions = map[string]string{
	"build": "build and package artifacts",
	"clean": "cleanup build output",
	"test":  "run tests for all packages",
}

func main() {
	flag.BoolVar(&strip, "strip", false, "run resulting executables through upx (if installed)")
	flag.StringVar(&goos, "os", "windows:linux", "systems to build for (separated by column, e.g. `windows:linux:mac`)")
	flag.StringVar(&outDir, "out", "build", "specifies build output `directory`")
	flag.BoolVar(&help, "help", false, "prints usage")
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage: go run build.go [OPTIONS] <COMMAND>\n")
		fmt.Printf("\nCommands:\n")
		for param, desc := range descriptions {
			fmt.Printf("  %s\n    \t%s\n", param, desc)
		}
		fmt.Printf("\nFlags:\n")
		flag.PrintDefaults()
	}

	command := flag.Arg(0)
	if command == "" {
		flag.Usage()
		return
	}

	if invoked, ok := commands[command]; ok {
		invoked()
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		flag.Usage()
	}
}

func must(err error) {
	if err != nil {
		fmt.Printf("err: %s\n", err)
		panic(err)
	}
}

// Builds binaries for the various OS's
func build() {
	vars := map[string]string{
		"windows": "windows",
		"linux":   "linux",
		"mac":     "darwin",
	}

	osList := strings.Split(goos, ":")
	for _, osValue := range osList {
		if v, ok := vars[osValue]; ok {
			buildArtifact(v)

			if strip {
				upx(v)
			}
		} else {
			log.Printf("Invalid OS supplied: %s, skipping...\n", osValue)
		}
	}
}

func buildArtifact(system string) {
	os.Setenv("GOOS", system)
	os.Setenv("GOARCH", "amd64")

	fileExt := ""
	ldflags := "-s -w"
	if system == "windows" {
		fileExt = ".exe"
		ldflags += " -H windowsgui"
	}

	outFile := path.Join(outDir, fmt.Sprintf("paletter-%s", system))
	outFile += fileExt

	args := []string{
		"build",
		"-o", outFile,
		"-ldflags", ldflags,
		"./cmd/paletter",
	}

	buildCmd := exec.Command("go", args...)
	log.Printf("Building package for %s %v\n", system, buildCmd.Args)
	err := buildCmd.Run()
	must(err)
}

func upx(system string) {
	outFile := path.Join(outDir, fmt.Sprintf("paletter-%s", system))
	if system == "windows" {
		outFile += ".exe"
	}

	args := []string{
		"--brute",
		outFile,
	}

	upxCmd := exec.Command("upx", args...)
	upxCmd.Stdout = os.Stdout
	err := upxCmd.Run()
	if err != nil {
		log.Println(err.Error())
		log.Println("upx error, skipping...")
		return
	}
}

// Removes built artifacts and deployment directory (playground)
func clean() {
	paths := []string{
		outDir,
	}

	for _, p := range paths {
		err := os.RemoveAll(p)
		must(err)
	}
}

// Runs tests
func test() {
	args := []string{
		"test",
		"-v",
		"-benchmem",
		"./...",
	}

	testCmd := exec.Command("go", args...)
	testCmd.Stdout = os.Stdout
	err := testCmd.Run()
	must(err)
}

package main

////////////////////////////////////////////////////////////////////////////////

import (
	"log"
	"os"
	"path"
)

////////////////////////////////////////////////////////////////////////////////

const (
	rundir     = "/src/github.com/sabhiram/trade-bot/server"
	escpath    = "/bin/esc"
	staticpath = "/static/static.go"
)

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Remove the existing static.go file from the static dir.
	p := path.Join(os.Getenv("GOPATH"), rundir, staticpath)
	if err := os.Remove(p); err != nil {
		log.Printf("Warning: Unable to remove %s!\n", p)
	}
}

////////////////////////////////////////////////////////////////////////////////`

func init() {
	// The generated code will get created under the current path, so make sure
	// it's being run from the correct directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current working directory")
	}

	runpath := path.Join(os.Getenv("GOPATH"), rundir)
	if cwd != runpath {
		log.Fatalf("must be run from directory: %s\n", runpath)
	}
	// Validate that the ESC tool has been installed.
	// `go install github.com/mjibson/esc`
	p := path.Join(os.Getenv("GOPATH"), escpath)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		log.Fatalf("esc utility missing. Run: `go install github.com/mjibson/esc`")
	}
}

////////////////////////////////////////////////////////////////////////////////

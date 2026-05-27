package main

import (
	"os"

	"github.com/aaron70/decoy-cli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"

	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/cmd"
)

const (
	exitFailure = 1
)

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitFailure)
		}
	}
}

package main

import (
	"github.com/dredzone/mongobackup/internal/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

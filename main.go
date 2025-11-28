package main

import (
	"os"

	"github.com/evarela/kq/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

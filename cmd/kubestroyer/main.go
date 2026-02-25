package main

import (
	"os"

	"github.com/Rolix44/Kubestroyer/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

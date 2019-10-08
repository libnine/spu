package main

import (
	"os"

	"../../src/cli"
)

func main() {
	args := os.Args[1:]
	cli.Run(args)
}

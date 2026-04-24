package main

import (
	"fmt"
	"os"

	"github.com/nath-sign/godu/internal/cli"
)

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scan":
		if err := cli.RunScan(os.Args[2:]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		usage()
		os.Exit(0)
	default:
		usage()
		os.Exit(1)
	}

}

func usage() {
	fmt.Println("Usage: godu <command> <paths...>")
	fmt.Println("Commands:")
	fmt.Println("  scan - scan the filesystem and compute the total size of all regular files contained within the given path")
	cli.ScanUsage()
	fmt.Println("  -h, --help - show this help message and exit")

}

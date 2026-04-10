package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nath-sign/godu/internal/scan"
)

const (
	// Path flag
	FlagPath        = "p"
	FlagPathMessage = "the path to scan, can be a file or directory, default is the current directory"

	// Symlinks flag
	FlagSymlinks        = "s"
	FlagSymlinksMessage = "follow symlinks in the scan, default is false"

	// Skip errors flag
	FlagSkipErrors        = "e"
	FlagSkipErrorsMessage = "skip errors in the scan, default is false"

	// Verbose flag
	FlagVerbose        = "v"
	FlagVerboseMessage = "verbose output (print subdir results), default is false"
)

func RunScan(args []string) error {
	var (
		path       string
		symlinks   bool
		skipErrors bool
		verbose    bool
	)
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	fs.StringVar(&path, FlagPath, ".", FlagPathMessage)
	fs.BoolVar(&symlinks, FlagSymlinks, false, FlagSymlinksMessage)
	fs.BoolVar(&skipErrors, FlagSkipErrors, false, FlagSkipErrorsMessage)
	fs.BoolVar(&verbose, FlagVerbose, false, FlagVerboseMessage)
	fs.Parse(args)

	if path == "" {
		return fmt.Errorf("path is required")
	}

	config := scan.NewConfig(filepath.Clean(path), symlinks, skipErrors, verbose)
	fmt.Println("Config:", config)

	result, err := scan.Scan(config)
	if err != nil {
		return err
	}
	fmt.Println("Result:", result)

	return nil
}

func ScanUsage() {
	fmt.Fprintf(os.Stdout, "   -%s <path> - [REQUIRED] %s\n", FlagPath, FlagPathMessage)
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagSymlinks, FlagSymlinksMessage)
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagSkipErrors, FlagSkipErrorsMessage)
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagVerbose, FlagVerboseMessage)
}

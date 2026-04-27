package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nath-sign/godu/internal/scan"
)

const (
	// FlagSymlinks enables following symlinks while walking paths.
	FlagSymlinks        = "s"
	FlagSymlinksMessage = "follow symlinks in the scan, default is false"

	// FlagSkipErrors continues the scan when read/stat operations fail.
	FlagSkipErrors        = "e"
	FlagSkipErrorsMessage = "skip errors in the scan, default is false"

	// FlagVerbose enables per-path verbose output during scanning.
	FlagVerbose        = "v"
	FlagVerboseMessage = "verbose output (print subdir results), default is false"
)

// RunScan parses scan command flags and executes scans for each root path.
func RunScan(args []string) error {
	var (
		symlinks   bool
		skipErrors bool
		verbose    bool
	)
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	fs.BoolVar(&symlinks, FlagSymlinks, false, FlagSymlinksMessage)
	fs.BoolVar(&skipErrors, FlagSkipErrors, false, FlagSkipErrorsMessage)
	fs.BoolVar(&verbose, FlagVerbose, false, FlagVerboseMessage)
	fs.Parse(args)

	roots := fs.Args()
	if len(roots) == 0 {
		return fmt.Errorf("path is required")
	}

	results := scan.InitResult("Total")
	for _, path := range roots {
		config := scan.NewConfig(filepath.Clean(path), symlinks, skipErrors, verbose)
		//fmt.Println("Config:", config)

		result, err := scan.Scan(config)
		if err != nil {
			return err
		}
		if config.Verbose {
			scan.PrintVerbose(result.TotalSize, result.RootPath)
		}
		results.AddFromResult(result)
	}
	fmt.Println("Result:", results)
	return nil
}

// ScanUsage prints the scan command flags and their descriptions.
func ScanUsage() {
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagSymlinks, FlagSymlinksMessage)
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagSkipErrors, FlagSkipErrorsMessage)
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagVerbose, FlagVerboseMessage)
}

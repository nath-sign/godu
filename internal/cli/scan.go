package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nath-sign/godu/internal/scan"
)

const (

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

func ScanUsage() {
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagSymlinks, FlagSymlinksMessage)
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagSkipErrors, FlagSkipErrorsMessage)
	fmt.Fprintf(os.Stdout, "   -%s <bool> - %s\n", FlagVerbose, FlagVerboseMessage)
}

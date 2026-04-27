package scan

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Walker struct {
	Result *Result
	cfg    *Config
}

// Scan walks the configured root path and returns aggregated metrics.
func Scan(config *Config) (*Result, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	result := InitResult(config.RootPath)
	walker := &Walker{
		Result: result,
		cfg:    config,
	}

	if err := walker.WalkPath(config.RootPath); err != nil {
		return nil, fmt.Errorf("failed to walk path: %w", err)
	}

	return result, nil
}

// PrintVerbose prints a single verbose scan output line.
func PrintVerbose(size int64, path string, writer ...io.Writer) {
	var w io.Writer
	if len(writer) == 0 {
		w = os.Stdout
	} else {
		w = writer[0]
	}
	fmt.Fprintf(w, "%-*d\t %s\n", 8, size, path)
}

// WalkSymlink resolves and walks a symlink target.
func (w *Walker) WalkSymlink(path string) error {
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		if !w.cfg.SkipErrors {
			return fmt.Errorf("failed to evaluate symlink: %w", err)
		}
		w.Result.AddError()
	} else {
		return w.WalkPath(realPath)
	}
	return nil
}

// WalkDirectory recursively walks directory entries and merges local results.
func (w *Walker) WalkDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		if !w.cfg.SkipErrors {
			return fmt.Errorf("failed to read directory: %w", err)
		}
		w.Result.AddError()
	} else {
		localResult := InitResult(path)
		localWalker := &Walker{
			Result: localResult,
			cfg:    w.cfg,
		}
		for _, entry := range entries {
			if err := localWalker.WalkPath(filepath.Join(path, entry.Name())); err != nil {
				return err
			}
		}
		w.Result.AddDirectory()
		w.Result.AddFromResult(localResult)
		if w.cfg.Verbose && filepath.Clean(path) != filepath.Clean(w.cfg.RootPath) {
			PrintVerbose(localResult.TotalSize, localResult.RootPath)
		}
	}
	return nil
}

// WalkRegularFile adds the size of a regular file to the current result.
func (w *Walker) WalkRegularFile(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if !w.cfg.SkipErrors {
			return fmt.Errorf("failed to stat path: %w", err)
		}
		w.Result.AddError()
	} else {
		w.Result.AddFile(fi.Size())
	}
	return nil
}

// WalkPath dispatches walking logic based on file mode.
func (w *Walker) WalkPath(path string) error {
	fi, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("failed to lstat path: %w", err)
	}
	mode := fi.Mode()
	switch {
	case mode&os.ModeSymlink != 0:
		if !w.cfg.FollowSymlinks {
			w.Result.AddSymlinkSkipped()
		} else {
			return w.WalkSymlink(path)
		}
	case mode.IsDir():
		return w.WalkDirectory(path)
	case mode.IsRegular():
		return w.WalkRegularFile(path)
	default:
		w.Result.AddOtherCount()
	}
	return nil
}

package scan

import (
	"fmt"
	"os"
	"path/filepath"
)

type Walker struct {
	Result *Result
	cfg    *Config
}

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

func PrintVerbose(size int64, path string) {
	fmt.Fprintf(os.Stdout, "%-*d\t %s\n", 8, size, path)
}

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

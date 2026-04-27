package scan

import "fmt"

// Config controls how Scan walks the filesystem from a root path.
type Config struct {
	RootPath       string
	FollowSymlinks bool
	SkipErrors     bool
	Verbose        bool
}

// NewConfig builds a Config for a single scan root.
func NewConfig(rootPath string, followSymlinks bool, skipErrors bool, verbose bool) *Config {
	return &Config{
		RootPath:       rootPath,
		FollowSymlinks: followSymlinks,
		SkipErrors:     skipErrors,
		Verbose:        verbose,
	}
}

// Validate checks that required Config fields are set.
func (c *Config) Validate() error {
	if c.RootPath == "" {
		return fmt.Errorf("root path is required")
	}
	return nil
}

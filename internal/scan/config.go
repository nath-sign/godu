package scan

import "fmt"

type Config struct {
	RootPath       string
	FollowSymlinks bool
	SkipErrors     bool
	Verbose        bool
}

func NewConfig(rootPath string, followSymlinks bool, skipErrors bool, verbose bool) *Config {
	return &Config{
		RootPath:       rootPath,
		FollowSymlinks: followSymlinks,
		SkipErrors:     skipErrors,
		Verbose:        verbose,
	}
}

func (c *Config) Validate() error {
	if c.RootPath == "" {
		return fmt.Errorf("root path is required")
	}
	return nil
}

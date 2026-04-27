package scan

import "fmt"

// Result stores aggregated scan metrics for a root path.
type Result struct {
	RootPath             string
	TotalSize            int64
	FilesCount           int64
	DirectoriesCount     int64
	SymlinksSkippedCount int64
	ErrorsCount          int64
	OtherCount           int64
}

// InitResult creates an empty Result for the provided root path label.
func InitResult(rootPath string) *Result {
	return &Result{
		RootPath:             rootPath,
		TotalSize:            0,
		FilesCount:           0,
		DirectoriesCount:     0,
		SymlinksSkippedCount: 0,
		ErrorsCount:          0,
		OtherCount:           0,
	}
}

// AddFile increments the file count and accumulates its size.
func (r *Result) AddFile(size int64) {
	r.TotalSize += size
	r.FilesCount++
}

// AddDirectory increments the directory count.
func (r *Result) AddDirectory() {
	r.DirectoriesCount++
}

// AddSymlinkSkipped increments the skipped symlink count.
func (r *Result) AddSymlinkSkipped() {
	r.SymlinksSkippedCount++
}

// AddError increments the error count.
func (r *Result) AddError() {
	r.ErrorsCount++
}

// AddOtherCount increments the count for non-regular/non-directory entries.
func (r *Result) AddOtherCount() {
	r.OtherCount++
}

// AddFromResult merges counters and totals from another Result.
func (r *Result) AddFromResult(other *Result) {
	r.TotalSize += other.TotalSize
	r.FilesCount += other.FilesCount
	r.DirectoriesCount += other.DirectoriesCount
	r.SymlinksSkippedCount += other.SymlinksSkippedCount
	r.ErrorsCount += other.ErrorsCount
	r.OtherCount += other.OtherCount
}

// ToGB converts TotalSize from bytes to gibibytes.
func (r *Result) ToGB() float64 {
	return float64(r.TotalSize) / 1024 / 1024 / 1024
}

// String formats Result as a human-readable summary line.
func (r *Result) String() string {
	return fmt.Sprintf("RootPath: %s, TotalSize: %d (%.2f GB), FilesCount: %d, DirectoriesCount: %d, SymlinksSkippedCount: %d, ErrorsCount: %d, OtherCount: %d", r.RootPath, r.TotalSize, r.ToGB(), r.FilesCount, r.DirectoriesCount, r.SymlinksSkippedCount, r.ErrorsCount, r.OtherCount)
}

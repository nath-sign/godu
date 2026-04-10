package scan

import "fmt"

type Result struct {
	RootPath             string
	TotalSize            int64
	FilesCount           int64
	DirectoriesCount     int64
	SymlinksSkippedCount int64
	ErrorsCount          int64
	OtherCount           int64
}

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

func (r *Result) AddFile(size int64) {
	r.TotalSize += size
	r.FilesCount++
}

func (r *Result) AddDirectory() {
	r.DirectoriesCount++
}

func (r *Result) AddSymlinkSkipped() {
	r.SymlinksSkippedCount++
}

func (r *Result) AddError() {
	r.ErrorsCount++
}

func (r *Result) AddOtherCount() {
	r.OtherCount++
}

func (r *Result) AddFromResult(other *Result) {
	r.TotalSize += other.TotalSize
	r.FilesCount += other.FilesCount
	r.DirectoriesCount += other.DirectoriesCount
	r.SymlinksSkippedCount += other.SymlinksSkippedCount
	r.ErrorsCount += other.ErrorsCount
	r.OtherCount += other.OtherCount
}

func (r *Result) ToGB() float64 {
	return float64(r.TotalSize) / 1024 / 1024 / 1024
}

func (r *Result) String() string {
	return fmt.Sprintf("RootPath: %s, TotalSize: %d (%.2f GB), FilesCount: %d, DirectoriesCount: %d, SymlinksSkippedCount: %d, ErrorsCount: %d, OtherCount: %d", r.RootPath, r.TotalSize, r.ToGB(), r.FilesCount, r.DirectoriesCount, r.SymlinksSkippedCount, r.ErrorsCount, r.OtherCount)
}

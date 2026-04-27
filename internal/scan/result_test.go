package scan

import (
	"testing"
)

func TestResultAddFile(t *testing.T) {
	result := InitResult("test")
	result.AddFile(100)
	if result.TotalSize != 100 {
		t.Fatalf("expected total size to be 100, got %d", result.TotalSize)
	}
	if result.FilesCount != 1 {
		t.Fatalf("expected files count to be 1, got %d", result.FilesCount)
	}
}

func TestResultAddDirectory(t *testing.T) {
	result := InitResult("test")
	result.AddDirectory()
	if result.DirectoriesCount != 1 {
		t.Fatalf("expected directories count to be 1, got %d", result.DirectoriesCount)
	}
}

func TestResultAddSymlinkSkipped(t *testing.T) {
	result := InitResult("test")
	result.AddSymlinkSkipped()
	if result.SymlinksSkippedCount != 1 {
		t.Fatalf("expected symlinks skipped count to be 1, got %d", result.SymlinksSkippedCount)
	}
}

func TestResultAddError(t *testing.T) {
	result := InitResult("test")
	result.AddError()
	if result.ErrorsCount != 1 {
		t.Fatalf("expected errors count to be 1, got %d", result.ErrorsCount)
	}
}

func TestResultAddOtherCount(t *testing.T) {
	result := InitResult("test")
	result.AddOtherCount()
	if result.OtherCount != 1 {
		t.Fatalf("expected other count to be 1, got %d", result.OtherCount)
	}
}

func TestResultAddFromResult(t *testing.T) {
	result := InitResult("test")
	result.AddFile(100)
	result2 := InitResult("test2")
	result2.AddFile(100)
	result2.AddDirectory()
	result2.AddSymlinkSkipped()
	result2.AddError()
	result2.AddOtherCount()
	result.AddFromResult(result2)
	if result.TotalSize != 200 {
		t.Fatalf("expected total size to be 200, got %d", result.TotalSize)
	}
	if result.FilesCount != 2 {
		t.Fatalf("expected files count to be 2, got %d", result.FilesCount)
	}
	if result.DirectoriesCount != 1 {
		t.Fatalf("expected directories count to be 1, got %d", result.DirectoriesCount)
	}
	if result.SymlinksSkippedCount != 1 {
		t.Fatalf("expected symlinks skipped count to be 1, got %d", result.SymlinksSkippedCount)
	}
	if result.ErrorsCount != 1 {
		t.Fatalf("expected errors count to be 1, got %d", result.ErrorsCount)
	}
	if result.OtherCount != 1 {
		t.Fatalf("expected other count to be 1, got %d", result.OtherCount)
	}
}

func TestResultToGB(t *testing.T) {
	result := InitResult("test")
	result.AddFile(1024 * 1024 * 1024)
	if result.ToGB() != 1 {
		t.Fatalf("expected to gb to be 1, got %f", result.ToGB())
	}
}

func TestResultString(t *testing.T) {
	result := InitResult("test")
	result.AddFile(100)
	if result.String() != "RootPath: test, TotalSize: 100 (0.00 GB), FilesCount: 1, DirectoriesCount: 0, SymlinksSkippedCount: 0, ErrorsCount: 0, OtherCount: 0" {
		t.Fatalf("expected string to be 'test: 100 bytes', got %s", result.String())
	}
}

package scan

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}

func TestPrintVerbose(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	PrintVerbose(100, "test", buf)
	str := fmt.Sprintf("%-*d\t %s\n", 8, 100, "test")
	if buf.String() != str {
		t.Fatalf("expected print verbose to be '%s', got %s", str, buf.String())
	}
}

func TestScanRegularFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filePath := filepath.Join(root, "a.txt")
	content := "hello world"
	mustWriteFile(t, filePath, content)

	result, err := Scan(NewConfig(filePath, false, false, false))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.TotalSize != int64(len(content)) {
		t.Fatalf("expected total size %d, got %d", len(content), result.TotalSize)
	}
	if result.FilesCount != 1 {
		t.Fatalf("expected files count 1, got %d", result.FilesCount)
	}
	if result.DirectoriesCount != 0 {
		t.Fatalf("expected directories count 0, got %d", result.DirectoriesCount)
	}
	if result.SymlinksSkippedCount != 0 {
		t.Fatalf("expected symlinks skipped count 0, got %d", result.SymlinksSkippedCount)
	}
}

func TestScanDirectoryRecursively(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("create nested dir: %v", err)
	}

	contentA := "abc"
	contentB := "12345"
	mustWriteFile(t, filepath.Join(root, "a.txt"), contentA)
	mustWriteFile(t, filepath.Join(nested, "b.txt"), contentB)

	result, err := Scan(NewConfig(root, false, false, false))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedSize := int64(len(contentA) + len(contentB))
	if result.TotalSize != expectedSize {
		t.Fatalf("expected total size %d, got %d", expectedSize, result.TotalSize)
	}
	if result.FilesCount != 2 {
		t.Fatalf("expected files count 2, got %d", result.FilesCount)
	}
	// root and nested are both counted.
	if result.DirectoriesCount != 2 {
		t.Fatalf("expected directories count 2, got %d", result.DirectoriesCount)
	}
}

func TestScanSymlinkSkippedWhenDisabled(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	target := filepath.Join(root, "target.txt")
	mustWriteFile(t, target, "content")

	link := filepath.Join(root, "target-link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink not supported on this environment: %v", err)
	}

	result, err := Scan(NewConfig(link, false, false, false))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.SymlinksSkippedCount != 1 {
		t.Fatalf("expected symlinks skipped count 1, got %d", result.SymlinksSkippedCount)
	}
	if result.FilesCount != 0 {
		t.Fatalf("expected files count 0, got %d", result.FilesCount)
	}
	if result.TotalSize != 0 {
		t.Fatalf("expected total size 0, got %d", result.TotalSize)
	}
}

func TestScanFollowsSymlinkWhenEnabled(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	target := filepath.Join(root, "target.txt")
	content := "content"
	mustWriteFile(t, target, content)

	link := filepath.Join(root, "target-link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink not supported on this environment: %v", err)
	}

	result, err := Scan(NewConfig(link, true, false, false))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.SymlinksSkippedCount != 0 {
		t.Fatalf("expected symlinks skipped count 0, got %d", result.SymlinksSkippedCount)
	}
	if result.FilesCount != 1 {
		t.Fatalf("expected files count 1, got %d", result.FilesCount)
	}
	if result.TotalSize != int64(len(content)) {
		t.Fatalf("expected total size %d, got %d", len(content), result.TotalSize)
	}
}

func TestScanRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	_, err := Scan(NewConfig("", false, false, false))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid config") {
		t.Fatalf("expected invalid config error, got %v", err)
	}
}

func TestScanFailsForMissingPath(t *testing.T) {
	t.Parallel()

	_, err := Scan(NewConfig("/definitely/not/found/path", false, false, false))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to walk path") {
		t.Fatalf("expected walk path error, got %v", err)
	}
}

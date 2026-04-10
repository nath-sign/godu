# godu — Disk Usage Analyzer (v1 Specification)

## Overview

`godu` is a command-line disk usage analyzer written in Go. Its purpose is to scan a given filesystem path and compute the total size of all regular files contained within that path.

Version 1 focuses on correctness, robustness, and clean architecture. It deliberately avoids advanced features such as concurrency, tree views, duplicate detection, or interactive interfaces.

---

## Goals

### Primary Goals
- Accurately compute total disk usage for a given path
- Handle real-world filesystem edge cases safely
- Provide predictable and consistent behavior
- Establish a clean foundation for future features

### Non-Goals (v1)
- No concurrency
- No interactive UI (TUI)
- No duplicate file detection
- No advanced filtering (exclude patterns, depth limits)
- No persistent storage or caching

---

## CLI Specification

### Command

`godu scan <path>`


### Arguments

- `<path>`: Required. The filesystem path to scan.
  - Can be a file or directory
  - Must exist

### Behavior

- If `<path>` is a file:
  - Return the size of that file

- If `<path>` is a directory:
  - Recursively scan all contents
  - Aggregate sizes of all regular files

### Exit Codes

- `0`: Success (even if some files were skipped due to errors)
- `1`: Fatal error (invalid path, cannot access root, etc.)

---

## Filesystem Semantics

### Included

- Regular files:
  - Their logical size is added to the total

- Directories:
  - Traversed recursively
  - Do not contribute size directly

### Excluded

- Symbolic links:
  - Not followed
  - Not included in size
  - Counted separately

- Special files:
  - Sockets
  - Named pipes (FIFOs)
  - Device files
  - Ignored entirely

---

## Error Handling

### Principles

- The scan should continue whenever possible
- Errors should not abort the entire scan unless critical
- Errors should be recorded and reported

### Fatal Errors (terminate execution)

- Path does not exist
- Root path cannot be accessed (e.g. permission denied on root)

### Non-Fatal Errors (continue scanning)

- Cannot read a subdirectory
- Cannot stat a file
- Broken or inaccessible entries

### Error Reporting

- Total number of errors is reported in the final output
- Optional verbose reporting may be added later

---

## Symlink Policy

- Symbolic links are never followed
- They are counted separately as "skipped"
- Prevents:
  - Infinite recursion (cycles)
  - Double counting
  - Traversing outside the intended root

---

## Size Semantics

### Definition

- Total size is defined as the sum of:
  - `size` of all regular files encountered

### Notes

- This is logical file size, not disk block usage
- Sparse files are counted by logical size
- Directory metadata size is not included

---

## Scan Behavior

### High-Level Algorithm

1. Validate input path
2. Inspect root:
   - If file → return its size
   - If directory → begin traversal
3. For each entry:
   - Determine type
   - Apply rules:
     - File → add size
     - Directory → recurse
     - Symlink → skip
     - Other → ignore
4. Accumulate results
5. Return final aggregated result

---

## Data Model (Conceptual)

### Scan Result

The scan produces a result containing:

- Root path
- Total size in bytes
- Number of files
- Number of directories visited
- Number of symlinks skipped
- Number of errors encountered

### Notes

- Individual file entries are not stored in memory
- Aggregation is done incrementally during traversal
- Memory usage should remain low regardless of directory size

---

## Output Specification

### Default Output (Human-readable)

Example:
```
Path: /home/user/projects
Total Size: 3.4 GB
Files: 1245
Directories: 132
Symlinks Skipped: 12
Errors: 3
```


### Formatting Rules

- Sizes should be human-readable (KB, MB, GB)
- Counts are displayed as integers
- Output is plain text

---

## Performance Expectations

### v1 Priorities

- Correctness over speed
- Stability on large directory trees
- Reasonable memory usage

### Known Limitations

- Sequential traversal (no concurrency)
- Performance limited by filesystem I/O
- Not optimized for extremely large datasets

---

## Constraints and Assumptions

- Designed primarily for Unix-like systems
- Should work cross-platform where possible
- Does not require elevated privileges
- Assumes standard filesystem semantics

---

## Testing Strategy

### Core Test Cases

#### Basic
- Empty directory
- Single file
- Nested directories

#### Edge Cases
- Non-existent path
- File as root input
- Symlinks (file and directory)
- Broken symlink

#### Error Handling
- Unreadable directory
- Partial scan with errors

#### Correctness
- Known file sizes
- Mixed content structures

---

## Project Structure (Initial)
```
cmd/godu/ # CLI entrypoint
internal/scan/ # scanning logic
internal/format/ # output formatting
```


### Responsibilities

- `cmd/godu`:
  - Argument parsing
  - Command dispatch
  - Output printing

- `internal/scan`:
  - Filesystem traversal
  - Aggregation
  - Error handling

- `internal/format`:
  - Human-readable formatting
  - Size conversion

---

## Future Extensions (Not in v1)

- Concurrent scanning
- Top largest files/directories
- Tree view output
- Duplicate file detection
- JSON/CSV export
- Exclude patterns
- Depth limiting
- Interactive TUI

---

## Definition of Done (v1)

The implementation is considered complete when:

- A path can be scanned reliably
- Total size is correct
- Symlinks are skipped safely
- Errors are handled without crashing
- Output is clear and consistent
- All defined test cases pass

---

## Summary

Version 1 of `godu` is a focused, reliable disk usage scanner. It establishes a strong architectural foundation and clear behavior model, enabling future expansion into a full-featured analysis tool.
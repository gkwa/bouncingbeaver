package testutil

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("test.update", false, "update golden files")

// Golden returns the golden file content. If the `test.update` flag is specified,
// it updates the file with the current output and returns it.
func Golden(t *testing.T, actual []byte, filename string) []byte {
	t.Helper()

	golden := filepath.Join("testdata", filename)

	if *update {
		dir := filepath.Dir(golden)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(golden, actual, 0o644); err != nil {
			t.Fatalf("Failed to update golden file %s: %v", golden, err)
		}
	}

	expected, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("Failed to read golden file %s: %v", golden, err)
	}

	return expected
}

// GoldenSafe is like Golden but creates the file if it doesn't exist without failing
func GoldenSafe(t *testing.T, actual []byte, filename string) ([]byte, error) {
	t.Helper()

	golden := filepath.Join("testdata", filename)

	if *update {
		dir := filepath.Dir(golden)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(golden, actual, 0o644); err != nil {
			return nil, fmt.Errorf("failed to update golden file %s: %v", golden, err)
		}
		return actual, nil
	}

	// Check if file exists
	if _, err := os.Stat(golden); os.IsNotExist(err) {
		// Create the directory and file
		dir := filepath.Dir(golden)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(golden, actual, 0o644); err != nil {
			return nil, fmt.Errorf("failed to create golden file %s: %v", golden, err)
		}

		return nil, fmt.Errorf("created golden file %s", golden)
	}

	expected, err := os.ReadFile(golden)
	if err != nil {
		return nil, fmt.Errorf("failed to read golden file %s: %v", golden, err)
	}

	return expected, nil
}

// ToJSON converts any value to pretty-printed JSON bytes with trailing newline
func ToJSON(t *testing.T, v interface{}) []byte {
	t.Helper()

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	// Add trailing newline for consistency with text editors
	data = append(data, '\n')

	return data
}

// simpleDiff creates a simple line-by-line diff
func simpleDiff(expected, actual string) string {
	expectedLines := strings.Split(strings.TrimRight(expected, "\n"), "\n")
	actualLines := strings.Split(strings.TrimRight(actual, "\n"), "\n")

	var diff strings.Builder
	maxLines := len(expectedLines)
	if len(actualLines) > maxLines {
		maxLines = len(actualLines)
	}

	for i := 0; i < maxLines; i++ {
		var expectedLine, actualLine string

		if i < len(expectedLines) {
			expectedLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actualLine = actualLines[i]
		}

		if expectedLine != actualLine {
			if expectedLine != "" {
				diff.WriteString(fmt.Sprintf("- %s\n", expectedLine))
			}
			if actualLine != "" {
				diff.WriteString(fmt.Sprintf("+ %s\n", actualLine))
			}
		}
	}

	return diff.String()
}

// AssertGoldenMatch provides clear diff output for golden file mismatches
func AssertGoldenMatch(t *testing.T, actual, expected []byte, filename string) {
	t.Helper()

	diffOutput := simpleDiff(string(expected), string(actual))
	if diffOutput == "" {
		diffOutput = "Files differ but no line differences detected"
	}

	t.Errorf("Golden file mismatch for %s\n\nDifferences:\n%s\nTo update golden file, run:\n  go test -test.update",
		filename, diffOutput)
}

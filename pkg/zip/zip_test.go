package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestCompress_Directory(t *testing.T) {
	inputPath, err := filepath.Abs("test_data")
	if err != nil {
		t.Fatal(err)
	}

	outputPath := filepath.Join(t.TempDir(), "output.zip")

	err = Compress(inputPath, outputPath)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	r, err := zip.OpenReader(outputPath)
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	defer r.Close()

	expected := []string{"1.txt", "2.txt"}
	var got []string
	for _, f := range r.File {
		got = append(got, f.Name)
	}
	sort.Strings(got)
	sort.Strings(expected)

	if len(got) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(got), got)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("expected file %q, got %q", expected[i], got[i])
		}
	}
}

func TestCompress_Directory_FileContents(t *testing.T) {
	inputPath, err := filepath.Abs("test_data")
	if err != nil {
		t.Fatal(err)
	}

	outputPath := filepath.Join(t.TempDir(), "output.zip")

	err = Compress(inputPath, outputPath)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	r, err := zip.OpenReader(outputPath)
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("failed to open %s in zip: %v", f.Name, err)
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("failed to read %s: %v", f.Name, err)
		}

		originalPath := filepath.Join(inputPath, f.Name)
		originalContent, err := os.ReadFile(originalPath)
		if err != nil {
			t.Fatalf("failed to read original %s: %v", originalPath, err)
		}

		if string(content) != string(originalContent) {
			t.Errorf("content mismatch for %s: got %q, want %q", f.Name, content, originalContent)
		}
	}
}

func TestCompress_InvalidInput(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "output.zip")
	err := Compress("/nonexistent/path", outputPath)
	if err == nil {
		t.Error("expected error for nonexistent input path, got nil")
	}
}

func TestCompress_InvalidOutput(t *testing.T) {
	inputPath, err := filepath.Abs("test_data")
	if err != nil {
		t.Fatal(err)
	}

	err = Compress(inputPath, "/nonexistent/dir/output.zip")
	if err == nil {
		t.Error("expected error for invalid output path, got nil")
	}
}

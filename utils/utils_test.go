package utils

import (
	"bytes"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
)

/* Test Utils */

// Test DeleteEqual utility function
func TestDeleteEqual(t *testing.T) {
	n := 100
	s := make([]uint, 0)
	for range n {
		s = append(s, rand.UintN(10))
	}

	el := s[0]
	for _, val := range DeleteEqual(s, el) {
		if val == el {
			t.Fatal("Didn't remove all elements from slice")
		}
	}
}

// Test FileExists utility function
func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test")

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create a file: %s", err)
	}
	file.Close()

	if !FileExists(path) {
		t.Fatal("Created file doesn't exist")
	}
}

// Test DirExists utility function
func TestDirExists(t *testing.T) {
	dir := t.TempDir()
	if !DirExists(dir) {
		t.Fatal("Created directory doesn't exist")
	}
}

// Test TempFilePath utility function
func TestTempFilePath(t *testing.T) {
	dir := t.TempDir()
	path, err := TempFilePath(dir)
	if err != nil {
		t.Fatalf("Failed to create a temporary file")
	}

	if !FileExists(path) {
		t.Fatal("Created temporary file doesn't exist")
	}
}

// Test Relink utility function
func TestRelink(t *testing.T) {
	dir := t.TempDir()
	pathFrom := filepath.Join(dir, "from")
	pathTo := filepath.Join(dir, "to")

	err := os.WriteFile(pathFrom, []byte("Link"), 0644)
	if err != nil {
		t.Fatalf("Failed to a file: %s", err)
	}

	Relink(pathTo, pathFrom)

	buf, err := os.ReadFile(pathTo)
	if err != nil {
		t.Fatalf("Failed to read a linked file")
	}

	if !bytes.Equal([]byte("Link"), buf) {
		t.Fatal("Linked files don't have the same content")
	}
}

// Test GetHardLinkID utility function
func TestGetHardLinkID(t *testing.T) {
	dir := t.TempDir()
	pathFrom := filepath.Join(dir, "from")
	pathTo := filepath.Join(dir, "to")

	err := os.WriteFile(pathFrom, []byte("Link"), 0644)
	if err != nil {
		t.Fatalf("Failed to a file: %s", err)
	}

	Relink(pathTo, pathFrom)

	buf, err := os.ReadFile(pathTo)
	if err != nil {
		t.Fatalf("Failed to read a linked file")
	}

	if !bytes.Equal([]byte("Link"), buf) {
		t.Fatal("Linked files don't have the same content")
	}

	HLIDFrom, err := GetHardLinkID(pathFrom)
	if err != nil {
		t.Fatalf("Failed to get hardlink ID for the original file: %s", err)
	}

	HLIDTo, err := GetHardLinkID(pathTo)
	if err != nil {
		t.Fatalf("Failed to get hardlink ID for the linked file: %s", err)
	}

	if HLIDFrom != HLIDTo {
		t.Fatal("Hardlink ID of the original and linked files are not equal")
	}
}

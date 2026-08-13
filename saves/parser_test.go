package saves

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

var (
	save1 = Save1{}
	save2 = Save2{}
)

// Test save file's bytes parsing into Save object
func TestParseSaveBytesChapter1(t *testing.T) {
	save1bytes, err := Save2Bytes(&save1)
	if err != nil {
		t.Fatalf("Failed to parse chapter 1 save object to bytes: %s", err)
	}

	parsedSave, err := ParseSaveBytes(save1bytes, 1)
	if err != nil {
		t.Fatalf("Failed to parse chapter 1 save bytes to save object: %s", err)
	}

	saveRes, ok := parsedSave.(*Save1)
	if !ok {
		t.Fatalf("Failed to convert save object to chapter 1 save object")
	}

	if *saveRes != save1 {
		t.Fatalf("Resulting save object isn't equal to the parsed one")
	}
}

// Test save file's bytes parsing into Save object
func TestParseSaveBytesChapter2(t *testing.T) {
	save2bytes, err := Save2Bytes(&save2)
	if err != nil {
		t.Fatalf("Failed to parse chapter 2 save object to bytes: %s", err)
	}

	parsedSave, err := ParseSaveBytes(save2bytes, 2)
	if err != nil {
		t.Fatalf("Failed to parse chapter 2 save bytes to save object: %s", err)
	}

	saveRes, ok := parsedSave.(*Save2)
	if !ok {
		t.Fatalf("Failed to convert save object to chapter 2 save object")
	}

	if *saveRes != save2 {
		t.Fatalf("Resulting save object isn't equal to the parsed one")
	}
}

// Test save file's path parsing into Save object
func TestLoadSaveChapter1(t *testing.T) {
	save1bytes, err := Save2Bytes(&save1)
	if err != nil {
		t.Fatalf("Failed to parse chapter 1 save object to bytes: %s", err)
	}

	path := filepath.Join(t.TempDir(), "testloadsave")
	err = os.WriteFile(path, save1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write chapter 1 save bytes to file: %s", err)
	}

	parsedSave, err := LoadSave(path, 1)
	if err != nil {
		t.Fatalf("Failed to parse chapter 1 save bytes to save object: %s", err)
	}

	saveRes, ok := parsedSave.(*Save1)
	if !ok {
		t.Fatalf("Failed to convert save object to chapter 1 save object")
	}

	if *saveRes != save1 {
		t.Fatalf("Resulting save object isn't equal to the parsed one")
	}
}

// Test save file's path parsing into Save object
func TestLoadSaveChapter2(t *testing.T) {
	save2bytes, err := Save2Bytes(&save2)
	if err != nil {
		t.Fatalf("Failed to parse chapter 2 save object to bytes: %s", err)
	}

	path := filepath.Join(t.TempDir(), "testloadsave")
	err = os.WriteFile(path, save2bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write chapter 2 save bytes to file: %s", err)
	}

	parsedSave, err := LoadSave(path, 2)
	if err != nil {
		t.Fatalf("Failed to parse chapter 2 save bytes to save object: %s", err)
	}

	saveRes, ok := parsedSave.(*Save2)
	if !ok {
		t.Fatalf("Failed to convert save object to chapter 2 save object")
	}

	if *saveRes != save2 {
		t.Fatalf("Resulting save object isn't equal to the parsed one")
	}
}

// Test save file's Reader object parsing into Save object
func TestParseSaveReaderChapter1(t *testing.T) {
	save1bytes, err := Save2Bytes(&save1)
	if err != nil {
		t.Fatalf("Failed to parse chapter 1 save object to bytes: %s", err)
	}

	buf := bytes.NewBuffer(save1bytes)
	parsedSave, err := ParseSaveReader(buf, 1)
	if err != nil {
		t.Fatalf("Failed to parse chapter 1 save bytes to save object: %s", err)
	}

	saveRes, ok := parsedSave.(*Save1)
	if !ok {
		t.Fatalf("Failed to convert save object to chapter 1 save object")
	}

	if *saveRes != save1 {
		t.Fatalf("Resulting save object isn't equal to the parsed one")
	}
}

// Test save file's Reader object parsing into Save object
func TestParseSaveReaderChapter2(t *testing.T) {
	save2bytes, err := Save2Bytes(&save2)
	if err != nil {
		t.Fatalf("Failed to parse chapter 2 save object to bytes: %s", err)
	}

	buf := bytes.NewBuffer(save2bytes)
	parsedSave, err := ParseSaveReader(buf, 2)
	if err != nil {
		t.Fatalf("Failed to parse chapter 2 save bytes to save object: %s", err)
	}

	saveRes, ok := parsedSave.(*Save2)
	if !ok {
		t.Fatalf("Failed to convert save object to chapter 2 save object")
	}

	if *saveRes != save2 {
		t.Fatalf("Resulting save object isn't equal to the parsed one")
	}
}

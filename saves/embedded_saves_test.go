package saves

import (
	"bytes"
	"testing"
)

// Test chapter 1 embedded bytes on validity using parser
func TestChapter1EmbeddedBytes(t *testing.T) {
	_, err := ParseSaveBytes(chapter1bytes, 1)
	if err != nil {
		t.Fatalf("Chapter 1 example save file is invalid: %s", err)
	}
}

// Test chapter 2 embedded bytes on validity using parser
func TestChapter2EmbeddedBytes(t *testing.T) {
	_, err := ParseSaveBytes(chapter2bytes, 2)
	if err != nil {
		t.Fatalf("Chapter 2 example save file is invalid: %s", err)
	}
}

// Test chapter 3 embedded bytes on validity using parser
func TestChapter3EmbeddedBytes(t *testing.T) {
	_, err := ParseSaveBytes(chapter3bytes, 3)
	if err != nil {
		t.Fatalf("Chapter 3 example save file is invalid: %s", err)
	}
}

// Test chapter 4 embedded bytes on validity using parser
func TestChapter4EmbeddedBytes(t *testing.T) {
	_, err := ParseSaveBytes(chapter4bytes, 4)
	if err != nil {
		t.Fatalf("Chapter 4 example save file is invalid: %s", err)
	}
}

// Test chapter 5 embedded bytes on validity using parser
func TestChapter5EmbeddedBytes(t *testing.T) {
	_, err := ParseSaveBytes(chapter5bytes, 5)
	if err != nil {
		t.Fatalf("Chapter 5 example save file is invalid: %s", err)
	}
}

// Test getExampleSaveBytesForChapter function
func TestGetExampleSaveBytesForChapter(t *testing.T) {
	embeddedFiles := [5][]byte{
		chapter1bytes,
		chapter2bytes,
		chapter3bytes,
		chapter4bytes,
		chapter5bytes,
	}
	for i := range 5 {
		saveBytes, err := GetExampleSaveBytesForChapter(i + 1)
		if err != nil {
			t.Errorf("Failed to receive bytes of chapter %d example save file bytes: %s", i+1, err)
			continue
		}
		if !bytes.Equal(saveBytes, embeddedFiles[i]) {
			t.Errorf("Received chapter %d example save file bytes are invalid", i+1)
		}
	}
}

package saves

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* Test Config */

var defaultDrContent = []byte(`
[G1]
Name="[EMPTY]"
[G2_9]
Name="[EMPTY]"
[G2_2]
Name="G2_2"
[G2_5]
Name="G2_5"
[G5_1]
Name="G5_1"
[G0]
Name="[EMPTY]"
`)

/* Utils */

func createDrINI(t *testing.T, content ...[]byte) (DrINI, error) {
	if len(content) > 1 {
		return DrINI{}, errors.New("Too many arguments")
	}
	if len(content) == 0 {
		content = append(content, defaultDrContent)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "dr.ini")
	if err := os.WriteFile(path, content[0], 0644); err != nil {
		return DrINI{}, err
	}
	return NewDrINI(path)
}

/* Utils Tests */

func TestChapterIndexToKey(t *testing.T) {
	for _, i := range []int{0, 1, 2, 3, 4, 5, 9} {
		if chapterIndexToKey(1, i) != fmt.Sprintf("G%d", i) {
			t.Fatalf("Wrong selector name for chapter 1 index %d", i)
		}
	}
	for _, i := range []int{2, 3, 4, 5, 6, 7} {
		for _, j := range []int{0, 1, 2, 3, 4, 5, 9} {
			if chapterIndexToKey(i, j) != fmt.Sprintf("G%d_%d", i, j) {
				t.Fatalf("Wrong selector name for chapter %d index %d", i, j)
			}
		}
	}
}

/* Dr Tests */

// Load data from dr.ini
func TestLoadDrIni(t *testing.T) {
	if _, err := createDrINI(t); err != nil {
		t.Fatalf("Failed to load a dr.ini file: %s", err)
	}
}

// Clean empty save slots
func TestCleanDrIni(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load a dr.ini file: %s", err)
	}
	dr.Clean()

	a := dr.Sections()
	b := []string{"G2_2", "G2_5", "G5_1"}

	if len(a) != len(b) {
		t.Fatalf("Wrong Dr contents after cleaning %s", dr.Sections())
	}

	freq := make(map[string]int)
	for _, v := range a {
		freq[v]++
	}

	for _, v := range b {
		if freq[v] == 0 {
			t.Fatalf("Wrong Dr contents after cleaning %s", dr.Sections())
		}
		freq[v]--
	}
}

// Get a save slot section from dr.ini
func TestDrIniGetSlot(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load a dr.ini file: %s", err)
	}

	slot, ok := dr.GetSlot(2, 2)
	if !ok || slot == nil {
		t.Fatal("Failed to get a save slot")
	}

	if val, ok := (*slot)["Name"]; !ok || !strings.Contains(val, "G2_2") {
		t.Fatal("Returned save slot is invalid")
	}
}

// Get all save slots sections from dr.ini
func TestDrIniGetSlots(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load a dr.ini file: %s", err)
	}

	slots := dr.GetSlots()
	expected := map[string]map[string]string{
		"G1": {
			"Name": `"[EMPTY]"`,
		},
		"G2_9": {
			"Name": `"[EMPTY]"`,
		},
		"G2_2": {
			"Name": `"G2_2"`,
		},
		"G2_5": {
			"Name": `"G2_5"`,
		},
		"G5_1": {
			"Name": `"G5_1"`,
		},
		"G0": {
			"Name": `"[EMPTY]"`,
		},
	}
	res := make(map[string]map[string]string)
	for label, slot := range slots {
		res[label] = *slot
	}

	if !maps.EqualFunc(expected, res, func(v1, v2 map[string]string) bool {
		return maps.Equal(v1, v2)
	}) {
		t.Fatal("Returned slots are invalid")
	}
}

// Set a save slot to dr.ini
func TestDrIniSetSlot(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load dr.ini: %v", err)
	}

	// Create a new SlotDr (map) with necessary fields.
	slot := SlotDr{
		"Chapter": "2",
		"Name":    `"NewSlot"`,
		"Time":    `"123.45"`,
		"Room":    `"678.9"`,
	}

	// Set a new slot at chapter 2, index 3 (should succeed).
	ok := dr.SetSlot(slot, 3)
	if !ok {
		t.Fatal("SetSlot returned false for a new slot")
	}

	// Verify the slot exists in memory.
	got, ok := dr.GetSlot(2, 3)
	if !ok {
		t.Fatal("GetSlot failed after SetSlot")
	}
	if !maps.Equal(map[string]string(*got), map[string]string(slot)) {
		t.Errorf("Slot content mismatch: got %v, want %v", *got, slot)
	}

	// Try to set an existing slot (chapter 2, index 2 already exists).
	ok = dr.SetSlot(slot, 2)
	if ok {
		t.Fatal("SetSlot succeeded when slot already exists")
	}
}

// Set a save slot from a save object to dr.ini
func TestDrIniSetSlotFromSave(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load dr.ini: %v", err)
	}

	save1 := &Save1{
		PlayerName: "Kris",
		Time:       42.0,
		Room:       1.5,
	}
	save2 := &Save2{
		PlayerName: "Susie",
		Time:       99.9,
		Room:       2.7,
	}

	// 1. Set from Save1 at a free slot (chapter 2, index 3).
	ok := dr.SetSlotFromSave(save1, 2, 3, false)
	if !ok {
		t.Fatal("SetSlotFromSave failed for new slot")
	}
	got, ok := dr.GetSlot(2, 3)
	if !ok {
		t.Fatal("Slot not found after SetSlotFromSave")
	}
	expected := map[string]string{
		"Name": `"Kris"`,
		"Time": `"42.000000"`,
		"Room": `"1.500000"`,
	}
	if !maps.Equal(map[string]string(*got), expected) {
		t.Errorf("Slot content mismatch: got %v, want %v", *got, expected)
	}

	// 2. Set from Save2 at an existing slot with replace=false (should fail).
	ok = dr.SetSlotFromSave(save2, 2, 2, false) // slot G2_2 exists
	if ok {
		t.Fatal("SetSlotFromSave with replace=false succeeded on existing slot")
	}

	// 3. Set from Save2 at an existing slot with replace=true (should overwrite).
	ok = dr.SetSlotFromSave(save2, 2, 2, true)
	if !ok {
		t.Fatal("SetSlotFromSave with replace=true failed")
	}
	got, ok = dr.GetSlot(2, 2)
	if !ok {
		t.Fatal("Slot not found after overwrite")
	}
	expected = map[string]string{
		"Name": `"Susie"`,
		"Time": `"99.900000"`,
		"Room": `"2.700000"`,
	}
	if !maps.Equal(map[string]string(*got), expected) {
		t.Errorf("Overwritten slot content mismatch: got %v, want %v", *got, expected)
	}
}

// Copy a save slot in dr.ini
func TestDrIniCopySlot(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load dr.ini: %v", err)
	}

	// Copy from G2_2 (exists) to G2_3 (free).
	ok := dr.CopySlot(2, 2, 3)
	if !ok {
		t.Fatal("CopySlot failed when destination is free")
	}

	// Verify both source and destination exist and have same content.
	src, _ := dr.GetSlot(2, 2)
	dst, _ := dr.GetSlot(2, 3)
	if !maps.Equal(map[string]string(*src), map[string]string(*dst)) {
		t.Errorf("Copied slot content differs: src=%v, dst=%v", *src, *dst)
	}

	// Try to copy to an existing destination (should fail).
	ok = dr.CopySlot(2, 3, 2)
	if ok {
		t.Fatal("CopySlot succeeded when destination exists")
	}
}

// Move a save slot in dr.ini
func TestDrIniMoveSlot(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load dr.ini: %v", err)
	}

	// Move G5_1 to G5_2 (free).
	ok := dr.MoveSlot(5, 1, 2)
	if !ok {
		t.Fatal("MoveSlot failed")
	}

	// Check source is removed, destination has the content.
	_, ok = dr.GetSlot(5, 1)
	if ok {
		t.Error("Source slot still exists after move")
	}
	dst, ok := dr.GetSlot(5, 2)
	if !ok {
		t.Fatal("Destination slot not found after move")
	}
	expected := map[string]string{"Name": `"G5_1"`} // original content
	if !maps.Equal(map[string]string(*dst), expected) {
		t.Errorf("Moved content mismatch: got %v, want %v", *dst, expected)
	}

	// Try to move to an existing destination (should fail).
	ok = dr.MoveSlot(5, 2, 2)
	if ok {
		t.Fatal("MoveSlot succeeded when destination exists")
	}
}

// Swap save slots in dr.ini
func TestDrIniSwapSlots(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load dr.ini: %v", err)
	}

	// Swap G2_2 and G2_5
	ok := dr.SwapSlots(2, 2, 5)
	if !ok {
		t.Fatal("SwapSlots failed")
	}

	// Verify swapped contents.
	s1, _ := dr.GetSlot(2, 2)
	s2, _ := dr.GetSlot(2, 5)
	if !maps.Equal(map[string]string(*s1), map[string]string{"Name": `"G2_5"`}) {
		t.Errorf("After swap, G2_2 content is wrong: got %v", *s1)
	}
	if !maps.Equal(map[string]string(*s2), map[string]string{"Name": `"G2_2"`}) {
		t.Errorf("After swap, G2_5 content is wrong: got %v", *s2)
	}

	// Try to swap with a non-existing slot (should fail).
	ok = dr.SwapSlots(2, 2, 0) // chapter 3 index 0 doesn't exist
	if ok {
		t.Fatal("SwapSlots succeeded with non-existing slot")
	}
}

// Remove a save slot in dr.ini
func TestDrIniRemoveSlot(t *testing.T) {
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load dr.ini: %v", err)
	}

	// Remove G2_2.
	dr.RemoveSlot(2, 2)
	_, ok := dr.GetSlot(2, 2)
	if ok {
		t.Error("Slot still exists after RemoveSlot")
	}
}

// Write Dr content to dr.ini file
func TestDrIniWrite(t *testing.T) {
	// We'll create a fresh DrINI, modify it, and then read the file back.
	dr, err := createDrINI(t)
	if err != nil {
		t.Fatalf("Failed to load dr.ini: %v", err)
	}

	// Add a new slot.
	slot := SlotDr{
		"Chapter": "3",
		"Name":    `"New"`,
		"Time":    `"1.0"`,
		"Room":    `"2.0"`,
	}
	ok := dr.SetSlot(slot, 1)
	if !ok {
		t.Fatal("SetSlot failed")
	}
	// Write is called inside SetSlot? No, SetSlot does not call Write, so we must call it explicitly.
	// Actually SetSlot does not write; only Copy/Move/Swap/Remove do. So we call Write.
	err = dr.Write()
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Reload the file and verify the new slot exists.
	reloaded, err := NewDrINI(dr.Path)
	if err != nil {
		t.Fatalf("Failed to reload dr.ini: %v", err)
	}
	got, ok := reloaded.GetSlot(3, 1)
	if !ok {
		t.Fatal("New slot not found after reload")
	}
	expected := map[string]string{
		"Chapter": `3`,
		"Name":    `"New"`,
		"Time":    `"1.0"`,
		"Room":    `"2.0"`,
	}
	if !maps.Equal(map[string]string(*got), expected) {
		t.Errorf("Reloaded slot content mismatch: got %v, want %v", *got, expected)
	}
}

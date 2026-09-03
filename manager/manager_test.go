package manager

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/LammoGit/Deltarune-Save-Manager/saves"
	"github.com/LammoGit/Deltarune-Save-Manager/utils"
)

// Test SaveID string conversion
func TestSaveIDStringConversion(t *testing.T) {
	saveID := SaveID{"Save", 1, false}
	if saveID.String() != "1_a_Save" {
		t.Fatalf("Invalid SaveID to string conversion: %s", saveID.String())
	}
}

// Test SaveID side B string conversion
func TestSaveIDSideBStringConversion(t *testing.T) {
	saveID := SaveID{"Save", 5, true}
	if saveID.String() != "5_b_Save" {
		t.Fatalf("Invalid SaveID side B to string conversion: %s", saveID.String())
	}
}

// Test SlotID string conversion
func TestSlotIDStringConversion(t *testing.T) {
	slotID := SlotID{1, 2, false}
	if slotID.String() != "filech1_2" {
		t.Fatalf("Invalid SlotID to string conversion: %s", slotID.String())
	}
}

// Test SlotID side B string conversion
func TestSlotIDSideBStringConversion(t *testing.T) {
	slotID := SlotID{1, 2, true}
	if slotID.String() != "filech1_2_b" {
		t.Fatalf("Invalid SlotID side B to string conversion: %s", slotID.String())
	}
}

// Test loadSaves function
func TestLoadSaves(t *testing.T) {
	savesPath := t.TempDir()

	chapter1bytes, err := saves.GetExampleSaveBytesForChapter(1)
	if err != nil {
		t.Fatalf("Failed to get bytes for example chapter 1 save file: %s", err)
	}

	expectedSaveGeneral, err := saves.ParseSaveBytes(chapter1bytes, 1)
	if err != nil {
		t.Fatalf("Failed to parse save file bytes into save object: %s", err)
	}
	expectedSave, ok := expectedSaveGeneral.(*saves.Save1)
	if !ok {
		t.Fatal("Expected save object to be of type *Save1")
	}

	path1 := filepath.Join(savesPath, "1_a_Save1")
	path2 := filepath.Join(savesPath, "1_a_Save2")
	path3 := filepath.Join(savesPath, "1_a_Save3")

	err = os.WriteFile(path1, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write bytes to a file: %s", err)
	}
	err = os.WriteFile(path2, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write bytes to a file: %s", err)
	}
	err = os.Link(path2, path3)
	if err != nil {
		t.Fatalf("Failed to create a hard link: %s", err)
	}

	savesMap, saveLinks, err := loadSaves(savesPath)
	if err != nil {
		t.Fatalf("loadSaves returned error: %s", err)
	}

	// Check that all three saves are present and correct
	expectedIDs := []SaveID{
		{"Save1", 1, false},
		{"Save2", 1, false},
		{"Save3", 1, false},
	}
	for _, id := range expectedIDs {
		saveGeneral, ok := savesMap[id]
		if !ok {
			t.Errorf("SaveID %v not found in saves map", id)
			continue
		}
		save, ok := saveGeneral.(*saves.Save1)
		if !ok {
			t.Errorf("Save for %v is not *Save1", id)
			continue
		}
		if *save != *expectedSave {
			t.Errorf("Save data for %v does not match expected", id)
		}
	}

	// Check hardlink grouping: path2 and path3 should share a hardlink, path1 separate
	hardLinkID1, err := utils.GetHardLinkID(path1)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	hardLinkID2, err := utils.GetHardLinkID(path2)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	hardLinkID3, err := utils.GetHardLinkID(path3)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	if hardLinkID2 != hardLinkID3 {
		t.Error("Hardlink IDs for path2 and path3 differ, but they are linked")
	}
	if hardLinkID1 == hardLinkID2 {
		t.Error("Hardlink IDs for path1 and path2 are same, but they are not linked")
	}

	// Check saveLinks map
	ids1 := saveLinks[hardLinkID1]
	if len(ids1) != 1 || ids1[0] != (SaveID{"Save1", 1, false}) {
		t.Errorf("saveLinks for hardlink %v has unexpected values: %v", hardLinkID1, ids1)
	}
	ids2 := saveLinks[hardLinkID2]
	if len(ids2) != 2 {
		t.Errorf("saveLinks for hardlink %v should have 2 entries, got %d", hardLinkID2, len(ids2))
	} else {
		// Check that both Save2 and Save3 are present (order may vary)
		found2, found3 := false, false
		for _, id := range ids2 {
			if id == (SaveID{"Save2", 1, false}) {
				found2 = true
			}
			if id == (SaveID{"Save3", 1, false}) {
				found3 = true
			}
		}
		if !found2 || !found3 {
			t.Errorf("saveLinks for hardlink %v missing Save2 or Save3", hardLinkID2)
		}
	}
}

// Test loadSaves function with non-existant path
func TestLoadSavesWithInvalidPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "abc")
	_, _, err := loadSaves(path)
	if err == nil {
		t.Fatal("Loading saves from non-existant path didn't return an error")
	}
}

// Test loadSlots function
func TestLoadSlots(t *testing.T) {
	slotsPath := t.TempDir()

	chapter1bytes, err := saves.GetExampleSaveBytesForChapter(1)
	if err != nil {
		t.Fatalf("Failed to get bytes for example chapter 1 save file: %s", err)
	}

	expectedSaveGeneral, err := saves.ParseSaveBytes(chapter1bytes, 1)
	if err != nil {
		t.Fatalf("Failed to parse save file bytes into save object: %s", err)
	}
	expectedSave, ok := expectedSaveGeneral.(*saves.Save1)
	if !ok {
		t.Fatal("Expected save object to be of type *Save1")
	}

	path1 := filepath.Join(slotsPath, "filech1_1")
	path2 := filepath.Join(slotsPath, "filech1_2")
	path3 := filepath.Join(slotsPath, "filech1_3_b")

	err = os.WriteFile(path1, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write bytes to a file: %s", err)
	}
	err = os.WriteFile(path2, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write bytes to a file: %s", err)
	}
	err = os.Link(path2, path3)
	if err != nil {
		t.Fatalf("Failed to create a hard link: %s", err)
	}

	slots, slotLinks, err := loadSlots(slotsPath)
	if err != nil {
		t.Fatalf("loadSlots returned error: %s", err)
	}

	// Check all slots are present
	expectedSlotIDs := []SlotID{
		{1, 1, false},
		{1, 2, false},
		{1, 3, true},
	}
	for _, id := range expectedSlotIDs {
		saveGeneral, ok := slots[id]
		if !ok {
			t.Errorf("SlotID %v not found in slots map", id)
			continue
		}
		save, ok := saveGeneral.(*saves.Save1)
		if !ok {
			t.Errorf("Slot for %v is not *Save1", id)
			continue
		}
		if *save != *expectedSave {
			t.Errorf("Slot data for %v does not match expected", id)
		}
	}

	// Check hardlink grouping
	hardLinkID1, err := utils.GetHardLinkID(path1)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	hardLinkID2, err := utils.GetHardLinkID(path2)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	hardLinkID3, err := utils.GetHardLinkID(path3)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	if hardLinkID2 != hardLinkID3 {
		t.Error("Hardlink IDs for path2 and path3 differ, but they are linked")
	}
	if hardLinkID1 == hardLinkID2 {
		t.Error("Hardlink IDs for path1 and path2 are same, but they are not linked")
	}

	// Check slotLinks map
	ids1 := slotLinks[hardLinkID1]
	if len(ids1) != 1 || ids1[0] != (SlotID{1, 1, false}) {
		t.Errorf("slotLinks for hardlink %v has unexpected values: %v", hardLinkID1, ids1)
	}
	ids2 := slotLinks[hardLinkID2]
	if len(ids2) != 2 {
		t.Errorf("slotLinks for hardlink %v should have 2 entries, got %d", hardLinkID2, len(ids2))
	} else {
		found2, found3 := false, false
		for _, id := range ids2 {
			if id == (SlotID{1, 2, false}) {
				found2 = true
			}
			if id == (SlotID{1, 3, true}) {
				found3 = true
			}
		}
		if !found2 || !found3 {
			t.Errorf("slotLinks for hardlink %v missing Slot2 or Slot3", hardLinkID2)
		}
	}
}

// Test loadSlots function with non-existant path
func TestLoadSlotsWithInvalidPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "abc")
	_, _, err := loadSlots(path)
	if err == nil {
		t.Fatal("Loading slots from non-existant path didn't return an error")
	}
}

// Test creation of a new SaveManager object
func TestNewSaveManager(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()

	// Create some saves and slots to load
	chapter1bytes, err := saves.GetExampleSaveBytesForChapter(1)
	if err != nil {
		t.Fatalf("Failed to get bytes for example chapter 1 save file: %s", err)
	}
	savePath := filepath.Join(managerPath, "1_a_TestSave")
	err = os.WriteFile(savePath, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write save file: %s", err)
	}
	slotPath := filepath.Join(slotsPath, "filech1_1")
	err = os.WriteFile(slotPath, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write slot file: %s", err)
	}
	// Create dr.ini (can be empty)
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err = os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}

	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager returned error: %s", err)
	}
	if sm == nil {
		t.Fatal("NewSaveManager returned nil")
	}
	if sm.ManagerPath != managerPath {
		t.Errorf("ManagerPath mismatch: got %s, want %s", sm.ManagerPath, managerPath)
	}
	if sm.SlotsPath != slotsPath {
		t.Errorf("SlotsPath mismatch: got %s, want %s", sm.SlotsPath, slotsPath)
	}
	if len(sm.Saves) != 1 {
		t.Errorf("Expected 1 save, got %d", len(sm.Saves))
	}
	if len(sm.Slots) != 1 {
		t.Errorf("Expected 1 slot, got %d", len(sm.Slots))
	}
	if len(sm.SaveLinks) != 1 {
		t.Errorf("Expected 1 save hardlink group, got %d", len(sm.SaveLinks))
	}
	if len(sm.SlotLinks) != 1 {
		t.Errorf("Expected 1 slot hardlink group, got %d", len(sm.SlotLinks))
	}
	// Check dr.ini loaded
	if sm.Dr.Path == "" {
		t.Error("DrINI not loaded (Path is empty)")
	}
}

// Test creation of a new SaveManager object with non-existant save files' directory path
func TestNewSaveManagerWithInvalidSavesPath(t *testing.T) {
	slotsPath := t.TempDir()
	managerPath := filepath.Join(t.TempDir(), "nonexistent")
	_, err := NewSaveManager(managerPath, slotsPath)
	if err == nil {
		t.Fatal("NewSaveManager with invalid saves path did not return error")
	}
}

// Test creation of a new SaveManager object with non-existant slots' directory path
func TestNewSaveManagerWithInvalidSlotsPath(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := filepath.Join(t.TempDir(), "nonexistent")
	_, err := NewSaveManager(managerPath, slotsPath)
	if err == nil {
		t.Fatal("NewSaveManager with invalid slots path did not return error")
	}
}

// Test hardlink string representation from SaveID
func TestHardLinkIDFromSaveID(t *testing.T) {
	// Setup a manager with a save
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create a save
	err = sm.Create("TestSave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	saveID := SaveID{"TestSave", 1, false}
	hardLinkID, ok := sm.hardLinkIDFromSaveID(saveID)
	if !ok {
		t.Fatal("hardLinkIDFromSaveID returned false for existing save")
	}
	if hardLinkID == "" {
		t.Error("hardLinkIDFromSaveID returned empty string")
	}
	// Verify that the hardlink ID corresponds to the file
	path := filepath.Join(managerPath, saveID.String())
	expectedID, err := utils.GetHardLinkID(path)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	if hardLinkID != expectedID {
		t.Errorf("hardLinkID mismatch: got %s, want %s", hardLinkID, expectedID)
	}
}

// Test hardlink string representation from SaveID for non-existant save
func TestHardLinkIDFromInvalidSaveID(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	saveID := SaveID{"NonExistent", 1, false}
	_, ok := sm.hardLinkIDFromSaveID(saveID)
	if ok {
		t.Fatal("hardLinkIDFromSaveID returned true for non-existent save")
	}
}

// Test hardlink string representation from SlotID
func TestHardLinkIDFromSlotID(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create a slot by setting a save to it
	err = sm.Create("TestSave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	err = sm.SetSlot("TestSave", 1, 1, false, false)
	if err != nil {
		t.Fatalf("SetSlot failed: %s", err)
	}
	slotID := SlotID{1, 1, false}
	hardLinkID, ok := sm.hardLinkIDFromSlotID(slotID)
	if !ok {
		t.Fatal("hardLinkIDFromSlotID returned false for existing slot")
	}
	if hardLinkID == "" {
		t.Error("hardLinkIDFromSlotID returned empty string")
	}
	path := filepath.Join(slotsPath, slotID.String())
	expectedID, err := utils.GetHardLinkID(path)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	if hardLinkID != expectedID {
		t.Errorf("hardLinkID mismatch: got %s, want %s", hardLinkID, expectedID)
	}
}

// Test hardlink string representation from SlotID for non-existant slot
func TestHardLinkIDFromInvalidSlotID(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	slotID := SlotID{1, 99, false}
	_, ok := sm.hardLinkIDFromSlotID(slotID)
	if ok {
		t.Fatal("hardLinkIDFromSlotID returned true for non-existent slot")
	}
}

// Test new save file creation
func TestCreateSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}

	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	saveID := SaveID{"MySave", 1, false}
	path := filepath.Join(managerPath, saveID.String())
	if !utils.FileExists(path) {
		t.Error("Save file not created")
	}
	if _, ok := sm.Saves[saveID]; !ok {
		t.Error("Save not added to Saves map")
	}
	hardLinkID, err := utils.GetHardLinkID(path)
	if err != nil {
		t.Fatalf("GetHardLinkID failed: %s", err)
	}
	if ids, ok := sm.SaveLinks[hardLinkID]; !ok || len(ids) != 1 || ids[0] != saveID {
		t.Errorf("SaveLinks not updated correctly: %v", sm.SaveLinks[hardLinkID])
	}
}

// Test new save file creation with taken save name
func TestCreateSaveWithTakenSaveName(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	err = sm.Create("MySave", 1, false)
	if err == nil {
		t.Fatal("Create with taken name did not return error")
	}
}

// Test new save file creation with invalid chapter
func TestCreateSaveWithInvalidChapter(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Create("MySave", utils.MAXCHAPTER+1, false)
	if err == nil {
		t.Fatal("Create with invalid chapter did not return error")
	}
}

// Test swapping two save files
func TestSwapSaves(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Create("SaveA", 1, false)
	if err != nil {
		t.Fatalf("Create SaveA failed: %s", err)
	}
	err = sm.Create("SaveB", 1, false)
	if err != nil {
		t.Fatalf("Create SaveB failed: %s", err)
	}
	// Write distinct content to each file so we can verify swap
	pathA := filepath.Join(managerPath, SaveID{"SaveA", 1, false}.String())
	pathB := filepath.Join(managerPath, SaveID{"SaveB", 1, false}.String())
	err = os.WriteFile(pathA, []byte("contentA"), 0644)
	if err != nil {
		t.Fatalf("Write to SaveA failed: %s", err)
	}
	err = os.WriteFile(pathB, []byte("contentB"), 0644)
	if err != nil {
		t.Fatalf("Write to SaveB failed: %s", err)
	}
	// Reload saves to update objects? Actually we need to reload because we wrote directly.
	// Instead, we can just read the files after swap.
	err = sm.Swap("SaveA", "SaveB", 1, false)
	if err != nil {
		t.Fatalf("Swap failed: %s", err)
	}
	// Check file contents swapped
	contentA, err := os.ReadFile(pathA)
	if err != nil {
		t.Fatalf("Read SaveA failed: %s", err)
	}
	contentB, err := os.ReadFile(pathB)
	if err != nil {
		t.Fatalf("Read SaveB failed: %s", err)
	}
	if string(contentA) != "contentB" || string(contentB) != "contentA" {
		t.Errorf("Swap did not swap contents: A=%q, B=%q", contentA, contentB)
	}
	// Check Saves map and SaveLinks were updated (we can check hardlinks)
	// After swap, SaveA should now point to the file that was SaveB, and vice versa.
	// The hardlink IDs of the files are the same as before, but the mapping in SaveLinks should have swapped.
	// The hardlink IDs should be the same as before (since file contents changed but inode same).
	// We can check that SaveA's hardlink ID equals the old hardlink ID of SaveB, and vice versa.
	newHardLinkIDA, _ := utils.GetHardLinkID(pathA)
	newHardLinkIDB, _ := utils.GetHardLinkID(pathB)
	if newHardLinkIDA == newHardLinkIDB {
		t.Error("After swap, files should not be hardlinked")
	}
	// The save objects should have been swapped in Saves map; but we don't compare deep equality, just existence.
	// Verify SaveLinks: each hardlink should have the correct SaveID.
	idsA, ok := sm.SaveLinks[newHardLinkIDA]
	if !ok || len(idsA) != 1 || idsA[0] != (SaveID{"SaveA", 1, false}) {
		t.Errorf("SaveLinks for hardlink ID %v has incorrect saves: %v", newHardLinkIDA, idsA)
	}
	idsB, ok := sm.SaveLinks[newHardLinkIDB]
	if !ok || len(idsB) != 1 || idsB[0] != (SaveID{"SaveB", 1, false}) {
		t.Errorf("SaveLinks for hardlink ID %v has incorrect saves: %v", newHardLinkIDB, idsB)
	}
}

// Test swapping two save files with save files being linked
func TestSwapLinkedSaves(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create SaveA and SaveB as hardlinks
	err = sm.Create("SaveA", 1, false)
	if err != nil {
		t.Fatalf("Create SaveA failed: %s", err)
	}
	pathA := filepath.Join(managerPath, SaveID{"SaveA", 1, false}.String())
	pathB := filepath.Join(managerPath, SaveID{"SaveB", 1, false}.String())
	err = os.Link(pathA, pathB)
	if err != nil {
		t.Fatalf("Link failed: %s", err)
	}
	// Now load the manager again? Actually we need to refresh the SaveManager to recognize the new hardlink.
	// But we can just create a new manager from scratch.
	sm2, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Now SaveA and SaveB are hardlinked. Swap should be a no-op.
	err = sm2.Swap("SaveA", "SaveB", 1, false)
	if err != nil {
		t.Fatalf("Swap failed: %s", err)
	}
	// Check that files are still hardlinked and no changes
	hardLinkIDA, _ := utils.GetHardLinkID(pathA)
	hardLinkIDB, _ := utils.GetHardLinkID(pathB)
	if hardLinkIDA != hardLinkIDB {
		t.Error("Files are not hardlinked after swap (should remain linked)")
	}
	// Check SaveLinks: both SaveA and SaveB should be in the same group
	ids, ok := sm2.SaveLinks[hardLinkIDA]
	if !ok || len(ids) != 2 {
		t.Errorf("Expected 2 saves in hardlink group, got %d", len(ids))
	} else {
		foundA, foundB := false, false
		for _, id := range ids {
			if id == (SaveID{"SaveA", 1, false}) {
				foundA = true
			}
			if id == (SaveID{"SaveB", 1, false}) {
				foundB = true
			}
		}
		if !foundA || !foundB {
			t.Errorf("Hardlink group missing SaveA or SaveB")
		}
	}
}

// Test swapping two save files with first save being non-existant
func TestSwapSavesWithInvalidFirstSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Create("SaveB", 1, false)
	if err != nil {
		t.Fatalf("Create SaveB failed: %s", err)
	}
	err = sm.Swap("SaveA", "SaveB", 1, false)
	if err == nil {
		t.Fatal("Swap with non-existent SaveA did not return error")
	}
}

// Test swapping two save files with second save being non-existant
func TestSwapSavesWithInvalidSecondSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Create("SaveA", 1, false)
	if err != nil {
		t.Fatalf("Create SaveA failed: %s", err)
	}
	err = sm.Swap("SaveA", "SaveB", 1, false)
	if err == nil {
		t.Fatal("Swap with non-existent SaveB did not return error")
	}
}

// Test setting save file to a slot
func TestSetSlot(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create a save
	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	// Set to slot 1
	err = sm.SetSlot("MySave", 1, 1, false, false)
	if err != nil {
		t.Fatalf("SetSlot failed: %s", err)
	}
	slotID := SlotID{1, 1, false}
	slotPath := filepath.Join(slotsPath, slotID.String())
	if !utils.FileExists(slotPath) {
		t.Error("Slot file not created")
	}
	// Check that slot is hardlinked to save
	savePath := filepath.Join(managerPath, SaveID{"MySave", 1, false}.String())
	hardLinkSave, _ := utils.GetHardLinkID(savePath)
	hardLinkSlot, _ := utils.GetHardLinkID(slotPath)
	if hardLinkSave != hardLinkSlot {
		t.Errorf("Slot and save are not hardlinked: save=%s, slot=%s", hardLinkSave, hardLinkSlot)
	}
	// Check Slots map
	if _, ok := sm.Slots[slotID]; !ok {
		t.Error("Slot not added to Slots map")
	}
	// Check SlotLinks
	if ids, ok := sm.SlotLinks[hardLinkSave]; !ok || len(ids) != 1 || ids[0] != slotID {
		t.Errorf("SlotLinks not updated correctly: %v", sm.SlotLinks[hardLinkSave])
	}
	// Check dr.ini (we can just check that the file exists and has some content)
	drContent, err := os.ReadFile(drIniPath)
	if err != nil {
		t.Fatalf("Failed to read dr.ini: %s", err)
	}
	if len(drContent) == 0 {
		t.Error("dr.ini seems empty after SetSlot")
	}
}

// Test setting non-existant save file to a slot
func TestSetSlotInvalidSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.SetSlot("NonExistent", 1, 1, false, false)
	if err == nil {
		t.Fatal("SetSlot with non-existent save did not return error")
	}
}

// Test setting save file to a slot taken by managed save
func TestSetManagedSlot(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create two saves and set first to slot
	err = sm.Create("SaveA", 1, false)
	if err != nil {
		t.Fatalf("Create SaveA failed: %s", err)
	}
	err = sm.Create("SaveB", 1, false)
	if err != nil {
		t.Fatalf("Create SaveB failed: %s", err)
	}
	err = sm.SetSlot("SaveA", 1, 1, false, false)
	if err != nil {
		t.Fatalf("SetSlot SaveA failed: %s", err)
	}
	// Now set SaveB to same slot (should succeed)
	err = sm.SetSlot("SaveB", 1, 1, false, false)
	if err != nil {
		t.Fatalf("SetSlot SaveB failed: %s", err)
	}
	// Check that slot now points to SaveB
	slotPath := filepath.Join(slotsPath, SlotID{1, 1, false}.String())
	saveBPath := filepath.Join(managerPath, SaveID{"SaveB", 1, false}.String())
	hardLinkSlot, _ := utils.GetHardLinkID(slotPath)
	hardLinkSaveB, _ := utils.GetHardLinkID(saveBPath)
	if hardLinkSlot != hardLinkSaveB {
		t.Errorf("Slot not hardlinked to SaveB")
	}
	// Check that SaveA no longer has the slot in its group
	hardLinkSaveA, _ := utils.GetHardLinkID(filepath.Join(managerPath, SaveID{"SaveA", 1, false}.String()))
	if ids, ok := sm.SlotLinks[hardLinkSaveA]; ok && len(ids) > 0 {
		t.Errorf("SaveA's hardlink still has slots: %v", ids)
	}
	// Check SaveB's group has the slot
	if ids, ok := sm.SlotLinks[hardLinkSaveB]; !ok || len(ids) != 1 || ids[0] != (SlotID{1, 1, false}) {
		t.Errorf("SaveB's hardlink missing slot: %v", sm.SlotLinks[hardLinkSaveB])
	}
}

// Test setting save file to a slot taken by unmanaged save
func TestSetUnmanagedSlot(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	// Create an unmanaged slot: just write a file in slotsPath
	slotID := SlotID{1, 1, false}
	slotPath := filepath.Join(slotsPath, slotID.String())
	chapter1bytes, err := saves.GetExampleSaveBytesForChapter(1)
	if err != nil {
		t.Fatalf("Failed to get example bytes: %s", err)
	}
	err = os.WriteFile(slotPath, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write unmanaged slot: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create a save
	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	// Attempt to set slot (eraseUnmanaged=false) should error
	err = sm.SetSlot("MySave", 1, 1, false, false)
	if err == nil {
		t.Fatal("SetSlot with unmanaged slot and eraseUnmanaged=false did not return error")
	}
}

// Test setting save file to a slot taken by unmanaged save with eraseUnmanaged allowed
func TestSetUnmanageedSlotWithErasingPermission(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	// Create unmanaged slot
	slotID := SlotID{1, 1, false}
	slotPath := filepath.Join(slotsPath, slotID.String())
	chapter1bytes, err := saves.GetExampleSaveBytesForChapter(1)
	if err != nil {
		t.Fatalf("Failed to get example bytes: %s", err)
	}
	err = os.WriteFile(slotPath, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write unmanaged slot: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create a save
	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	// Set slot with eraseUnmanaged=true
	err = sm.SetSlot("MySave", 1, 1, false, true)
	if err != nil {
		t.Fatalf("SetSlot with eraseUnmanaged=true failed: %s", err)
	}
	// Check slot now points to save
	savePath := filepath.Join(managerPath, SaveID{"MySave", 1, false}.String())
	hardLinkSave, _ := utils.GetHardLinkID(savePath)
	hardLinkSlot, _ := utils.GetHardLinkID(slotPath)
	if hardLinkSave != hardLinkSlot {
		t.Errorf("Slot not hardlinked to save after set")
	}
}

// Test clearing a save slot taken by managed save
func TestUnsetSlot(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create and set slot
	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	err = sm.SetSlot("MySave", 1, 1, false, false)
	if err != nil {
		t.Fatalf("SetSlot failed: %s", err)
	}
	slotID := SlotID{1, 1, false}
	slotPath := filepath.Join(slotsPath, slotID.String())
	if !utils.FileExists(slotPath) {
		t.Fatal("Slot file does not exist before unset")
	}
	// Unset
	err = sm.UnsetSlot(1, 1, false)
	if err != nil {
		t.Fatalf("UnsetSlot failed: %s", err)
	}
	if utils.FileExists(slotPath) {
		t.Error("Slot file still exists after unset")
	}
	if _, ok := sm.Slots[slotID]; ok {
		t.Error("Slot still in Slots map")
	}
	// Check SlotLinks
	hardLinkSave, _ := utils.GetHardLinkID(filepath.Join(managerPath, SaveID{"MySave", 1, false}.String()))
	if ids, ok := sm.SlotLinks[hardLinkSave]; ok && len(ids) > 0 {
		t.Errorf("SlotLinks still contains slot after unset: %v", ids)
	}
}

// Test clearing an empty save slot
func TestUnsetEmptySlot(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Unset on empty slot should succeed (no-op)
	err = sm.UnsetSlot(1, 1, false)
	if err != nil {
		t.Fatalf("UnsetSlot on empty slot returned error: %s", err)
	}
	// Ensure no file created
	slotPath := filepath.Join(slotsPath, SlotID{1, 1, false}.String())
	if utils.FileExists(slotPath) {
		t.Error("Slot file was created unexpectedly")
	}
}

// Test clearing a save slot taken by unmanaged save
func TestUnsetUnmanagedSlot(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	// Create unmanaged slot
	slotID := SlotID{1, 1, false}
	slotPath := filepath.Join(slotsPath, slotID.String())
	chapter1bytes, err := saves.GetExampleSaveBytesForChapter(1)
	if err != nil {
		t.Fatalf("Failed to get example bytes: %s", err)
	}
	err = os.WriteFile(slotPath, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write unmanaged slot: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Unset with eraseUnmanaged=false should error
	err = sm.UnsetSlot(1, 1, false)
	if err == nil {
		t.Fatal("UnsetSlot on unmanaged slot with eraseUnmanaged=false did not error")
	}
}

// Test clearing a save slot taken by unmanaged save with eraseUnmanaged allowed
func TestUnsetUnmanagedSlotWithErasingPermission(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	// Create unmanaged slot
	slotID := SlotID{1, 1, false}
	slotPath := filepath.Join(slotsPath, slotID.String())
	chapter1bytes, err := saves.GetExampleSaveBytesForChapter(1)
	if err != nil {
		t.Fatalf("Failed to get example bytes: %s", err)
	}
	err = os.WriteFile(slotPath, chapter1bytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write unmanaged slot: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Unset with eraseUnmanaged=true
	err = sm.UnsetSlot(1, 1, true)
	if err != nil {
		t.Fatalf("UnsetSlot with eraseUnmanaged=true failed: %s", err)
	}
	if utils.FileExists(slotPath) {
		t.Error("Unmanaged slot file was not removed")
	}
}

// Test saving a save slot to a save file
func TestSaveSlot(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create a save and set it to slot
	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	err = sm.SetSlot("MySave", 1, 1, false, false)
	if err != nil {
		t.Fatalf("SetSlot failed: %s", err)
	}
	// Now call SaveSlot to create a managed save from the slot
	err = sm.SaveSlot("NewSave", 1, 1, false)
	if err != nil {
		t.Fatalf("SaveSlot failed: %s", err)
	}
	newSaveID := SaveID{"NewSave", 1, false}
	newSavePath := filepath.Join(managerPath, newSaveID.String())
	if !utils.FileExists(newSavePath) {
		t.Error("New save file not created")
	}
	// Check that it's hardlinked to the slot
	slotPath := filepath.Join(slotsPath, SlotID{1, 1, false}.String())
	hardLinkSlot, _ := utils.GetHardLinkID(slotPath)
	hardLinkNew, _ := utils.GetHardLinkID(newSavePath)
	if hardLinkSlot != hardLinkNew {
		t.Errorf("New save not hardlinked to slot")
	}
	// Check Saves map
	if _, ok := sm.Saves[newSaveID]; !ok {
		t.Error("New save not in Saves map")
	}
	// Check SaveLinks
	if ids, ok := sm.SaveLinks[hardLinkSlot]; !ok || !slices.Contains(ids, newSaveID) {
		t.Errorf("SaveLinks missing new save: %v", sm.SaveLinks[hardLinkSlot])
	}
}

// Test saving a non-existant save slot to a save file
func TestSaveInvalidSlot(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// SaveSlot on non-existent slot should do nothing (no error)
	err = sm.SaveSlot("NewSave", 1, 99, false)
	if err != nil {
		t.Fatalf("SaveSlot on non-existent slot returned error: %s", err)
	}
	// Ensure no save file created
	savePath := filepath.Join(managerPath, SaveID{"NewSave", 1, false}.String())
	if utils.FileExists(savePath) {
		t.Error("Save file created for non-existent slot")
	}
}

// Test saving a save slot to a taken save file
func TestSaveSlotTakenSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create a save and set it to slot
	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	err = sm.SetSlot("MySave", 1, 1, false, false)
	if err != nil {
		t.Fatalf("SetSlot failed: %s", err)
	}
	// Create another save with the name we want to use for SaveSlot
	err = sm.Create("NewSave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	// Now try SaveSlot with same name
	err = sm.SaveSlot("NewSave", 1, 1, false)
	if err == nil {
		t.Fatal("SaveSlot with taken name did not return error")
	}
}

// Test removing a save
func TestRemoveSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	// Create and set a save to a slot (so we have both)
	err = sm.Create("MySave", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	err = sm.SetSlot("MySave", 1, 1, false, false)
	if err != nil {
		t.Fatalf("SetSlot failed: %s", err)
	}
	saveID := SaveID{"MySave", 1, false}
	savePath := filepath.Join(managerPath, saveID.String())
	slotID := SlotID{1, 1, false}
	slotPath := filepath.Join(slotsPath, slotID.String())
	// Remove save without removing slots
	err = sm.Remove("MySave", 1, false, false)
	if err != nil {
		t.Fatalf("Remove failed: %s", err)
	}
	if utils.FileExists(savePath) {
		t.Error("Save file still exists after removal")
	}
	if _, ok := sm.Saves[saveID]; ok {
		t.Error("Save still in Saves map")
	}
	// Slot should still exist
	if !utils.FileExists(slotPath) {
		t.Error("Slot file was removed unexpectedly")
	}
	if _, ok := sm.Slots[slotID]; !ok {
		t.Error("Slot removed from Slots map unexpectedly")
	}
	// SlotLinks should still have the slot with its hardlink
	hardLinkSlot, _ := utils.GetHardLinkID(slotPath)
	if ids, ok := sm.SlotLinks[hardLinkSlot]; !ok || len(ids) != 1 || ids[0] != slotID {
		t.Errorf("SlotLinks missing slot: %v", sm.SlotLinks[hardLinkSlot])
	}
}

// Test removing a non-existant save
func TestRemoveInvalidSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Remove("NonExistent", 1, false, false)
	if err == nil {
		t.Fatal("Remove on non-existent save did not return error")
	}
}

// Test renaming a save
func TestRenameSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Create("OldName", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	// Rename
	err = sm.Rename("OldName", "NewName", 1, false)
	if err != nil {
		t.Fatalf("Rename failed: %s", err)
	}
	oldID := SaveID{"OldName", 1, false}
	newID := SaveID{"NewName", 1, false}
	oldPath := filepath.Join(managerPath, oldID.String())
	newPath := filepath.Join(managerPath, newID.String())
	if utils.FileExists(oldPath) {
		t.Error("Old file still exists after rename")
	}
	if !utils.FileExists(newPath) {
		t.Error("New file not created after rename")
	}
	if _, ok := sm.Saves[oldID]; ok {
		t.Error("Old save still in Saves map")
	}
	if _, ok := sm.Saves[newID]; !ok {
		t.Error("New save not in Saves map")
	}
	// Check SaveLinks: the hardlink should have new ID instead of old
	hardLinkID, _ := utils.GetHardLinkID(newPath)
	if ids, ok := sm.SaveLinks[hardLinkID]; !ok || len(ids) != 1 || ids[0] != newID {
		t.Errorf("SaveLinks not updated: %v", sm.SaveLinks[hardLinkID])
	}
}

// Test renaming a save to a taken name
func TestRenameSaveToTakenName(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Create("SaveA", 1, false)
	if err != nil {
		t.Fatalf("Create SaveA failed: %s", err)
	}
	err = sm.Create("SaveB", 1, false)
	if err != nil {
		t.Fatalf("Create SaveB failed: %s", err)
	}
	err = sm.Rename("SaveA", "SaveB", 1, false)
	if err == nil {
		t.Fatal("Rename to taken name did not return error")
	}
}

// Test renaming a non-existant save
func TestRenameInvalidSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Rename("NonExistent", "NewName", 1, false)
	if err == nil {
		t.Fatal("Rename on non-existent save did not return error")
	}
}

// Test copying a save
func TestCopySave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Create("Original", 1, false)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	origPath := filepath.Join(managerPath, SaveID{"Original", 1, false}.String())
	// Read original content to compare later
	origContent, err := os.ReadFile(origPath)
	if err != nil {
		t.Fatalf("Read original failed: %s", err)
	}

	// Copy the save
	err = sm.Copy("Original", "Copy", 1, false)
	if err != nil {
		t.Fatalf("Copy failed: %s", err)
	}
	copyID := SaveID{"Copy", 1, false}
	copyPath := filepath.Join(managerPath, copyID.String())
	if !utils.FileExists(copyPath) {
		t.Error("Copy file not created")
	}
	// Check content copied
	copyContent, err := os.ReadFile(copyPath)
	if err != nil {
		t.Fatalf("Read copy failed: %s", err)
	}
	if string(copyContent) != string(origContent) {
		t.Errorf("Copy content mismatch: expected %q, got %q", origContent, copyContent)
	}
	// Check Saves map
	if _, ok := sm.Saves[copyID]; !ok {
		t.Error("Copy not in Saves map")
	}
	// Check SaveLinks: copy should have its own hardlink
	hardLinkCopy, _ := utils.GetHardLinkID(copyPath)
	if ids, ok := sm.SaveLinks[hardLinkCopy]; !ok || len(ids) != 1 || ids[0] != copyID {
		t.Errorf("SaveLinks missing copy: %v", sm.SaveLinks[hardLinkCopy])
	}
	// Original should still have its own hardlink
	hardLinkOrig, _ := utils.GetHardLinkID(origPath)
	if hardLinkOrig == hardLinkCopy {
		t.Error("Original and copy are hardlinked, but should be separate")
	}
}

// Test copying a non-existant save
func TestCopyInvalidSave(t *testing.T) {
	managerPath := t.TempDir()
	slotsPath := t.TempDir()
	drIniPath := filepath.Join(slotsPath, "dr.ini")
	err := os.WriteFile(drIniPath, []byte("[General]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dr.ini: %s", err)
	}
	sm, err := NewSaveManager(managerPath, slotsPath)
	if err != nil {
		t.Fatalf("NewSaveManager failed: %s", err)
	}
	err = sm.Copy("NonExistent", "Copy", 1, false)
	if err == nil {
		t.Fatal("Copy on non-existent save did not return error")
	}
}

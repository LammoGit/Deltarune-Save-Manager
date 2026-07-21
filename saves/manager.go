package saves

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"path/filepath"
)

type SaveID struct {
	Name    string
	Chapter int
	SideB   bool
}

func (id SaveID) String() string {
	if id.SideB {
		return fmt.Sprintf("%d_a_%s", id.Chapter, id.Name)
	} else {
		return fmt.Sprintf("%d_b_%s", id.Chapter, id.Name)
	}
}

type SlotID struct {
	Chapter int
	Slot    int
	SideB   bool
}

func (id SlotID) String() string {
	if id.SideB {
		return fmt.Sprintf("filech%d_%d_b", id.Chapter, id.Slot)
	} else {
		return fmt.Sprintf("filech%d_%d", id.Chapter, id.Slot)
	}
}

type SaveManager struct {
	ManagerPath string
	SlotsPath   string
	Saves       map[SaveID]Save
	Slots       map[SlotID]Save
	SaveLinks   map[string][]SaveID
	SlotLinks   map[string][]SlotID
	Dr          DrINI
}

func loadSaves(dirPath string) (map[SaveID]Save, map[string][]SaveID, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}
	
	saves := make(map[SaveID]Save)
	links := make(map[string][]SaveID)

	for _, entry := range entries {
		if entry.IsDir() { continue }

		savePath := filepath.Join(dirPath, entry.Name())

		match := saveRegex.FindStringSubmatch(entry.Name())

		if len(match) != 4 {
			continue
		}

		chapter, err := strconv.Atoi(match[1])
		if err != nil { continue }

		sideB := match[2] != "a"
		name  := match[3]

		save, err := LoadSave(savePath, chapter)
		if err != nil { continue }

		saveID := SaveID{name, chapter, sideB}
		hardLinkID, err := getHardLinkID(savePath)
		if err != nil { continue }

		saves[saveID] = save
		links[hardLinkID] = append(links[hardLinkID], saveID)
	}

	return saves, links, nil
}

func loadSlots(dirPath string) (map[SlotID]Save, map[string][]SlotID, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}

	slots := make(map[SlotID]Save)
	links := make(map[string][]SlotID)

	for _, entry := range entries {
		if entry.IsDir() { continue }

		slotPath := filepath.Join(dirPath, entry.Name())
		
		match := slotRegex.FindStringSubmatch(entry.Name())
		
		if len(match) != 4 { 
			continue 
		}
		
		chapter, err := strconv.Atoi(match[1])
		if err != nil { continue }
		
		slot, err := strconv.Atoi(match[2])
		if err != nil { continue }

		sideB := match[3] != ""
		
		save, err := LoadSave(slotPath, chapter)
		if err != nil { continue }
		
		slotID := SlotID{ chapter, slot, sideB }
		hardLinkID, err := getHardLinkID(slotPath)
		if err != nil { continue }

		slots[slotID] = save
		links[hardLinkID] = append(links[hardLinkID], slotID)
	}

	return slots, links, nil
}

func NewSaveManager(managerPath, slotsPath string) (sm *SaveManager, err error) {
	saves, saveLinks, err := loadSaves(managerPath)
	if err != nil {
		return
	}

	slots, slotLinks, err := loadSlots(slotsPath)
	if err != nil {
		return
	}
	
	dr, err := NewDrINI(filepath.Join(slotsPath, "dr.ini"))
	if err != nil {
		return
	}
	
	sm = &SaveManager {
		ManagerPath: managerPath,
		SlotsPath:   slotsPath,
		Saves:       saves,
		Slots:       slots,
		SaveLinks:   saveLinks,
		SlotLinks:   slotLinks,
		Dr:          dr,
	}
	return
}

/* Utils */
func (sm *SaveManager) hardLinkIDFromSaveID(id SaveID) (string, bool) {
	for hardLinkID, saveID := range sm.SaveLinks {
		if slices.Contains(saveID, id) {
			return hardLinkID, true
		}
	}
	return "", false
}

func (sm *SaveManager) hardLinkIDFromSlotID(id SlotID) (string, bool) {
	for hardLinkID, slotID := range sm.SlotLinks {
		if slices.Contains(slotID, id) {
			return hardLinkID, true
		}
	}
	return "", false
}

/* Manage Saves */
func (sm *SaveManager) Create(name string, chapter int) error {
	if chapter > MAX_CHAPTER {
		return ErrChapterNotSupported
	}
	
	if name == "" {
		return ErrEmptySaveName
	}
	
	saveID := SaveID{Name: name, Chapter: chapter}
	path := filepath.Join(sm.ManagerPath, saveID.String())
	file, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := getExampleSaveBytesForChapter(chapter)
	if err != nil {
		return err
	}
	
	_, err = file.Write(content)
	if err != nil {
		return err
	}

	hardLinkID, err := getHardLinkID(path)
	if err != nil {
		return err
	}

	sm.SaveLinks[hardLinkID] = append(sm.SaveLinks[hardLinkID], saveID)
	sm.Saves[saveID], err = LoadSave(path, chapter)
	if err != nil {
		return err
	}

	return nil
}

func (sm *SaveManager) Swap(name1, name2 string, chapter int) error {
	saveID1 := SaveID{Name: name1, Chapter: chapter}
	saveID2 := SaveID{Name: name2, Chapter: chapter}

	path1 := filepath.Join(sm.ManagerPath, saveID1.String())
	path2 := filepath.Join(sm.ManagerPath, saveID2.String())

	// If at least one doesn't exist return an error
	if !fileExists(path1) || !fileExists(path2) {
		return ErrSaveNotExist
	}


	hardLinkID1, _ := sm.hardLinkIDFromSaveID(saveID1)
	hardLinkID2, _ := sm.hardLinkIDFromSaveID(saveID2)

	// No need to swap if they point to the same memory
	if hardLinkID1 == hardLinkID2 {
		return nil
	}

	// Get a temporary path
	tmpPath, err := tempFilePath(sm.ManagerPath)
	if err != nil {
		return err
	}

	// Move first save to temporary path
	err = os.Rename(path1, tmpPath)
	if err != nil {
		return err
	}

	// Move second save to first save path
	err = os.Rename(path2, path1)
	if err != nil {
		_ = os.Rename(tmpPath, path1)
		return err
	}

	// Move first save to second save path
	err = os.Rename(tmpPath, path2)
	if err != nil {
		_ = os.Rename(path1, path2)
		_ = os.Rename(tmpPath, path1)
		return err
	}

	// Swap save objects inside of the save map
	tmpSave := sm.Saves[saveID1]
	sm.Saves[saveID1] = sm.Saves[saveID2]
	sm.Saves[saveID2] = tmpSave

	// Remove old save identifiers from saves' hard links map
	sm.SaveLinks[hardLinkID1] = deleteEqual(sm.SaveLinks[hardLinkID1], saveID1)
	sm.SaveLinks[hardLinkID2] = deleteEqual(sm.SaveLinks[hardLinkID2], saveID2)

	// Add new save identifiers to saves' hard links map
	sm.SaveLinks[hardLinkID1] = append(sm.SaveLinks[hardLinkID1], saveID2)
	sm.SaveLinks[hardLinkID2] = append(sm.SaveLinks[hardLinkID2], saveID1)

	return nil
}

func (sm *SaveManager) SetSlot(name string, chapter, slot int, eraseUnmanaged bool) error {
	saveID := SaveID{Name: name, Chapter: chapter}
	savePath := filepath.Join(sm.ManagerPath, saveID.String())

	slotID := SlotID{Chapter: chapter, Slot: slot}
	slotPath := filepath.Join(sm.SlotsPath, slotID.String())

	saveHardLink, ok := sm.hardLinkIDFromSaveID(saveID)
	if !ok {
		return ErrSaveNotExist
	}

	slotHardLink, slotExist := sm.hardLinkIDFromSlotID(slotID)
	// if slotExist && saveHardLink == slotHardLink {
	// 	return nil
	// }

	// Get saves pointing to the same space in memory
	saves, ok := sm.SaveLinks[slotHardLink]
	
	// If slot exists, doesn't have linked saves and
	// it's not allowed to delete unlinked slots, than return error
	if slotExist && (!ok || len(saves) == 0) && !eraseUnmanaged {
		return fmt.Errorf("Slot is already taken by an unmanaged save")
	}

	// Make slot link to the given save file
	err := relink(slotPath, savePath)
	if err != nil {
		return err
	}
	// Update slot object
	sm.Slots[slotID], _ = LoadSave(slotPath, chapter)
	// Remove slot identifier from the old hard link
	sm.SlotLinks[slotHardLink] = deleteEqual(sm.SlotLinks[slotHardLink], slotID)
	// Add slot identifier to the new hard link
	sm.SlotLinks[saveHardLink] = append(sm.SlotLinks[saveHardLink], slotID)
	// Update dr.ini
	sm.Dr.SetSlotFromSave(
		sm.Slots[slotID],
		chapter,
		slot,
		true,
	)

	return nil
}

func (sm *SaveManager) UnsetSlot(chapter, slot int, eraseUnmanaged bool) error {
	slotID := SlotID{Chapter: chapter, Slot: slot}
	slotPath := filepath.Join(sm.SlotsPath, slotID.String())

	slotHardLink, ok := sm.hardLinkIDFromSlotID(slotID)
	// If slot file doesn't exist just exit
	if !ok {
		return nil
	}

	// Get saves pointing to the same space in memory
	saves, ok := sm.SaveLinks[slotHardLink]

	// If slot exist, doessn't have linked saves and
	// it's not allowed to delete unlinked slots, than return error
	if (!ok || len(saves) == 0) && !eraseUnmanaged {
		return fmt.Errorf("Slot is taken by an unmanaged save")
	}

	// Remove the slot file
	err := os.Remove(slotPath)
	if err != nil {
		return err
	}
	// Remove the slot object from the slots map
	delete(sm.Slots, slotID)

	// Remove slot identifier from the hard link map
	sm.SlotLinks[slotHardLink] = deleteEqual(sm.SlotLinks[slotHardLink], slotID)

	// Update dr.ini
	sm.Dr.RemoveSlot(
		chapter,
		slot,
	)

	return nil
}

func (sm *SaveManager) SaveSlot(name string, chapter, slot int) error {
	saveID := SaveID{Name: name, Chapter: chapter}
	savePath := filepath.Join(sm.ManagerPath, saveID.String())

	slotID := SlotID{Chapter: chapter, Slot: slot}
	slotPath := filepath.Join(sm.SlotsPath, slotID.String())


	slotHardLink, ok := sm.hardLinkIDFromSlotID(slotID)
	// If slot file doesn't exist just exit
	if !ok {
		return nil
	}

	// If save file name is already taken, then return error
	if fileExists(savePath) {
		return ErrSaveNameIsTaken
	}
	
	// Create a new save file linking to the same memory space
	// as the slot
	err := os.Link(slotPath, savePath)
	if err != nil {
		return err
	}

	// Add new save object to the saves map
	sm.Saves[saveID], _ = LoadSave(savePath, chapter)

	// Add save idetifier to the hard link map
	sm.SaveLinks[slotHardLink] = append(sm.SaveLinks[slotHardLink], saveID)

	return nil
}

func (sm *SaveManager) Remove(name string, chapter int, removeSlots bool) error {
	saveID := SaveID{Name: name, Chapter: chapter}
	savePath := filepath.Join(sm.ManagerPath, saveID.String())

	saveHardLink, ok := sm.hardLinkIDFromSaveID(saveID)
	// If save file doesn't exist just exit
	if !ok {
		return ErrSaveNotExist
	}

	// Remove the save file
	err := os.Remove(savePath)
	if err != nil {
		return err
	}

	// Remove the save object from the saves map
	delete(sm.Saves, saveID)

	// Remove save identifier from the hard link map
	sm.SaveLinks[saveHardLink] = deleteEqual(sm.SaveLinks[saveHardLink], saveID)

	// Remove all linked slots if required
	if removeSlots {
		for _, slotID := range sm.SlotLinks[saveHardLink] {
			// Get path to the slot file
			slotPath := filepath.Join(sm.SlotsPath, slotID.String())
			// Remove the slot file
			os.Remove(slotPath)
			// Remove the slot file's object
			delete(sm.Slots, slotID)
			// update dr.ini
			sm.Dr.RemoveSlot(
				chapter,
				slotID.Slot,
			)
		}
		// Remove all slots with the save's hard link
		delete(sm.SlotLinks, saveHardLink)
	}
	return nil
}

func (sm *SaveManager) Rename(nameFrom, nameTo string, chapter int) error {
	saveIDFrom := SaveID{Name: nameFrom, Chapter: chapter}
	saveIDTo := SaveID{Name: nameTo, Chapter: chapter}

	savePathFrom := filepath.Join(sm.ManagerPath, saveIDFrom.String())
	savePathTo := filepath.Join(sm.ManagerPath, saveIDTo.String())

	saveHardLink, ok := sm.hardLinkIDFromSaveID(saveIDFrom)
	// If save file doesn't exist return an error
	if !ok {
		return ErrSaveNotExist
	}

	// If save name is already taken return an error
	if fileExists(savePathTo) {
		return ErrSaveNameIsTaken
	}

	// Rename the save file
	err := os.Rename(savePathFrom, savePathTo)
	if err != nil {
		return err
	}

	// Update save identifier for the save object in the saves map
	sm.Saves[saveIDTo] = sm.Saves[saveIDFrom]
	delete(sm.Saves, saveIDFrom)

	// Remove old save identifier from the hard link map
	sm.SaveLinks[saveHardLink] = deleteEqual(sm.SaveLinks[saveHardLink], saveIDFrom)

	// Add new save identifier to the hard link map
	sm.SaveLinks[saveHardLink] = append(sm.SaveLinks[saveHardLink], saveIDTo)

	return nil
}

func (sm *SaveManager) Copy(nameFrom, nameTo string, chapter int) error {
	saveIDFrom := SaveID{Name: nameFrom, Chapter: chapter}
	saveIDTo := SaveID{Name: nameTo, Chapter: chapter}

	savePathFrom := filepath.Join(sm.ManagerPath, saveIDFrom.String())
	savePathTo := filepath.Join(sm.ManagerPath, saveIDTo.String())

	// If save file doesn't exist return an error
	if !fileExists(savePathFrom) {
		return ErrSaveNotExist
	}

	// If save name is already taken return an error
	if fileExists(savePathTo) {
		return ErrSaveNameIsTaken
	}

	// Create a copy save file
	fileFrom, err := os.Open(savePathFrom)
	if err != nil {
		return err
	}
	defer fileFrom.Close()

	fileTo, err := os.Create(savePathTo)
	if err != nil {
		return err
	}
	defer fileTo.Close()

	fileFrom.WriteTo(fileTo)
	
	// Load save object
	_, err 	= fileFrom.Seek(0, 0)
	if err != nil {
		return err
	}

	save, err := ParseSaveReader(fileFrom, chapter)
	if err != nil {
		return err
	}
	
	// Add new save object to the saves map
	sm.Saves[saveIDTo] = save
	
	// Add new save identifier to the hard link map
	saveHardLink, err := getHardLinkID(savePathTo)
	if err != nil {
		return err
	}
	sm.SaveLinks[saveHardLink] = append(sm.SaveLinks[saveHardLink], saveIDTo)

	return nil
}

func (sm *SaveManager) Edit(name string, chapter int, props map[string]string) error {
	saveID := SaveID{Name: name, Chapter: chapter}
	save, ok := sm.Saves[saveID]
	if !ok {
		return ErrSaveNotExist
	}
	return Edit(&save, props)
}

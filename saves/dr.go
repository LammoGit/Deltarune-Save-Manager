package saves

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"dsm/ini"
	"dsm/utils"
)

// SlotDr type represents a save slot section in dr.ini
type SlotDr ini.INISection

// DrINI type represents the dr.ini file
// it holds data of the dr.ini and path to it
type DrINI struct {
	ini.INI
	Path string
}

// chapterIndexToKey generates a save slot section label
// from its chapter and index
func chapterIndexToKey(chapter, index int) string {
	if chapter == 1 {
		return fmt.Sprintf("G%d", index)
	}
	return fmt.Sprintf("G%d_%d", chapter, index)
}

// NewDrINI creates a new DrINI object
func NewDrINI(path string) (dr DrINI, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	ini, err := ini.Parse(bytes.NewBuffer(content))
	if err != nil {
		return
	}
	dr = DrINI{ini, path}
	dr.Clean()
	return
}

// Clean removes all empty saves
func (dr *DrINI) Clean() {
	for label, section := range dr.INI {
		if name, ok := section["Name"]; ok && name == "[EMPTY]" {
			delete(dr.INI, label)
		}
	}
}

// GetSlot returns a save slot section from the dr.ini file
// given its chapter and index
func (dr *DrINI) GetSlot(chapter, index int) (slot *SlotDr, ok bool) {
	key := chapterIndexToKey(chapter, index)
	val, ok := dr.INI[key]
	if ok {
		slot = (*SlotDr)(&val)
	}
	return
}

// GetSlots returns all save slot sections from the dr.ini file
func (dr *DrINI) GetSlots() (slots []*SlotDr) {
	for label, slot := range dr.INI {
		if utils.SlotSectionLabelRegex.MatchString(label) {
			slots = append(slots, (*SlotDr)(&slot))
		}
	}
	return
}

// SetSlot sets a save slot with given slot and index
func (dr *DrINI) SetSlot(slot SlotDr, index int) bool {
	chapterS, ok := slot["Chapter"]
	if !ok {
		return false
	}

	chapter, err := strconv.Atoi(chapterS)
	if err != nil {
		return false
	}

	if _, ok := dr.GetSlot(chapter, index); ok {
		return false
	}

	dr.INI[chapterIndexToKey(chapter, index)] = ini.INISection(slot)
	return true
}

// SetSlotFromSave sets a save slot at given chapter and index by setting sector data
// using given save. Replaces previous save when replace is set to true
func (dr *DrINI) SetSlotFromSave(save Save, chapter int, index int, replace bool) bool {
	if !replace {
		if _, ok := dr.GetSlot(chapter, index); ok {
			return false
		}
	}

	f := func(name string, time any, room float64) {
		dr.INI[chapterIndexToKey(chapter, index)] = ini.INISection(map[string]string{
			"Name": fmt.Sprintf(`"%s"`, name),
			"Time": fmt.Sprintf(`"%f"`, time),
			"Room": fmt.Sprintf(`"%f"`, room),
		})
	}

	switch s := save.(type) {
	case Save1:
		f(s.PlayerName, s.Time, s.Room)
	case Save2:
		f(s.PlayerName, s.Time, s.Room)
	case Save3:
		f(s.PlayerName, s.Time, s.Room)
	case Save4:
		f(s.PlayerName, s.Time, s.Room)
	case Save5:
		f(s.PlayerName, s.Time, s.Room)
	default:
		return false
	}

	_ = dr.Write()
	return true
}

// CopySlot copies a save slot data from one slot to another
func (dr *DrINI) CopySlot(chapter, indexFrom, indexTo int) bool {
	keyTo := chapterIndexToKey(chapter, indexTo)
	if _, ok := dr.INI[keyTo]; ok {
		return false
	}

	keyFrom := chapterIndexToKey(chapter, indexFrom)
	slot, ok := dr.INI[keyFrom]
	if !ok {
		return false
	}

	dr.INI[keyTo] = slot
	dr.Write()
	return true
}

// MoveSlot moves a save slot data from one slot to another
func (dr *DrINI) MoveSlot(chapter, indexFrom, indexTo int) bool {
	keyTo := chapterIndexToKey(chapter, indexTo)
	if _, ok := dr.INI[keyTo]; ok {
		return false
	}

	keyFrom := chapterIndexToKey(chapter, indexFrom)
	slot, ok := dr.INI[keyFrom]
	if !ok {
		return false
	}

	dr.INI[keyTo] = slot
	delete(dr.INI, keyFrom)
	dr.Write()
	return true
}

// SwapSlots swaps save slot data of two slots
func (dr *DrINI) SwapSlots(chapter, index1, index2 int) bool {
	key1 := chapterIndexToKey(chapter, index1)
	slot1, ok := dr.INI[key1]
	if !ok {
		return false
	}

	key2 := chapterIndexToKey(chapter, index2)
	slot2, ok := dr.INI[key2]
	if !ok {
		return false
	}

	dr.INI[key1] = slot2
	dr.INI[key2] = slot1
	dr.Write()
	return true
}

// RemoveSlot removes a save slot from the dr.ini
func (dr *DrINI) RemoveSlot(chapter, index int) {
	key := chapterIndexToKey(chapter, index)
	delete(dr.INI, key)
	dr.Write()
}

// Write writes current contents of the DrINI object to the dr.ini file
func (dr *DrINI) Write() error {
	file, err := os.OpenFile(dr.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	err = dr.INI.Write(file, true)
	return err
}

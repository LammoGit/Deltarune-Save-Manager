package saves

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"dsm/ini"
)

var slotectionLabelRegex = regexp.MustCompile(`G(?:[2-7]_)?\d+`)

type SlotDr ini.INISection

type DrINI struct {
	ini.INI
	Path string
}

func NewDrINI(path string) (dr DrINI, err error){
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

func (dr *DrINI) Clean() {
	for label, section := range dr.INI {
		if name, ok := section["Name"];ok && name == "[EMPTY]" {
			delete(dr.INI, label)
		}
	}
}

func (dr *DrINI) GetSlot(chapter, index int) (slot *SlotDr, ok bool) {
	key := chapterIndexToKey(chapter, index)
	val, ok := dr.INI[key]
	if ok {
		slot = (*SlotDr)(&val)
	}
	return
}

func (dr *DrINI) GetSlots() (slots []*SlotDr) {
	for label, slot := range dr.INI {
		if slotectionLabelRegex.MatchString(label) {
			slots = append(slots, (*SlotDr)(&slot))
		}
	}
	return
}

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

func (dr *DrINI) RemoveSlot(chapter, index int) {
	key := chapterIndexToKey(chapter, index)
	delete(dr.INI, key)
	dr.Write()
}

func (dr *DrINI) Write() error {
	file, err := os.OpenFile(dr.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	err = dr.INI.Write(file, true)
	return err
}

func chapterIndexToKey(chapter, index int) string {
	if chapter == 1 {
		return fmt.Sprintf("G%d", index)
	}
	return fmt.Sprintf("G%d_%d", chapter, index)
}

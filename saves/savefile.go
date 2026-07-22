package saves

import (
	"bytes"
	"dsm/utils"
	_ "embed"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/* Loading chapters' default save files */

//go:embed 1_savefile
var chapter1bytes []byte

//go:embed 2_savefile
var chapter2bytes []byte

//go:embed 3_savefile
var chapter3bytes []byte

//go:embed 4_savefile
var chapter4bytes []byte

//go:embed 5_savefile
var chapter5bytes []byte

/* Type defenitions */

//go:generate stringer -type=Character -trimprefix=Character
type Character int

const (
	CharacterRalsei Character = iota
	CharacterKris
	CharacterSusie
	CharacterRalseiWithHat
	CharacterNoelle
)

type Weapon int

const ()

type Armor int

const ()

type Spell int

const ()

type Element int

const ()

type ItemStats struct {
	Attack    int
	Defence   int
	Magic     int
	Bolts     int
	Grazeamt  int
	GrazeSize int
	BoltSpeed int
	Special   int
}

type CharacterStats struct {
	HP          int
	MaxHP       int
	Attack      int
	Defence     int
	Magic       int
	Guts        int
	Weapon      Weapon
	Armor       [2]Armor
	WeaponStyle string
	ItemsStats  [4]ItemStats
	Spells      [12]Spell
}

type InventorySlot struct {
	Item    int
	KeyItem int
	Weapon  int
	Armor   int
}

type ItemAndPhoneSlot struct {
	Item  int
	Phone int
}

type LightWorldStats struct {
	Weapon         int
	Armor          int
	XP             int
	Level          int
	Gold           int
	HP             int
	MaxHP          int
	Attack         int
	Defence        int
	WeaponStrength int
	ArmorDefence   int
	Inventory      [8]ItemAndPhoneSlot
}

type GlobalFlags struct {
	Flags [2500]string
}

type Save1 struct {
	PlayerName        string
	CharName          string
	OtherNames        [5]string
	Characters        [3]Character
	Gold              int
	XP                int
	Level             int
	Inv               int
	Invc              int
	Darkzone          bool
	CharactersStats   [4]CharacterStats
	BoltSpeed         int
	Grazeamt          int
	GrazeSize         int
	Inventory         [13]InventorySlot
	Tension           int
	MaxTension        int
	LightWorldStats   LightWorldStats
	GlobalFlags       GlobalFlags
	UnusedGlobalFlags [7499]string
	Plot              float64
	Room              float64
	Time              float64
}

type ItemStats2 struct {
	Attack        int
	Defence       int
	Magic         int
	Bolts         int
	Grazeamt      int
	GrazeSize     int
	BoltSpeed     int
	Special       int
	Element       Element
	ElementAmount float64
}

type CharacterStats2 struct {
	HP          int
	MaxHP       int
	Attack      int
	Defence     int
	Magic       int
	Guts        int
	Weapon      Weapon
	Armor       [2]Armor
	WeaponStyle string
	ItemsStats  [4]ItemStats2
	Spells      [12]Spell
}

type Inventory2 struct {
	ItemsAndKeyItems [13]struct{ Item, KeyItem int }
	WeaponsAndArmors [48]struct {
		Weapon Weapon
		Armor  Armor
	}
	PocketItems [72]int
}

type Save2 struct {
	PlayerName      string
	CharName        string
	OtherNames      [5]string
	Characters      [3]Character
	Gold            int
	XP              int
	Level           int
	Inv             int
	Invc            int
	Darkzone        int
	CharactersStats [5]CharacterStats2
	BoltSpeed       int
	Grazeamt        int
	GrazeSize       int
	Inventory       Inventory2
	Tension         int
	MaxTension      int
	LightWorldStats LightWorldStats
	GlobalFlags     GlobalFlags
	Plot            float64
	Room            float64
	Time            float64
}

type Save3 Save2

type Save4 Save2

type Save5 Save2

type Save any

func parseSaveLine(lines []string, cur int, kind reflect.Kind, v reflect.Value) (int, error) {
	switch kind {
	case reflect.String:
		if cur >= len(lines) {
			return cur, utils.ErrShortSaveFile
		}
		if !v.CanSet() {
			return cur, utils.ErrValueCannotBeSet
		}
		v.SetString(lines[cur])
		cur++
	case reflect.Bool:
		if cur >= len(lines) {
			return cur, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return cur, utils.ErrValueCannotBeSet
		}
		num, err := strconv.Atoi(strings.Trim(lines[cur], " "))
		if err != nil {
			return cur, fmt.Errorf("%w: failed to parse to a boolean value %q", utils.ErrWrongLineType, lines[cur])
		}

		v.SetBool(num != 0)
		cur++
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if cur >= len(lines) {
			return cur, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return cur, utils.ErrValueCannotBeSet
		}

		var bitSize int
		switch kind {
		case reflect.Int:
			bitSize = 0
		case reflect.Int8:
			bitSize = 8
		case reflect.Int16:
			bitSize = 16
		case reflect.Int32:
			bitSize = 32
		case reflect.Int64:
			bitSize = 64
		}

		num, err := strconv.ParseInt(strings.Trim(lines[cur], " "), 10, bitSize)
		if err != nil {
			return cur, fmt.Errorf("%w: failed to parse to a signed integer value %q", utils.ErrWrongLineType, lines[cur])
		}

		v.SetInt(int64(num))
		cur++
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if cur >= len(lines) {
			return cur, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return cur, utils.ErrValueCannotBeSet
		}

		var bitSize int
		switch kind {
		case reflect.Uint:
			bitSize = 0
		case reflect.Uint8:
			bitSize = 8
		case reflect.Uint16:
			bitSize = 16
		case reflect.Uint32:
			bitSize = 32
		case reflect.Uint64:
			bitSize = 64
		}

		num, err := strconv.ParseUint(strings.Trim(lines[cur], " "), 10, bitSize)
		if err != nil {
			return cur, fmt.Errorf("%w: failed to parse to an unsigned integer value %q", utils.ErrWrongLineType, lines[cur])
		}

		v.SetUint(uint64(num))
		cur++
	case reflect.Float32, reflect.Float64:
		if cur >= len(lines) {
			return cur, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return cur, utils.ErrValueCannotBeSet
		}

		var bitSize int
		switch kind {
		case reflect.Float32:
			bitSize = 32
		case reflect.Float64:
			bitSize = 64
		}

		num, err := strconv.ParseFloat(strings.Trim(lines[cur], " "), bitSize)
		if err != nil {
			return cur, fmt.Errorf("%w: failed to parse to a floating point number value %q", utils.ErrWrongLineType, lines[cur])
		}

		v.SetFloat(float64(num))
		cur++
	case reflect.Struct:
		for _, fieldValue := range v.Fields() {
			fieldKind := fieldValue.Kind()

			var err error
			cur, err = parseSaveLine(lines, cur, fieldKind, fieldValue)
			if err != nil {
				return cur, err
			}
		}
	case reflect.Array:
		elemKind := v.Type().Elem().Kind()
		for i := 0; i < v.Len(); i++ {
			var err error
			cur, err = parseSaveLine(lines, cur, elemKind, v.Index(i))
			if err != nil {
				return cur, err
			}
		}
	}
	return cur, nil
}

func ParseSaveBytes(buf []byte, chapter int) (Save, error) {
	buf = bytes.ReplaceAll(buf, []byte("\r\n"), []byte("\n"))
	text := string(bytes.Trim(buf, "\n"))

	lines := strings.Split(text, "\n")

	switch chapter {
	case 1:
		save := Save1{}
		saveValue := reflect.ValueOf(&save).Elem()
		_, err := parseSaveLine(lines, 0, saveValue.Kind(), saveValue)
		return save, err
	default:
		save := Save2{}
		saveValue := reflect.ValueOf(&save).Elem()
		_, err := parseSaveLine(lines, 0, saveValue.Kind(), saveValue)
		if err != nil {
			fmt.Printf("%d, %s\n", chapter, err)
		}
		return save, err
	}
}

func LoadSave(path string, chapter int) (save Save, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	save, err = ParseSaveBytes(content, chapter)
	return
}

func ParseSaveReader(r io.Reader, chapter int) (save Save, err error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return
	}
	save, err = ParseSaveBytes(buf, chapter)
	return
}

func getExampleSaveBytesForChapter(chapter int) ([]byte, error) {
	switch chapter {
	case 1:
		return chapter1bytes, nil
	case 2:
		return chapter2bytes, nil
	case 3:
		return chapter3bytes, nil
	case 4:
		return chapter4bytes, nil
	case 5:
		return chapter5bytes, nil
	}
	return []byte{}, utils.ErrChapterNotSupported
}

func Edit(s *Save, props map[string]string) error {
	return nil
}

package saves

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/LammoGit/Deltarune-Save-Manager/utils"
)

// nextLine fetches the next trimmed line from the scanner
func nextLine(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", utils.ErrShortSaveFile
	}
	return strings.TrimSpace(scanner.Text()), nil
}

// parseSaveField parses individual Save object's field
func parseSaveField(w io.Writer, kind reflect.Kind, v reflect.Value) error {
	switch kind {
	case reflect.String:
		fmt.Fprintf(w, "%s\n", v.String())
	case reflect.Bool:
		if v.Bool() {
			fmt.Fprintln(w, "1")
		} else {
			fmt.Fprintln(w, "0")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(w, "%d\n", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(w, "%d\n", v.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(w, "%f\n", v.Float())
	case reflect.Struct:
		for _, fieldValue := range v.Fields() {
			fieldKind := fieldValue.Kind()

			err := parseSaveField(w, fieldKind, fieldValue)
			if err != nil {
				return err
			}
		}
	case reflect.Array:
		elemKind := v.Type().Elem().Kind()
		for i := 0; i < v.Len(); i++ {
			err := parseSaveField(w, elemKind, v.Index(i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Save2Bytes returns bytes from a given Save object
func Save2Bytes(save Save) ([]byte, error) {
	var buf bytes.Buffer

	switch save.(type) {
	case *Save1:
	case *Save2:
	case nil:
		return nil, errors.New("save is a nil pointer")
	default:
		return nil, utils.ErrChapterNotSupported
	}

	saveValue := reflect.ValueOf(save).Elem()
	err := parseSaveField(&buf, saveValue.Kind(), saveValue)
	return buf.Bytes(), err
}

// ParseSaveBytes returns a Save object from a save file bytes and save's chapter
func ParseSaveBytes(buf []byte, chapter int) (Save, error) {
	reader := bytes.NewReader(buf)
	scanner := bufio.NewScanner(reader)

	var save Save
	switch chapter {
	case 1:
		return ParseSave1Generated(scanner)
	case 2, 3, 4, 5:
		return ParseSave2Generated(scanner)
	default:
		return save, utils.ErrChapterNotSupported
	}
}

// LoadSave returns a Save object from a save file path and save's chapter
func LoadSave(path string, chapter int) (Save, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseSaveReader(file, chapter)
}

// ParseSaveReader returns a Save object from a save file Reader object and save's chapter
func ParseSaveReader(r io.Reader, chapter int) (Save, error) {
	scanner := bufio.NewScanner(r)
	var save Save
	switch chapter {
	case 1:
		return ParseSave1Generated(scanner)
	case 2, 3, 4, 5:
		return ParseSave2Generated(scanner)
	default:
		return save, utils.ErrChapterNotSupported
	}
}

func ParseSave1Generated(scanner *bufio.Scanner) (*Save1, error) {
	s := &Save1{}
	var text string
	var err error
	var num int
	var num64 int64
	var f64 float64

	// Line Parsing: s.PlayerName
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	s.PlayerName = text

	// Line Parsing: s.CharName
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	s.CharName = text

	// Line Parsing: s.OtherNames
	for i0 := range 5 {
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		s.OtherNames[i0] = text
	}

	// Line Parsing: s.Characters
	for i0 := range 3 {
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num, err = strconv.Atoi(text); err != nil {
			return nil, err
		}
		s.Characters[i0] = Character(num)
	}

	// Line Parsing: s.Gold
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Gold = int(num64)

	// Line Parsing: s.XP
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.XP = int(num64)

	// Line Parsing: s.Level
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Level = int(num64)

	// Line Parsing: s.Inv
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Inv = int(num64)

	// Line Parsing: s.Invc
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Invc = int(num64)

	// Line Parsing: s.Darkzone
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num, err = strconv.Atoi(text); err != nil {
		return nil, err
	}
	s.Darkzone = num != 0

	// Line Parsing: s.CharactersStats
	for i0 := range 4 {
		// Line Parsing: s.CharactersStats[i0].HP
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].HP = int(num64)

		// Line Parsing: s.CharactersStats[i0].MaxHP
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].MaxHP = int(num64)

		// Line Parsing: s.CharactersStats[i0].Attack
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Attack = int(num64)

		// Line Parsing: s.CharactersStats[i0].Defence
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Defence = int(num64)

		// Line Parsing: s.CharactersStats[i0].Magic
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Magic = int(num64)

		// Line Parsing: s.CharactersStats[i0].Guts
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Guts = int(num64)

		// Line Parsing: s.CharactersStats[i0].Weapon
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Weapon = Weapon(num64)

		// Line Parsing: s.CharactersStats[i0].Armor
		for i1 := range 2 {
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num, err = strconv.Atoi(text); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].Armor[i1] = Armor(num)
		}

		// Line Parsing: s.CharactersStats[i0].WeaponStyle
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].WeaponStyle = text

		// Line Parsing: s.CharactersStats[i0].ItemsStats
		for i1 := range 4 {
			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Attack
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Attack = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Defence
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Defence = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Magic
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Magic = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Bolts
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Bolts = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].GrazeAmount
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].GrazeAmount = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].GrazeSize
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].GrazeSize = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].BoltSpeed
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].BoltSpeed = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Special
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Special = int(num64)

		}

		// Line Parsing: s.CharactersStats[i0].Spells
		for i1 := range 12 {
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num, err = strconv.Atoi(text); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].Spells[i1] = Spell(num)
		}

	}

	// Line Parsing: s.BoltSpeed
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.BoltSpeed = int(num64)

	// Line Parsing: s.GrazeAmount
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.GrazeAmount = int(num64)

	// Line Parsing: s.GrazeSize
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.GrazeSize = int(num64)

	// Line Parsing: s.Inventory
	for i0 := range 13 {
		// Line Parsing: s.Inventory[i0].Item
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.Inventory[i0].Item = Item(num64)

		// Line Parsing: s.Inventory[i0].KeyItem
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.Inventory[i0].KeyItem = KeyItem(num64)

		// Line Parsing: s.Inventory[i0].Weapon
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.Inventory[i0].Weapon = Weapon(num64)

		// Line Parsing: s.Inventory[i0].Armor
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.Inventory[i0].Armor = Armor(num64)

	}

	// Line Parsing: s.Tension
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Tension = int(num64)

	// Line Parsing: s.MaxTension
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.MaxTension = int(num64)

	// Line Parsing: s.LightWorldStats
	// Line Parsing: s.LightWorldStats.Weapon
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Weapon = LItem(num64)

	// Line Parsing: s.LightWorldStats.Armor
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Armor = LItem(num64)

	// Line Parsing: s.LightWorldStats.XP
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.XP = int(num64)

	// Line Parsing: s.LightWorldStats.Level
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Level = int(num64)

	// Line Parsing: s.LightWorldStats.Gold
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Gold = int(num64)

	// Line Parsing: s.LightWorldStats.HP
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.HP = int(num64)

	// Line Parsing: s.LightWorldStats.MaxHP
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.MaxHP = int(num64)

	// Line Parsing: s.LightWorldStats.Attack
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Attack = int(num64)

	// Line Parsing: s.LightWorldStats.Defence
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Defence = int(num64)

	// Line Parsing: s.LightWorldStats.WeaponStrength
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.WeaponStrength = int(num64)

	// Line Parsing: s.LightWorldStats.ArmorDefence
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.ArmorDefence = int(num64)

	// Line Parsing: s.LightWorldStats.Inventory
	for i1 := range 8 {
		// Line Parsing: s.LightWorldStats.Inventory[i1].Item
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.LightWorldStats.Inventory[i1].Item = LItem(num64)

		// Line Parsing: s.LightWorldStats.Inventory[i1].Phone
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.LightWorldStats.Inventory[i1].Phone = Phone(num64)

	}

	// Line Parsing: s.GlobalFlags
	for i0 := range 2500 {
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		s.GlobalFlags[i0] = text
	}

	// Line Parsing: s.ExtraGlobalFlags
	for i0 := range 7499 {
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		s.ExtraGlobalFlags[i0] = text
	}

	// Line Parsing: s.Plot
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if f64, err = strconv.ParseFloat(text, 64); err != nil {
		return nil, err
	}
	s.Plot = float64(f64)

	// Line Parsing: s.Room
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if f64, err = strconv.ParseFloat(text, 64); err != nil {
		return nil, err
	}
	s.Room = float64(f64)

	// Line Parsing: s.Time
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if f64, err = strconv.ParseFloat(text, 64); err != nil {
		return nil, err
	}
	s.Time = float64(f64)

	return s, nil
}

func ParseSave2Generated(scanner *bufio.Scanner) (*Save2, error) {
	s := &Save2{}
	var text string
	var err error
	var num int
	var num64 int64
	var f64 float64

	// Line Parsing: s.PlayerName
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	s.PlayerName = text

	// Line Parsing: s.CharName
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	s.CharName = text

	// Line Parsing: s.OtherNames
	for i0 := range 5 {
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		s.OtherNames[i0] = text
	}

	// Line Parsing: s.Characters
	for i0 := range 3 {
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num, err = strconv.Atoi(text); err != nil {
			return nil, err
		}
		s.Characters[i0] = Character(num)
	}

	// Line Parsing: s.Gold
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Gold = int(num64)

	// Line Parsing: s.XP
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.XP = int(num64)

	// Line Parsing: s.Level
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Level = int(num64)

	// Line Parsing: s.Inv
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Inv = int(num64)

	// Line Parsing: s.Invc
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Invc = int(num64)

	// Line Parsing: s.Darkzone
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num, err = strconv.Atoi(text); err != nil {
		return nil, err
	}
	s.Darkzone = num != 0

	// Line Parsing: s.CharactersStats
	for i0 := range 5 {
		// Line Parsing: s.CharactersStats[i0].HP
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].HP = int(num64)

		// Line Parsing: s.CharactersStats[i0].MaxHP
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].MaxHP = int(num64)

		// Line Parsing: s.CharactersStats[i0].Attack
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Attack = int(num64)

		// Line Parsing: s.CharactersStats[i0].Defence
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Defence = int(num64)

		// Line Parsing: s.CharactersStats[i0].Magic
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Magic = int(num64)

		// Line Parsing: s.CharactersStats[i0].Guts
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Guts = int(num64)

		// Line Parsing: s.CharactersStats[i0].Weapon
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].Weapon = Weapon(num64)

		// Line Parsing: s.CharactersStats[i0].Armor
		for i1 := range 2 {
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num, err = strconv.Atoi(text); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].Armor[i1] = Armor(num)
		}

		// Line Parsing: s.CharactersStats[i0].WeaponStyle
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		s.CharactersStats[i0].WeaponStyle = text

		// Line Parsing: s.CharactersStats[i0].ItemsStats
		for i1 := range 4 {
			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Attack
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Attack = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Defence
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Defence = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Magic
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Magic = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Bolts
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Bolts = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].GrazeAmount
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].GrazeAmount = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].GrazeSize
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].GrazeSize = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].BoltSpeed
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].BoltSpeed = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Special
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Special = int(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].Element
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].Element = Element(num64)

			// Line Parsing: s.CharactersStats[i0].ItemsStats[i1].ElementAmount
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if f64, err = strconv.ParseFloat(text, 64); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].ItemsStats[i1].ElementAmount = float64(f64)

		}

		// Line Parsing: s.CharactersStats[i0].Spells
		for i1 := range 12 {
			if text, err = nextLine(scanner); err != nil {
				return nil, err
			}
			if num, err = strconv.Atoi(text); err != nil {
				return nil, err
			}
			s.CharactersStats[i0].Spells[i1] = Spell(num)
		}

	}

	// Line Parsing: s.BoltSpeed
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.BoltSpeed = int(num64)

	// Line Parsing: s.GrazeAmount
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.GrazeAmount = int(num64)

	// Line Parsing: s.GrazeSize
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.GrazeSize = int(num64)

	// Line Parsing: s.Inventory
	// Line Parsing: s.Inventory.ItemsAndKeyItems
	for i1 := range 13 {
		// Line Parsing: s.Inventory.ItemsAndKeyItems[i1].Item
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.Inventory.ItemsAndKeyItems[i1].Item = Item(num64)

		// Line Parsing: s.Inventory.ItemsAndKeyItems[i1].KeyItem
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.Inventory.ItemsAndKeyItems[i1].KeyItem = KeyItem(num64)

	}

	// Line Parsing: s.Inventory.WeaponsAndArmors
	for i1 := range 48 {
		// Line Parsing: s.Inventory.WeaponsAndArmors[i1].Weapon
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.Inventory.WeaponsAndArmors[i1].Weapon = Weapon(num64)

		// Line Parsing: s.Inventory.WeaponsAndArmors[i1].Armor
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.Inventory.WeaponsAndArmors[i1].Armor = Armor(num64)

	}

	// Line Parsing: s.Inventory.PocketItems
	for i1 := range 72 {
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num, err = strconv.Atoi(text); err != nil {
			return nil, err
		}
		s.Inventory.PocketItems[i1] = Item(num)
	}

	// Line Parsing: s.Tension
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.Tension = int(num64)

	// Line Parsing: s.MaxTension
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.MaxTension = int(num64)

	// Line Parsing: s.LightWorldStats
	// Line Parsing: s.LightWorldStats.Weapon
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Weapon = LItem(num64)

	// Line Parsing: s.LightWorldStats.Armor
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Armor = LItem(num64)

	// Line Parsing: s.LightWorldStats.XP
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.XP = int(num64)

	// Line Parsing: s.LightWorldStats.Level
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Level = int(num64)

	// Line Parsing: s.LightWorldStats.Gold
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Gold = int(num64)

	// Line Parsing: s.LightWorldStats.HP
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.HP = int(num64)

	// Line Parsing: s.LightWorldStats.MaxHP
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.MaxHP = int(num64)

	// Line Parsing: s.LightWorldStats.Attack
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Attack = int(num64)

	// Line Parsing: s.LightWorldStats.Defence
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.Defence = int(num64)

	// Line Parsing: s.LightWorldStats.WeaponStrength
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.WeaponStrength = int(num64)

	// Line Parsing: s.LightWorldStats.ArmorDefence
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
		return nil, err
	}
	s.LightWorldStats.ArmorDefence = int(num64)

	// Line Parsing: s.LightWorldStats.Inventory
	for i1 := range 8 {
		// Line Parsing: s.LightWorldStats.Inventory[i1].Item
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.LightWorldStats.Inventory[i1].Item = LItem(num64)

		// Line Parsing: s.LightWorldStats.Inventory[i1].Phone
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		if num64, err = strconv.ParseInt(text, 10, 64); err != nil {
			return nil, err
		}
		s.LightWorldStats.Inventory[i1].Phone = Phone(num64)

	}

	// Line Parsing: s.GlobalFlags
	for i0 := range 2500 {
		if text, err = nextLine(scanner); err != nil {
			return nil, err
		}
		s.GlobalFlags[i0] = text
	}

	// Line Parsing: s.Plot
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if f64, err = strconv.ParseFloat(text, 64); err != nil {
		return nil, err
	}
	s.Plot = float64(f64)

	// Line Parsing: s.Room
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if f64, err = strconv.ParseFloat(text, 64); err != nil {
		return nil, err
	}
	s.Room = float64(f64)

	// Line Parsing: s.Time
	if text, err = nextLine(scanner); err != nil {
		return nil, err
	}
	if f64, err = strconv.ParseFloat(text, 64); err != nil {
		return nil, err
	}
	s.Time = float64(f64)

	return s, nil
}

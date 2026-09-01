//go:generate go run ../cmd/generator/main.go
package saves

import (
	"fmt"
	"strings"
	"time"
)

// ItemStats1 represents stats of a dark world item in chapter 1
type ItemStats1 struct {
	Attack      int
	Defence     int
	Magic       int
	Bolts       int
	GrazeAmount int
	GrazeSize   int
	BoltSpeed   int
	Special     int
}

// ItemStats2 represents stats of a dark world item in chapter 2 and further
type ItemStats2 struct {
	Attack        int
	Defence       int
	Magic         int
	Bolts         int
	GrazeAmount   int
	GrazeSize     int
	BoltSpeed     int
	Special       int
	Element       Element
	ElementAmount float64
}

// CharacterStats1 represents stats of a character in a dark world in chapter 1
type CharacterStats1 struct {
	HP          int
	MaxHP       int
	Attack      int
	Defence     int
	Magic       int
	Guts        int
	Weapon      Weapon
	Armor       [2]Armor
	WeaponStyle string
	ItemsStats  [4]ItemStats1
	Spells      [12]Spell
}

// CharacterStats2 represents stats of a character in a dark world in chapter 2 and further
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

// Inventory1 represents inventory contents in a dark world in chapter 1
type Inventory1 [13]struct {
	Item    Item
	KeyItem KeyItem
	Weapon  Weapon
	Armor   Armor
}

// Inventory2 represents inventory contents in a dark world in chapter 2 and further
type Inventory2 struct {
	ItemsAndKeyItems [13]struct {
		Item    Item
		KeyItem KeyItem
	}
	WeaponsAndArmors [48]struct {
		Weapon Weapon
		Armor  Armor
	}
	PocketItems [72]Item // Items in the storage
}

// LightInventory represents inventory contents in the light world
type LightInventory [8]struct {
	Item  LItem
	Phone Phone
}

// LightWorldStats represents stats of Kris is the light world
type LightWorldStats struct {
	Weapon         LItem
	Armor          LItem
	XP             int
	Level          int
	Gold           int
	HP             int
	MaxHP          int
	Attack         int
	Defence        int
	WeaponStrength int
	ArmorDefence   int
	Inventory      LightInventory
}

// GlobalFlags represents used game flags
type GlobalFlags [2500]string

// Save1 represents chapter 1 save file format
type Save1 struct {
	PlayerName       string
	CharName         string
	OtherNames       [5]string
	Characters       [3]Character
	Gold             int
	XP               int
	Level            int
	Inv              int
	Invc             int
	Darkzone         bool
	CharactersStats  [4]CharacterStats1
	BoltSpeed        int
	GrazeAmount      int
	GrazeSize        int
	Inventory        Inventory1
	Tension          int
	MaxTension       int
	LightWorldStats  LightWorldStats
	GlobalFlags      GlobalFlags
	ExtraGlobalFlags [7499]string
	Plot             float64
	Room             float64
	Time             float64
}

// Save2 represents chapter 2 and further save file format
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
	Darkzone        bool
	CharactersStats [5]CharacterStats2
	BoltSpeed       int
	GrazeAmount     int
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

// Save3 has the same structure as chapter 2 save file
type Save3 Save2

// Save4 has the same structure as chapter 2 save file
type Save4 Save2

// Save5 has the same structure as chapter 2 save file
type Save5 Save2

// Save must implement the String() method
type Save interface {
	String() string
}

// String returns a string representation of ItemStats1 structure
func (st ItemStats1) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintf(bp, "Attack: %d\n", st.Attack)
	fmt.Fprintf(bp, "Defence: %d\n", st.Defence)
	fmt.Fprintf(bp, "Magic: %d\n", st.Magic)
	fmt.Fprintf(bp, "Bolts: %d\n", st.Bolts)
	fmt.Fprintf(bp, "GrazeAmount: %d\n", st.GrazeAmount)
	fmt.Fprintf(bp, "GrazeSize: %d\n", st.GrazeSize)
	fmt.Fprintf(bp, "BoltSpeed: %d\n", st.BoltSpeed)
	fmt.Fprintf(bp, "Special: %d\n", st.Special)

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (st ItemStats2) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintf(bp, "Attack: %d\n", st.Attack)
	fmt.Fprintf(bp, "Defence: %d\n", st.Defence)
	fmt.Fprintf(bp, "Magic: %d\n", st.Magic)
	fmt.Fprintf(bp, "Bolts: %d\n", st.Bolts)
	fmt.Fprintf(bp, "GrazeAmount: %d\n", st.GrazeAmount)
	fmt.Fprintf(bp, "GrazeSize: %d\n", st.GrazeSize)
	fmt.Fprintf(bp, "BoltSpeed: %d\n", st.BoltSpeed)
	fmt.Fprintf(bp, "Special: %d\n", st.Special)
	fmt.Fprintf(bp, "Element: %f %s\n", st.ElementAmount, st.Element)

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (st CharacterStats1) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintf(bp, "HP: %d\n", st.HP)
	fmt.Fprintf(bp, "Max HP: %d\n", st.MaxHP)
	fmt.Fprintf(bp, "Attack: %d\n", st.Attack)
	fmt.Fprintf(bp, "Defence: %d\n", st.Defence)
	fmt.Fprintf(bp, "Magic: %d\n", st.Magic)
	fmt.Fprintf(bp, "Guts: %d\n", st.Guts)
	fmt.Fprintf(bp, "Weapon: %s\n", st.Weapon)
	fmt.Fprintf(bp, "First Armor: %s\n", st.Armor[0])
	fmt.Fprintf(bp, "Second Armor: %s\n", st.Armor[1])
	fmt.Fprintf(bp, "Weapon Style: %s\n", st.WeaponStyle)

	fmt.Fprintln(bp, "Item stats:")
	for i, itemStats := range st.ItemsStats {
		fmt.Fprintf(bp, "Item%d:\n", i+1)
		fmt.Fprintln(bp, itemStats)
	}

	fmt.Fprintln(bp, "Spells:")
	for i, spell := range st.Spells {
		if spell == SpellEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, spell)
	}

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (st CharacterStats2) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintf(bp, "HP: %d\n", st.HP)
	fmt.Fprintf(bp, "Max HP: %d\n", st.MaxHP)
	fmt.Fprintf(bp, "Attack: %d\n", st.Attack)
	fmt.Fprintf(bp, "Defence: %d\n", st.Defence)
	fmt.Fprintf(bp, "Magic: %d\n", st.Magic)
	fmt.Fprintf(bp, "Guts: %d\n", st.Guts)
	fmt.Fprintf(bp, "Weapon: %s\n", st.Weapon)
	fmt.Fprintf(bp, "First Armor: %s\n", st.Armor[0])
	fmt.Fprintf(bp, "Second Armor: %s\n", st.Armor[1])
	fmt.Fprintf(bp, "Weapon Style: %s\n", st.WeaponStyle)

	fmt.Fprintln(bp, "Item stats:")
	for i, itemStats := range st.ItemsStats {
		if i == 3 {
			continue
		}
		if i == 0 && st.Weapon == WeaponEmpty {
			continue
		}
		if i == 1 && st.Armor[0] == ArmorEmpty {
			continue
		}
		if i == 2 && st.Armor[1] == ArmorEmpty {
			continue
		}
		fmt.Fprintf(bp, "Item%d:\n", i+1)
		fmt.Fprintln(bp, itemStats)
	}

	fmt.Fprintln(bp, "Spells:")
	for i, spell := range st.Spells {
		if spell == SpellEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, spell)
	}

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (in Inventory1) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintln(bp, "Items:")
	for i, slot := range in {
		if slot.Item == ItemEmpty {
			continue
		}
		if slot.Item == ItemEndOfInventory {
			break
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.Item)
	}

	fmt.Fprintln(bp, "Key Items:")
	for i, slot := range in {
		if slot.KeyItem == KeyItemEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.KeyItem)
	}

	fmt.Fprintln(bp, "Weapons:")
	for i, slot := range in {
		if slot.Weapon == WeaponEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.Weapon)
	}

	fmt.Fprintln(bp, "Armors:")
	for i, slot := range in {
		if slot.Armor == ArmorEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.Armor)
	}

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (in Inventory2) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintln(bp, "Items:")
	for i, slot := range in.ItemsAndKeyItems {
		if slot.Item == ItemEmpty {
			continue
		}
		if slot.Item == ItemEndOfInventory {
			break
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.Item)
	}

	fmt.Fprintln(bp, "Key Items:")
	for i, slot := range in.ItemsAndKeyItems {
		if slot.KeyItem == KeyItemEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.KeyItem)
	}

	fmt.Fprintln(bp, "Weapons:")
	for i, slot := range in.WeaponsAndArmors {
		if slot.Weapon == WeaponEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.Weapon)
	}

	fmt.Fprintln(bp, "Armors:")
	for i, slot := range in.WeaponsAndArmors {
		if slot.Armor == ArmorEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.Armor)
	}

	fmt.Fprintln(bp, "Storage Items:")
	for i, item := range in.PocketItems {
		if item == ItemEmpty {
			continue
		}
		if item == ItemEndOfInventory {
			break
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, item)
	}

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (linv LightInventory) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintln(bp, "Items:")
	for i, slot := range linv {
		if slot.Item == LItemEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.Item)
	}

	fmt.Fprintln(bp, "Phone Numbers:")
	for i, slot := range linv {
		if slot.Phone == PhoneEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, slot.Phone)
	}

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (st LightWorldStats) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintf(bp, "Weapon: %s\n", st.Weapon)
	fmt.Fprintf(bp, "Armor: %s\n", st.Armor)
	fmt.Fprintf(bp, "Exp: %d\n", st.XP)
	fmt.Fprintf(bp, "Level: %d\n", st.Level)
	fmt.Fprintf(bp, "Gold: %d\n", st.Gold)
	fmt.Fprintf(bp, "HP: %d\n", st.HP)
	fmt.Fprintf(bp, "Max HP: %d\n", st.MaxHP)
	fmt.Fprintf(bp, "Attack: %d\n", st.Attack)
	fmt.Fprintf(bp, "Defence: %d\n", st.Defence)
	fmt.Fprintf(bp, "Weapon Strength: %d\n", st.WeaponStrength)
	fmt.Fprintf(bp, "Armor Defence: %d\n", st.ArmorDefence)
	fmt.Fprint(bp, st.Inventory.String())

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (gf GlobalFlags) String() string {
	var b strings.Builder
	bp := &b

	// SideB progression flag
	fmt.Fprintf(bp, "SideB Active: %t\n", gf[916] == "0")
	fmt.Fprintf(bp, "SideB Progression: %s\n", gf[915])

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (s *Save1) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintf(bp, "Player name: %s\n", s.PlayerName)
	fmt.Fprintf(bp, "Vessel name: %s\n", s.CharName)
	fmt.Fprintln(bp, "Other vessel names:")
	for _, name := range s.OtherNames {
		fmt.Fprintf(bp, "- %s\n", name)
	}
	fmt.Fprintln(bp, "Characters:")
	for i, character := range s.Characters {
		if character == CharacterEmpty {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, character)
	}

	fmt.Fprintf(bp, "Gold: %d\n", s.Gold)
	fmt.Fprintf(bp, "Exp: %d\n", s.XP)
	fmt.Fprintf(bp, "Level: %d\n", s.Level)
	fmt.Fprintf(bp, "Inv: %d\n", s.Inv)
	fmt.Fprintf(bp, "Invc: %d\n", s.Invc)
	fmt.Fprintf(bp, "Darkworld: %t", s.Darkzone)

	for i, character := range s.Characters {
		fmt.Fprintf(bp, "%s Stats\n", character)
		fmt.Fprintln(bp, s.CharactersStats[i+1])
	}

	fmt.Fprintf(bp, "Soul Speed: %d\n", s.BoltSpeed)
	fmt.Fprintf(bp, "Graze Amount: %d\n", s.GrazeAmount)
	fmt.Fprintf(bp, "Graze Size: %d\n", s.GrazeSize)

	fmt.Fprint(bp, s.Inventory.String())

	fmt.Fprintf(bp, "Tension: %d\n", s.Tension)
	fmt.Fprintf(bp, "Max Tension: %d\n", s.MaxTension)

	fmt.Fprintln(bp, "Light World Stats:")
	fmt.Fprintln(bp, s.LightWorldStats)

	fmt.Fprintln(bp, "Global Flags:")
	fmt.Fprintln(bp, s.GlobalFlags)
	fmt.Fprintf(bp, "Plot: %f\n", s.Plot)
	fmt.Fprintf(bp, "Room: %f\n", s.Room)
	fmt.Fprintf(bp, "Time: %s\n", time.Duration(uint64(s.Time/30))*time.Second)

	return b.String()
}

// String returns a string representation of ItemStats1 structure
func (s *Save2) String() string {
	var b strings.Builder
	bp := &b

	fmt.Fprintf(bp, "Player name: %s\n", s.PlayerName)
	fmt.Fprintf(bp, "Vessel name: %s\n", s.CharName)
	fmt.Fprintln(bp, "Other vessel names:")
	for _, name := range s.OtherNames {
		fmt.Fprintf(bp, "- %s\n", name)
	}
	fmt.Fprintln(bp, "Characters:")
	for i, character := range s.Characters {
		if character == 0 {
			continue
		}
		fmt.Fprintf(bp, "%d. %s\n", i+1, character)
	}

	fmt.Fprintf(bp, "Gold: %d\n", s.Gold)
	fmt.Fprintf(bp, "Exp: %d\n", s.XP)
	fmt.Fprintf(bp, "Level: %d\n", s.Level)
	fmt.Fprintf(bp, "Inv: %d\n", s.Inv)
	fmt.Fprintf(bp, "Invc: %d\n", s.Invc)
	fmt.Fprintf(bp, "Darkworld: %t\n\n", s.Darkzone)

	for i, character := range s.Characters {
		fmt.Fprintf(bp, "%s Stats\n", character)
		fmt.Fprintln(bp, s.CharactersStats[i+1])
	}

	fmt.Fprintf(bp, "Soul Speed: %d\n", s.BoltSpeed)
	fmt.Fprintf(bp, "Graze Amount: %d\n", s.GrazeAmount)
	fmt.Fprintf(bp, "Graze Size: %d\n", s.GrazeSize)

	fmt.Fprint(bp, s.Inventory.String())

	fmt.Fprintf(bp, "Tension: %d\n", s.Tension)
	fmt.Fprintf(bp, "Max Tension: %d\n", s.MaxTension)

	fmt.Fprintln(bp, "Light World Stats:")
	fmt.Fprintln(bp, s.LightWorldStats)

	fmt.Fprintln(bp, "Global Flags:")
	fmt.Fprintln(bp, s.GlobalFlags)

	fmt.Fprintf(bp, "Plot: %f\n", s.Plot)
	fmt.Fprintf(bp, "Room: %f\n", s.Room)
	fmt.Fprintf(bp, "Time: %s\n", time.Duration(uint64(s.Time/30))*time.Second)

	return b.String()
}

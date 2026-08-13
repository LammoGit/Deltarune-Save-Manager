package saves

import (
	"fmt"
	"strings"
	"testing"
)

// Test ItemStats1 structure to string conversion
func TestItemStatsChapter1StringConversion(t *testing.T) {
	is := ItemStats1{
		Attack:      1,
		Defence:     2,
		Magic:       3,
		Bolts:       4,
		GrazeAmount: 5,
		GrazeSize:   6,
		BoltSpeed:   7,
		Special:     8,
	}

	expected := `Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
`

	if is.String() != expected {
		t.Fatal("Chapter 1 item stats conversion to string doesn't match expected")
	}
}

// Test ItemStats2 structure to string conversion
func TestItemStatsChapter2StringConversion(t *testing.T) {
	is := ItemStats2{
		Attack:        1,
		Defence:       2,
		Magic:         3,
		Bolts:         4,
		GrazeAmount:   5,
		GrazeSize:     6,
		BoltSpeed:     7,
		Special:       8,
		Element:       ElementPuppetCat,
		ElementAmount: 0.1,
	}

	expected := `Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat
`

	if is.String() != expected {
		t.Fatal("Chapter 2 item stats conversion to string doesn't match expected")
	}
}

// Test CharacterStats1 structure to string conversion
func TestCharacterStatsChapter1StringConversion(t *testing.T) {
	is := ItemStats1{
		Attack:      1,
		Defence:     2,
		Magic:       3,
		Bolts:       4,
		GrazeAmount: 5,
		GrazeSize:   6,
		BoltSpeed:   7,
		Special:     8,
	}

	cs := CharacterStats1{
		HP:          1,
		MaxHP:       2,
		Attack:      3,
		Defence:     4,
		Magic:       5,
		Guts:        6,
		Weapon:      WeaponSpookysword,
		Armor:       [2]Armor{ArmorAmberCard, ArmorWhiteRibbon},
		WeaponStyle: "Normal",
		ItemsStats:  [4]ItemStats1{is, is, is, is},
		Spells:      [12]Spell{SpellACT, SpellEmpty, SpellRudeBuster},
	}

	var b strings.Builder

	fmt.Fprintln(&b, `HP: 1
Max HP: 2
Attack: 3
Defence: 4
Magic: 5
Guts: 6
Weapon: Spookysword
First Armor: AmberCard
Second Armor: WhiteRibbon
Weapon Style: Normal
Item stats:`)

	for i := range 4 {
		fmt.Fprintf(&b, "Item%d:\n", i+1)
		fmt.Fprintln(&b, is.String())
	}

	fmt.Fprintln(&b, `Spells:
1. ACT
3. RudeBuster`)

	if cs.String() != b.String() {
		t.Fatal("Chapter 1 character stats conversion to string doesn't match expected")
	}
}

// Test CharacterStats2 structure to string conversion
func TestCharacterStatsChapter2StringConversion(t *testing.T) {
	is := ItemStats2{
		Attack:        1,
		Defence:       2,
		Magic:         3,
		Bolts:         4,
		GrazeAmount:   5,
		GrazeSize:     6,
		BoltSpeed:     7,
		Special:       8,
		Element:       ElementPuppetCat,
		ElementAmount: 0.1,
	}

	cs := CharacterStats2{
		HP:          1,
		MaxHP:       2,
		Attack:      3,
		Defence:     4,
		Magic:       5,
		Guts:        6,
		Weapon:      WeaponSpookysword,
		Armor:       [2]Armor{ArmorAmberCard, ArmorWhiteRibbon},
		WeaponStyle: "Normal",
		ItemsStats:  [4]ItemStats2{is, is, is, is},
		Spells:      [12]Spell{SpellACT, SpellEmpty, SpellRudeBuster},
	}

	var b strings.Builder

	fmt.Fprintln(&b, `HP: 1
Max HP: 2
Attack: 3
Defence: 4
Magic: 5
Guts: 6
Weapon: Spookysword
First Armor: AmberCard
Second Armor: WhiteRibbon
Weapon Style: Normal
Item stats:`)

	for i := range 3 {
		fmt.Fprintf(&b, "Item%d:\n", i+1)
		fmt.Fprintln(&b, is.String())
	}

	fmt.Fprintln(&b, `Spells:
1. ACT
3. RudeBuster`)

	if cs.String() != b.String() {
		t.Fatal("Chapter 2 character stats conversion to string doesn't match expected")
	}
}

// Test Inventory1 structure to string conversion
func TestInventory1StringConversion(t *testing.T) {
	inv := Inventory1{
		{
			ItemCDBagel,
			KeyItemEgg,
			WeaponAbsorbAx,
			ArmorBlueRibbon,
		},
		{
			ItemEmpty,
			KeyItemEmpty,
			WeaponEmpty,
			ArmorEmpty,
		},
		{
			ItemDogDollar,
			KeyItemLancer,
			WeaponManeAx,
			ArmorDogWindow,
		},
	}

	expected := `Items:
1. CDBagel
3. DogDollar
Key Items:
1. Egg
3. Lancer
Weapons:
1. AbsorbAx
3. ManeAx
Armors:
1. BlueRibbon
3. DogWindow
`

	if inv.String() != expected {
		t.Fatal("Chapter 1 inventory conversion to string doesn't match expected")
	}
}

// Test Inventory2 structure to string conversion
func TestInventory2StringConversion(t *testing.T) {
	inv := Inventory2{
		[13]struct {
			Item    Item
			KeyItem KeyItem
		}{
			{ItemCDBagel, KeyItemEgg},
			{ItemEmpty, KeyItemEmpty},
			{ItemDogDollar, KeyItemLancer},
		},
		[48]struct {
			Weapon Weapon
			Armor  Armor
		}{
			{WeaponAbsorbAx, ArmorBlueRibbon},
			{WeaponEmpty, ArmorEmpty},
			{WeaponManeAx, ArmorDogWindow},
		},
		[72]Item{ItemCDBagel, ItemEmpty, ItemDogDollar},
	}

	expected := `Items:
1. CDBagel
3. DogDollar
Key Items:
1. Egg
3. Lancer
Weapons:
1. AbsorbAx
3. ManeAx
Armors:
1. BlueRibbon
3. DogWindow
Storage Items:
1. CDBagel
3. DogDollar
`

	if inv.String() != expected {
		t.Fatal("Chapter 2 inventory conversion to string doesn't match expected")
	}
}

// Test LightInventory structure to string conversion
func TestLightInventoryStringConversion(t *testing.T) {
	linv := LightInventory{
		{LItemHotChocolate, PhoneHome},
		{LItemEmpty, PhoneEmpty},
		{LItemCards, PhoneSans},
	}

	expected := `Items:
1. HotChocolate
3. Cards
Phone Numbers:
1. Home
3. Sans
`

	if linv.String() != expected {
		t.Fatal("Light world inventory conversion to string doesn't match expected")
	}
}

// Test LightWorldStats structure to string conversion
func TestLightWorldStatsStringConversion(t *testing.T) {
	linv := LightInventory{
		{LItemHotChocolate, PhoneHome},
		{LItemEmpty, PhoneEmpty},
		{LItemCards, PhoneSans},
	}

	lws := LightWorldStats{
		Weapon:         LItemPencil,
		Armor:          LItemBandage,
		XP:             1,
		Level:          2,
		Gold:           3,
		HP:             4,
		MaxHP:          5,
		Attack:         6,
		Defence:        7,
		WeaponStrength: 8,
		ArmorDefence:   9,
		Inventory:      linv,
	}

	var b strings.Builder
	fmt.Fprintln(&b, `Weapon: Pencil
Armor: Bandage
Exp: 1
Level: 2
Gold: 3
HP: 4
Max HP: 5
Attack: 6
Defence: 7
Weapon Strength: 8
Armor Defence: 9`)
	fmt.Fprint(&b, linv.String())

	if lws.String() != b.String() {
		t.Fatal("Light world stats conversion to string doesn't match expected")
	}
}

// Test GlobalFlags structure to string conversion
func TestGlobalFlagsStringConversion(t *testing.T) {
	gf := GlobalFlags{}

	gf[916] = "0"
	gf[915] = "20"

	expected := `SideB Active: true
SideB Progression: 20
`

	if gf.String() != expected {
		t.Fatal("Global flags conversion to string doesn't match expected")
	}
}

// Test Save1 structure to string conversion
func TestSaveChapter1StringConversion(t *testing.T) {
	is := ItemStats1{
		Attack:      1,
		Defence:     2,
		Magic:       3,
		Bolts:       4,
		GrazeAmount: 5,
		GrazeSize:   6,
		BoltSpeed:   7,
		Special:     8,
	}

	cs := CharacterStats1{
		HP:          1,
		MaxHP:       2,
		Attack:      3,
		Defence:     4,
		Magic:       5,
		Guts:        6,
		Weapon:      WeaponSpookysword,
		Armor:       [2]Armor{ArmorAmberCard, ArmorWhiteRibbon},
		WeaponStyle: "Normal",
		ItemsStats:  [4]ItemStats1{is, is, is, is},
		Spells:      [12]Spell{SpellACT, SpellEmpty, SpellRudeBuster},
	}

	inv := Inventory1{
		{
			ItemCDBagel,
			KeyItemEgg,
			WeaponAbsorbAx,
			ArmorBlueRibbon,
		},
		{
			ItemEmpty,
			KeyItemEmpty,
			WeaponEmpty,
			ArmorEmpty,
		},
		{
			ItemDogDollar,
			KeyItemLancer,
			WeaponManeAx,
			ArmorDogWindow,
		},
	}

	linv := LightInventory{
		{LItemHotChocolate, PhoneHome},
		{LItemEmpty, PhoneEmpty},
		{LItemCards, PhoneSans},
	}

	lws := LightWorldStats{
		Weapon:         LItemPencil,
		Armor:          LItemBandage,
		XP:             1,
		Level:          2,
		Gold:           3,
		HP:             4,
		MaxHP:          5,
		Attack:         6,
		Defence:        7,
		WeaponStrength: 8,
		ArmorDefence:   9,
		Inventory:      linv,
	}

	gf := GlobalFlags{}

	save := Save1{
		PlayerName:       "Player",
		CharName:         "Vessel",
		OtherNames:       [5]string{},
		Characters:       [3]Character{CharacterKris, CharacterSusie, CharacterRalsei},
		Gold:             1,
		XP:               2,
		Level:            3,
		Inv:              4,
		Invc:             5,
		Darkzone:         true,
		CharactersStats:  [4]CharacterStats1{cs, cs, cs, cs},
		BoltSpeed:        6,
		GrazeAmount:      7,
		GrazeSize:        8,
		Inventory:        inv,
		Tension:          9,
		MaxTension:       10,
		LightWorldStats:  lws,
		GlobalFlags:      gf,
		ExtraGlobalFlags: [7499]string{},
		Plot:             11.0,
		Room:             12.0,
		Time:             30.0 * 13,
	}

	expected := `Player name: Player
Vessel name: Vessel
Other vessel names:
- 
- 
- 
- 
- 
Characters:
1. Kris
2. Susie
3. Ralsei
Gold: 1
Exp: 2
Level: 3
Inv: 4
Invc: 5
Darkworld: trueKris Stats
HP: 1
Max HP: 2
Attack: 3
Defence: 4
Magic: 5
Guts: 6
Weapon: Spookysword
First Armor: AmberCard
Second Armor: WhiteRibbon
Weapon Style: Normal
Item stats:
Item1:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item2:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item3:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item4:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Spells:
1. ACT
3. RudeBuster

Susie Stats
HP: 1
Max HP: 2
Attack: 3
Defence: 4
Magic: 5
Guts: 6
Weapon: Spookysword
First Armor: AmberCard
Second Armor: WhiteRibbon
Weapon Style: Normal
Item stats:
Item1:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item2:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item3:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item4:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Spells:
1. ACT
3. RudeBuster

Ralsei Stats
HP: 1
Max HP: 2
Attack: 3
Defence: 4
Magic: 5
Guts: 6
Weapon: Spookysword
First Armor: AmberCard
Second Armor: WhiteRibbon
Weapon Style: Normal
Item stats:
Item1:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item2:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item3:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Item4:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8

Spells:
1. ACT
3. RudeBuster

Soul Speed: 6
Graze Amount: 7
Graze Size: 8
Items:
1. CDBagel
3. DogDollar
Key Items:
1. Egg
3. Lancer
Weapons:
1. AbsorbAx
3. ManeAx
Armors:
1. BlueRibbon
3. DogWindow
Tension: 9
Max Tension: 10
Light World Stats:
Weapon: Pencil
Armor: Bandage
Exp: 1
Level: 2
Gold: 3
HP: 4
Max HP: 5
Attack: 6
Defence: 7
Weapon Strength: 8
Armor Defence: 9
Items:
1. HotChocolate
3. Cards
Phone Numbers:
1. Home
3. Sans

Global Flags:
SideB Active: false
SideB Progression: 

Plot: 11.000000
Room: 12.000000
Time: 13s
`

	if save.String() != expected {
		t.Fatal("Chapter 1 save conversion to string doesn't match expected")
	}
}

// Test Save2 structure to string conversion
func TestSaveChapter2StringConversion(t *testing.T) {
	is := ItemStats2{
		Attack:        1,
		Defence:       2,
		Magic:         3,
		Bolts:         4,
		GrazeAmount:   5,
		GrazeSize:     6,
		BoltSpeed:     7,
		Special:       8,
		Element:       ElementPuppetCat,
		ElementAmount: 0.1,
	}

	cs := CharacterStats2{
		HP:          1,
		MaxHP:       2,
		Attack:      3,
		Defence:     4,
		Magic:       5,
		Guts:        6,
		Weapon:      WeaponSpookysword,
		Armor:       [2]Armor{ArmorAmberCard, ArmorWhiteRibbon},
		WeaponStyle: "Normal",
		ItemsStats:  [4]ItemStats2{is, is, is, is},
		Spells:      [12]Spell{SpellACT, SpellEmpty, SpellRudeBuster},
	}

	inv := Inventory2{
		[13]struct {
			Item    Item
			KeyItem KeyItem
		}{
			{ItemCDBagel, KeyItemEgg},
			{ItemEmpty, KeyItemEmpty},
			{ItemDogDollar, KeyItemLancer},
		},
		[48]struct {
			Weapon Weapon
			Armor  Armor
		}{
			{WeaponAbsorbAx, ArmorBlueRibbon},
			{WeaponEmpty, ArmorEmpty},
			{WeaponManeAx, ArmorDogWindow},
		},
		[72]Item{ItemCDBagel, ItemEmpty, ItemDogDollar},
	}

	linv := LightInventory{
		{LItemHotChocolate, PhoneHome},
		{LItemEmpty, PhoneEmpty},
		{LItemCards, PhoneSans},
	}

	lws := LightWorldStats{
		Weapon:         LItemPencil,
		Armor:          LItemBandage,
		XP:             1,
		Level:          2,
		Gold:           3,
		HP:             4,
		MaxHP:          5,
		Attack:         6,
		Defence:        7,
		WeaponStrength: 8,
		ArmorDefence:   9,
		Inventory:      linv,
	}

	gf := GlobalFlags{}

	save := Save2{
		PlayerName:      "Player",
		CharName:        "Vessel",
		OtherNames:      [5]string{},
		Characters:      [3]Character{CharacterKris, CharacterSusie, CharacterRalsei},
		Gold:            1,
		XP:              2,
		Level:           3,
		Inv:             4,
		Invc:            5,
		Darkzone:        true,
		CharactersStats: [5]CharacterStats2{cs, cs, cs, cs, cs},
		BoltSpeed:       6,
		GrazeAmount:     7,
		GrazeSize:       8,
		Inventory:       inv,
		Tension:         9,
		MaxTension:      10,
		LightWorldStats: lws,
		GlobalFlags:     gf,
		Plot:            11.0,
		Room:            12.0,
		Time:            30.0 * 13,
	}

	expected := `Player name: Player
Vessel name: Vessel
Other vessel names:
- 
- 
- 
- 
- 
Characters:
1. Kris
2. Susie
3. Ralsei
Gold: 1
Exp: 2
Level: 3
Inv: 4
Invc: 5
Darkworld: true

Kris Stats
HP: 1
Max HP: 2
Attack: 3
Defence: 4
Magic: 5
Guts: 6
Weapon: Spookysword
First Armor: AmberCard
Second Armor: WhiteRibbon
Weapon Style: Normal
Item stats:
Item1:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Item2:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Item3:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Spells:
1. ACT
3. RudeBuster

Susie Stats
HP: 1
Max HP: 2
Attack: 3
Defence: 4
Magic: 5
Guts: 6
Weapon: Spookysword
First Armor: AmberCard
Second Armor: WhiteRibbon
Weapon Style: Normal
Item stats:
Item1:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Item2:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Item3:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Spells:
1. ACT
3. RudeBuster

Ralsei Stats
HP: 1
Max HP: 2
Attack: 3
Defence: 4
Magic: 5
Guts: 6
Weapon: Spookysword
First Armor: AmberCard
Second Armor: WhiteRibbon
Weapon Style: Normal
Item stats:
Item1:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Item2:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Item3:
Attack: 1
Defence: 2
Magic: 3
Bolts: 4
GrazeAmount: 5
GrazeSize: 6
BoltSpeed: 7
Special: 8
Element: 0.100000 PuppetCat

Spells:
1. ACT
3. RudeBuster

Soul Speed: 6
Graze Amount: 7
Graze Size: 8
Items:
1. CDBagel
3. DogDollar
Key Items:
1. Egg
3. Lancer
Weapons:
1. AbsorbAx
3. ManeAx
Armors:
1. BlueRibbon
3. DogWindow
Storage Items:
1. CDBagel
3. DogDollar
Tension: 9
Max Tension: 10
Light World Stats:
Weapon: Pencil
Armor: Bandage
Exp: 1
Level: 2
Gold: 3
HP: 4
Max HP: 5
Attack: 6
Defence: 7
Weapon Strength: 8
Armor Defence: 9
Items:
1. HotChocolate
3. Cards
Phone Numbers:
1. Home
3. Sans

Global Flags:
SideB Active: false
SideB Progression: 

Plot: 11.000000
Room: 12.000000
Time: 13s
`

	if save.String() != expected {
		t.Fatal("Chapter 2 save conversion to string doesn't match expected")
	}
}

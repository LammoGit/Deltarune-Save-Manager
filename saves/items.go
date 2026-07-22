package saves

//go:generate stringer -type=Item -trimprefix=Item
type Item int

const (
	ItemEmpty Item = iota
	ItemDarkerCandy
	ItemReviveMint
	ItemGlowshard
	ItemManual
	ItemBrokenCake
	ItemTopCake
	ItemSpinCake
	ItemDarkburger
	ItemLancerCookie
	ItemGigaSalad
	ItemClubswich
	ItemHeartsDonut
	ItemChocDiamond
	ItemFavSandwich
	ItemRouxlsRoux
	ItemCDBagel
	ItemMannequin
	ItemRottenTea
	ItemRottenTea19
	ItemRottenTea20
	ItemRottenTea21
	ItemDDBurger
	ItemLightCandy
	ItemButJuice
	ItemSpagettiCode
	ItemJavaCookie
	ItemTensionBit
	ItemTensionGem
	ItemTensionMax
	ItemReviveDust
	ItemReviveBrite
	ItemSPoison
	ItemDogDollar
	ItemTVDinner
	ItemPipis
	ItemFlatSoda
	ItemTVSlop
	ItemExecBuffet
	ItemDeluxeDinner
	ItemPunchBowl
	ItemFlavigne
	ItemGreenTea
	ItemOrangeJuice
	ItemAncientSweet = iota + 16
	ItemRhapsotea
	ItemScarlixir
	ItemBitterTear
	ItemSchadenbrot
	ItemTreeCake
	ItemSPotion
	ItemRawMoon
	ItemPhanta
	ItemFlowerSoda
	ItemShikacola
)

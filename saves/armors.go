package saves

//go:generate stringer -type=Armor -trimprefix=Armor
type Armor int

const (
	ArmorEmpty Armor = iota
	ArmorAmberCard
	ArmorDiceBrace
	ArmorPinkRibbon
	ArmorWhiteRibbon
	ArmorIronShackle
	ArmorMouseToken
	ArmorJevilstail
	ArmorSilverCard
	ArmorTwinRibbon
	ArmorGlowWrist
	ArmorChainMail
	ArmorBShotBowtie
	ArmorSpikeBand
	ArmorSilverWatch
	ArmorTensionBow
	ArmorMannequin
	ArmorDarkGoldBand
	ArmorSkyMantle
	ArmorSpikeShackle
	ArmorFrayedBowtie
	ArmorDealmaker
	ArmorRoyalPin
	ArmorShadowMantle
	ArmorLodeStone
	ArmorGingerGuard
	ArmorBlueRibbon
	ArmorTennaTie
	ArmorMonarchRBN = iota + 2
	ArmorTrueTie
	ArmorDogWindow
	ArmorRedRibbon
	ArmorNetskieHat
	ArmorSethSpecs
	ArmorYellowHat
	ArmorOGlove
	ArmorGreenApron
	ArmorWaferguard = iota + 13
	ArmorMysticBand
	ArmorPowerBand
	ArmorPrincessRBN
	ArmorGoldWidow
)

type ArmorInfo struct {
	Name          string
	Description   string
	SusieMessage  string
	RalseiMessage string
	NoelleMessage string
	Attack        int
	Defence       int
	Magic         int
	BoltSpeed     int
	Grazeamt      int
	GrazeSize     int
	KrisCanWear   bool
	SusieCanWear  bool
	RalseiCanWear bool
	NoelleCanWear bool
	Ability       string
	AbilityIcon   int
	Icon          int
	Price         int
	Element       Element
	ElementAmount float64
}

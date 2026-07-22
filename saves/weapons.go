package saves

//go:generate stringer -type=Weapon -trimprefix=Weapon
type Weapon int

const (
	WeaponEmpty Weapon = iota
	WeaponWoodBlade
	WeaponManeAx
	WeaponRedScarf
	WeaponEverybodyWeapon
	WeaponSpookysword
	WeaponBraveAx
	WeaponDevilsknife
	WeaponTrefoil
	WeaponRagger
	WeaponDaintyScarf
	WeaponTwistedSwd
	WeaponSnowRing
	WeaponThornRing
	WeaponBounceBlade
	WeaponCheerScarf
	WeaponMechaSaber
	WeaponAutoAxe
	WeaponFiberScarf
	WeaponRagger2
	WeaponBrokenSwd
	WeaponPuppetScarf
	WeaponFreezeRing
	WeaponSaber10
	WeaponToxicAxe
	WeaponFlexScarf
	WeaponBlackShard
	WeaponWoodBlade2 = iota + 3
	WeaponThatchet
	WeaponBlueShoes
	WeaponAquaKnife
	WeaponFloweryScarf
	WeaponBrokenScarf
	WeaponGildedRose
	WeaponMistleWP
	WeaponJingleBlade = iota + 15
	WeaponScarfMark
	WeaponJusticeAxe
	WeaponWinglade
	WeaponAbsorbAx
)

type WeaponInfo struct {
	Name           string
	Description    string
	SusieMessage   string
	RalseiMessage  string
	NoelleMessagge string
	Attack         int
	Defence        int
	Magic          int
	Style          string
	Grazeamt       int
	GrazeSize      int
	KrisCanWear    bool
	SusieCanWear   bool
	RalseiCanWear  bool
	NoelleCanWear  bool
	AbilityIcon    int
	Ability        string
	Price          int
}

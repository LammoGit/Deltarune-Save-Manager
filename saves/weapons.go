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
)
const (
	WeaponWoodBlade2 Weapon = iota + 30
	WeaponThatchet
	WeaponBlueShoes
	WeaponAquaKnife
	WeaponFloweryScarf
	WeaponBrokenScarf
	WeaponGildedRose
	WeaponMistleWP
)

const (
	WeaponJingleBlade Weapon = iota + 50
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
	GrazeAmount    int
	GrazeSize      int
	KrisCanWear    bool
	SusieCanWear   bool
	RalseiCanWear  bool
	NoelleCanWear  bool
	AbilityIcon    int
	Ability        string
	Price          int
}

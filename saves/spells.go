package saves

//go:generate stringer -type=Spell -trimprefix=Spell
type Spell int

const (
	SpellEmpty Spell = iota
	SpellRudeSword
	SpellHealPlayer
	SpellPacify
	SpellRudeBuster
	SpellRedBuster
	SpellDualHeal
	SpellACT
	SpellSleepMist
	SpellIceShock
	SpellSnowGrave
	SpellHeal
	SpellReviveSong
	SpellScythemare
)

//go:generate stringer -type=SpellTarget -trimprefix=SpellTarget
type SpellTarget int

const (
	SpellTargetAll SpellTarget = iota
	SpellTargetAlly
	SpellTargetEnemy
)

type SpellInfo struct {
	Name              string
	BattleName        string
	Description       string
	BattleDescription string
	Target            SpellTarget
	Cost              int
	SpellUsable       int
}

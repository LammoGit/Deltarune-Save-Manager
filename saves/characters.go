package saves

//go:generate stringer -type=Character -trimprefix=Character
type Character int

const (
	CharacterEmpty Character = iota
	CharacterKris
	CharacterSusie
	CharacterRalsei
	CharacterNoelle
)

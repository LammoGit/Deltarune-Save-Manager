package saves

//go:generate stringer -type=Element -trimprefix=Element
type Element int

const (
	ElementEmpty       Element = 0
	ElementElectroHoly         = 1
	ElementElectro
	ElementHoly
	ElementDarkStar = 5
	ElementDark
	ElementStar
	ElementPuppetCat = 6
	ElementPuppet
	ElementCat
	ElementMouse = 7
)

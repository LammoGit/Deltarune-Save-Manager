package saves

//go:generate stringer -type=Element -trimprefix=Element
type Element int

const (
	ElementEmpty       Element = 0
	ElementElectroHoly Element = 1
	ElementElectro
	ElementHoly
	ElementDarkStar Element = 5
	ElementDark
	ElementStar
	ElementPuppetCat Element = 6
	ElementPuppet
	ElementCat
	ElementMouse Element = 7
)

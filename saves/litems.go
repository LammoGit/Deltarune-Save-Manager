package saves

//go:generate stringer -type=LItem -trimprefix=LItem
type LItem int

const (
	LItemEmpty LItem = iota
	LItemHotChocolate
	LItemPencil
	LItemBandage
	LItemBouquet
	LItemBallOfJunk
	LItemHalloweenPencil
	LItemLuckyPencil
	LItemEgg
	LItemCards
	LItemBoxOfHeartCandy
	LItemGlass
	LItemEraser
	LItemMechPencil
	LItemWristwatch
	LItemHolidayPencil
	LItemCactusNeedle
	LItemBlackShard
	LItemQuillPen
	LItemHoneyToast
	LItemBread
	LItemSeeds
	LItemPencil2
	LItemPetal
)

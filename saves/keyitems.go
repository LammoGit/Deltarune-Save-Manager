package saves

//go:generate stringer -type=KeyItem -trimprefix=KeyItem
type KeyItem int

const (
	KeyItemEmpty KeyItem = iota
	KeyItemPhone
	KeyItemEgg
	KeyItemBrokenCake
	KeyItemBrokenKeyA
	KeyItemDoorKey
	KeyItemBrokenKeyB
	KeyItemBrokenKeyC
	KeyItemLancer
	KeyItemRouxlsKaard
	KeyItemEmptyDisk
	KeyItemLoadedDisk
	KeyItemKeyGen
	KeyItemShadowCrystal
	KeyItemStarwalker
	KeyItemPureCrystal
	KeyItemOddController
	KeyItemBackstagePass
	KeyItemTripTicket
	KeyItemLancerCon
	KeyItemScissors
	KeyItemYellowShred
	KeyItemBootOil
	KeyItemRedSplatter
	KeyItemBromideR
	KeyItemPetalFeather
	KeyItemPerpBook
	KeyItemBlueString
	KeyItemTrainPlan
	KeyItemYellowKey
	KeyItemSheetMusic
	KeyItemClaimbClaws
	KeyItemMysteryKey
	KeyItemBromideF
)

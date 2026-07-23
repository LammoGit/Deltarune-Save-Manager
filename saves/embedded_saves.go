package saves

import (
	_ "embed"

	"github.com/LammoGit/Deltarune-Save-Manager/utils"
)

//go:embed 1_savefile
var chapter1bytes []byte

//go:embed 2_savefile
var chapter2bytes []byte

//go:embed 3_savefile
var chapter3bytes []byte

//go:embed 4_savefile
var chapter4bytes []byte

//go:embed 5_savefile
var chapter5bytes []byte

func getExampleSaveBytesForChapter(chapter int) ([]byte, error) {
	switch chapter {
	case 1:
		return chapter1bytes, nil
	case 2:
		return chapter2bytes, nil
	case 3:
		return chapter3bytes, nil
	case 4:
		return chapter4bytes, nil
	case 5:
		return chapter5bytes, nil
	}
	return []byte{}, utils.ErrChapterNotSupported
}

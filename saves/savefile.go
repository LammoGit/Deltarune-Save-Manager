package saves

import (
	"bytes"
	"strings"
	"strconv"
	"io"
	"os"
)

type Save []string

func (s Save) Name() string {
	return s[0]
}

func (s Save) CharName() string {
	return s[1]
}

func (s *Save) Edit(props map[string]string) error {
	return nil
}

func ParseSaveBytes(buf []byte) (Save, error) {
	buf = bytes.ReplaceAll(buf, []byte("\r\n"), []byte("\n"))
	text := string(bytes.Trim(buf, "\n"))
	return Save(strings.Split(text, "\n")), nil
}

func LoadSave(path string) (save Save, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	save, err = ParseSaveBytes(content)
	return
}

func ParseSaveReader(r io.Reader) (save Save, err error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return 
	}
	save, err = ParseSaveBytes(buf)
	return
}

func getExampleSaveBytesForChapter(chapter int) []byte {
	return []byte(strconv.Itoa(chapter))
}

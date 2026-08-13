package saves

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/LammoGit/Deltarune-Save-Manager/utils"
)

// parseSaveLine parses individual line in a parsed save file
func parseSaveLine(scanner *bufio.Scanner, lineNum int, kind reflect.Kind, v reflect.Value) (int, error) {
	switch kind {
	case reflect.String:
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return lineNum, err
			}
			return lineNum, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return lineNum, utils.ErrValueCannotBeSet
		}

		v.SetString(scanner.Text())
		lineNum++
	case reflect.Bool:
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return lineNum, err
			}
			return lineNum, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return lineNum, utils.ErrValueCannotBeSet
		}

		text := strings.TrimSpace(scanner.Text())
		num, err := strconv.Atoi(text)
		if err != nil {
			return lineNum, fmt.Errorf("%w: failed to parse to a boolean value %q on line %d", utils.ErrWrongLineType, text, lineNum)
		}
		v.SetBool(num != 0)
		lineNum++
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return lineNum, err
			}
			return lineNum, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return lineNum, utils.ErrValueCannotBeSet
		}

		var bitSize int
		switch kind {
		case reflect.Int:
			bitSize = 0
		case reflect.Int8:
			bitSize = 8
		case reflect.Int16:
			bitSize = 16
		case reflect.Int32:
			bitSize = 32
		case reflect.Int64:
			bitSize = 64
		}

		text := strings.TrimSpace(scanner.Text())
		num, err := strconv.ParseInt(text, 10, bitSize)
		if err != nil {
			return lineNum, fmt.Errorf("%w: failed to parse to a signed integer value %q on line %d", utils.ErrWrongLineType, text, lineNum)
		}

		v.SetInt(int64(num))
		lineNum++
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return lineNum, err
			}
			return lineNum, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return lineNum, utils.ErrValueCannotBeSet
		}

		var bitSize int
		switch kind {
		case reflect.Uint:
			bitSize = 0
		case reflect.Uint8:
			bitSize = 8
		case reflect.Uint16:
			bitSize = 16
		case reflect.Uint32:
			bitSize = 32
		case reflect.Uint64:
			bitSize = 64
		}

		text := strings.TrimSpace(scanner.Text())
		num, err := strconv.ParseUint(text, 10, bitSize)
		if err != nil {
			return lineNum, fmt.Errorf("%w: failed to parse to an unsigned integer value %q on line %d", utils.ErrWrongLineType, text, lineNum)
		}

		v.SetUint(uint64(num))
		lineNum++
	case reflect.Float32, reflect.Float64:
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return lineNum, err
			}
			return lineNum, utils.ErrShortSaveFile
		}

		if !v.CanSet() {
			return lineNum, utils.ErrValueCannotBeSet
		}

		var bitSize int
		switch kind {
		case reflect.Float32:
			bitSize = 32
		case reflect.Float64:
			bitSize = 64
		}

		text := strings.TrimSpace(scanner.Text())
		num, err := strconv.ParseFloat(text, bitSize)
		if err != nil {
			return lineNum, fmt.Errorf("%w: failed to parse to a floating point number value %q on line %d", utils.ErrWrongLineType, text, lineNum)
		}

		v.SetFloat(float64(num))
		lineNum++
	case reflect.Struct:
		for _, fieldValue := range v.Fields() {
			fieldKind := fieldValue.Kind()

			var err error
			lineNum, err = parseSaveLine(scanner, lineNum, fieldKind, fieldValue)
			if err != nil {
				return lineNum, err
			}
		}
	case reflect.Array:
		elemKind := v.Type().Elem().Kind()
		for i := 0; i < v.Len(); i++ {
			var err error
			lineNum, err = parseSaveLine(scanner, lineNum, elemKind, v.Index(i))
			if err != nil {
				return lineNum, err
			}
		}
	}
	return lineNum, nil
}

// parseSaveField parses individual Save object's field
func parseSaveField(w io.Writer, kind reflect.Kind, v reflect.Value) error {
	switch kind {
	case reflect.String:
		fmt.Fprintf(w, "%s\n", v.String())
	case reflect.Bool:
		if v.Bool() {
			fmt.Fprintln(w, "1")
		} else {
			fmt.Fprintln(w, "0")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(w, "%d\n", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(w, "%d\n", v.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(w, "%f\n", v.Float())
	case reflect.Struct:
		for _, fieldValue := range v.Fields() {
			fieldKind := fieldValue.Kind()

			err := parseSaveField(w, fieldKind, fieldValue)
			if err != nil {
				return err
			}
		}
	case reflect.Array:
		elemKind := v.Type().Elem().Kind()
		for i := 0; i < v.Len(); i++ {
			err := parseSaveField(w, elemKind, v.Index(i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Save2Bytes returns bytes from a given Save object
func Save2Bytes(save Save) ([]byte, error) {
	var buf bytes.Buffer

	switch save.(type) {
	case *Save1:
	case *Save2:
	case nil:
		return nil, errors.New("save is a nil pointer")
	default:
		return nil, utils.ErrChapterNotSupported
	}

	saveValue := reflect.ValueOf(save).Elem()
	err := parseSaveField(&buf, saveValue.Kind(), saveValue)
	return buf.Bytes(), err
}

// ParseSaveBytes returns a Save object from a save file bytes and save's chapter
func ParseSaveBytes(buf []byte, chapter int) (Save, error) {
	reader := bytes.NewReader(buf)
	scanner := bufio.NewScanner(reader)

	var save Save
	var err error
	switch chapter {
	case 1:
		save = &Save1{}
	case 2, 3, 4, 5:
		save = &Save2{}
	default:
		return save, utils.ErrChapterNotSupported
	}
	saveValue := reflect.ValueOf(save).Elem()
	_, err = parseSaveLine(scanner, 0, saveValue.Kind(), saveValue)
	return save, err
}

// LoadSave returns a Save object from a save file path and save's chapter
func LoadSave(path string, chapter int) (Save, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseSaveReader(file, chapter)
}

// ParseSaveReader returns a Save object from a save file Reader object and save's chapter
func ParseSaveReader(r io.Reader, chapter int) (Save, error) {
	scanner := bufio.NewScanner(r)
	var save Save
	var err error
	switch chapter {
	case 1:
		save = &Save1{}
	case 2, 3, 4, 5:
		save = &Save2{}
	default:
		return save, utils.ErrChapterNotSupported
	}
	saveValue := reflect.ValueOf(save).Elem()
	_, err = parseSaveLine(scanner, 0, saveValue.Kind(), saveValue)
	return save, err
}

package saves

import (
	"bufio"
	"bytes"
	"dsm/utils"
	_ "embed"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

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

func LoadSave(path string, chapter int) (Save, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseSaveReader(file, chapter)
}

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

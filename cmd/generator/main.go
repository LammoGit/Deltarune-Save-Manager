package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/LammoGit/Deltarune-Save-Manager/saves"
)

var parserCode = `package saves

import (
	"bufio"
	"strconv"
	"strings"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/LammoGit/Deltarune-Save-Manager/utils"
)

// nextLine fetches the next trimmed line from the scanner
func nextLine(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", utils.ErrShortSaveFile
	}
	return strings.TrimSpace(scanner.Text()), nil
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
	switch chapter {
	case 1:
		return ParseSave1Generated(scanner)
	case 2, 3, 4, 5:
		return ParseSave2Generated(scanner)
	default:
		return save, utils.ErrChapterNotSupported
	}
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
	switch chapter {
	case 1:
		return ParseSave1Generated(scanner)
	case 2, 3, 4, 5:
		return ParseSave2Generated(scanner)
	default:
		return save, utils.ErrChapterNotSupported
	}
}
`

func main() {
	var buf bytes.Buffer

	buf.WriteString(parserCode)

	generateParserForType(&buf, "Save1", reflect.TypeFor[saves.Save1]())
	generateParserForType(&buf, "Save2", reflect.TypeFor[saves.Save2]())

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(buf.String())
		log.Fatalf("failed to format generated code: %v", err)
	}

	err = os.WriteFile("parser_gen.go", formatted, 0644)
	if err != nil {
		log.Fatalf("failed to write generated file: %v", err)
	}
	fmt.Println("Successfully generated parsers")
}

// generateParserForType walks a struct recursively or sequentially and generates pure Go code
func generateParserForType(w *bytes.Buffer, structName string, t reflect.Type) {
	fmt.Fprintf(
		w,
		"\nfunc Parse%sGenerated(scanner *bufio.Scanner) (*%s, error) {\n",
		structName,
		structName,
	)
	fmt.Fprintf(w, "\ts := &%s{}\n", structName)
	w.WriteString("\tvar text string\n\tvar err error\n\tvar num int\n\tvar num64 int64\n\tvar f64 float64\n\n")

	writeFields(w, t, "s.")

	w.WriteString("\treturn s, nil\n}\n")
}

// writeFields processes struct fields and maps them into native sequential Go code
func writeFields(w *bytes.Buffer, t reflect.Type, prefix string, depthVar ...int) {
	depth := 0
	if len(depthVar) > 0 {
		depth = depthVar[0]
	}

	for field := range t.Fields() {
		fieldName := prefix + field.Name
		kind := field.Type.Kind()

		fmt.Fprintf(w, "\t// Line Parsing: %s\n", fieldName)

		switch kind {
		case reflect.String:
			fmt.Fprintf(w, "\tif text, err = nextLine(scanner); err != nil { return nil, err }\n")
			fmt.Fprintf(w, "\t%s = text\n", fieldName)

		case reflect.Bool:
			fmt.Fprintf(w, "\tif text, err = nextLine(scanner); err != nil { return nil, err }\n")
			fmt.Fprintf(w, "\tif num, err = strconv.Atoi(text); err != nil { return nil, err }\n")
			fmt.Fprintf(w, "\t%s = num != 0\n", fieldName)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Fprintf(w, "\tif text, err = nextLine(scanner); err != nil { return nil, err }\n")
			fmt.Fprintf(w, "\tif num64, err = strconv.ParseInt(text, 10, 64); err != nil { return nil, err }\n")
			fmt.Fprintf(w, "\t%s = %s(num64)\n", fieldName, strings.TrimPrefix(field.Type.String(), "saves."))

		case reflect.Float32, reflect.Float64:
			fmt.Fprintf(w, "\tif text, err = nextLine(scanner); err != nil { return nil, err }\n")
			fmt.Fprintf(w, "\tif f64, err = strconv.ParseFloat(text, 64); err != nil { return nil, err }\n")
			fmt.Fprintf(w, "\t%s = %s(f64)\n", fieldName, strings.TrimPrefix(field.Type.String(), "saves."))

		case reflect.Struct:
			writeFields(w, field.Type, fieldName+".", depth+1)

		case reflect.Array:
			elemKind := field.Type.Elem().Kind()
			arrayLen := field.Type.Len()

			fmt.Fprintf(w, "\tfor i%d := range %d {\n", depth, arrayLen)

			switch elemKind {
			case reflect.String:
				fmt.Fprintf(w, "\t\tif text, err = nextLine(scanner); err != nil { return nil, err }\n")
				fmt.Fprintf(w, "\t\t%s[i%d] = text\n", fieldName, depth)
			case reflect.Int:
				fmt.Fprintf(w, "\t\tif text, err = nextLine(scanner); err != nil { return nil, err }\n")
				fmt.Fprintf(w, "\t\tif num, err = strconv.Atoi(text); err != nil { return nil, err }\n")
				fmt.Fprintf(w, "\t\t%s[i%d] = %s(num)\n", fieldName, depth, strings.TrimPrefix(field.Type.Elem().String(), "saves."))
			case reflect.Struct:
				writeFields(w, field.Type.Elem(), fmt.Sprintf("%s[i%d].", fieldName, depth), depth+1)
			default:
				if field.Type.Elem().Kind() == reflect.Struct {
					writeFields(w, field.Type.Elem(), fmt.Sprintf("%s[i%d].", fieldName, depth), depth+1)
				}
			}
			fmt.Fprintf(w, "\t}\n")
		}
		w.WriteString("\n")
	}
}

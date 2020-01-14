package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	codeCharacters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type sourcesFlag []string

func (i *sourcesFlag) String() string {
	return strings.Join(*i, ",")
}

func (i *sourcesFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type errorInfo struct {
	Code        uint
	Prefix      string
	Label       string
	Description string
	Source      string
}

// This was copied from https://golang.org/pkg/sort/ and modified.

type lessFunc func(p1, p2 *errorInfo) bool

type multiSorter struct {
	errors []errorInfo
	less   []lessFunc
}

func (ms *multiSorter) Sort(errors []errorInfo) {
	ms.errors = errors
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (ms *multiSorter) Len() int {
	return len(ms.errors)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.errors[i], ms.errors[j] = ms.errors[j], ms.errors[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.errors[i], &ms.errors[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}

func sortByPrefix(e1, e2 *errorInfo) bool {
	return e1.Prefix < e2.Prefix
}

func sortByLabel(e1, e2 *errorInfo) bool {
	return e1.Label < e2.Label
}

func sortByCode(e1, e2 *errorInfo) bool {
	return e1.Code < e2.Code
}

func main() {
	var sources sourcesFlag
	var pkg string
	var out string
	var skipValidate bool
	var lineNumbersMustMatch bool
	var testOut string

	flag.BoolVar(&skipValidate, "skip-validate", false, "Skip validation of sources.")
	flag.StringVar(&pkg, "package", "errors", "The package of the source file(s).")
	flag.StringVar(&out, "out", "generated.go", "The name of the source file generated by the program.")
	flag.Var(&sources, "source", "The source file to read errors from.")
	flag.BoolVar(&lineNumbersMustMatch, "match-line-numbers", false, "Line numbers must batch codes in packages.")
	flag.StringVar(&testOut, "test-out", "generated_test.go", "The name of the test source file generated by the program.")

	flag.Parse()

	var errors []errorInfo
	allErrors := map[string]errorInfo{}

	for _, source := range sources {
		sourceErrors, err := collectErrors(source, lineNumbersMustMatch)
		if err != nil {
			log.Fatalf("error processing source %s: %v\n", source, err)
		}

		for _, sourceError := range sourceErrors {

			if !skipValidate {
				foundErr, ok := allErrors[sourceError.Serialized()]
				if ok {
					panic(fmt.Sprintf("Duplicate error found: %s[%s-%d] %s[%s-%d]",
						foundErr.Source, foundErr.Prefix, foundErr.Code,
						sourceError.Source, sourceError.Prefix, sourceError.Code,
					))
				}
				allErrors[sourceError.Serialized()] = sourceError
			}

			errors = append(errors, sourceError)
		}
	}

	OrderedBy(sortByPrefix, sortByCode).Sort(errors)
	if err := writeGeneratedSource(out, pkg, errors); err != nil {
		panic(err)
	}
	if len(testOut) > 0 {
		if err := writeGeneratedSourceTest(testOut, pkg, errors); err != nil {
			panic(err)
		}
	}
}

func collectErrors(source string, lineNumbersMustMatch bool) ([]errorInfo, error) {
	var errors []errorInfo

	reader, err := os.Open(source)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(reader)
	csvReader.Comment = '#'
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = 4

	lineNumber := 0

	for {
		record, err := csvReader.Read()
		lineNumber = lineNumber + 1
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		code, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid record %d: %s", lineNumber, strings.Join(record, ","))
		}

		if lineNumbersMustMatch && code != lineNumber {
			return nil, fmt.Errorf("line number and record code do not match: %d != %d", lineNumber, code)
		}

		errors = append(errors, errorInfo{
			Code:        uint(code),
			Prefix:      record[1],
			Label:       record[2],
			Description: record[3],
			Source:      source,
		})

	}

	return errors, nil
}

func writeGeneratedSource(out, pkg string, errors []errorInfo) error {
	dest := os.Stdout
	if out != "stdout" {
		var err error
		dest, err = os.Create(out)
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := dest.Close(); closeErr != nil {
				log.Println("error closing source file", out, ":", closeErr.Error())
			}
		}()
	}

	return packageTemplate.Execute(dest, struct {
		Package   string
		Timestamp time.Time
		Codes     []errorInfo
	}{
		Package:   pkg,
		Timestamp: time.Now(),
		Codes:     errors,
	})
}

func writeGeneratedSourceTest(out, pkg string, errors []errorInfo) error {
	dest := os.Stdout
	if out != "stdout" {
		var err error
		dest, err = os.Create(out)
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := dest.Close(); closeErr != nil {
				log.Println("error closing source file", out, ":", closeErr.Error())
			}
		}()
	}

	return testTemplate.Execute(dest, struct {
		Package   string
		Timestamp time.Time
		Codes     []errorInfo
	}{
		Package:   pkg,
		Timestamp: time.Now(),
		Codes:     errors,
	})
}

func (e errorInfo) Serialized() string {
	return fmt.Sprintf("%s%s", e.Prefix, SerializeCode(e.Code))
}

// SerializeCode serializes a number into a english-friendly string.
func SerializeCode(code uint) string {
	if code <= uint(len(codeCharacters)) {
		return Left(string(codeCharacters[code]), 8, string(codeCharacters[0]))
	}
	places := make([]string, 50)
	num := int(code)
	place := 0
	for {
		remainder := num % len(codeCharacters)
		digit := string(codeCharacters[remainder])
		places[place] = digit
		place++
		num = num / len(codeCharacters)
		if num < 1 {
			break
		}
	}
	for i, j := 0, len(places)-1; i < j; i, j = i+1, j-1 {
		places[i], places[j] = places[j], places[i]
	}
	return Left(strings.Join(places, ""), 8, string(codeCharacters[0]))
}

func times(str string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(str, n)
}

// Left left-pads the string with pad up to len runes
// len may be exceeded if
func Left(str string, length int, pad string) string {
	return times(pad, length-len(str)) + str
}

var packageTemplate = template.Must(template.New("").Funcs(template.FuncMap{
	"lower": strings.ToLower,
}).Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by herr at {{ .Timestamp }}
package {{ .Package }}

import (
    "fmt"
)

type CodedError interface {
    Code() int
    Description() string
    Prefix() string
    error
}

{{ range .Codes }}
type {{ .Label }}Error struct {
    Err error
}
{{ end }}
{{ range .Codes }}var _ CodedError = {{ .Label }}Error{}
{{ end }}
// ErrorFromCode returns the CodedError for a serialized coded error string. 
func ErrorFromCode(code string) (bool, error) {
    switch code {
{{- range .Codes }}
    case "{{ .Serialized }}":
        return true, {{ .Label }}Error{}
{{- end }}
    default:
        return false, fmt.Errorf("unknown error code: %s", code)
    }
}
{{range .Codes }}
func (e {{ .Label }}Error) Error() string {
    return "{{ .Serialized }}"
}

func (e {{ .Label }}Error) Unwrap() error {
	return e.Err
}

func (e {{ .Label }}Error) Is(target error) bool {
    t, ok := target.({{ .Label }}Error)
    if !ok {
        return false
    }
    return t.Prefix() == "{{ .Prefix }}" && t.Code() == {{ .Code }}
}

func (e {{ .Label }}Error) Code() int {
    return {{ .Code }}
}

func (e {{ .Label }}Error) Description() string {
    return "{{ .Description }}"
}

func (e {{ .Label }}Error) Prefix() string {
    return "{{ .Prefix }}"
}

func (e {{ .Label }}Error) String() string {
    return "{{ .Serialized }} {{ .Description }}"
}
{{ end }}

`))

var testTemplate = template.Must(template.New("").Funcs(template.FuncMap{
	"lower": strings.ToLower,
}).Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by herr at {{ .Timestamp }}
package {{ .Package }}

import (
    "fmt"
	"testing"
	"errors"
)

{{ range .Codes }}
func Test{{ .Label }} (t *testing.T) {
    err1 := {{ .Label }}Error{}
    if err1.Prefix() != "{{ .Prefix }}" {
		t.Errorf("Assertion failed on {{ .Label }}: %s != {{ .Prefix }}", err1.Prefix())
    }
    if err1.Code() != {{ .Code }} {
		t.Errorf("Assertion failed on {{ .Label }}: %d != {{ .Code }}", err1.Code())
    }
    if err1.Description() != "{{ .Description }}" {
		t.Errorf("Assertion failed on {{ .Label }}: %s != {{ .Description }}", err1.Description())
    }

	errNotFound := fmt.Errorf("not found")
	errThingNotFound := fmt.Errorf("thing: %w", errNotFound)
	err2 := {{ .Label }}Error{ Err: errThingNotFound }
	errNestErr2 := fmt.Errorf("oh snap: %w", err2)
    if err2.Code() != {{ .Code }} {
		t.Errorf("Assertion failed on {{ .Label }}: %d != {{ .Code }}", err1.Code())
    }
    if !errors.Is(err2, errNotFound) {
		t.Errorf("Assertion failed on {{ .Label }}: errNotFound not unwrapped correctly")
    }
    if !errors.Is(err2, errThingNotFound) {
		t.Errorf("Assertion failed on {{ .Label }}: errThingNotFound not unwrapped correctly")
    }
    if !errors.Is(err2, {{ .Label }}Error{}) {
		t.Errorf("Assertion failed on {{ .Label }}: {{ .Label }}Error{} not identified correctly")
    }
    if !errors.Is(errNestErr2, {{ .Label }}Error{}) {
		t.Errorf("Assertion failed on {{ .Label }}: {{ .Label }}Error{} not identified correctly")
    }
}
{{ end }}

`))

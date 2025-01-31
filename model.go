package keg

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rwxrob/choose"
	"github.com/rwxrob/json"
	"github.com/rwxrob/term"
)

const IsoDateFmt = `2006-01-02 15:04:05Z`
const IsoDateExpStr = `\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\dZ`

// Local contains a name to full path mapping for kegs stored locally.
type Local struct {
	Name string
	Path string
}

// DexEntry represents a single line in an index (usually the latest.md
// or nodes.tsv file). All three fields are always required.
type DexEntry struct {
	U time.Time // updated
	T string    // title
	N int       // node id (also see ID)
}

// MarshalJSON produces JSON text that contains one DexEntry per line
// that has not been HTML escaped (unlike the default) and that uses
// a consistent DateTime format. Note that the (broken) encoding/json
// encoder is not used at all.
func (e *DexEntry) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	buf.WriteRune('{')
	buf.WriteString(`"U":"` + e.U.Format(IsoDateFmt) + `",`)
	buf.WriteString(`"N":` + strconv.Itoa(e.N) + `,`)
	buf.WriteString(`"T":"` + json.Escape(e.T) + `"`)
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (e DexEntry) TSV() string {
	return fmt.Sprintf("%v\t%v\t%v", e.N, e.U.Format(IsoDateFmt), e.T)
}

// ID returns the node identifier as a string instead of an integer.
// Returns an empty string if unable to parse the integer.
func (e DexEntry) ID() string { return strconv.Itoa(e.N) }

// MD returns the entry as a single Markdown list item for inclusion in
// the dex/nodex.md file:
//
//     1. Second last changed in UTC in ISO8601 (RFC3339)
//     2. Current title (always first line of README.md)
//     2. Unique node integer identifier
//
// Note that the second of last change is based on *any* file within the
// node directory changing, not just the README.md or meta files.
func (e DexEntry) MD() string {
	return fmt.Sprintf(
		"* %v [%v](/%v)",
		e.U.Format(IsoDateFmt),
		e.T, e.N,
	)
}

// String implements fmt.Stringer interface as MD.
func (e DexEntry) String() string { return e.MD() }

// Asinclude returns a KEGML include link list item without the time
// suitable for creating include blocks in node files.
func (e DexEntry) AsInclude() string {
	return fmt.Sprintf("* [%v](/%v)", e.T, e.N)
}

// Dex is a collection of DexEntry structs. This allows mapping methods
// for its serialization to different output formats.
type Dex []DexEntry

// MarshalJSON produces JSON text that contains one DexEntry per line
// that has not been HTML escaped (unlike the default).
func (d *Dex) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	buf.WriteString("[")
	for _, entry := range *d {
		byt, _ := entry.MarshalJSON()
		buf.Write(byt)
		buf.WriteString(",\n")
	}
	byt := buf.Bytes()
	byt[len(byt)-2] = ']'
	return byt, nil
}

// String fulfills the fmt.Stringer interface as JSON. Any error returns
// a "null" string.
func (e Dex) String() string { return e.TSV() }

// MD renders the entire Dex as a Markdown list suitable for the
// standard dex/latest.md file.
func (e Dex) MD() string {
	var str string
	for _, entry := range e {
		str += entry.MD() + "\n"
	}
	return str
}

// AsIncludes renders the entire Dex as a KEGML include list (markdown
// bulleted list) and cab be useful from within editing sessions to
// include from the current keg without leaving the terminal editor.
func (e Dex) AsIncludes() string {
	var str string
	for _, entry := range e {
		str += entry.AsInclude() + "\n"
	}
	return str
}

// TSV renders the entire Dex as a loadable tab-separated values file.
func (e Dex) TSV() string {
	var str string
	for _, entry := range e {
		str += entry.TSV() + "\n"
	}
	return str
}

// Highest returns the highest integer value identifier.
func (d Dex) Highest() int {
	var highest int
	for _, e := range d {
		if e.N > highest {
			highest = e.N
		}
	}
	return highest
}

// Highest returns Highest as string.
func (d Dex) HighestString() string { return strconv.Itoa(d.Highest()) }

// HighestWidth returns width of highest integer identifier.
func (d Dex) HighestWidth() int { return len(d.HighestString()) }

// Pretty returns a string with pretty color string with time stamps
// rendered in more readable way.
func (d Dex) Pretty() string {
	var str string
	nwidth := d.HighestWidth()
	for _, e := range d {
		str += fmt.Sprintf(
			"%v%v %v%-"+strconv.Itoa(nwidth)+"v %v%v%v\n",
			term.Black, e.U.Format(`2006-01-02 15:03Z`),
			term.Green, e.N,
			term.White, e.T,
			term.Reset,
		)
	}
	return str
}

// PrettyLines returns Pretty but each line separate and without line
// return.
func (d Dex) PrettyLines() []string {
	lines := make([]string, 0, len(d))
	nwidth := d.HighestWidth()
	for _, e := range d {
		lines = append(lines, fmt.Sprintf(
			"%v%v %v%-"+strconv.Itoa(nwidth)+"v %v%v%v",
			term.Black, e.U.Format(`2006-01-02 15:03Z`),
			term.Green, e.N,
			term.White, e.T,
			term.Reset,
		))
	}
	return lines
}

// ByID orders the Dex from lowest to highest node ID integer.
func (e Dex) ByID() Dex {
	sort.Slice(e, func(i, j int) bool {
		return e[i].N < e[j].N
	})
	return e
}

// WithTitleText filters all nodes with titles that do not contain the text
// substring in the title.
func (e Dex) WithTitleText(keyword string) Dex {
	dex := Dex{}
	for _, d := range e {
		if strings.Index(strings.ToLower(d.T), strings.ToLower(keyword)) >= 0 {
			dex = append(dex, d)
		}
	}
	return dex
}

// ChooseWithTitleText returns a single *DexEntry for the keyword
// passed. If there are more than one then user is prompted to choose
// from list sent to the terminal.
func (d Dex) ChooseWithTitleText(key string) *DexEntry {
	hits := d.WithTitleText(key)
	switch len(hits) {
	case 1:
		return &hits[0]
	case 0:
		return nil
	default:
		i, _, err := choose.From(hits.PrettyLines())
		if err != nil {
			return nil
		}
		if i < 0 {
			return nil
		}
		return &hits[i]
	}
}

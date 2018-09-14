package pretty

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

// TODO rename WriteTag to WriteKey
// TODO rename WriteTagValue to WriteKeyValue

// ESCAPING:
//
//   The string starts with a special character:
//   One of !#%@&*`?|>{[ or -.
//   The string starts or ends with whitespace characters.
//   The string contains : or # character sequences.
//   The string ends with a colon.
//   The value looks like a number or boolean (123, 1.23, true, false, null) but should be a string.

// Implement state pooling. See:
//     https://golang.org/src/fmt/print.go#L131

type writer interface {
	// Indent appends padding to the buffer.
	Indent()

	// IndentIncrease inreases indentation.
	IndentIncrease()

	// IndentDecrese decreases indentation.
	IndentDecrease()

	// Write appends the contents of p to the buffer.
	Write(p []byte) (n int, err error)

	// WriteByte appends the contents of b to the buffer.
	WriteByte(b byte) error

	// WriteString appends the contents of s to the buffer.
	WriteString(s string) (n int, err error)

	// WriteTag appends the key to the buffer.
	WriteTag(v interface{})

	// WriteTag appends the keypair to the buffer.
	WriteTagValue(k, v interface{})
}

//
// node writer
//

type baseWriter struct {
	bytes.Buffer
	depth int
}

func (w *baseWriter) Indent() {
	for i := 0; i < w.depth; i++ {
		w.WriteString("  ")
	}
}

func (w *baseWriter) IndentIncrease() {
	w.depth++
}

func (w *baseWriter) IndentDecrease() {
	w.depth--
}

func (w *baseWriter) WriteTag(v interface{}) {
	w.WriteByte('\n')
	w.Indent()
	writeValue(w, v)
	w.WriteByte(':')
}

func (w *baseWriter) WriteTagValue(k, v interface{}) {
	if isZero(v) {
		return
	}
	w.WriteTag(k)
	if isPrimative(v) {
		w.WriteByte(' ')
		writeValue(w, v)
	} else if isSlice(v) {
		w.WriteByte('\n')
		w.Indent()
		writeValue(w, v)
	} else {
		w.depth++
		w.WriteByte('\n')
		w.Indent()
		writeValue(w, v)
		w.depth--
	}
}

//
// sequence writer
//

type indexWriter struct {
	writer
	index int
}

func (w *indexWriter) WriteTag(v interface{}) {
	w.WriteByte('\n')
	if w.index == 0 {
		w.IndentDecrease()
		w.Indent()
		w.IndentIncrease()
		w.WriteByte('-')
		w.WriteByte(' ')
	} else {
		w.Indent()
	}
	writeValue(w, v)
	w.WriteByte(':')
	w.index++
}

func (w *indexWriter) WriteTagValue(k, v interface{}) {
	if isZero(v) {
		return
	}
	w.WriteTag(k)
	if isPrimative(v) {
		w.WriteByte(' ')
		writeValue(w, v)
	} else if isSlice(v) {
		w.WriteByte('\n')
		w.Indent()
		writeValue(w, v)
	} else {
		w.IndentIncrease()
		w.WriteByte('\n')
		w.Indent()
		writeValue(w, v)
		w.IndentDecrease()
	}
}

//
// helper functions
//

func writeBool(w writer, v bool) {
	w.WriteString(
		strconv.FormatBool(v),
	)
}

func writeFloat(w writer, v float64) {
	w.WriteString(
		strconv.FormatFloat(v, 'g', -1, 64),
	)
}

func writeInt(w writer, v int) {
	w.WriteString(
		strconv.Itoa(v),
	)
}

func writeEncode(w writer, v string) {
	if len(v) == 0 {
		w.WriteByte('"')
		w.WriteByte('"')
		return
	}
	for _, b := range v {
		if isQuoted(b) {
			fmt.Fprintf(w, "%q", v)
			return
		}
	}
	w.WriteString(v)
}

func writeValue(w writer, v interface{}) {
	if v == nil {
		w.WriteByte('~')
		return
	}
	switch v := v.(type) {
	case bool, int, float64, string:
		writeScalar(w, v)
	case []interface{}:
		writeSequence(w, v)
	case []string:
		writeSequenceStr(w, v)
	case map[interface{}]interface{}:
		writeMapping(w, v)
	case map[string]string:
		writeMappingStr(w, v)
	}
}

func writeScalar(w writer, v interface{}) {
	switch v := v.(type) {
	case bool:
		writeBool(w, v)
	case int:
		writeInt(w, v)
	case float64:
		writeFloat(w, v)
	case string:
		writeEncode(w, v)
	}
}

func writeSequence(w writer, v []interface{}) {
	if len(v) == 0 {
		w.WriteByte('[')
		w.WriteByte(']')
		return
	}
	for i, v := range v {
		if i != 0 {
			w.WriteByte('\n')
			w.Indent()
		}
		w.WriteByte('-')
		w.WriteByte(' ')
		w.IndentIncrease()
		writeValue(w, v)
		w.IndentDecrease()
	}
}

func writeSequenceStr(w writer, v []string) {
	if len(v) == 0 {
		w.WriteByte('[')
		w.WriteByte(']')
		return
	}
	for i, v := range v {
		if i != 0 {
			w.WriteByte('\n')
			w.Indent()
		}
		w.WriteByte('-')
		w.WriteByte(' ')
		writeEncode(w, v)
	}
}

func writeMapping(w writer, v map[interface{}]interface{}) {
	if len(v) == 0 {
		w.WriteByte('{')
		w.WriteByte('}')
		return
	}
	var keys []string
	for k := range v {
		s := fmt.Sprint(k)
		keys = append(keys, s)
	}
	sort.Strings(keys)
	for i, k := range keys {
		v := v[k]
		if i != 0 {
			w.WriteByte('\n')
			w.Indent()
		}
		writeEncode(w, k)
		w.WriteByte(':')
		if v == nil || isPrimative(v) || isZero(v) {
			w.WriteByte(' ')
			writeValue(w, v)
		} else {
			slice := isSlice(v)
			if !slice {
				w.IndentIncrease()
			}
			w.WriteByte('\n')
			w.Indent()
			writeValue(w, v)
			if !slice {
				w.IndentDecrease()
			}
		}
	}
}

func writeMappingStr(w writer, v map[string]string) {
	if len(v) == 0 {
		w.WriteByte('{')
		w.WriteByte('}')
		return
	}
	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		v := v[k]
		if i != 0 {
			w.WriteByte('\n')
			w.Indent()
		}
		writeEncode(w, k)
		w.WriteByte(':')
		w.WriteByte(' ')
		writeEncode(w, v)
	}
}

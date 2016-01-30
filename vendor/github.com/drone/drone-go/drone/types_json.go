package drone

import (
	"encoding/json"
	"strconv"
)

// StringSlice representes a string or an array of strings.
type StringSlice struct {
	parts []string
}

// UnmarshalJSON unmarshals bytes into a StringSlice.
func (e *StringSlice) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	p := make([]string, 0, 1)
	if err := json.Unmarshal(b, &p); err != nil {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		p = append(p, s)
	}

	e.parts = p
	return nil
}

// Len returns the number of strings in a StringSlice.
func (e *StringSlice) Len() int {
	if e == nil {
		return 0
	}
	return len(e.parts)
}

// Slice returns the slice of strings in a StringSlice.
func (e *StringSlice) Slice() []string {
	if e == nil {
		return nil
	}
	return e.parts
}

// NewStringSlice takes a slice of strings and returns a StringSlice.
func NewStringSlice(parts []string) StringSlice {
	return StringSlice{parts}
}

// StringMap representes a string or a map of strings.
type StringMap struct {
	parts map[string]string
}

// UnmarshalJSON unmarshals bytes into a StringMap.
func (e *StringMap) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	p := map[string]string{}
	if err := json.Unmarshal(b, &p); err != nil {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		p[""] = s
	}

	e.parts = p
	return nil
}

// Len returns the number of elements in the StringMap.
func (e *StringMap) Len() int {
	if e == nil {
		return 0
	}
	return len(e.parts)
}

// String satisfies the fmt.Stringer interface.
func (e *StringMap) String() (str string) {
	if e == nil {
		return
	}
	for _, val := range e.parts {
		return val // returns the first string value
	}
	return
}

// Map returns StringMap as a map of string to string.
func (e *StringMap) Map() map[string]string {
	if e == nil {
		return nil
	}
	return e.parts
}

// NewStringMap takes a map of string to string and returns a StringMap.
func NewStringMap(parts map[string]string) StringMap {
	return StringMap{parts}
}

// StringInt representes a string or an integer value.
type StringInt struct {
	value string
}

// UnmarshalJSON unmarshals bytes into a StringInt.
func (e *StringInt) UnmarshalJSON(b []byte) error {
	var num int
	err := json.Unmarshal(b, &num)
	if err == nil {
		e.value = strconv.Itoa(num)
		return nil
	}
	return json.Unmarshal(b, &e.value)
}

func (e StringInt) String() string {
	return e.value
}

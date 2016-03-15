package template

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/aymerick/raymond"
	"github.com/drone/drone-go/drone"
	"net/url"
)

func init() {
	raymond.RegisterHelpers(funcs)
}

// Render parses and executes a template, returning the results
// in string format.
func Render(template string, playload *drone.Payload) (string, error) {
	if strings.HasPrefix(template, "http") {
		resp, err := http.Get(template)

		if err != nil {
			return "", fmt.Errorf("Failed to fetch remote template: %s", err)
		}

		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", fmt.Errorf("Failed to read remote template: %s", err)
		}

		template = string(content)
	}

	return raymond.Render(template, normalize(playload))
}

// RenderTrim parses and executes a template, returning the results
// in string format. The result is trimmed to remove left and right
// padding and newlines that may be added unintentially in the
// template markup.
func RenderTrim(template string, playload *drone.Payload) (string, error) {
	out, err := Render(template, playload)
	return strings.Trim(out, " \n"), err
}

// Write parses and executes a template, writing the results to
// writer w.
func Write(w io.Writer, template string, playload *drone.Payload) error {
	out, err := Render(template, playload)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, out)
	return err
}

var funcs = map[string]interface{}{
	"uppercasefirst": uppercaseFirst,
	"uppercase":      strings.ToUpper,
	"lowercase":      strings.ToLower,
	"duration":       toDuration,
	"datetime":       toDatetime,
	"success":        isSuccess,
	"failure":        isFailure,
	"truncate":       truncate,
	"urlencode":      urlencode,
}

// truncate is a helper function that truncates a string by a particular length.
func truncate(s string, len int) string {
	if utf8.RuneCountInString(s) <= len {
		return s
	}
	runes := []rune(s)
	return string(runes[:len])

}

// uppercaseFirst is a helper function that takes a string and capitalizes
// the first letter.
func uppercaseFirst(s string) string {
	a := []rune(s)
	a[0] = unicode.ToUpper(a[0])
	s = string(a)
	return s
}

// toDuration is a helper function that calculates a duration for a start and
// and end time, and returns the duration in string format.
func toDuration(started, finished float64) string {
	return fmt.Sprintln(time.Duration(finished-started) * time.Second)
}

// toDatetime is a helper function that converts a unix timestamp to a string.
func toDatetime(timestamp float64, layout, zone string) string {
	if len(zone) == 0 {
		return time.Unix(int64(timestamp), 0).Format(layout)
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		fmt.Printf("Error parsing timezone, defaulting to local timezone. %s\n", err)
		return time.Unix(int64(timestamp), 0).Local().Format(layout)
	}
	return time.Unix(int64(timestamp), 0).In(loc).Format(layout)
}

// isSuccess is a helper function that executes a block iff the status
// is success, else it executes the else block.
func isSuccess(conditional bool, options *raymond.Options) string {
	if !conditional {
		return options.Inverse()
	}

	switch options.ParamStr(0) {
	case "success":
		return options.Fn()
	default:
		return options.Inverse()
	}
}

// isFailure is a helper function that executes a block iff the status
// is a form of failure, else it executes the else block.
func isFailure(conditional bool, options *raymond.Options) string {
	if !conditional {
		return options.Inverse()
	}

	switch options.ParamStr(0) {
	case "failure", "error", "killed":
		return options.Fn()
	default:
		return options.Inverse()
	}
}

// normalize takes a Go representation of the variable, marshals
// to json and then unmarshals to a map[string]interfacce{}. This
// is important because it let's us use the JSON variable names
// in our template
func normalize(in interface{}) map[string]interface{} {
	data, _ := json.Marshal(in) // we own the types, so this should never fail

	out := map[string]interface{}{}
	json.Unmarshal(data, &out)
	return out
}

// urlencode is a helper function that encode a block
// to url safe string.
func urlencode(options *raymond.Options) string {
	return url.QueryEscape(options.Fn())
}

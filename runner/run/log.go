package main

import (
	"bytes"
	"fmt"

	"github.com/Sirupsen/logrus"
)

type formatter struct {
	nocolor bool
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("\033[2m")
	buf.WriteString("[drone]")

	for k, v := range entry.Data {
		if k != "exit_code" {
			continue
		}

		if v == 0 {
			buf.WriteString("\033[32m \u2713\033[0m")
		} else {
			buf.WriteString("\033[31m \u2717\033[0m")
		}
	}

	buf.WriteByte(' ')
	buf.WriteString(entry.Message)
	buf.WriteByte(' ')

	for k, v := range entry.Data {
		buf.WriteString(
			fmt.Sprintf("%s=%v", k, v),
		)
	}

	buf.WriteString("\033[0m")
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

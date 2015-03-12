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

	buf.WriteString("[drone]")
	buf.WriteByte(' ')
	buf.WriteString(entry.Message)
	buf.WriteByte(' ')

	for k, v := range entry.Data {
		buf.WriteString(
			fmt.Sprintf("%s=%v", k, v),
		)
	}

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

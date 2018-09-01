// Copyright 2018 Drone.IO Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

// A Logger represents an active logging object that generates
// lines of output to an io.Writer.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
}

// Discard returns a no-op logger
func Discard() Logger {
	return &discard{}
}

type discard struct{}

func (*discard) Debug(args ...interface{})                 {}
func (*discard) Debugf(format string, args ...interface{}) {}
func (*discard) Debugln(args ...interface{})               {}
func (*discard) Error(args ...interface{})                 {}
func (*discard) Errorf(format string, args ...interface{}) {}
func (*discard) Errorln(args ...interface{})               {}
func (*discard) Info(args ...interface{})                  {}
func (*discard) Infof(format string, args ...interface{})  {}
func (*discard) Infoln(args ...interface{})                {}
func (*discard) Warn(args ...interface{})                  {}
func (*discard) Warnf(format string, args ...interface{})  {}
func (*discard) Warnln(args ...interface{})                {}

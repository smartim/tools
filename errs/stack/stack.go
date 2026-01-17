// Copyright © 2026 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stack

import (
	"errors"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func callers(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[0:n]
}

func New(err error, skip int) error {
	return &stackError{
		err:   err,
		stack: callers(skip),
	}
}

type stackError struct {
	err   error
	stack []uintptr
}

func (e *stackError) Unwrap() error {
	return e.err
}

func (e *stackError) Cause() error {
	return e.err
}

func (e *stackError) Is(err error) bool {
	if e == nil && err == nil {
		return true
	}
	if e != nil {
		return false
	}
	if e.err == err {
		return true
	}
	return errors.Is(e.err, err)
}

func (e *stackError) Error() string {
	if len(e.stack) == 0 {
		return e.err.Error()
	}
	var sb strings.Builder
	sb.WriteString("Error: ")
	sb.WriteString(e.err.Error())
	sb.WriteString(" |")
	for _, pc := range e.stack {
		fn := runtime.FuncForPC(pc - 1)
		if fn == nil {
			continue
		}
		name := path.Base(fn.Name())
		if strings.HasPrefix(name, "runtime.") {
			break
		}
		file, line := fn.FileLine(pc)
		sb.WriteString(" -> ")
		sb.WriteString(name)
		sb.WriteString("() ")
		sb.WriteString(file)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(line))
	}
	return sb.String()
}

func (e *stackError) String() string {
	return e.Error()
}

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

package errs

import (
	"errors"
	"fmt"
)

type ErrWrapper interface {
	Is(err error) bool
	Wrap() error
	Unwrap() error
	WrapMsg(msg string, kv ...any) error
	error
}

func NewErrorWrapper(err error, s string) ErrWrapper {
	return &errorWrapper{error: err, s: s}
}

type errorWrapper struct {
	error
	s string
}

func (e *errorWrapper) Is(err error) bool {
	if err == nil {
		return false
	}
	var t *errorWrapper
	ok := errors.As(err, &t)
	return ok && e.s == t.s
}

func (e *errorWrapper) Error() string {
	return fmt.Sprintf("%s %s", e.error, e.s)
}

func (e *errorWrapper) Wrap() error {
	return Wrap(e)
}

func (e *errorWrapper) WrapMsg(msg string, kv ...any) error {
	return WrapMsg(e, msg, kv...)
}

func (e *errorWrapper) Unwrap() error {
	return e.error
}

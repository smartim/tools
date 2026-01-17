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

package specialerror

import (
	"errors"

	"github.com/smartim/tools/errs"
)

var handlers []func(err error) errs.CodeError

func AddErrHandler(h func(err error) errs.CodeError) (err error) {
	if h == nil {
		return errs.New("nil handler")
	}
	handlers = append(handlers, h)
	return nil
}

func AddReplace(target error, codeErr errs.CodeError) error {
	handler := func(err error) errs.CodeError {
		if errors.Is(err, target) {
			return codeErr
		}
		return nil
	}

	if err := AddErrHandler(handler); err != nil {
		return err
	}

	return nil
}

func ErrCode(err error) errs.CodeError {
	var codeErr errs.CodeError
	if errors.As(err, &codeErr) {
		return codeErr
	}
	for i := 0; i < len(handlers); i++ {
		if codeErr := handlers[i](err); codeErr != nil {
			return codeErr
		}
	}
	return nil
}

func ErrString(err error) errs.Error {
	var codeErr errs.Error
	if errors.As(err, &codeErr) {
		return codeErr
	}
	return nil
}

func ErrWrapper(err error) errs.ErrWrapper {
	var codeErr errs.ErrWrapper
	if errors.As(err, &codeErr) {
		return codeErr
	}
	return nil
}

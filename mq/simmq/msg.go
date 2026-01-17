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

package simmq

import "context"

type message struct {
	ctx   context.Context
	key   string
	value []byte
}

func (m *message) Context() context.Context {
	return m.ctx
}

func (m *message) Key() string {
	return m.key
}

func (m *message) Value() []byte {
	return m.value
}

func (m *message) Mark() {
}

func (m *message) Commit() {
}

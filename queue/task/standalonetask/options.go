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

package standalonetask

type Options[T any, K comparable] func(*QueueManager[T, K])

func WithStrategy[T any, K comparable](s strategy) Options[T, K] {
	return func(tm *QueueManager[T, K]) {
		tm.assignStrategy = getStrategy[T, K](s)
	}
}

func WithAfterProcessPushFunc[T any, K comparable](fs ...func(key K, data T)) Options[T, K] {
	return func(tm *QueueManager[T, K]) {
		tm.afterProcessPushFunc = fs
	}
}

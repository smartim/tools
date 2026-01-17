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

package redistimer

import "time"

// Option configures a RedisTimer
type Option[T any] func(*RedisTimer[T])

// WithNamespace sets the Redis key namespace
func WithNamespace[T any](namespace string) Option[T] {
	return func(rt *RedisTimer[T]) {
		rt.namespace = namespace
	}
}

// WithPollInterval sets the polling interval for checking expired timers
func WithPollInterval[T any](interval time.Duration) Option[T] {
	return func(rt *RedisTimer[T]) {
		rt.pollInterval = interval
	}
}

// WithBatchSize sets the batch size for processing expired timers
func WithBatchSize[T any](size int) Option[T] {
	return func(rt *RedisTimer[T]) {
		rt.batchSize = size
	}
}

// WithMarshal sets custom marshal function
func WithMarshal[T any](fn func(T) ([]byte, error)) Option[T] {
	return func(rt *RedisTimer[T]) {
		rt.marshal = fn
	}
}

// WithUnmarshal sets custom unmarshal function
func WithUnmarshal[T any](fn func([]byte) (T, error)) Option[T] {
	return func(rt *RedisTimer[T]) {
		rt.unmarshal = fn
	}
}

// WithSelfContainedKey indicates that the key itself contains all the data.
// This option is only valid for string types (T must be string).
// When enabled:
//   - No additional hash storage is used
//   - The ZSET member (key) is the data itself
//   - More efficient for simple string timers
//
// Example usage:
//
//	manager := New[string](ctx, config, func(s string) string { return s }, handlers, WithSelfContainedKey[string]())
func WithSelfContainedKey[T any]() Option[T] {
	return func(rt *RedisTimer[T]) {
		// Type check: only allow for string type
		var zero T
		if _, ok := any(zero).(string); !ok {
			panic("WithSelfContainedKey can only be used with string type")
		}
		rt.selfContainedKey = true
	}
}

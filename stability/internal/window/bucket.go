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

package window

import (
	"sync/atomic"
)

// Bucket is a basic structure that stores multiple float64 data points
type Bucket struct {
	Points []float64 // stores individual data points
	Count  int64     // number of data points in this bucket
	next   *Bucket   // points to the next bucket, forming a circular structure
}

// NewBucket creates a new bucket
func NewBucket() *Bucket {
	return &Bucket{
		Points: make([]float64, 0),
	}
}

// Append adds the given value to the bucket
func (b *Bucket) Append(val float64) {
	b.Points = append(b.Points, val)
	atomic.AddInt64(&b.Count, 1)
}

// Add adds the value at the specified offset
func (b *Bucket) Add(offset int, val float64) {
	b.Points[offset] += val
	atomic.AddInt64(&b.Count, 1)
}

// Reset clears the bucket
func (b *Bucket) Reset() {
	b.Points = b.Points[:0]
	atomic.StoreInt64(&b.Count, 0)
}

// Next return next bucket
func (b *Bucket) Next() *Bucket {
	return b.next
}

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

import "fmt"

// Iterator returns an iterator over count buckets starting from the specified offset
func (w *Window) Iterator(offset int, count int) Iterator {
	return Iterator{
		count: count,
		cur:   &w.buckets[offset%w.size],
	}
}

// Iterator is an iterator for traversing buckets in the window
type Iterator struct {
	count         int     // total number of buckets to iterate
	iteratedCount int     // number of buckets already iterated
	cur           *Bucket // pointer to the current bucket
}

// Next returns true if there are still buckets left to iterate
func (i *Iterator) Next() bool {
	return i.count != i.iteratedCount
}

// Bucket returns the current bucket
func (i *Iterator) Bucket() Bucket {
	if !(i.Next()) {
		panic(fmt.Errorf("window: iteration out of range iteratedCount: %d count: %d", i.iteratedCount, i.count))
	}
	bucket := *i.cur
	i.iteratedCount++
	i.cur = i.cur.Next()
	return bucket
}

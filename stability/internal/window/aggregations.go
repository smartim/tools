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

// Sum computes the total sum of values within the window
func Sum(iterator Iterator) float64 {
	var sum float64
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			sum += p
		}
	}
	return sum
}

// Avg computes the average of values within the window
func Avg(iterator Iterator) float64 {
	var sum float64
	var count int
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			sum += p
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// Max computes the maximum value within the window
func Max(iterator Iterator) float64 {
	var max float64
	var initialized bool
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			if !initialized {
				max = p
				initialized = true
				continue
			}
			if p > max {
				max = p
			}
		}
	}
	return max
}

// Min computes the minimum value within the window
func Min(iterator Iterator) float64 {
	var min float64
	var initialized bool
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			if !initialized {
				min = p
				initialized = true
				continue
			}
			if p < min {
				min = p
			}
		}
	}
	return min
}

// Count computes the total number of data points within the window
func Count(iterator Iterator) float64 {
	var count int64
	for iterator.Next() {
		bucket := iterator.Bucket()
		count += bucket.Count
	}
	return float64(count)
}

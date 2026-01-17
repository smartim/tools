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

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestName(t *testing.T) {
	const topic = "test"
	ctx := context.Background()

	done := make(chan struct{})

	go func() {
		defer func() {
			t.Log("consumer end")
			close(done)
		}()
		fn := func(ctx context.Context, key string, value []byte) error {
			t.Logf("consumer key: %s, value: %s", key, value)
			return nil
		}
		c := GetTopicConsumer(topic)
		for {
			if err := c.Subscribe(ctx, fn); err != nil {
				t.Log("subscribe err", err)
				return
			}
		}
	}()

	var wg sync.WaitGroup

	var count atomic.Int64
	for i := 0; i < 4; i++ {
		wg.Add(1)
		p := GetTopicProducer(topic)
		key := fmt.Sprintf("go_%d", i+1)

		go func() {
			defer func() {
				t.Log("producer end", key)
				wg.Done()
			}()
			for i := 1; i <= 10; i++ {
				value := fmt.Sprintf("value_%d", count.Add(1))
				if err := p.SendMessage(ctx, key, []byte(value)); err != nil {
					t.Log("send message err", key, value, err)
					return
				}
			}
		}()
	}

	wg.Wait()
	_ = GetTopicProducer(topic).Close()
	<-done
}

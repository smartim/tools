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
	"errors"
	"sync"
	"sync/atomic"

	"github.com/smartim/tools/mq"
)

var (
	errClosed = errors.New("memory mq closed")
)

func NewMemory(size int) (mq.Producer, mq.Consumer) {
	m := newMemory(size, nil)
	return m, m
}

func newMemory(size int, fn func()) *memory {
	return &memory{
		ch:   make(chan *message, size),
		done: make(chan struct{}),
		fn:   fn,
	}
}

type memory struct {
	lock   sync.RWMutex
	closed atomic.Bool
	ch     chan *message
	done   chan struct{}
	fn     func()
}

func (x *memory) Subscribe(ctx context.Context, fn mq.Handler) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg, ok := <-x.ch:
		if !ok {
			return errClosed
		}
		if err := fn(msg); err != nil {
			return err
		}
		return nil
	}
}

func (x *memory) SendMessage(ctx context.Context, key string, value []byte) error {
	if x.closed.Load() {
		return errClosed
	}
	msg := &message{
		ctx:   context.WithoutCancel(ctx),
		key:   key,
		value: value,
	}
	x.lock.RLock()
	defer x.lock.RUnlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-x.done:
		return errClosed
	case x.ch <- msg:
		return nil
	}
}

func (x *memory) Close() error {
	if !x.closed.CompareAndSwap(false, true) {
		return nil
	}
	close(x.done)
	if x.fn != nil {
		x.fn()
	}
	x.lock.Lock()
	defer x.lock.Unlock()
	close(x.ch)
	return nil
}

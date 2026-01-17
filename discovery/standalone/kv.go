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

package standalone

import (
	"context"
	"sync"

	"github.com/smartim/tools/discovery"
)

type keyValue struct {
	lock sync.RWMutex
	kv   map[string][]byte
}

func (x *keyValue) SetKey(ctx context.Context, key string, data []byte) error {
	tmp := make([]byte, len(data))
	copy(tmp, data)
	x.lock.Lock()
	if x.kv == nil {
		x.kv = make(map[string][]byte)
	}
	x.kv[key] = tmp
	x.lock.Unlock()
	return nil
}

func (x *keyValue) SetWithLease(ctx context.Context, key string, val []byte, ttl int64) error {
	return discovery.ErrNotSupported
}

func (x *keyValue) GetKey(ctx context.Context, key string) ([]byte, error) {
	x.lock.RLock()
	defer x.lock.RUnlock()
	if x.kv != nil {
		if v, ok := x.kv[key]; ok {
			tmp := make([]byte, len(v))
			copy(tmp, v)
			return tmp, nil
		}
	}
	return nil, nil
}

func (x *keyValue) GetKeyWithPrefix(ctx context.Context, key string) ([][]byte, error) {
	return nil, discovery.ErrNotSupported
}

func (x *keyValue) DelData(ctx context.Context, key string) error {
	x.lock.Lock()
	if x.kv != nil {
		delete(x.kv, key)
	}
	x.lock.Unlock()
	return nil
}

func (x *keyValue) WatchKey(ctx context.Context, key string, fn discovery.WatchKeyHandler) error {
	return discovery.ErrNotSupported
}

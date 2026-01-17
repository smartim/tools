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
	"sync"

	"github.com/smartim/tools/mq"
)

const defaultMqSize = 1024 * 16

var (
	topicMq   map[string]*memory
	topicLock sync.Mutex
)

func getTopicMemory(topic string) *memory {
	topicLock.Lock()
	defer topicLock.Unlock()
	if topicMq == nil {
		topicMq = make(map[string]*memory)
	}
	val, ok := topicMq[topic]
	if !ok {
		val = newMemory(defaultMqSize, func() {
			topicLock.Lock()
			delete(topicMq, topic)
			topicLock.Unlock()
		})
		topicMq[topic] = val
	}
	return val
}

func GetTopicProducer(topic string) mq.Producer {
	return getTopicMemory(topic)
}

func GetTopicConsumer(topic string) mq.Consumer {
	return getTopicMemory(topic)
}

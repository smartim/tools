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

package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/smartim/tools/mq"
)

func NewKafkaProducerV2(config *Config, addr []string, topic string) (mq.Producer, error) {
	conf, err := BuildProducerConfig(*config)
	if err != nil {
		return nil, err
	}
	producer, err := NewProducer(conf, addr)
	if err != nil {
		return nil, err
	}
	return &mqProducer{
		topic:    topic,
		producer: producer,
	}, nil
}

type mqProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func (x *mqProducer) SendMessage(ctx context.Context, key string, value []byte) error {
	headers, err := GetMQHeaderWithContext(ctx)
	if err != nil {
		return err
	}
	kMsg := &sarama.ProducerMessage{
		Topic:   x.topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(value),
		Headers: headers,
	}
	if _, _, err := x.producer.SendMessage(kMsg); err != nil {
		return err
	}
	return nil
}

func (x *mqProducer) Close() error {
	return x.producer.Close()
}

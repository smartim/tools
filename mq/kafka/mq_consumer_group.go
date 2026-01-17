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
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/smartim/tools/log"
	"github.com/smartim/tools/mcontext"
	"github.com/smartim/tools/mq"
)

func NewMConsumerGroupV2(ctx context.Context, conf *Config, groupID string, topics []string, autoCommitEnable bool) (mq.Consumer, error) {
	config, err := BuildConsumerGroupConfig(conf, sarama.OffsetNewest, autoCommitEnable)
	if err != nil {
		return nil, err
	}
	group, err := NewConsumerGroup(config, conf.Addr, groupID)
	if err != nil {
		return nil, err
	}
	mcg := &mqConsumerGroup{
		topics:   topics,
		groupID:  groupID,
		consumer: group,
		msg:      make(chan *consumerMessage, 64),
	}
	mcg.ctx, mcg.cancel = context.WithCancel(ctx)
	go mcg.loopConsume()
	return mcg, nil
}

type consumerMessage struct {
	Msg     *sarama.ConsumerMessage
	Session sarama.ConsumerGroupSession
}

type mqConsumerGroup struct {
	topics   []string
	groupID  string
	consumer sarama.ConsumerGroup
	ctx      context.Context
	cancel   context.CancelFunc
	msg      chan *consumerMessage
	once     sync.Once
}

func (*mqConsumerGroup) Setup(sarama.ConsumerGroupSession) error { return nil }

func (*mqConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (x *mqConsumerGroup) closeMsgChan() {
	x.once.Do(func() {
		x.cancel()
		close(x.msg)
	})
}

func (x *mqConsumerGroup) loopConsume() {
	defer x.closeMsgChan()
	ctx := mcontext.SetOperationID(x.ctx, fmt.Sprintf("consumer_group_%s_%s_%d", strings.Join(x.topics, "_"), x.groupID, rand.Uint32()))
	for {
		if err := x.consumer.Consume(x.ctx, x.topics, x); err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				return
			case errors.Is(err, sarama.ErrClosedConsumerGroup):
				return
			}
			log.ZWarn(ctx, "consume err", err, "topic", x.topics, "groupID", x.groupID)
		}
	}
}

func (x *mqConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer func() {
		_ = recover()
	}()

	msg := claim.Messages()
	for {
		select {
		case <-x.ctx.Done():
			return context.Canceled
		case val, ok := <-msg:
			if !ok {
				return nil
			}

			select {
			case <-x.ctx.Done():
				return context.Canceled
			case x.msg <- &consumerMessage{Msg: val, Session: session}:
			}
		}
	}
}

func (x *mqConsumerGroup) Subscribe(ctx context.Context, fn mq.Handler) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case msg, ok := <-x.msg:
		if !ok {
			return sarama.ErrClosedConsumerGroup
		}
		ctx := GetContextWithMQHeader(msg.Msg.Headers)
		if err := fn(kafkaMessage{ctx: ctx, msg: msg}); err != nil {
			return err
		}
		return nil
	}
}

func (x *mqConsumerGroup) Close() error {
	x.cancel()
	return x.consumer.Close()
}

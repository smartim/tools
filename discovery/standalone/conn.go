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
	"time"

	"github.com/smartim/tools/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newClientConn() *clientConn {
	return &clientConn{
		registry:   newRegistry(),
		serializer: newProtoSerializer(),
	}
}

type clientConn struct {
	registry   *registry
	serializer serializer
}

func (x *clientConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	handler, err := x.registry.getMethod(ctx, method)
	if err != nil {
		return err
	}
	log.ZInfo(ctx, "standalone rpc server request", "method", method, "req", args)
	start := time.Now()
	resp, err := handler(ctx, args, nil)
	if err == nil {
		log.ZInfo(ctx, "standalone rpc server response success", "method", method, "cost", time.Since(start), "req", args, "resp", resp)
	} else {
		log.ZError(ctx, "standalone rpc server response error", err, "method", method, "cost", time.Since(start), "req", args)
	}
	if err != nil {
		return err
	}
	tmp, err := x.serializer.Marshal(resp)
	if err != nil {
		return err
	}
	return x.serializer.Unmarshal(tmp, reply)
}

func (x *clientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Errorf(codes.Unimplemented, "method stream not implemented")
}

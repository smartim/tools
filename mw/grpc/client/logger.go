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

package client

import (
	"context"
	"time"

	"github.com/smartim/tools/log"
	"google.golang.org/grpc"
)

func GrpcClientLogger() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.ZInfo(ctx, "rpc client request", "method", method, "req", req)
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			log.ZInfo(ctx, "rpc client response success", "method", method, "cost", time.Since(start), "req", req, "resp", reply)
		} else {
			log.ZError(ctx, "rpc client response error", err, "method", method, "cost", time.Since(start), "req", req)
		}
		return err
	})
}

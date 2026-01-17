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

package server

import (
	"context"
	"time"

	"github.com/smartim/tools/log"
	"google.golang.org/grpc"
)

func GrpcServerLogger() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.ZInfo(ctx, "rpc server request", "method", info.FullMethod, "req", req)
		start := time.Now()
		resp, err = handler(ctx, req)
		if err == nil {
			log.ZInfo(ctx, "rpc server response success", "method", info.FullMethod, "cost", time.Since(start), "req", req, "resp", resp)
		} else {
			log.ZError(ctx, "rpc server response error", err, "method", info.FullMethod, "cost", time.Since(start), "req", req)
		}
		return
	})
}

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
	"errors"

	"github.com/smartim/protocol/errinfo"
	"github.com/smartim/tools/errs"
	"github.com/smartim/tools/log"
	"github.com/smartim/tools/mw/specialerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcServerErrorConvert() grpc.ServerOption {
	type grpcError interface {
		error
		GRPCStatus() *status.Status
	}
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return
		}
		var grpcErr grpcError
		if errors.As(err, &grpcErr) {
			return
		}
		err = codeErrorToGrpcError(ctx, getCodeError(err))
		return
	})
}

func getCodeError(err error) errs.CodeError {
	if codeErr := specialerror.ErrCode(err); codeErr != nil {
		return codeErr
	}
	return errs.ErrInternalServer.WithDetail(errs.Unwrap(err).Error())
}

func codeErrorToGrpcError(ctx context.Context, codeErr errs.CodeError) error {
	grpcStatus := status.New(codes.Code(codeErr.Code()), codeErr.Msg())
	if detail := codeErr.Detail(); detail != "" {
		errInfo := &errinfo.ErrorInfo{Cause: detail}
		details, err := grpcStatus.WithDetails(errInfo)
		if err == nil {
			return details.Err()
		} else {
			log.ZError(ctx, "rpc server response WithDetails failed", err, "codeErr", codeErr)
		}
	}
	return grpcStatus.Err()
}

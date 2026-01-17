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

	"github.com/smartim/tools/discovery"
	"google.golang.org/grpc"
)

var global *svcDiscoveryRegistry

func init() {
	conn := newDiscoveryConn()
	global = &svcDiscoveryRegistry{
		Conn:             conn,
		ServiceRegistrar: conn.conn.registry,
	}
}

func GetDiscoveryConn() discovery.Conn {
	return global
}

func GetServiceRegistrar() grpc.ServiceRegistrar {
	return global
}

func GetKeyValue() discovery.KeyValue {
	return global
}

func GetSvcDiscoveryRegistry() discovery.SvcDiscoveryRegistry {
	return global
}

type svcDiscoveryRegistry struct {
	discovery.Conn
	grpc.ServiceRegistrar
	keyValue
}

func (x *svcDiscoveryRegistry) AddOption(opts ...grpc.DialOption) {}

func (x *svcDiscoveryRegistry) Register(ctx context.Context, serviceName, host string, port int, opts ...grpc.DialOption) error {
	return nil
}

func (x *svcDiscoveryRegistry) Close() {}

func (x *svcDiscoveryRegistry) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	return "", nil
}

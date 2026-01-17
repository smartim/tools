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

import "google.golang.org/protobuf/proto"

type serializer interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}

func newProtoSerializer() serializer {
	return protoSerializer{}
}

type protoSerializer struct{}

func (protoSerializer) Marshal(in any) ([]byte, error) {
	return proto.Marshal(in.(proto.Message))
}

func (protoSerializer) Unmarshal(b []byte, out any) error {
	return proto.Unmarshal(b, out.(proto.Message))
}

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

package minio

import (
	"context"
	"path"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	conf := Config{
		Bucket:          "openim",
		AccessKeyID:     "root",
		SecretAccessKey: "openIM123",
		Endpoint:        "http://127.0.0.1:10005",
	}
	ctx := context.Background()
	m, err := NewMinio(ctx, nil, conf)
	if err != nil {
		panic(err)
	}
	t.Log(m.DeleteObject(ctx, "/openim/data/hash/6aeb6959cad0d0b2ef4a5d9f66ed394a"))
}

func TestName2(t *testing.T) {
	t.Log(strings.Trim(path.Base("openim/thumbnail/ae20fe3d6466fdb11bcf465386b51312/image_w640_h640.jpeg"), "."))

}

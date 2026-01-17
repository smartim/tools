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

package runtimeenv

import (
	"os"
	"strings"
)

const (
	Kubernetes = "kubernetes"
	Docker     = "docker"
	Source     = "source"
)

var runtimeEnv = runtimeEnvironment()

func isDocker() bool {
	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "docker")
}

func isKubernetes() bool {
	_, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount")
	return err == nil
}

func runtimeEnvironment() string {
	if isKubernetes() {
		return Kubernetes
	} else if isDocker() {
		return Docker
	} else {
		return Source
	}
}

func RuntimeEnvironment() string {
	return runtimeEnv
}

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

package cpu

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// Ensure PsutilCPU implements CPU interface
var _ CPU = (*PsutilCPU)(nil)

// PsutilCPU implements CPU monitoring using gopsutil library
type PsutilCPU struct {
	interval time.Duration
	cores    int
}

// NewPsutilCPU creates a new CPU monitor based on gopsutil
func NewPsutilCPU(interval time.Duration) (CPU, error) {
	cores, err := cpu.Counts(true)
	if err != nil {
		cores = runtime.NumCPU() // Fallback to runtime
	}

	cpu := &PsutilCPU{
		interval: interval,
		cores:    cores,
	}

	// Test usage retrieval
	_, err = cpu.Usage()
	if err != nil {
		return nil, err
	}

	return cpu, nil
}

// Usage returns current CPU usage percentage (0-1000)
func (pc *PsutilCPU) Usage() (uint64, error) {
	percents, err := cpu.Percent(pc.interval, false)
	if err != nil {
		return 0, err
	}

	if len(percents) == 0 {
		return 0, nil
	}

	// Convert percentage to value between 0-1000
	// Ensure usage doesn't exceed 1000
	usage := min(uint64(percents[0]*10), 1000)

	return usage, nil
}

// Info returns CPU information
func (pc *PsutilCPU) Info() Info {
	info := Info{
		Cores: pc.cores,
		Quota: float64(pc.cores),
	}

	// Try to get CPU frequency
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		// Convert MHz to Hz
		info.Frequency = uint64(cpuInfo[0].Mhz * 1000000)
	}

	return info
}

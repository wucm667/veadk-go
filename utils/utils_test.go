// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd. and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"testing"
	"time"
)

func TestConvertTimeMillToTime(t *testing.T) {
	tests := []struct {
		name     string
		timeMill int64
		expected time.Time
	}{
		{
			name:     "zero timestamp",
			timeMill: 0,
			expected: time.Unix(0, 0),
		},
		{
			name:     "seconds without milliseconds",
			timeMill: 1609459200000, // 2021-01-01 00:00:00 UTC
			expected: time.Unix(1609459200, 0),
		},
		{
			name:     "with milliseconds",
			timeMill: 1609459200123, // 2021-01-01 00:00:00.123 UTC
			expected: time.Unix(1609459200, 123000000),
		},
		{
			name:     "negative timestamp",
			timeMill: -1609459200000, // 1918-12-31 00:00:00 UTC
			expected: time.Unix(-1609459200, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTimeMillToTime(tt.timeMill)
			if !result.Equal(tt.expected) {
				t.Errorf("ConvertTimeMillToTime(%d) = %v, want %v", tt.timeMill, result, tt.expected)
			}
		})
	}
}

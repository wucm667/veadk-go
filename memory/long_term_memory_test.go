// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd. and/or its affiliates.
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

package memory

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/volcengine/veadk-go/memory/long_term_memory_backends"
	"google.golang.org/adk/memory"
)

type mockMemoryServiceImpl struct {
	memory.Service
}

func TestNewLongTermMemory(t *testing.T) {
	tests := []struct {
		name    string
		backend LongTermBackendType
		wantErr bool
	}{
		{
			name:    "has user config",
			backend: "viking",
			wantErr: false,
		},
		{
			name:    "default config",
			backend: "",
			wantErr: false,
		},
		{
			name:    "unsupported backend",
			backend: "test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		mockey.PatchConvey(tt.name, t, func() {
			t.Run(tt.name, func(t *testing.T) {
				mockey.Mock(long_term_memory_backends.NewVikingDbMemoryBackend).Return(&mockMemoryServiceImpl{}, nil).Build()
				memoryService, err := NewLongTermMemoryService(tt.backend, nil)
				assert.True(t, tt.wantErr == (err != nil))
				if err == nil {
					assert.NotNil(t, memoryService)
				}
			})
		})
	}
}

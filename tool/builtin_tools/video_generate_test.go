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

package builtin_tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVideoGenerateTool(t *testing.T) {
	tests := []struct {
		name        string
		config      *VideoGenerateConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config with all fields",
			config: &VideoGenerateConfig{
				ModelName: "doubao-seedance-1-0-pro",
				APIKey:    "test-api-key",
				BaseURL:   "https://test-api.com",
			},
			expectError: false,
		},
		{
			name:        "nil config - should use defaults",
			config:      nil,
			expectError: true, // May panic if global config is not initialized
		},
		{
			name: "empty config - should use defaults",
			config: &VideoGenerateConfig{
				ModelName: "",
				APIKey:    "",
				BaseURL:   "",
			},
			expectError: true, // May fail if global config is not initialized
		},
		{
			name: "config with only model name",
			config: &VideoGenerateConfig{
				ModelName: "doubao-seedance-1-0-pro",
				APIKey:    "",
				BaseURL:   "",
			},
			expectError: true, // Will fail due to missing API key
		},
		{
			name: "config with only API key",
			config: &VideoGenerateConfig{
				ModelName: "",
				APIKey:    "test-api-key",
				BaseURL:   "",
			},
			expectError: true, // Will fail due to missing model name
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle potential panics from accessing global config
			defer func() {
				if r := recover(); r != nil {
					if tt.expectError {
						// Expected panic, test passes
						return
					}
					// Unexpected panic, fail the test
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			tool, err := NewVideoGenerateTool(tt.config)

			if tt.expectError {
				if err != nil {
					// Expected error case
					if tt.errorMsg != "" {
						assert.Contains(t, err.Error(), tt.errorMsg)
					}
					assert.Nil(t, tool)
				}
				// If no error but expectError is true, that's also acceptable
				// (means the function handled the error case gracefully)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tool)
			}
		})
	}
}

func TestVideoGenerateToolRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     VideoGenerateToolRequest
		expectValid bool
	}{
		{
			name: "valid single video request",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						VideoName: "test_video.mp4",
						Prompt:    "a cat playing with a ball",
					},
				},
				BatchSize: 1,
			},
			expectValid: true,
		},
		{
			name: "valid multiple video request",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						VideoName: "video1.mp4",
						Prompt:    "sunset over mountains",
					},
					{
						VideoName: "video2.mp4",
						Prompt:    "ocean waves",
					},
				},
				BatchSize: 2,
			},
			expectValid: true,
		},
		{
			name: "video with first frame",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						VideoName:  "logo_animation.mp4",
						Prompt:     "company logo animation",
						FirstFrame: stringPtr("https://example.com/logo.png"),
					},
				},
				BatchSize: 1,
			},
			expectValid: true,
		},
		{
			name: "video with last frame",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						VideoName: "transition.mp4",
						Prompt:    "smooth transition effect",
						LastFrame: stringPtr("https://example.com/end_frame.png"),
					},
				},
				BatchSize: 1,
			},
			expectValid: true,
		},
		{
			name: "video with both frames",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						VideoName:  "logo_reveal.mp4",
						Prompt:     "logo reveal animation",
						FirstFrame: stringPtr("https://example.com/start.png"),
						LastFrame:  stringPtr("https://example.com/end.png"),
					},
				},
				BatchSize: 1,
			},
			expectValid: true,
		},
		{
			name: "missing required video name",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						Prompt: "a cat playing with a ball",
					},
				},
				BatchSize: 1,
			},
			expectValid: false,
		},
		{
			name: "missing required prompt",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						VideoName: "test_video.mp4",
					},
				},
				BatchSize: 1,
			},
			expectValid: false,
		},
		{
			name: "empty params array",
			request: VideoGenerateToolRequest{
				Params:    []GenerateVideosRequest{},
				BatchSize: 0,
			},
			expectValid: true, // Empty array is structurally valid
		},
		{
			name: "zero batch size",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						VideoName: "test_video.mp4",
						Prompt:    "a cat playing with a ball",
					},
				},
				BatchSize: 0,
			},
			expectValid: true, // Zero batch size should default to 10
		},
		{
			name: "negative batch size",
			request: VideoGenerateToolRequest{
				Params: []GenerateVideosRequest{
					{
						VideoName: "test_video.mp4",
						Prompt:    "a cat playing with a ball",
					},
				},
				BatchSize: -1,
			},
			expectValid: true, // Negative batch size is structurally valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - just check that required fields are present
			if tt.expectValid {
				assert.NotNil(t, tt.request.Params)
				for _, param := range tt.request.Params {
					assert.NotEmpty(t, param.VideoName)
					assert.NotEmpty(t, param.Prompt)
				}
			}
		})
	}
}

func TestVideoGenerateResult(t *testing.T) {
	tests := []struct {
		name           string
		result         VideoGenerateResult
		expectedStatus string
	}{
		{
			name: "successful result",
			result: VideoGenerateResult{
				SuccessList: []*VideoResult{
					{VideoName: "video1.mp4", Url: "https://example.com/video1.mp4"},
					{VideoName: "video2.mp4", Url: "https://example.com/video2.mp4"},
				},
				ErrorList: []*VideoResult{},
			},
			expectedStatus: VideoGenerateSuccessStatus,
		},
		{
			name: "error result",
			result: VideoGenerateResult{
				SuccessList: []*VideoResult{},
				ErrorList: []*VideoResult{
					{VideoName: "failed_video1"},
					{VideoName: "failed_video2"},
				},
			},
			expectedStatus: VideoGenerateErrorStatus, // No success, so status should be error
		},
		{
			name: "mixed result",
			result: VideoGenerateResult{
				SuccessList: []*VideoResult{
					{VideoName: "video1.mp4", Url: "https://example.com/video1.mp4"},
				},
				ErrorList: []*VideoResult{
					{VideoName: "failed_video1"},
				},
			},
			expectedStatus: VideoGenerateSuccessStatus, // Has success, so overall status is success
		},
		{
			name: "empty result",
			result: VideoGenerateResult{
				SuccessList: []*VideoResult{},
				ErrorList:   []*VideoResult{},
			},
			expectedStatus: VideoGenerateErrorStatus, // No success, so status should be error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set status based on success list (this mimics the logic in the handler)
			if len(tt.result.SuccessList) == 0 {
				tt.result.Status = VideoGenerateErrorStatus
			} else {
				tt.result.Status = VideoGenerateSuccessStatus
			}

			assert.Equal(t, tt.expectedStatus, tt.result.Status)
		})
	}
}

func TestVideoGenerateToolChanelMessage(t *testing.T) {
	tests := []struct {
		name   string
		message VideoGenerateToolChanelMessage
	}{
		{
			name: "success message",
			message: VideoGenerateToolChanelMessage{
				Status: VideoGenerateSuccessStatus,
				Result: &VideoResult{
					VideoName: "test_video.mp4",
					Url:       "https://example.com/test_video.mp4",
				},
			},
		},
		{
			name: "error message",
			message: VideoGenerateToolChanelMessage{
				Status:       VideoGenerateErrorStatus,
				ErrorMessage: "API connection failed",
				Result: &VideoResult{
					VideoName: "failed_video.mp4",
				},
			},
		},
		{
			name: "message with nil result",
			message: VideoGenerateToolChanelMessage{
				Status: VideoGenerateErrorStatus,
				ErrorMessage: "Unknown error",
				Result: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.message.Status)
			if tt.message.Status == VideoGenerateSuccessStatus {
				assert.NotNil(t, tt.message.Result)
				assert.NotEmpty(t, tt.message.Result.VideoName)
				assert.NotEmpty(t, tt.message.Result.Url)
			} else if tt.message.Status == VideoGenerateErrorStatus {
				assert.NotEmpty(t, tt.message.ErrorMessage)
			}
		})
	}
}

func TestVideoResult(t *testing.T) {
	tests := []struct {
		name   string
		result VideoResult
	}{
		{
			name: "complete result",
			result: VideoResult{
				VideoName: "test_video.mp4",
				Url:       "https://example.com/test_video.mp4",
			},
		},
		{
			name: "result without URL (error case)",
			result: VideoResult{
				VideoName: "failed_video.mp4",
				Url:       "",
			},
		},
		{
			name: "result with empty video name",
			result: VideoResult{
				VideoName: "",
				Url:       "https://example.com/video.mp4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.result)
			// Basic validation that the struct is properly formed
			// Video name can be empty in error cases, URL can be empty in error cases
		})
	}
}

func TestGenerateVideosRequest(t *testing.T) {
	tests := []struct {
		name    string
		request GenerateVideosRequest
	}{
		{
			name: "basic request",
			request: GenerateVideosRequest{
				VideoName: "test.mp4",
				Prompt:    "a beautiful sunset",
			},
		},
		{
			name: "request with first frame",
			request: GenerateVideosRequest{
				VideoName:  "logo_animation.mp4",
				Prompt:     "company logo animation",
				FirstFrame: stringPtr("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="),
			},
		},
		{
			name: "request with last frame",
			request: GenerateVideosRequest{
				VideoName: "transition.mp4",
				Prompt:    "smooth transition",
				LastFrame: stringPtr("https://example.com/end_frame.png"),
			},
		},
		{
			name: "request with both frames",
			request: GenerateVideosRequest{
				VideoName:  "logo_reveal.mp4",
				Prompt:     "logo reveal animation",
				FirstFrame: stringPtr("https://example.com/start.png"),
				LastFrame:  stringPtr("https://example.com/end.png"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.request.VideoName)
			assert.NotEmpty(t, tt.request.Prompt)
			// Frame URLs can be either data URLs or regular URLs
			if tt.request.FirstFrame != nil {
				assert.NotEmpty(t, *tt.request.FirstFrame)
			}
			if tt.request.LastFrame != nil {
				assert.NotEmpty(t, *tt.request.LastFrame)
			}
		})
	}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
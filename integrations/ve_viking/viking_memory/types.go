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

package viking_memory

type CollectionInfoRequest struct {
	CollectionName string `json:"collection_name"`
	ProjectName    string `json:"project_name"`
	ResourceId     string `json:"resource_id"`
}

type CollectionCreateRequest struct {
	CollectionName      string
	ProjectName         string
	Description         string
	ResourceId          string
	BuiltinEventTypes   []string `json:"BuiltinEventTypes,omitempty"`
	BuiltinProfileTypes []string `json:"BuiltinProfileTypes,omitempty"`
}

type CollectionSearchMemoryRequest struct {
	CollectionName string `json:"collection_name,omitempty"`
	ProjectName    string `json:"project_name,omitempty"`
	ResourceId     string `json:"resource_id,omitempty"`
	Query          string `json:"query,omitempty"`
	Filter         Filter `json:"filter,omitempty"`
	Limit          int    `json:"limit,omitempty"`
}

type Filter struct {
	UserId     []string `json:"user_id,omitempty"`
	MemoryType []string `json:"memory_type,omitempty"`
}

type AddSessionRequest struct {
	CollectionName string     `json:"collection_name"`
	ProjectName    string     `json:"project_name"`
	ResourceId     string     `json:"resource_id" `
	SessionId      string     `json:"session_id"`
	Messages       []*Message `json:"messages,omitempty"`
	Metadata       Metadata   `json:"metadata,omitempty"`
	Profiles       []Profile  `json:"profiles,omitempty"`
}

type Message struct {
	Role     string `json:"role,omitempty"`
	Content  string `json:"content,omitempty"`
	RoleId   string `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	Time     int64  `json:"time,omitempty"`
}

type Metadata struct {
	DefaultUserId        string `json:"default_user_id,omitempty"`
	DefaultUserName      string `json:"default_user_name,omitempty"`
	DefaultAssistantId   string `json:"default_assistant_id,omitempty"`
	DefaultAssistantName string `json:"default_assistant_name,omitempty"`
	Time                 int64  `json:"time,omitempty"`
	GroupId              string `json:"group_id,omitempty"`
}

type Profile struct {
	ProfileType  string
	ProfileScope map[string]interface{}
}

type CollectionSearchMemoryResponse struct {
	Code      int64                               `json:"code,omitempty"`
	Message   string                              `json:"message,omitempty"`
	Data      *CollectionSearchMemoryResponseData `json:"data,omitempty"`
	RequestID string                              `json:"request_id,omitempty"`
}
type CollectionSearchMemoryResponseData struct {
	CollectionName string                          `json:"collection_name,omitempty"`
	Count          int32                           `json:"count,omitempty"`
	ResultList     []*CollectionSearchResponseItem `json:"result_list,omitempty"`
}

type CollectionSearchResponseItem struct {
	MemoryInfo *MemoryInfo `json:"memory_info,omitempty"`
	Time       int64       `json:"time,omitempty"`
}

type MemoryInfo struct {
	Summary string `json:"summary,omitempty"`
}

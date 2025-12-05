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

package knowledge_base

import (
	"fmt"

	"github.com/volcengine/veadk-go/knowledgebase"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

// knowledgeBaseTool is a tool for loading the common knowledge base.
// In the future, we will support multiple knowledge bases for different users.
type knowledgeBaseTool struct {
	name        string
	description string
	backend     knowledgebase.KnowledgeBase
}

func NewKnowledgeBaseTool(name, description string, backend knowledgebase.KnowledgeBase) tool.Tool {
	return &knowledgeBaseTool{
		name:        name,
		description: description,
		backend:     backend,
	}
}

func (k *knowledgeBaseTool) Name() string {
	return k.name
}

func (k *knowledgeBaseTool) Description() string {
	return k.description
}

func (k *knowledgeBaseTool) IsLongRunning() bool {
	return false
}

// Declaration returns the GenAI FunctionDeclaration for the knowledgeBase tool.
//
// This declaration allows the LLM to understand and call the tool
// by specifying the function name, a detailed description of its
// purpose, and the required input parameters (schema).
func (k *knowledgeBaseTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        k.name,
		Description: k.description,
		Parameters: &genai.Schema{
			Title:       "query",
			Type:        genai.TypeString,
			Description: "The query for loading the knowledge base for this tool.",
		},
	}
}

// Run implements tool.Tool.
func (k *knowledgeBaseTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	m, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected args type, got: %T", args)
	}

	chunks, err := k.backend.Search(m)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"knowledge": chunks,
	}, nil
}

// ProcessRequest processes the LLM request. It packs the tool, appends initial
// instructions, and processes any load artifacts function calls.
func (k *knowledgeBaseTool) ProcessRequest(ctx tool.Context, req *model.LLMRequest) error {
	return nil
}

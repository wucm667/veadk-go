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

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	veagent "github.com/volcengine/veadk-go/agent/llmagent"
	"github.com/volcengine/veadk-go/common"
	"github.com/volcengine/veadk-go/tool/builtin_tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	cfg := veagent.Config{
		ModelName:    common.DEFAULT_MODEL_AGENT_NAME,
		ModelAPIBase: common.DEFAULT_MODEL_AGENT_API_BASE,
		ModelAPIKey:  os.Getenv("MODEL_API_KEY"),
	}
	imageGenerate, err := builtin_tools.NewImageGenerateTool(&builtin_tools.ImageGenerateConfig{
		ModelName: common.DEFAULT_MODEL_IMAGE_NAME,
		BaseURL:   common.DEFAULT_MODEL_IMAGE_API_BASE,
		APIKey:    os.Getenv("MODEL_API_KEY"),
	})
	if err != nil {
		fmt.Printf("NewLLMAgent failed: %v", err)
		return
	}

	cfg.Tools = []tool.Tool{imageGenerate}

	a, err := veagent.New(&cfg)
	if err != nil {
		fmt.Printf("NewLLMAgent failed: %v", err)
		return
	}

	sessionService := session.InMemoryService()
	runner, err := runner.New(runner.Config{
		AppName:        "image_generator",
		Agent:          a,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatal(err)
	}

	session, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: "image_generator",
		UserID:  "user1234",
	})
	if err != nil {
		log.Fatal(err)
	}

	runImageGenerate(ctx, runner, session.Session.ID(), "给我生成两组图片，每组图片各两张，第一组充满活力的特写编辑肖像，模特眼神犀利，头戴雕塑感帽子，色彩拼接丰富，眼部焦点锐利，景深较浅，具有Vogue杂志封面的美学风格，采用中画幅拍摄，工作室灯光效果强烈；第二组阳光男孩的图片。")

}

func runImageGenerate(ctx context.Context, r *runner.Runner, sessionID string, prompt string) {
	fmt.Printf("\n> %s\n", prompt)
	events := r.Run(
		ctx,
		"user1234",
		sessionID,
		genai.NewContentFromText(prompt, genai.RoleUser),
		agent.RunConfig{
			StreamingMode: agent.StreamingModeNone,
		},
	)
	for event, err := range events {
		if err != nil {
			log.Fatalf("ERROR during agent execution: %v", err)
		}

		for _, v := range event.Content.Parts {
			if v.Text != "" {
				fmt.Printf("Agent Response: %s\n", v.Text)
			}
		}
	}
}

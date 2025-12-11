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
	videoGenerate, err := builtin_tools.NewVideoGenerateTool(&builtin_tools.VideoGenerateConfig{
		ModelName: common.DEFAULT_MODEL_VIDEO_NAME,
		BaseURL:   common.DEFAULT_MODEL_VIDEO_API_BASE,
		APIKey:    os.Getenv("MODEL_API_KEY"),
	})
	if err != nil {
		fmt.Printf("NewLLMAgent failed: %v", err)
		return
	}

	cfg.Tools = []tool.Tool{videoGenerate}

	a, err := veagent.New(&cfg)
	if err != nil {
		fmt.Printf("NewLLMAgent failed: %v", err)
		return
	}
	sessionService := session.InMemoryService()
	runner, err := runner.New(runner.Config{
		AppName:        "video_generator",
		Agent:          a,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatal(err)
	}

	session, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: "video_generator",
		UserID:  "user1234",
	})
	if err != nil {
		log.Fatal(err)
	}

	runImageGenerate(ctx, runner, session.Session.ID(), "多个镜头。一名侦探进入一间光线昏暗的房间。他检查桌上的线索，手里拿起桌上的某个物品。镜头转向他正在思索。 --ratio 16:9")

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

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

package knowledgebase

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/volcengine/veadk-go/integrations/ve_viking_knowledge"
	"github.com/volcengine/veadk-go/log"
)

const (
	VikingKnowledgeBaseIndexNotExistCode = 1000005
	VikingKnowledgeBaseSuccessCode       = 0
	DefaultTopK                          = 5
	DefaultChunkDiffusionCount           = 3
)

var NewVikingKnowledgeBaseErr = errors.New("NewVikingKnowledgeBase error")
var VikingKnowledgeBaseSearchErr = errors.New("VikingKnowledgeBase search error")

type Config struct {
	AK                  string
	SK                  string
	Index               string
	Project             string
	Region              string
	CreateIfNotExist    bool
	TopK                int32
	ChunkDiffusionCount int32
}

type VikingKnowledgeBackend struct {
	client *ve_viking_knowledge.Client
	config *Config
	index  string
}

func NewVikingKnowledgeBackend(cfg *Config) (KnowledgeBase, error) {
	client, err := ve_viking_knowledge.New(&ve_viking_knowledge.Client{
		Index:   cfg.Index,
		Project: cfg.Project,
		Region:  cfg.Region,
		AK:      cfg.AK,
		SK:      cfg.SK,
	})
	if err != nil {
		return nil, fmt.Errorf("%w : %w", NewVikingKnowledgeBaseErr, err)
	}

	collectionInfo, err := client.CollectionInfo()
	if err != nil {
		return nil, fmt.Errorf("%w : get viking collection info error: %w", NewVikingKnowledgeBaseErr, err)
	}

	if collectionInfo.Code == VikingKnowledgeBaseIndexNotExistCode {
		if cfg.CreateIfNotExist {
			_, err := client.CollectionCreate("")
			if err != nil {
				return nil, fmt.Errorf("%w : create viking knowledge index error: %w", NewVikingKnowledgeBaseErr, err)
			}
			log.Info("Create viking knowledge index", cfg.Index, "successfully")
		} else {
			return nil, fmt.Errorf("%w : viking index not exist", NewVikingKnowledgeBaseErr)
		}
	}

	if cfg.TopK <= 0 {
		cfg.TopK = DefaultTopK
	}
	if cfg.ChunkDiffusionCount <= 0 {
		cfg.ChunkDiffusionCount = DefaultChunkDiffusionCount
	}
	return &VikingKnowledgeBackend{
		client: client,
		config: cfg,
		// todo create new index
		index: cfg.Index,
	}, nil
}

func (v *VikingKnowledgeBackend) Search(query string) ([]KnowledgeEntry, error) {
	chunks, err := v.client.SearchKnowledge(
		query, v.config.TopK, nil, true, v.config.ChunkDiffusionCount)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", VikingKnowledgeBaseSearchErr, err)
	}

	if chunks.Code != VikingKnowledgeBaseSuccessCode {
		return nil, fmt.Errorf("%w : with bad code %d, message:%s", VikingKnowledgeBaseSearchErr, chunks.Code, chunks.Message)
	}

	var entries []KnowledgeEntry
	for _, item := range chunks.Data.ResultList {
		var metadata = make(map[string]any)
		if item.DocInfo.DocMeta != "" {
			err = json.Unmarshal([]byte(item.DocInfo.DocMeta), &metadata)
			if err != nil {
				return nil, fmt.Errorf("%w : Unmarshal DocMeta error:%w", VikingKnowledgeBaseSearchErr, err)
			}
		}
		entries = append(entries, KnowledgeEntry{
			Content:  item.Content,
			Metadata: metadata,
		})
	}

	return entries, nil
}

func (v *VikingKnowledgeBackend) Index() string {
	return v.index
}

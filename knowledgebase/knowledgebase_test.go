package knowledgebase

import (
    "errors"
    "reflect"
    "testing"

    "github.com/volcengine/veadk-go/knowledgebase/ktypes"
)

type dummyBackend struct{}

func (d *dummyBackend) Index() string { return "dummy" }
func (d *dummyBackend) Search(query string, opts ...map[string]any) ([]ktypes.KnowledgeEntry, error) {
    return []ktypes.KnowledgeEntry{}, nil
}
func (d *dummyBackend) AddFromText(text []string, opts ...map[string]any) error { return nil }
func (d *dummyBackend) AddFromFiles(files []string, opts ...map[string]any) error { return nil }
func (d *dummyBackend) AddFromDirectory(directory string, opts ...map[string]any) error { return nil }

func TestNewKnowledgeBase_WithBackendInstance_Defaults(t *testing.T) {
    backend := &dummyBackend{}
    kb, err := NewKnowledgeBase(backend)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if kb.Backend != backend {
        t.Fatalf("backend not set")
    }
    if kb.Name != DefaultName {
        t.Fatalf("default name not applied: got %s", kb.Name)
    }
    if kb.Description != DefaultDescription {
        t.Fatalf("default description not applied")
    }
}

func TestNewKnowledgeBase_WithOptions(t *testing.T) {
    cfg := map[string]int{"x": 1}
    kb, err := NewKnowledgeBase(&dummyBackend{}, WithName("kb"), WithDescription("desc"), WithBackendConfig(cfg))
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if kb.Name != "kb" {
        t.Fatalf("name option not applied: %s", kb.Name)
    }
    if kb.Description != "desc" {
        t.Fatalf("description option not applied: %s", kb.Description)
    }
    if !reflect.DeepEqual(kb.BackendConfig, cfg) {
        t.Fatalf("backend config not applied")
    }
}

func TestGetKnowledgeBackend_InvalidBackend(t *testing.T) {
    _, err := getKnowledgeBackend(ktypes.RedisBackend, nil)
    if !errors.Is(err, InvalidKnowledgeBackendErr) {
        t.Fatalf("expected InvalidKnowledgeBackendErr, got %v", err)
    }
}

func TestGetKnowledgeBackend_VikingInvalidConfig(t *testing.T) {
    b, err := getKnowledgeBackend(ktypes.VikingBackend, struct{}{})
    if err == nil || b != nil {
        t.Fatalf("expected error for invalid config type")
    }
}

func TestNewKnowledgeBase_StringBackend_Invalid(t *testing.T) {
    _, err := NewKnowledgeBase(ktypes.RedisBackend)
    if !errors.Is(err, InvalidKnowledgeBackendErr) {
        t.Fatalf("expected InvalidKnowledgeBackendErr, got %v", err)
    }
}

func TestNewKnowledgeBase_NilBackend_Error(t *testing.T) {
    _, err := NewKnowledgeBase(nil)
    if !errors.Is(err, InvalidKnowledgeBackendErr) {
        t.Fatalf("expected InvalidKnowledgeBackendErr for nil backend, got %v", err)
    }
}

func TestNewKnowledgeBase_UnsupportedType(t *testing.T) {
    _, err := NewKnowledgeBase(123)
    if err == nil {
        t.Fatalf("expected error for unsupported backend type")
    }
}

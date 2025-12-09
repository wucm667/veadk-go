package utils

import (
	"errors"
	"testing"
)

func TestExtractOptsValue_Success(t *testing.T) {
	opts := map[string]any{"a": 123, "b": "x"}
    ia := func(v any) (int, bool) { x, ok := v.(int); return x, ok }
    sa := func(v any) (string, bool) { s, ok := v.(string); return s, ok }

    gotInt, err := ExtractOptsValue[int](opts, "a", ia)
	if err != nil {
		t.Fatal(err)
	}
	if gotInt != 123 {
		t.Fatalf("got=%d want=123", gotInt)
	}

    gotStr, err := ExtractOptsValue[string](opts, "b", sa)
	if err != nil {
		t.Fatal(err)
	}
	if gotStr != "x" {
		t.Fatalf("got=%s want=x", gotStr)
	}
}

func TestExtractOptsValue_Errors(t *testing.T) {
    ia := func(v any) (int, bool) { x, ok := v.(int); return x, ok }

    _, err := ExtractOptsValue[int](nil, "a", ia)
	if !errors.Is(err, OptsNIlErr) {
		t.Fatalf("want OptsNIlErr got=%v", err)
	}

    _, err = ExtractOptsValue[int](map[string]any{}, "a", ia)
	if !errors.Is(err, OptsInvalidKeyErr) {
		t.Fatalf("want OptsInvalidKeyErr got=%v", err)
	}

    _, err = ExtractOptsValue[int](map[string]any{"a": "x"}, "a", ia)
	if !errors.Is(err, OptsAssertTypeErr) {
		t.Fatalf("want OptsAssertTypeErr got=%v", err)
	}
}

func TestExtractOptsValueWithDefault_DefaultPaths(t *testing.T) {
    ia := func(v any) (int, bool) { x, ok := v.(int); return x, ok }
    ba := func(v any) (bool, bool) { b, ok := v.(bool); return b, ok }

    if got := ExtractOptsValueWithDefault[int](nil, "a", ia, 7); got != 7 {
		t.Fatalf("nil opts default got=%d", got)
	}
    if got := ExtractOptsValueWithDefault[int](map[string]any{}, "a", ia, 8); got != 8 {
		t.Fatalf("missing key default got=%d", got)
	}
    if got := ExtractOptsValueWithDefault[int](map[string]any{"a": "x"}, "a", ia, 9); got != 9 {
		t.Fatalf("assert fail default got=%d", got)
	}

    if got := ExtractOptsValueWithDefault[int](map[string]any{"a": 5}, "a", ia, 0); got != 5 {
		t.Fatalf("success got=%d want=5", got)
	}
    if got := ExtractOptsValueWithDefault[bool](map[string]any{"b": true}, "b", ba, false); got != true {
		t.Fatalf("bool success got=%v want=true", got)
	}
}


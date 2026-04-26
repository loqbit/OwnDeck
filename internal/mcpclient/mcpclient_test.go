package mcpclient

import (
	"context"
	"testing"
	"time"
)

func TestIntrospect_CommandNotFound(t *testing.T) {
	result := IntrospectWithTimeout(
		context.Background(),
		3*time.Second,
		"nonexistent-command-12345",
		nil, nil, "",
	)

	if result.Health != HealthUnreachable {
		t.Errorf("health = %q, want %q", result.Health, HealthUnreachable)
	}
	if result.Error == "" {
		t.Error("expected non-empty error message")
	}
	if len(result.Tools) != 0 {
		t.Errorf("expected 0 tools, got %d", len(result.Tools))
	}
}

func TestIntrospect_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	result := IntrospectWithTimeout(
		ctx,
		3*time.Second,
		"echo",
		[]string{"hello"}, nil, "",
	)

	if result.Health != HealthUnreachable {
		t.Errorf("health = %q, want %q", result.Health, HealthUnreachable)
	}
}

func TestIntrospect_Duration(t *testing.T) {
	result := IntrospectWithTimeout(
		context.Background(),
		1*time.Second,
		"nonexistent-command-12345",
		nil, nil, "",
	)

	if result.Duration == 0 {
		t.Error("expected non-zero duration")
	}
}

func TestBuildEnv(t *testing.T) {
	env := buildEnv(map[string]string{"MY_KEY": "my_value"})

	found := false
	for _, e := range env {
		if e == "MY_KEY=my_value" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected MY_KEY=my_value in env, got %v", env)
	}

	// Should also contain system PATH
	hasPath := false
	for _, e := range env {
		if len(e) > 5 && e[:5] == "PATH=" {
			hasPath = true
		}
	}
	if !hasPath {
		t.Error("expected PATH in env")
	}
}

func TestBuildEnv_NilExtra(t *testing.T) {
	env := buildEnv(nil)
	// Should still contain system vars
	if env == nil {
		t.Error("expected non-nil env even with nil extra")
	}
}

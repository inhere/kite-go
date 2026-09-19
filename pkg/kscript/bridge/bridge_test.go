package bridge

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	kscript2 "github.com/gookit/kscript"

	"github.com/inhere/kite-go/pkg/kscript"
)

func newLegacy(t *testing.T, scripts map[string]any) *kscript.Runner {
	t.Helper()
	runner := kscript.NewRunner()
	runner.PathResolver = func(path string) string { return path }
	runner.Scripts = scripts
	// The default script app dirs point at a configured path that does not
	// exist in a test environment.
	runner.ScriptAppDirs = nil
	runner.ScriptAppExts = nil
	return runner
}

func TestBridgeRunsLegacyTask(t *testing.T) {
	legacy := newLegacy(t, map[string]any{
		"tpl": map[string]any{
			"vars": map[string]any{"who": "bridged"},
			"run":  "echo ${who}",
		},
	})
	var out bytes.Buffer
	b := New(legacy, WithIO(kscript2.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}), WithBaseDir(t.TempDir()))
	found, err := b.TryRun(context.Background(), "tpl", nil, &kscript.RunCtx{Silent: true})
	if err != nil {
		t.Fatalf("TryRun: %v", err)
	}
	if !found {
		t.Fatal("task was not found")
	}
	if !strings.Contains(out.String(), "bridged") {
		t.Fatalf("legacy task variable was not rendered: %q", out.String())
	}
}

func TestBridgeReportsUnknownName(t *testing.T) {
	legacy := newLegacy(t, map[string]any{"known": "echo hi"})
	b := New(legacy, WithBaseDir(t.TempDir()))
	found, err := b.TryRun(context.Background(), "definitely-unknown", nil, nil)
	if err != nil || found {
		t.Fatalf("found=%v err=%v", found, err)
	}
}

func TestBridgeRunsExplicitShellTask(t *testing.T) {
	shell, script := "sh", `printf 'shell-ok'`
	if runtime.GOOS == "windows" {
		shell, script = "cmd", "echo shell-ok"
	}
	legacy := newLegacy(t, map[string]any{
		"typed": map[string]any{"type": shell, "run": script},
	})
	var out bytes.Buffer
	b := New(legacy, WithIO(kscript2.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}), WithBaseDir(t.TempDir()))
	found, err := b.TryRun(context.Background(), "typed", nil, &kscript.RunCtx{Silent: true})
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if !strings.Contains(out.String(), "shell-ok") {
		t.Fatalf("output=%q", out.String())
	}
}

func TestBridgePassesArgumentsToTask(t *testing.T) {
	legacy := newLegacy(t, map[string]any{
		"args": "echo ${1}-${2}",
	})
	var out bytes.Buffer
	b := New(legacy, WithIO(kscript2.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}), WithBaseDir(t.TempDir()))
	if _, err := b.TryRun(context.Background(), "args", []string{"one", "two"}, &kscript.RunCtx{Silent: true}); err != nil {
		t.Fatalf("TryRun: %v", err)
	}
	if !strings.Contains(out.String(), "one-two") {
		t.Fatalf("arguments were not passed: %q", out.String())
	}
}

func TestBridgeRunsScriptFile(t *testing.T) {
	dir := t.TempDir()
	ext, bin, body, marker := ".sh", "sh", "#!/bin/sh\nprintf file-ok\n", "file-ok"
	if runtime.GOOS == "windows" {
		ext, bin, body, marker = ".bat", "cmd /C", "@echo file-ok\r\n", "file-ok"
	}
	file := filepath.Join(dir, "hello"+ext)
	if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	legacy := newLegacy(t, map[string]any{})
	legacy.ScriptDirs = []string{dir}
	legacy.AllowedExt = []string{ext}
	legacy.ExtToBinMap = map[string]string{ext: bin}

	var out bytes.Buffer
	b := New(legacy, WithIO(kscript2.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}), WithBaseDir(t.TempDir()))
	found, err := b.TryRun(context.Background(), "hello", nil, &kscript.RunCtx{Silent: true})
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if !strings.Contains(out.String(), marker) {
		t.Fatalf("script file output=%q", out.String())
	}
}

func TestBridgePassesSettingsVarsAndRuntimeVars(t *testing.T) {
	legacy := newLegacy(t, map[string]any{
		"__settings": map[string]any{
			"vars": map[string]any{"project": "kite"},
			"env":  map[string]any{"KS_PROJECT": "kite"},
		},
		"tpl": "echo ${vars.project}-${gvs.app}",
	})
	var out bytes.Buffer
	b := New(legacy,
		WithIO(kscript2.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}),
		WithBaseDir(t.TempDir()),
		WithVarNames("gvs"),
	)
	found, err := b.TryRun(context.Background(), "tpl", nil, &kscript.RunCtx{
		Silent: true,
		Vars:   map[string]string{},
		AppendVarsFn: func(data map[string]any) map[string]any {
			data["gvs"] = map[string]any{"app": "kitego"}
			return data
		},
	})
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if !strings.Contains(out.String(), "kite-kitego") {
		t.Fatalf("settings and runtime vars not rendered: %q", out.String())
	}
}

func TestBridgeReportsConversionProblems(t *testing.T) {
	// A dependency on a task that does not exist cannot be converted.
	legacy := newLegacy(t, map[string]any{
		"broken": map[string]any{"deps": []any{"missing"}, "run": "echo hi"},
	})
	b := New(legacy, WithBaseDir(t.TempDir()))
	found, err := b.TryRun(context.Background(), "broken", nil, nil)
	if err == nil {
		t.Fatalf("expected a conversion error, found=%v", found)
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Fatalf("unhelpful error: %v", err)
	}
	// With the fallback enabled the legacy runner still runs the task.
	var out bytes.Buffer
	fallback := New(legacy,
		WithBaseDir(t.TempDir()),
		WithLegacyFallback(true),
		WithIO(kscript2.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}),
	)
	found, fallbackErr := fallback.TryRun(context.Background(), "broken", nil, &kscript.RunCtx{Silent: true})
	if !found {
		t.Fatal("fallback did not report the task as found")
	}
	// The legacy runner reports its own error, so the conversion error must not
	// be the one surfaced.
	if fallbackErr != nil && strings.Contains(fallbackErr.Error(), "invalid_definition") {
		t.Fatalf("conversion error leaked through the fallback: %v", fallbackErr)
	}
}

func TestBridgeDefinitionIsCachedPerShell(t *testing.T) {
	legacy := newLegacy(t, map[string]any{"plain": "echo hi"})
	b := New(legacy, WithBaseDir(t.TempDir()))
	first, err := b.runnerFor("", []string{"gvs"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := b.runnerFor("", []string{"gvs"})
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("converted runner was not cached")
	}
	third, err := b.runnerFor("sh", []string{"gvs"})
	if err != nil {
		t.Fatal(err)
	}
	if third == first {
		t.Fatal("a different shell must convert separately")
	}
}

func TestBridgeDryRunDoesNotExecute(t *testing.T) {
	legacy := newLegacy(t, map[string]any{"greet": "echo should-not-run"})
	var out bytes.Buffer
	b := New(legacy, WithIO(kscript2.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}), WithBaseDir(t.TempDir()))
	if _, err := b.TryRun(context.Background(), "greet", nil, &kscript.RunCtx{Silent: true, DryRun: true}); err != nil {
		t.Fatalf("TryRun: %v", err)
	}
	// Every printed line is a plan line; real execution would print the command
	// output without the prefix.
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "[dry-run] greet") {
		t.Fatalf("dry run did not print a plan: %q", out.String())
	}
	for _, line := range lines {
		if !strings.HasPrefix(line, "[dry-run]") {
			t.Fatalf("dry run executed the task: %q", out.String())
		}
	}
	if !strings.Contains(out.String(), "exec") {
		t.Fatalf("plan does not describe the action: %q", out.String())
	}
}

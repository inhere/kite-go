package bridge

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gookit/taskrun"

	"github.com/inhere/kite-go/pkg/kscript"
)

// shellCmd builds a legacy plain command line that runs a shell command, so the
// same fixture works for both engines on every platform. The legacy engine
// splits it with goutil's cmdline parser and the new engine splits it with
// formats.SplitCommandLine; both must produce the same argv.
func shellCmd(line string) string {
	if runtime.GOOS == "windows" {
		return `cmd /C "` + line + `"`
	}
	return `sh -c "` + line + `"`
}

// parityFixture exercises the semantics the migration must preserve: task
// settings, dependencies, task variables, command variables, task references,
// argument shorthands and shell expanded environment values.
func parityFixture() map[string]any {
	return map[string]any{
		"__settings": map[string]any{
			"vars": map[string]any{"project": "kite"},
		},
		"prep": shellCmd("echo prep >> trace.txt"),
		"collect": map[string]any{
			"deps": []any{"prep"},
			"vars": map[string]any{"who": "pack"},
			"run": []any{
				// The legacy renderer only substitutes the $1 form; ${1} stays
				// literal there. The new engine supports both, which is
				// recorded as an intentional fix in the migration document.
				shellCmd("echo A-${who}-$1 >> trace.txt"),
				shellCmd("echo B-${vars.project} >> trace.txt"),
				map[string]any{
					"name": "third",
					"vars": map[string]any{"inner": "deep"},
					"run":  shellCmd("echo C-${inner} >> trace.txt"),
				},
				// A command map with only "task" and no "run" was ignored by the
				// legacy loader, so the supported form is an @task: reference.
				map[string]any{"name": "ref", "run": "@task:prep"},
				shellCmd("echo D-$KS_PROJECT >> trace.txt"),
				shellCmd("echo E-$2 >> trace.txt"),
			},
		},
	}
}

// runLegacy executes the fixture with the in-tree runner and returns the trace
// file it produced.
func runLegacy(t *testing.T, fixture map[string]any, args []string) string {
	t.Helper()
	dir := t.TempDir()
	runner := kscript.NewRunner()
	runner.PathResolver = func(path string) string { return path }
	runner.Scripts = fixture
	runner.ScriptAppDirs = nil
	runner.ScriptAppExts = nil

	err := captureStdout(func() error {
		return runner.Run("collect", args, &kscript.RunCtx{
			Silent:  true,
			Workdir: dir,
			Env:     map[string]string{"KS_PROJECT": "kite"},
		})
	})
	if err != nil {
		t.Fatalf("legacy run failed: %v", err)
	}
	return readTrace(t, dir)
}

// runBridged executes the same fixture with the standalone library.
func runBridged(t *testing.T, fixture map[string]any, args []string) string {
	t.Helper()
	dir := t.TempDir()
	runner := kscript.NewRunner()
	runner.PathResolver = func(path string) string { return path }
	runner.Scripts = fixture
	runner.ScriptAppDirs = nil
	runner.ScriptAppExts = nil

	var out bytes.Buffer
	b := New(runner,
		WithBaseDir(dir),
		WithIO(taskrun.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}),
	)
	found, err := b.TryRun(context.Background(), "collect", args, &kscript.RunCtx{
		Silent:  true,
		Workdir: dir,
		Env:     map[string]string{"KS_PROJECT": "kite"},
	})
	if err != nil {
		t.Fatalf("bridged run failed: %v", err)
	}
	if !found {
		t.Fatal("bridged run did not find the task")
	}
	return readTrace(t, dir)
}

func TestEngineParityOnLegacyFixture(t *testing.T) {
	args := []string{"one", "two"}
	legacyTrace := runLegacy(t, parityFixture(), args)
	bridgedTrace := runBridged(t, parityFixture(), args)

	if legacyTrace != bridgedTrace {
		t.Fatalf("engine output differs:\nlegacy:\n%s\nbridged:\n%s", legacyTrace, bridgedTrace)
	}
	// Guard against a fixture that silently produced nothing.
	for _, want := range []string{
		"prep", "A-pack-one", "B-kite", "C-deep", "D-kite", "E-two",
	} {
		if !strings.Contains(legacyTrace, want) {
			t.Fatalf("trace is missing %q:\n%s", want, legacyTrace)
		}
	}
}

func TestEngineParityOnFailure(t *testing.T) {
	args := []string{"one", "two"}
	fixture := func() map[string]any {
		return map[string]any{
			"boom": shellCmd("exit 3"),
		}
	}
	legacyErr := captureLegacyError(t, fixture(), args)
	bridgeErr := captureBridgeError(t, fixture(), args)
	if (legacyErr == nil) != (bridgeErr == nil) {
		t.Fatalf("failure behavior differs: legacy=%v bridged=%v", legacyErr, bridgeErr)
	}
	if legacyErr == nil {
		t.Fatal("fixture should have failed on both engines")
	}
}

func captureLegacyError(t *testing.T, fixture map[string]any, args []string) error {
	t.Helper()
	dir := t.TempDir()
	runner := kscript.NewRunner()
	runner.PathResolver = func(path string) string { return path }
	runner.Scripts = fixture
	runner.ScriptAppDirs = nil
	runner.ScriptAppExts = nil
	var runErr error
	_ = captureStdout(func() error {
		runErr = runner.Run("boom", args, &kscript.RunCtx{Silent: true, Workdir: dir})
		return nil
	})
	return runErr
}

func captureBridgeError(t *testing.T, fixture map[string]any, args []string) error {
	t.Helper()
	dir := t.TempDir()
	runner := kscript.NewRunner()
	runner.PathResolver = func(path string) string { return path }
	runner.Scripts = fixture
	runner.ScriptAppDirs = nil
	runner.ScriptAppExts = nil
	var out bytes.Buffer
	b := New(runner, WithBaseDir(dir), WithIO(taskrun.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}))
	_, err := b.TryRun(context.Background(), "boom", args, &kscript.RunCtx{Silent: true, Workdir: dir})
	return err
}

func readTrace(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "trace.txt"))
	if err != nil {
		t.Fatalf("read trace: %v", err)
	}
	return string(data)
}

// captureStdout redirects the legacy engine's own logging to a pipe so test
// output stays readable.
func captureStdout(fn func() error) error {
	saved := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		return fn()
	}
	os.Stdout = writer
	drained := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, reader)
		close(drained)
	}()
	runErr := fn()
	_ = writer.Close()
	os.Stdout = saved
	<-drained
	_ = reader.Close()
	return runErr
}

package cmdbiz

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gookit/config/v2"

	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/kscript"
)

// newScriptRunner installs a script runner and config into the app box, the way
// boot does, so the engine switch can be exercised without starting the CLI.
func newScriptRunner(t *testing.T, scripts map[string]any, engine string) *kscript.Runner {
	t.Helper()
	cfg := config.New("test")
	cfg.Set("script_engine", engine)
	app.Set(app.ObjConf, cfg)

	runner := kscript.NewRunner()
	runner.PathResolver = func(path string) string { return path }
	runner.Scripts = scripts
	runner.ScriptAppDirs = nil
	runner.ScriptAppExts = nil
	app.Scripts = runner
	return runner
}

func TestScriptEngineSelection(t *testing.T) {
	newScriptRunner(t, map[string]any{"t": "echo hi"}, ScriptEngineLegacy)
	if got := ScriptEngine(); got != ScriptEngineLegacy {
		t.Fatalf("engine=%q", got)
	}
	newScriptRunner(t, map[string]any{"t": "echo hi"}, ScriptEngineTaskrun)
	if got := ScriptEngine(); got != ScriptEngineTaskrun {
		t.Fatalf("engine=%q", got)
	}
	// An unknown value falls back to the legacy runner.
	newScriptRunner(t, map[string]any{"t": "echo hi"}, "whatever")
	if got := ScriptEngine(); got != ScriptEngineLegacy {
		t.Fatalf("engine=%q", got)
	}
}

func TestRunScriptNameUsesKscriptEngine(t *testing.T) {
	newScriptRunner(t, map[string]any{
		"tpl": map[string]any{
			"vars": map[string]any{"who": "switched"},
			"run":  "echo ${who}",
		},
	}, ScriptEngineTaskrun)

	found, err := RunScriptName("tpl", nil, &kscript.RunCtx{Silent: true})
	if err != nil {
		t.Fatalf("RunScriptName: %v", err)
	}
	if !found {
		t.Fatal("script task was not found")
	}
}

func TestRunScriptNameReportsUnknownName(t *testing.T) {
	newScriptRunner(t, map[string]any{"known": "echo hi"}, ScriptEngineTaskrun)
	found, err := RunScriptName("definitely-unknown", nil, &kscript.RunCtx{Silent: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("unknown name must report found=false so Kite can fall back to a system command")
	}
}

func TestRunScriptOnlyUsesLegacyEngine(t *testing.T) {
	newScriptRunner(t, map[string]any{"t": "echo hi"}, ScriptEngineLegacy)
	if err := RunScriptOnly("definitely-unknown", nil, &kscript.RunCtx{Silent: true}); err == nil {
		t.Fatal("expected a not found error from the legacy runner")
	}
}

// The switched engine must really execute a task, and DryRun must not.
func TestRunScriptNameExecutesAndDryRunDoesNot(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	shell, script := "sh", "echo x > marker"
	if runtime.GOOS == "windows" {
		shell, script = "cmd", "echo x > marker"
	}
	newScriptRunner(t, map[string]any{
		"mark": map[string]any{"type": shell, "run": script},
	}, ScriptEngineTaskrun)

	if _, err := RunScriptName("mark", nil, &kscript.RunCtx{Silent: true, DryRun: true}); err != nil {
		t.Fatalf("dry run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "marker")); err == nil {
		t.Fatal("dry run executed the task")
	}

	if _, err := RunScriptName("mark", nil, &kscript.RunCtx{Silent: true}); err != nil {
		t.Fatalf("run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "marker")); err != nil {
		t.Fatalf("the switched engine did not execute the task: %v", err)
	}
}

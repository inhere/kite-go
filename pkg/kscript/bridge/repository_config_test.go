package bridge

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gookit/taskrun"

	"github.com/inhere/kite-go/pkg/kscript"
)

// repoRoot walks up until the module root is found.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("module root not found")
		}
		dir = parent
	}
}

// repositoryResolver mirrors the path aliases boot installs: $base is the
// repository root and $config is its config directory. Other aliases used by the
// shipped configuration ($os, $user) are expanded too.
func repositoryResolver(root string) func(string) string {
	configDir := filepath.Join(root, "config")
	home, _ := os.UserHomeDir()
	return func(path string) string {
		path = strings.ReplaceAll(path, "$os", runtime.GOOS)
		path = strings.ReplaceAll(path, "$user", filepath.Base(home))
		switch {
		case strings.HasPrefix(path, "$base/"):
			path = filepath.Join(root, strings.TrimPrefix(path, "$base/"))
		case path == "$base":
			path = root
		case strings.HasPrefix(path, "$config/"):
			path = filepath.Join(configDir, strings.TrimPrefix(path, "$config/"))
		case path == "$config":
			path = configDir
		}
		return path
	}
}

// TestRepositoryConfigConvertsCleanly converts the configuration this
// repository ships (config/module/scripts.yml and its per-OS optional file)
// and validates it with the standalone library. Nothing is executed: the check
// is conversion, validation and planning only.
func TestRepositoryConfigConvertsCleanly(t *testing.T) {
	root := repoRoot(t)
	resolve := repositoryResolver(root)
	scriptFile := resolve("$config/module/scripts.yml")
	if _, err := os.Stat(scriptFile); err != nil {
		t.Skipf("shipped script config is unavailable: %v", err)
	}

	runner := kscript.NewRunner()
	runner.PathResolver = resolve
	runner.DefineFiles = []string{"$config/module/scripts.yml", "?$config/module/scripts.$os.yml"}
	runner.ScriptDirs = []string{"$base/scripts"}
	// Auto discovery stays off so the test never scans the repository.
	runner.AutoTaskFiles = nil
	if err := runner.InitLoad(); err != nil {
		t.Fatalf("InitLoad: %v", err)
	}
	if len(runner.Scripts) == 0 {
		t.Fatal("the shipped config produced no script tasks")
	}

	b := New(runner, WithBaseDir(root), WithParseEnv(runner.ParseEnv))
	def, warnings, err := b.Definition("", b.runtimeVarNames(map[string]any{}), nil)
	if err != nil {
		t.Fatalf("conversion failed for the shipped config: %v", err)
	}
	if len(def.Tasks) == 0 {
		t.Fatal("no tasks were converted")
	}
	converted, err := taskrun.New(def)
	if err != nil {
		t.Fatalf("the converted definition did not validate: %v", err)
	}

	var requiresArgs int
	// Kite injects its runtime variables (including $@, $* and $1..$N) on every
	// run, so planning uses the same view.
	runtimeVars := b.runtimeVars(&kscript.RunCtx{}, nil, root)
	for name := range def.Tasks {
		plan, err := converted.Inspect(context.Background(), taskrun.Request{Task: name, Vars: runtimeVars})
		if err != nil {
			if errors.Is(err, taskrun.ErrInvalidRequest) {
				// A task whose arguments are required can only be planned with
				// them; that is correct behavior, not a conversion problem.
				requiresArgs++
				continue
			}
			t.Fatalf("planning task %q failed: %v", name, err)
		}
		if len(plan.Actions) == 0 {
			t.Fatalf("task %q planned no action", name)
		}
	}
	t.Logf("converted %d tasks (%d need arguments), %d warnings: %v",
		len(def.Tasks), requiresArgs, len(warnings), warnings)
}

// TestBridgeHandlesShippedConfigExpressions documents how the shipped commands
// are translated: $@ becomes the injected runtime variable and $? is left to
// the shell, exactly as the legacy renderer did.
func TestBridgeHandlesShippedConfigExpressions(t *testing.T) {
	root := repoRoot(t)
	runner := kscript.NewRunner()
	runner.PathResolver = repositoryResolver(root)
	runner.Scripts = map[string]any{
		"co": "git checkout $@",
		"br": "git branch $?",
		"st": "git status",
	}
	b := New(runner, WithBaseDir(root))
	var out strings.Builder
	b.io.Stdout = &out

	found, err := b.TryRun(context.Background(), "st", nil, &kscript.RunCtx{Silent: true, DryRun: true})
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if !strings.Contains(out.String(), "git") {
		t.Fatalf("dry run did not describe the command: %q", out.String())
	}

	def, _, err := b.Definition("", b.runtimeVarNames(map[string]any{}), nil)
	if err != nil {
		t.Fatal(err)
	}
	args := def.Tasks["co"].Steps[0].Exec.Args
	if strings.Join(args, " ") != "checkout ${vars.@}" {
		t.Fatalf("co args=%v", args)
	}
	brArgs := def.Tasks["br"].Steps[0].Exec.Args
	if strings.Join(brArgs, " ") != "branch $?" {
		t.Fatalf("br args=%v (the legacy renderer left $? untouched)", brArgs)
	}
}

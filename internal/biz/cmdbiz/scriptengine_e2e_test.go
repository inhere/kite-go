package cmdbiz

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gookit/config/v2"

	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/kscript"
)

// legacyTaskFile returns a task file in the shape Kite loads through
// DefineFiles: top level keys are task names plus the __settings entry.
func legacyTaskFile() string {
	shell := func(line string) string {
		if runtime.GOOS == "windows" {
			return `cmd /C "` + line + `"`
		}
		return `sh -c "` + line + `"`
	}
	return "__settings:\n" +
		"  vars:\n" +
		"    project: kite\n" +
		"  env:\n" +
		"    KS_PROJECT: kite\n" +
		"mark:\n" +
		"  desc: write a marker line\n" +
		"  run: '" + shell("echo mark >> out.txt") + "'\n" +
		"build:\n" +
		"  desc: build the project\n" +
		"  deps: [mark]\n" +
		"  vars:\n" +
		"    target: all\n" +
		"  run:\n" +
		"    - '" + shell("echo build-${target} >> out.txt") + "'\n" +
		"    - '" + shell("echo setting-${vars.project} >> out.txt") + "'\n" +
		"    - '@task:mark'\n"
}

func scriptFileName() (string, string) {
	if runtime.GOOS == "windows" {
		return "hello.bat", "@echo file-bat >> out.txt\r\n"
	}
	return "hello.sh", "#!/bin/sh\necho file-sh >> out.txt\n"
}

// newConfigRunner builds a runner from a real task file and script directory,
// the way boot configures it from DefineFiles and ScriptDirs.
func newConfigRunner(t *testing.T, engine string) (string, *kscript.Runner) {
	t.Helper()
	dir := t.TempDir()
	taskFile := filepath.Join(dir, "tasks.yml")
	if err := os.WriteFile(taskFile, []byte(legacyTaskFile()), 0o644); err != nil {
		t.Fatal(err)
	}
	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	file, body := scriptFileName()
	if err := os.WriteFile(filepath.Join(scriptsDir, file), []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	ext := filepath.Ext(file)
	bin := strings.TrimPrefix(ext, ".")
	if ext == ".bat" {
		bin = "cmd /C"
	}

	cfg := config.New("test")
	cfg.Set("script_engine", engine)
	app.Set(app.ObjConf, cfg)

	runner := kscript.NewRunner()
	runner.PathResolver = func(path string) string {
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(dir, path)
	}
	runner.DefineFiles = []string{taskFile}
	runner.ScriptDirs = []string{scriptsDir}
	runner.AllowedExt = []string{ext}
	runner.ExtToBinMap = map[string]string{ext: bin}
	app.Scripts = runner
	if err := runner.InitLoad(); err != nil {
		t.Fatalf("InitLoad: %v", err)
	}
	return dir, runner
}

func captureStdout(fn func()) {
	saved := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		fn()
		return
	}
	os.Stdout = writer
	drained := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, reader)
		close(drained)
	}()
	fn()
	_ = writer.Close()
	os.Stdout = saved
	<-drained
	_ = reader.Close()
}

func runTaskWithEngine(t *testing.T, engine, name string, args []string) string {
	t.Helper()
	dir, _ := newConfigRunner(t, engine)
	var runErr error
	captureStdout(func() {
		found, err := RunScriptName(name, args, &kscript.RunCtx{Silent: true, Workdir: dir})
		if err != nil {
			runErr = err
			return
		}
		if !found {
			t.Errorf("engine %s did not find %q", engine, name)
		}
	})
	if runErr != nil {
		t.Fatalf("engine %s: %v", engine, runErr)
	}
	data, err := os.ReadFile(filepath.Join(dir, "out.txt"))
	if err != nil {
		t.Fatalf("engine %s produced no output file: %v", engine, err)
	}
	return string(data)
}

// Both engines must produce the same observable result for the same task file
// loaded through DefineFiles.
func TestEnginesAgreeOnRealConfigFile(t *testing.T) {
	legacyOut := runTaskWithEngine(t, ScriptEngineLegacy, "build", nil)
	kscriptOut := runTaskWithEngine(t, ScriptEngineTaskrun, "build", nil)
	if legacyOut != kscriptOut {
		t.Fatalf("engines disagree:\nlegacy:\n%s\nkscript:\n%s", legacyOut, kscriptOut)
	}
	for _, want := range []string{"mark", "build-all", "setting-kite"} {
		if !strings.Contains(legacyOut, want) {
			t.Fatalf("output missing %q:\n%s", want, legacyOut)
		}
	}
	// The @task:mark reference runs the dependency a second time.
	if strings.Count(legacyOut, "mark") != 2 {
		t.Fatalf("task reference was not executed twice:\n%s", legacyOut)
	}
}

// Script files discovered from ScriptDirs are runnable by name through both
// engines.
func TestEnginesAgreeOnScriptFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		// The legacy engine runs BinName as the program and passes the file as
		// its first argument, so a Windows script file would need an
		// interpreter that works that way ("cmd" without /C opens an
		// interactive shell instead). The new engine supports structured
		// interpreter prefix arguments, which is covered by the bridge tests.
		t.Skip("legacy script-file execution on Windows has no usable interpreter form")
	}
	file, _ := scriptFileName()
	name := strings.TrimSuffix(file, filepath.Ext(file))
	legacyOut := runTaskWithEngine(t, ScriptEngineLegacy, name, nil)
	kscriptOut := runTaskWithEngine(t, ScriptEngineTaskrun, name, nil)
	if legacyOut != kscriptOut {
		t.Fatalf("engines disagree on a script file:\nlegacy:\n%s\nkscript:\n%s", legacyOut, kscriptOut)
	}
	if !strings.Contains(legacyOut, "file-") {
		t.Fatalf("script file produced no marker:\n%s", legacyOut)
	}
}

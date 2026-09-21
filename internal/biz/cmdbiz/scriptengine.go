package cmdbiz

import (
	"context"
	"os"
	"sync"

	"github.com/gookit/taskrun"

	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/kscript"
	"github.com/inhere/kite-go/pkg/kscript/bridge"
)

// Script engine selection.
//
// The legacy runner in pkg/kscript remains the default so this migration is
// reversible: set `script_engine: taskrun` in the config to run script tasks
// through the standalone github.com/gookit/taskrun library. Command aliases,
// extensions, plugins and system command fallback are unaffected either way and
// stay owned by Kite.
const (
	// ScriptEngineLegacy runs script tasks with the in-tree runner.
	ScriptEngineLegacy = "legacy"
	// ScriptEngineTaskrun runs script tasks through the standalone library.
	ScriptEngineTaskrun = "taskrun"
	// configScriptEngine is the configuration key that selects the engine.
	configScriptEngine = "script_engine"
)

var (
	scriptBridgeMu  sync.Mutex
	scriptBridge    *bridge.Bridge
	scriptBridgeFor *kscript.Runner
)

// ScriptEngine reports the configured script engine name.
func ScriptEngine() string {
	name := app.Cfg().String(configScriptEngine, ScriptEngineLegacy)
	switch name {
	case ScriptEngineLegacy, ScriptEngineTaskrun:
		return name
	default:
		return ScriptEngineLegacy
	}
}

// scriptBridgeInstance returns the migration bridge for the current script
// runner, rebuilding it when the runner is replaced. Conversion failures fall
// back to the legacy runner, so an unconvertible script never regresses.
func scriptBridgeInstance() *bridge.Bridge {
	scriptBridgeMu.Lock()
	defer scriptBridgeMu.Unlock()
	if scriptBridge != nil && scriptBridgeFor == app.Scripts {
		return scriptBridge
	}
	scriptBridge = bridge.New(app.Scripts,
		bridge.WithLegacyFallback(true),
		bridge.WithIO(taskrun.IO{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}),
	)
	scriptBridgeFor = app.Scripts
	return scriptBridge
}

// RunScriptName runs a script task or script file by name using the configured
// engine. found is false only when no script matched, so the caller can
// continue with system command fallback.
func RunScriptName(name string, args []string, ctx *kscript.RunCtx) (found bool, err error) {
	switch ScriptEngine() {
	case ScriptEngineTaskrun:
		return scriptBridgeInstance().TryRun(context.Background(), name, args, ctx)
	default:
		return app.Scripts.TryRun(name, args, ctx)
	}
}

// RunScriptOnly runs a script by name without any fallback, using the
// configured engine. It backs the explicit `kite run --type=script` path.
func RunScriptOnly(name string, args []string, ctx *kscript.RunCtx) error {
	switch ScriptEngine() {
	case ScriptEngineTaskrun:
		return scriptBridgeInstance().Run(context.Background(), name, args, ctx)
	default:
		return app.Scripts.Run(name, args, ctx)
	}
}

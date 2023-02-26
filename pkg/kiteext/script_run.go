package kiteext

import (
	"os"

	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/cliutil/cmdline"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
)

// ScriptItem struct
type ScriptItem struct {
	// Workdir for run script
	Workdir string

	// Name for script
	Name string
	// Desc message
	Desc string
	// Env setting for run
	Env map[string]string
	// Args script args definition.
	Args, ArgNames []string
	// Cmds commands define in ScriptRunner.DefineFiles
	Cmds []string

	// File script file path in ScriptRunner.ScriptDirs
	File string
	Bin  string
	Ext  string // eg: .go
}

// args type: string, strings
func (si *ScriptItem) loadArgsDefine(args any) error {
	if args == nil {
		return nil
	}

	switch typVal := args.(type) {
	case string: // desc
		si.Args = []string{typVal}
	case []string: // desc list
		si.Args = typVal
	case []any: // desc list
		si.Args = arrutil.SliceToStrings(typVal)
	// case map[string]string: // name with desc TODO map cannot be ordered
	// 	si.Args = typVal
	default:
		return errorx.Rawf("invalid args config for %q", si.Name)
	}
	return nil
}

// ScriptRunner struct
type ScriptRunner struct {
	// Scripts config and loaded from DefineFiles.
	//
	// format: {name: info, name2: info2, ...}
	Scripts map[string]any `json:"scripts"`
	// DefineFiles scripts define files, will read and add to Scripts
	DefineFiles []string `json:"define_files"`

	// ScriptDirs script file dirs, allow multi
	ScriptDirs []string `json:"script_dirs"`
	// AllowedExt allowed script file extensions. eg: .go, .sh
	AllowedExt []string `json:"allowed_ext"`
	// FindBinByExt on run a script file
	FindBinByExt bool `json:"find_bin_by_ext"`
	// ExtToBinMap settings
	ExtToBinMap map[string]string `json:"ext_to_bin_map"`
	// PathResolver handler
	PathResolver func(path string) string

	// loaded from ScriptDirs.
	//
	// format: {filename: filepath, ...}
	scriptFiles map[string]string
	// mark script loaded
	fileLoaded   bool
	defineLoaded bool
}

// AllowExt list
var AllowExt = []string{".sh", ".zsh", ".bash", ".php", ".go", ".gop", ".kts", ".java", ".gry", ".groovy"}

// ExtToBinMap data
//
// eg:
//
//	'#!/usr/bin/env bash'
//	'#!/usr/bin/env -S go run'
var ExtToBinMap = map[string]string{
	".sh":     "sh",
	".zsh":    "zsh",
	".bash":   "bash",
	".php":    "php",
	".gry":    "groovy",
	".groovy": "groovy",
	".go":     "go run",
}

// NewScriptRunner instance
func NewScriptRunner(fns ...func(sr *ScriptRunner)) *ScriptRunner {
	sr := &ScriptRunner{
		AllowedExt:   AllowExt,
		ExtToBinMap:  ExtToBinMap,
		PathResolver: sysutil.ExpandPath,
		scriptFiles:  map[string]string{},
	}

	for _, fn := range fns {
		fn(sr)
	}
	return sr
}

// InitLoad define scripts and script files.
func (sr *ScriptRunner) InitLoad() error {
	if err := sr.LoadDefineScripts(); err != nil {
		return err
	}
	return sr.LoadScriptFiles()
}

// LoadDefineScripts from DefineFiles
func (sr *ScriptRunner) LoadDefineScripts() error {
	if sr.defineLoaded {
		return nil
	}

	sr.defineLoaded = true
	for _, fpath := range sr.DefineFiles {
		fpath = sr.PathResolver(fpath)

		var sub map[string]any
		err := yaml.Decoder(fsutil.ReadAll(fpath), &sub)
		if err != nil {
			return err
		}

		sr.Scripts = maputil.SimpleMerge(sub, sr.Scripts)
	}

	return nil
}

// LoadScriptFiles from the ScriptDirs
func (sr *ScriptRunner) LoadScriptFiles() error {
	if sr.fileLoaded {
		return nil
	}
	sr.fileLoaded = true

	for _, dirPath := range sr.ScriptDirs {
		dirPath = sr.PathResolver(dirPath)
		des, err := os.ReadDir(dirPath)
		if err != nil {
			return err
		}

		for _, ent := range des {
			fName := ent.Name()
			if !ent.IsDir() {
				// add
				sr.scriptFiles[fName] = dirPath + "/" + fName
			}
		}
	}

	return nil
}

// Run script or script-file by name and with args
func (sr *ScriptRunner) Run(name string, args []string, beforeFn func()) error {
	found, err := sr.TryRun(name, args, beforeFn)
	if !found {
		return errorx.Rawf("script file %q is not exists", name)
	}
	return err
}

// TryRun script or script-file by name and with args
// TODO add param ctx map[string]string
func (sr *ScriptRunner) TryRun(name string, args []string, beforeFn func()) (found bool, err error) {
	if err := sr.InitLoad(); err != nil {
		return false, err
	}

	found = true

	// try check is script and run it.
	sItem, err := sr.ScriptItem(name)
	if err != nil {
		return found, err
	}
	if sItem != nil {
		return found, sr.doExecScript(sItem, name, args, beforeFn)
	}

	// try check and run script file.
	sItem, err = sr.ScriptFileItem(name)
	if err != nil {
		return found, err
	}

	if sItem != nil {
		return found, sr.doExecScript(sItem, name, args, beforeFn)
	}
	return false, nil
}

// RunDefinedScript by input name and with arguments
func (sr *ScriptRunner) RunDefinedScript(name string, args []string) error {
	if err := sr.InitLoad(); err != nil {
		return err
	}

	sItem, err := sr.ScriptItem(name)
	if err != nil {
		return err
	}

	if sItem != nil {
		return sr.doExecScript(sItem, name, args, nil)
	}
	return errorx.Rawf("script %q is not exists", name)
}

// RunScriptFile by input name and with arguments
func (sr *ScriptRunner) RunScriptFile(name string, args []string) error {
	if err := sr.InitLoad(); err != nil {
		return err
	}

	sItem, err := sr.ScriptFileItem(name)
	if err != nil {
		return err
	}

	if sItem != nil {
		return sr.doExecScript(sItem, name, args, nil)
	}
	return errorx.Rawf("script file %q is not exists", name)
}

func (sr *ScriptRunner) doExecScript(item *ScriptItem, name string, inArgs []string, beforeFn func()) error {
	if beforeFn != nil {
		beforeFn()
	}

	// run script
	if item.File == "" {
		ln := len(item.Cmds)
		if ln == 0 {
			return errorx.Rawf("empty cmd config for script %q", name)
		}

		// only one
		if ln == 1 {
			bin, args := cmdline.
				NewParser(item.Cmds[0]).
				WithParseEnv().
				BinAndArgs()

			return cmdr.NewCmd(bin, args...).AddArgs(inArgs).FlushRun()
		}

		cr := cmdr.NewRunner()
		for _, cmd := range item.Cmds {
			cr.AddCmdline(cmd)
		}
		return cr.Run()
	}

	// TODO
	return errorx.Rawf("TODO exec %q", name)
}

// DefinedScript info get
func (sr *ScriptRunner) DefinedScript(name string) (any, bool) {
	info, ok := sr.Scripts[name]
	return info, ok
}

// ScriptItem get script info as ScriptItem
func (sr *ScriptRunner) ScriptItem(name string) (*ScriptItem, error) {
	info, ok := sr.Scripts[name]
	if !ok {
		// return nil, errorx.Rawf("script %q is not exists", name)
		return nil, nil
	}

	return sr.newDefinedScriptItem(name, info)
}

func (sr *ScriptRunner) newDefinedScriptItem(name string, info any) (*ScriptItem, error) {
	si := &ScriptItem{Name: name}

	switch typVal := info.(type) {
	case string: // on command
		si.Cmds = []string{typVal}
	case []any: // as commands
		si.Cmds = arrutil.SliceToStrings(typVal)
	case map[string]any: // as commands
		data := maputil.Data(typVal)
		si.Desc = data.Str("desc")
		err := si.loadArgsDefine(data.Get("args"))
		if err != nil {
			return nil, err
		}
		si.Cmds = data.Strings("cmds")
	default:
		return nil, errorx.Rawf("invalid config of the script %q", name)
	}
	return si, nil
}

// ScriptFileItem info get
func (sr *ScriptRunner) ScriptFileItem(name string) (*ScriptItem, error) {
	if inExt := fsutil.FileExt(name); len(inExt) > 0 {
		fpath, ok := sr.scriptFiles[name]
		if !ok {
			return nil, nil
		}

		return newFileScriptItem(name, fpath, inExt)
	}

	for _, ext := range sr.AllowedExt {
		fpath, ok := sr.scriptFiles[name+ext]
		if !ok {
			return nil, nil
		}

		return newFileScriptItem(name, fpath, ext)
	}

	return nil, nil
}

func newFileScriptItem(name, fpath, ext string) (*ScriptItem, error) {
	return nil, nil
}

// IsDefinedScript name
func (sr *ScriptRunner) IsDefinedScript(name string) bool {
	_, ok := sr.Scripts[name]
	return ok
}

// DefinedScripts map
func (sr *ScriptRunner) DefinedScripts() map[string]any {
	return sr.Scripts
}

// ScriptFiles file map
func (sr *ScriptRunner) ScriptFiles() map[string]string {
	return sr.scriptFiles
}

func (sr *ScriptRunner) WithConfigFn(fn func(sr *ScriptRunner)) *ScriptRunner {
	fn(sr)
	return sr
}

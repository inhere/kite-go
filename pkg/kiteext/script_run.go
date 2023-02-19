package kiteext

// ScriptRunner struct
type ScriptRunner struct {
	PathResolver func(path string) string

	ScriptDirs  []string `json:"script_dirs"`
	DefineFiles []string `json:"define_files"`

	defineLoad, dirFileLoad bool

	// loaded from ScriptDirs
	scriptFiles map[string]string
	// loaded from DefineFiles
	scriptMap map[string]any
}

// DefineScripts map
func (sr *ScriptRunner) DefineScripts() {

}

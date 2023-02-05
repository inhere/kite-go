package kiteext

// ScriptRunner struct
type ScriptRunner struct {
	PathResolver func(path string) string

	ScriptDirs  []string `json:"script_dirs"`
	DefineFiles []string `json:"define_files"`

	// loaded from ScriptDirs
	scriptFiles map[string]string
	// loaded from DefineFiles
	scriptMap map[string]string
}

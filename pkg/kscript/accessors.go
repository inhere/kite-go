package kscript

// This file adds read-only accessors used by the migration bridge in
// pkg/kscript/bridge. The legacy implementation stays in place as the rollback
// point: nothing here changes existing behavior.

// SettingsData returns the loaded task settings in map form, so the migration
// bridge can pass them to the new library converter. InitLoad must have run.
func (r *Runner) SettingsData() map[string]any {
	data := map[string]any{}
	if len(r.taskSettings.Vars) > 0 {
		vars := make(map[string]any, len(r.taskSettings.Vars))
		for key, value := range r.taskSettings.Vars {
			vars[key] = value
		}
		data["vars"] = vars
	}
	if len(r.taskSettings.Env) > 0 {
		env := make(map[string]any, len(r.taskSettings.Env))
		for key, value := range r.taskSettings.Env {
			env[key] = value
		}
		data["env"] = env
	}
	if len(r.taskSettings.EnvPaths) > 0 {
		data["env_paths"] = append([]string(nil), r.taskSettings.EnvPaths...)
	}
	if len(r.taskSettings.Groups) > 0 {
		groups := make(map[string]any, len(r.taskSettings.Groups))
		for name, values := range r.taskSettings.Groups {
			inner := make(map[string]any, len(values))
			for key, value := range values {
				inner[key] = value
			}
			groups[name] = inner
		}
		data["groups"] = groups
	}
	if r.taskSettings.DefaultGroup != "" {
		data["default_group"] = r.taskSettings.DefaultGroup
	}
	return data
}

// ScriptFileMap returns the discovered script files as name to path, so the
// migration bridge can register them as structured File entries.
func (r *Runner) ScriptFileMap() map[string]string {
	out := make(map[string]string, len(r.scriptFiles))
	for name, path := range r.scriptFiles {
		out[name] = path
	}
	return out
}

// ExtToBin returns the configured extension to interpreter mapping.
func (r *Runner) ExtToBin() map[string]string {
	out := make(map[string]string, len(r.ExtToBinMap))
	for ext, bin := range r.ExtToBinMap {
		out[ext] = bin
	}
	return out
}

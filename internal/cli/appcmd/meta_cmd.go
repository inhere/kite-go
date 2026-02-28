package appcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/goccy/go-yaml"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/inhere/kite-go/internal/app"
)

var cmOpts = struct {
	format   cflag.EnumString
	strukt   cflag.EnumString
	output   string
	withFlag bool
}{
	format: cflag.NewEnumString("json", "yaml"),
	strukt: cflag.NewEnumString("both", "flat", "tree"),
}

// CommandMapCmd export all console commands info to JSON/YAML for smart search and matching
var CommandMapCmd = &gcli.Command{
	Name:    "cmd-map",
	Aliases: []string{"cmdmap"},
	Desc:    "export all console commands to JSON/YAML for smart search and AI matching",
	Examples: `
# export both flat and tree to current dir
{$fullCmd} -o .

# export only flat structure to stdout
{$fullCmd} --struct flat

# export tree as YAML to file
{$fullCmd} --struct tree --format yaml -o /tmp/commands.tree.yaml

# export with flags/args info
{$fullCmd} --flags -o .
`,
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&cmOpts.withFlag, "flags,flag", "also export command flags(options) and arguments")
		c.VarOpt(&cmOpts.format, "format", "fmt", "output format, allow: json, yaml", gflag.WithDefault("yaml"))
		c.VarOpt(&cmOpts.strukt, "struct", "s", "export structure type, allow: flat, tree, both", gflag.WithDefault("tree"))
		c.StrOpt2(&cmOpts.output, "output,o", "output file path or directory when --struct=both. default is stdout", gflag.WithDefault(""))
	},
	Func: func(c *gcli.Command, _ []string) error {
		cliApp := app.Cli
		format := cmOpts.format.String()
		strukt := cmOpts.strukt.String()
		outPath := cmOpts.output

		switch strukt {
		case "flat":
			items := buildFlatCmds(cliApp, cmOpts.withFlag)
			return writeExportData(items, format, outPath)

		case "tree":
			items := buildTreeCmds(cliApp, cmOpts.withFlag)
			return writeExportData(items, format, outPath)

		default: // both
			flat := buildFlatCmds(cliApp, cmOpts.withFlag)
			tree := buildTreeCmds(cliApp, cmOpts.withFlag)

			if outPath == "" {
				// stdout: wrap both in a single object
				return writeExportData(map[string]any{
					"flat": flat,
					"tree": tree,
				}, format, "")
			}

			// write to separate files in the output directory
			ext := fmtExt(format)
			flatFile := filepath.Join(outPath, "commands.flat"+ext)
			treeFile := filepath.Join(outPath, "commands.tree"+ext)

			if err := writeExportData(flat, format, flatFile); err != nil {
				return err
			}
			c.Successf("Written flat: %s\n", flatFile)

			if err := writeExportData(tree, format, treeFile); err != nil {
				return err
			}
			c.Successf("Written tree: %s\n", treeFile)
		}
		return nil
	},
}

// writeExportData marshal data and write to file or stdout
func writeExportData(data any, format, filePath string) error {
	var (
		bs  []byte
		err error
	)

	switch format {
	case "yaml":
		bs, err = yaml.Marshal(data)
	default: // json
		bs, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if filePath == "" {
		_, err = os.Stdout.Write(bs)
		return err
	}

	// ensure parent directory exists
	if dir := filepath.Dir(filePath); dir != "." {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create dir failed: %w", err)
		}
	}
	return os.WriteFile(filePath, bs, 0644)
}

// fmtExt returns file extension for the given format
func fmtExt(format string) string {
	if format == "yaml" || format == "yml" {
		return ".yaml"
	}
	return ".json"
}

// CmdFlagInfo flag/option info for export
type CmdFlagInfo struct {
	Name     string   `json:"name" yaml:"name"`
	Shorts   []string `json:"shorts,omitempty" yaml:"shorts,omitempty"`
	Desc     string   `json:"desc" yaml:"desc"`
	Default  any      `json:"default,omitempty" yaml:"default,omitempty"`
	Required bool     `json:"required,omitempty" yaml:"required,omitempty"`
	Hidden   bool     `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

// CmdArgInfo argument info for export
type CmdArgInfo struct {
	Name     string `json:"name" yaml:"name"`
	Desc     string `json:"desc" yaml:"desc"`
	Required bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Arrayed  bool   `json:"arrayed,omitempty" yaml:"arrayed,omitempty"`
}

// CmdFlatInfo flat command info, suitable for search/AI matching
type CmdFlatInfo struct {
	// Path is the full command path, e.g. "fs find"
	Path string `json:"path" yaml:"path"`
	// Name is the command name, e.g. "find"
	Name string `json:"name" yaml:"name"`
	// Group is the parent command name, e.g. "fs"; empty for top-level
	Group   string        `json:"group,omitempty" yaml:"group,omitempty"`
	Desc    string        `json:"desc" yaml:"desc"`
	Aliases []string      `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Hidden  bool          `json:"hidden,omitempty" yaml:"hidden,omitempty"`
	Flags   []CmdFlagInfo `json:"flags,omitempty" yaml:"flags,omitempty"`
	Args    []CmdArgInfo  `json:"args,omitempty" yaml:"args,omitempty"`
}

// CmdTreeInfo tree command info, preserves hierarchy
type CmdTreeInfo struct {
	Name    string         `json:"name" yaml:"name"`
	Desc    string         `json:"desc" yaml:"desc"`
	Aliases []string       `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Hidden  bool           `json:"hidden,omitempty" yaml:"hidden,omitempty"`
	Flags   []CmdFlagInfo  `json:"flags,omitempty" yaml:"flags,omitempty"`
	Args    []CmdArgInfo   `json:"args,omitempty" yaml:"args,omitempty"`
	Subs    []*CmdTreeInfo `json:"subs,omitempty" yaml:"subs,omitempty"`
}

// buildFlatCmds traverse all commands and return a flat list
func buildFlatCmds(cliApp *gcli.App, withFlags bool) []CmdFlatInfo {
	var list []CmdFlatInfo
	collectFlat(cliApp.Commands(), "", &list, withFlags)
	return list
}

func collectFlat(cmds map[string]*gcli.Command, parentPath string, list *[]CmdFlatInfo, withFlags bool) {
	for _, name := range sortedCmdKeys(cmds) {
		cmd := cmds[name]
		// skip alias entries, only process each command once by canonical name
		if cmd.Name != name {
			continue
		}

		var path string
		if parentPath == "" {
			path = name
		} else {
			path = parentPath + " " + name
		}

		// init to trigger Config callback and register flags/args
		cmd.Init()

		item := CmdFlatInfo{
			Path:    path,
			Name:    cmd.Name,
			Group:   parentPath,
			Desc:    cmd.Desc,
			Aliases: []string(cmd.Aliases),
			Hidden:  cmd.Hidden,
		}
		if withFlags {
			item.Flags = extractFlags(cmd)
			item.Args = extractArgs(cmd)
		}
		*list = append(*list, item)

		// recurse into sub-commands
		if subs := cmd.Commands(); len(subs) > 0 {
			collectFlat(subs, path, list, withFlags)
		}
	}
}

// buildTreeCmds build tree structure from all commands
func buildTreeCmds(cliApp *gcli.App, withFlags bool) []*CmdTreeInfo {
	return collectTree(cliApp.Commands(), withFlags)
}

func collectTree(cmds map[string]*gcli.Command, withFlags bool) []*CmdTreeInfo {
	var result []*CmdTreeInfo
	for _, name := range sortedCmdKeys(cmds) {
		cmd := cmds[name]
		if cmd.Name != name {
			continue
		}

		cmd.Init()

		node := &CmdTreeInfo{
			Name:    cmd.Name,
			Desc:    cmd.Desc,
			Aliases: []string(cmd.Aliases),
			Hidden:  cmd.Hidden,
		}
		if withFlags {
			node.Flags = extractFlags(cmd)
			node.Args = extractArgs(cmd)
		}
		if subs := cmd.Commands(); len(subs) > 0 {
			node.Subs = collectTree(subs, withFlags)
		}
		result = append(result, node)
	}
	return result
}

// extractFlags collect flag/option definitions from an initialized command
func extractFlags(cmd *gcli.Command) []CmdFlagInfo {
	opts := cmd.Opts()
	if len(opts) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var flags []CmdFlagInfo
	for _, opt := range opts {
		if seen[opt.Name] {
			continue
		}
		seen[opt.Name] = true
		flags = append(flags, CmdFlagInfo{
			Name:     opt.Name,
			Shorts:   opt.Shorts,
			Desc:     opt.Desc,
			Default:  opt.DefVal,
			Required: opt.Required,
			Hidden:   opt.Hidden,
		})
	}
	sort.Slice(flags, func(i, j int) bool { return flags[i].Name < flags[j].Name })
	return flags
}

// extractArgs collect argument definitions from an initialized command
func extractArgs(cmd *gcli.Command) []CmdArgInfo {
	args := cmd.Args()
	if len(args) == 0 {
		return nil
	}

	result := make([]CmdArgInfo, 0, len(args))
	for _, arg := range args {
		result = append(result, CmdArgInfo{
			Name:     arg.Name,
			Desc:     arg.Desc,
			Required: arg.Required,
			Arrayed:  arg.Arrayed,
		})
	}
	return result
}

func sortedCmdKeys(m map[string]*gcli.Command) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

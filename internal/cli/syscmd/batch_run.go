package syscmd

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

var btrOpts = struct {
	cmdbiz.CommonOpts
	cmdTpl  string
	inDirs  gflag.String
	allSub  bool
	exclude gflag.Strings
	// vars for command template
	cmdVars gflag.KVString
	// for range vars list, multi by comma
	forVars gflag.String
}{
	cmdVars: cflag.NewKVString(),
}

// NewBatchRunCmd instance
func NewBatchRunCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "brun",
		Aliases: []string{"batch-run"},
		Desc:    "batch run more commands at once, and allow with vars",
		Help: `
Build-in vars:
	{dir}  		  - current dir name
	{path} 		  - current dir full path
	{pwd}  		  - current workdir path
	{parent}  	  - parent dir name
	{parentPath}  - parent dir full path
`,
		Examples: `
# run command multi times in all subdir
{$fullCmd} --cmd "echo {dir}" --subdir
`,
		Config: func(c *gcli.Command) {
			btrOpts.BindCommonFlags(c)

			c.BoolOpt2(&btrOpts.allSub, "all-subdir, all-sub, subdir", "run command on the each subdir in the <cyan>--dirs</>")
			c.StrOpt2(&btrOpts.cmdTpl, "cmd, c", "want execute `command` line or template, allow use vars. eg: {dir}")
			c.VarOpt(&btrOpts.exclude, "exclude", "e", "exclude some subdir on with <cyan>--all-subdir</>")
			c.VarOpt(&btrOpts.inDirs, "dirs", "d", "run command on the each dir path, multi by comma. default is workdir")
			c.VarOpt2(&btrOpts.cmdVars, "vars,var,v", "sets template variables for render. format: `KEY=VALUE`")
			c.VarOpt2(&btrOpts.forVars, "foreach,range,for", "for range vars list, multi by comma. use: {item}")

			c.AddArg("cmd", "same of option <cyan>--cmd</>, set execute command line template, allow vars").WithAfterFn(func(a *gflag.CliArg) error {
				if btrOpts.cmdTpl == "" {
					btrOpts.cmdTpl = a.String()
					return nil
				}
				return errorx.Raw("cmd template has been bounded by option --cmd")
			})
		},
		Func: func(c *gcli.Command, _ []string) error {
			wkDirs := btrOpts.inDirs.Strings()
			if len(wkDirs) == 0 {
				wkDirs = []string{c.WorkDir()}
			}

			btrOpts.cmdTpl = strings.TrimSpace(btrOpts.cmdTpl)
			if btrOpts.cmdTpl == "" {
				return errors.New("please input command template")
			}

			for _, dir := range wkDirs {
				if btrOpts.allSub {
					err := fsutil.FindInDir(dir, func(path string, ent fs.DirEntry) error {
						// check exclude
						if arrutil.StringsHas(btrOpts.exclude, ent.Name()) {
							return nil
						}
						return runCmdInDir(path, c)
					}, fsutil.OnlyFindDir)
					if err != nil {
						return err
					}
				} else {
					if err := runCmdInDir(dir, c); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}

func runCmdInDir(dirPath string, c *gcli.Command) error {
	pDir := fsutil.Dir(dirPath)
	vars := btrOpts.cmdVars.Data()
	vars.Load(map[string]string{
		"pwd":        c.WorkDir(),
		"path":       dirPath,
		"dir":        fsutil.Name(dirPath),
		"parent":     fsutil.Name(pDir),
		"parentPath": pDir,
	})

	// render command template
	rpl := textutil.NewVarReplacer("{,}")
	str := rpl.ReplaceSMap(btrOpts.cmdTpl, vars)
	ers := errorx.Errors{}

	execCmd := cmdr.NewCmdline(str).
		WithWorkDir(btrOpts.Workdir).
		WithDryRun(btrOpts.DryRun).
		OutputToOS().
		PrintCmdline()

	if err := execCmd.Run(); err != nil {
		ers = append(ers, err)
	}
	return ers.ErrorOrNil()
}

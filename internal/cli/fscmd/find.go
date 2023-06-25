package fscmd

import (
	"fmt"
	"os"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/fsutil/finder"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

var ffOpts = struct {
	cmdbiz.CommonOpts
	Dirs string `flag:"desc=the find directory, multi by comma;shorts=in,dir"`

	Type    string `flag:"desc=the find type, allow: f/file, d/dir, b/both;shorts=t;default=file"`
	Name    string `flag:"desc=include file/dir name, multi by comma;shorts=n"`
	NotName string `flag:"desc=exclude file/dir names, multi by comma;shorts=N,nn"`
	Path    string `flag:"desc=contains file/dir path, multi by comma;shorts=p"`
	NotPath string `flag:"desc=exclude contains file/dir paths, multi by comma;shorts=P,np"`
	Like    string `flag:"desc=match file/dir name like, multi by comma;shorts=l"`
	NotLike string `flag:"desc=exclude match file/dir name like, multi by comma;shorts=nl"`
	Ext     string `flag:"desc=match file ext, multi by comma. eg: .md,.txt;shorts=e"`
	NotExt  string `flag:"desc=exclude match file ext, multi by comma;shorts=E,ne"`

	User  string `flag:"desc=match file/dir owner user name;shorts=U"`
	Group string `flag:"desc=match file/dir owner group name;shorts=G"`
	Atime string `flag:"desc=match file access time, format: 5m, 1h, 1d, 1w(TODO);shorts=a,at"`
	Mtime string `flag:"desc=match file modified time, format: 5m, 1h, 1d, 1w;shorts=m,mt"`
	Depth int    `flag:"desc=the find depth. if eq 1 like ls command;shorts=D"`
	Size  string `flag:"desc=match file size range, format: 20k, <20m, 10k-1m, 1m-20m, <1g;shorts=s"`

	Exec   string `flag:"desc=execute command for each file/dir;shorts=x"`
	Delete bool   `flag:"desc=delete matched files or dirs;shorts=del,rm"`

	Verb  bool `flag:"desc=show verbose info;shorts=v"`
	Clear bool `flag:"desc=output clear find result;shorts=c"`

	// NotRecursive find subdir
	NotRecursive bool `flag:"desc=not recursive find subdir. equals <mga>--depth=1</>;shorts=nr"`
	WithDotDir   bool `flag:"desc=include dot directories, start with <mga>.</>;shorts=dd"`
	WithDotFile  bool `flag:"desc=include dot files, start with <mga>.</>;shorts=df"`

	// runtime vars
	dirs []string
}{}

// FileFindCmd command
var FileFindCmd = &gcli.Command{
	Name:    "find",
	Desc:    "find files by match name or pattern, and support match contents",
	Aliases: []string{"glob", "search", "match"},
	Help: `
<cyan>### Mtime, Atime format</>:
  10m/<10m                    last 10 minutes
  >10m                        before 10 minutes
  1h/1hour/<1hour             last 1 hour
  1d/<1d/<24h                 last 1 day(24h)
  today                       today(00:00-24:00)
  yesterday                   yesterday(00:00-24:00)
  >24h              		  before yesterday
  >2d                         before 2 days
  // time range limit
  1h~10m                      last 1 hour to 10 minutes
  1d~1h                       last 1 day to 1 hour
  5h~1h                       last 5 hours to 1 hour

<cyan>### Variables in exec</>:
 {path} 	the file/dir path
 {name} 	the file/dir path name
 {dir}   	the directory path for {path}
`,
	Examples: `
# find files and run command
{$fullCmd} -t file --name "*.go" -x "cat {file}" .

# find and delete files
{$fullCmd} -t file --name "test,doc,[t|T]ests,[d|D]ocs" --del .

# list sub dirs and run command
{$fullCmd} -t dir --nr -x "ls -l {dir}" .
`,
	Config: func(c *gcli.Command) {
		ffOpts.BindWorkdirDryRun(c)
		c.MustFromStruct(&ffOpts, gflag.TagRuleNamed)
		c.AddArg("dirs", "the find directory, multi by comma or multi input. same as <mga>--dirs</>").
			SetArrayed().
			WithAfterFn(func(a *gflag.CliArg) error {
				ffOpts.dirs = a.Strings()
				return nil
			})
	},
	Func: func(c *gcli.Command, _ []string) error {
		if ffOpts.Dirs != "" {
			ffOpts.dirs = strutil.SplitValid(ffOpts.Dirs, ",")
		}
		if len(ffOpts.dirs) == 0 {
			return fmt.Errorf("please input find directory")
		}

		ff := buildFinder()

		if ffOpts.Verb {
			show.AList("Configuration:", ff.Config())
			// s := progress.RoundTripLoading(progress.GetCharTheme(-1), 500*time.Millisecond)
			// s.Start("[%s] Finding ... ...")
		}

		if !ffOpts.Clear {
			colorp.Warnln("Finding and results:")
		}

		spl := textutil.NewVarReplacer("{,}")
		ers := errorx.Errors{}

		ff.EachElem(func(el finder.Elem) {
			elPath := el.Path()
			if ffOpts.Clear {
				fmt.Println(elPath)
				return
			}

			if ffOpts.Delete {
				colorp.Warnf("Delete path: %s\n", elPath)
				if ffOpts.DryRun {
					colorp.Infoln("Dry run, skip delete")
				} else if err := os.RemoveAll(elPath); err != nil {
					ers = append(ers, err)
				}
				return
			}

			if ffOpts.Exec == "" {
				fmt.Println(el)
				return
			}

			// exec command
			vs := map[string]string{
				"path": elPath,
				"name": el.Name(),
				"dir":  fsutil.Dir(elPath),
			}

			execCmd := cmdr.NewCmdline(spl.RenderSimple(ffOpts.Exec, vs)).
				WithWorkDir(ffOpts.Workdir).
				WithDryRun(ffOpts.DryRun).
				OutputToOS().
				PrintCmdline()

			if err := execCmd.Run(); err != nil {
				ers = append(ers, err)
			}
		})
		// s.Stop("Find complete")

		if ffOpts.Clear {
			return ff.Err()
		}

		if ff.Num() > 0 {
			colorp.Successf("Total found %d paths\n", ff.Num())
		} else {
			colorp.Infoln("... Not found any paths")
		}
		return ff.Err()
	},
}

func buildFinder() *finder.Finder {
	cfg := finder.NewConfig(ffOpts.dirs...)
	cfg.ExcludeDotDir = !ffOpts.WithDotDir
	cfg.ExcludeDotFile = !ffOpts.WithDotFile

	if ffOpts.NotRecursive {
		ffOpts.Depth = 1
	}

	ff := finder.NewWithConfig(cfg)
	ff.WithMaxDepth(ffOpts.Depth).
		WithStrFlag(ffOpts.Type).
		WithNames(strutil.Split(ffOpts.Name, ",")).
		WithoutNames(strutil.Split(ffOpts.NotName, ",")).
		WithPaths(strutil.Split(ffOpts.Path, ",")).
		WithoutPaths(strutil.Split(ffOpts.NotPath, ",")).
		WithExts(strutil.Split(ffOpts.Ext, ",")).
		WithoutExts(strutil.Split(ffOpts.NotExt, ","))

	if ffOpts.Size != "" {
		ff.MatchFile(finder.HumanSize(ffOpts.Size))
	}
	if ffOpts.Mtime != "" {
		ff.MatchFile(finder.HumanModTime(ffOpts.Mtime))
	}
	if ffOpts.Like != "" {
		ff.With(finder.NameLikes(strutil.Split(ffOpts.Like, ",")))
	}
	if ffOpts.NotLike != "" {
		ff.Not(finder.NameLikes(strutil.Split(ffOpts.NotLike, ",")))
	}

	return ff
}

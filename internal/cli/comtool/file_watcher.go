package comtool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
)

var watcher *fsnotify.Watcher
var fwOpts = struct {
	Dir   string
	Ext   string
	Files gcli.Strings

	Config  string
	Exclude gcli.Strings

	handler func(event fsnotify.Event)
}{}

// FileWatcher command definition
func FileWatcher(handler func(event fsnotify.Event)) *gcli.Command {
	cmd := &gcli.Command{
		Name: "watch",
		Func: watch,

		Desc: "file system change notification, by fsnotify",

		Aliases: []string{"fwatch", "fswatch"},
		Examples: `watch a dir:
  {$fullCmd} -e .git -e .idea -d ./_examples --ext ".go|.md"
  watch a file(s):
  {$fullCmd} -f _examples/cliapp.go -f app.go
  open debug mode:
  {$binName} --verbose 4 {$cmd} -e .git -e .idea -d ./_examples --ext ".go|.md"
`,
	}

	cmd.StrOpt(&fwOpts.Dir, "dir", "d", "", "the want watched directory")
	cmd.StrOpt(&fwOpts.Ext, "ext", "", ".go", "the watched file extensions, multi split by '|'")
	cmd.VarOpt(&fwOpts.Files, "files", "f", "the want watched file paths")
	cmd.StrOpt(&fwOpts.Config, "config", "c", "", "load options from a json config")
	cmd.VarOpt(&fwOpts.Exclude, "exclude", "e", "the ignored directory or files")

	fwOpts.handler = handler

	return cmd
}

// test run:
// go run ./_examples/cliapp.go watch -e .git -e .idea -d ./_examples
func watch(c *gcli.Command, _ []string) (err error) {
	color.Info.Println("Work directory: ", c.WorkDir())

	if fwOpts.Dir == "" && len(fwOpts.Files) == 0 {
		return c.Errorf("watched directory or files cannot be empty")
	}

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				gcli.Logf(gcli.VerbInfo, "event: %s", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					gcli.Logf(gcli.VerbDebug, "modified file: %s", event.Name)
				}

				if fwOpts.handler != nil {
					fwOpts.handler(event)
				}
			case err := <-watcher.Errors:
				gcli.Logf(gcli.VerbError, "error: %s", err.Error())
			}
		}
	}()

	if len(fwOpts.Files) > 0 {
		if err = addWatchFiles(fwOpts.Files); err != nil {
			// <-done
			return err
		}
	}

	if fwOpts.Dir != "" {
		fmt.Println("- add watch dir: ", color.FgGreen.Render(fwOpts.Dir))

		if err = addWatchDir(fwOpts.Dir); err != nil {
			return err
		}
	}

	<-done
	return
}

func addWatchFiles(files []string) error {
	for _, path := range files {
		gcli.Logf(gcli.VerbDebug, "add watch file: %s", path)
		err := watcher.Add(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func addWatchDir(dir string) error {
	allowExt := ""
	if fwOpts.Ext != "" {
		// always wrap char "|". eg ".go|.md" -> "|.go|.md|"
		allowExt = "|" + fwOpts.Ext + "|"
	}

	// filepath.Match()
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil { // continue
			return err
		}

		// get base name.
		// /path/dir -> dir
		// /path/file.ext -> file.ext
		name := filepath.Base(path)
		if isExclude(name) { // skip
			return nil
		}

		if info.IsDir() {
			err = watcher.Add(path)
			gcli.Logf(gcli.VerbDebug, "add watch dir: %s", path)
			return err // continue OR err
		}

		// has ext limit
		if allowExt != "" {
			// get ext. eg ".go"
			ext := filepath.Ext(path)
			if strings.Contains(allowExt, "|"+ext+"|") {
				// add file watch
				err = watcher.Add(path)
				gcli.Logf(gcli.VerbDebug, "add watch file: %s", path)
			}
		} else { // add any file
			err = watcher.Add(path)
			gcli.Logf(gcli.VerbDebug, "add watch file: %s", path)
		}

		return err
	})
}

func isExclude(name string) bool {
	if len(fwOpts.Exclude) == 0 {
		return false
	}

	for _, exclude := range fwOpts.Exclude {
		if exclude == name {
			return true
		}
	}

	return false
}

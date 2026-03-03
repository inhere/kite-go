package aicmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/cli/aicmd/skills"
)

// delete skill options
var delSkillOpts = struct {
	Force bool `flag:"shorts=f;desc=Skip confirmation"`
}{}

var SkillsCmd = &gcli.Command{
	Name:    "skills",
	Desc:    "Manage local AI skills for extending capabilities",
	Aliases: []string{"skill"},
	Help: `
Manage local AI skills that extend Claude's capabilities.

Skills are stored in:
  - User level: ~/.claude/skills/<skill-name>/SKILL.md
  - Project level: ./.claude/skills/<skill-name>/SKILL.md

Each skill is a directory containing a SKILL.md file with YAML frontmatter
and markdown instructions.
`,
	Subs: []*gcli.Command{
		SkillsListCmd(),
		SkillsShowCmd(),
		SkillsCreateCmd(),
		SkillsEditCmd(),
		SkillsDeleteCmd(),
		SkillsPathCmd(),
		SkillsOpenCmd(),
	},
	Func: func(c *gcli.Command, args []string) error {
		return c.ShowHelp()
	},
}

// SkillsListCmd lists all available skills
func SkillsListCmd() *gcli.Command {
	var opts = struct {
		Scope string `flag:"name=scope,s;desc=Filter by scope: user, project, all;default=all"`
	}{}

	return &gcli.Command{
		Name:    "list",
		Desc:    "List all available skills",
		Aliases: []string{"ls"},
		Config: func(c *gcli.Command) {
			c.StrOpt2(&opts.Scope, "scope,s", "Filter by scope: user, project, all")
		},
		Func: func(c *gcli.Command, args []string) error {
			mgr := skills.NewManager()
			skillList, err := mgr.ScanSkills(opts.Scope)
			if err != nil {
				return err
			}

			if len(skillList) == 0 {
				c.Infof("No skills found.\n")
				c.Infof("Create one with: kite ai skills create <name>\n")
				return nil
			}

			c.Infof("Found %d skill(s):\n\n", len(skillList))

			// Group by scope
			userSkills := make([]*skills.Skill, 0)
			projectSkills := make([]*skills.Skill, 0)

			for _, s := range skillList {
				if s.Scope == "user" {
					userSkills = append(userSkills, s)
				} else {
					projectSkills = append(projectSkills, s)
				}
			}

			// Print user skills
			if ln := len(userSkills); ln > 0 {
				c.Printf("<mga>User Skills(%d)</> (~/.claude/skills/)\n", ln)
				for _, s := range userSkills {
					c.Printf("  <info>%-36s</>", s.Name)
					if s.Description != "" {
						c.Printf(" - %s", strutil.TextTruncate(s.Description, 86, "..."))
					}
					c.Printf("\n")
				}
				c.Printf("\n")
			}

			// Print project skills
			if ln := len(projectSkills); ln > 0 {
				c.Printf("<mga>Project Skills(%d)</> (./.claude/skills/)\n", ln)
				for _, s := range projectSkills {
					c.Printf("  <info>%-36s</>", s.Name)
					if s.Description != "" {
						c.Printf(" - %s", strutil.TextTruncate(s.Description, 86, "..."))
					}
					c.Printf("\n")
				}
			}

			return nil
		},
	}
}

// SkillsShowCmd shows details of a specific skill
func SkillsShowCmd() *gcli.Command {
	return &gcli.Command{
		Name: "show",
		Desc: "Show details of a specific skill",
		Config: func(c *gcli.Command) {
			c.AddArg("name", "Name of the skill to show", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()

			mgr := skills.NewManager()
			skill, err := mgr.GetSkill(name)
			if err != nil {
				return err
			}

			c.Printf("<mga>Skill: %s</>\n\n", skill.Name)
			c.Printf("  <mga>Path:</> %s\n", skill.Path)
			c.Printf("  <mga>Scope:</> %s\n", skill.Scope)

			if skill.Description != "" {
				c.Printf("  <mga>Description:</> %s\n", skill.Description)
			}

			// Print frontmatter
			if len(skill.Frontmatter) > 0 {
				c.Printf("\n<mga>Frontmatter:</>\n")
				for k, v := range skill.Frontmatter {
					c.Infof("  %s: %v\n", k, v)
				}
			}

			// Print content preview
			if skill.Content != "" {
				c.Printf("\n<mga>Content:</>\n")
				content := skill.Content
				if len(content) > 500 {
					content = content[:500] + "...\n(truncated)"
				}
				c.Infof("%s\n", content)
			}

			return nil
		},
	}
}

// SkillsCreateCmd creates a new skill
func SkillsCreateCmd() *gcli.Command {
	var opts = struct {
		Description string `flag:"name=desc,d;desc=Description for the skill"`
		Scope       string `flag:"name=scope,s;desc=Create in scope: user or project;default=user"`
	}{}

	return &gcli.Command{
		Name:    "create",
		Desc:    "Create a new skill",
		Aliases: []string{"new", "add"},
		Config: func(c *gcli.Command) {
			c.StrOpt2(&opts.Description, "desc,d", "Description for the skill")
			c.StrOpt2(&opts.Scope, "scope,s", "Create in scope: user or project")
			c.AddArg("name", "Name of the skill to create", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()

			mgr := skills.NewManager()
			if err := mgr.CreateSkill(name, opts.Description, opts.Scope); err != nil {
				return err
			}

			var skillPath string
			if opts.Scope == "project" {
				skillPath = fmt.Sprintf("./.claude/skills/%s/SKILL.md", name)
			} else {
				skillPath = fmt.Sprintf("~/.claude/skills/%s/SKILL.md", name)
			}

			c.Successf("Created skill: %s\n", name)
			c.Infof("Edit it at: %s\n", skillPath)

			return nil
		},
	}
}

// SkillsEditCmd opens a skill in the editor
func SkillsEditCmd() *gcli.Command {
	return &gcli.Command{
		Name: "edit",
		Desc: "Open a skill in your default editor",
		Config: func(c *gcli.Command) {
			c.AddArg("name", "Name of the skill to edit", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()

			mgr := skills.NewManager()
			skill, err := mgr.GetSkill(name)
			if err != nil {
				return err
			}

			// Get editor from environment
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = os.Getenv("VISUAL")
			}
			if editor == "" {
				// Default editors based on platform
				if runtime.GOOS == "windows" {
					editor = "notepad"
				} else {
					editor = "vim"
				}
			}

			c.Infof("Opening %s with %s...\n", skill.Path, editor)

			cmd := exec.Command(editor, skill.Path)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			return cmd.Run()
		},
	}
}

// SkillsDeleteCmd deletes a skill
func SkillsDeleteCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "delete",
		Desc:    "Delete a skill",
		Aliases: []string{"rm", "remove"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&delSkillOpts)
			c.AddArg("name", "Name of the skill to delete", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()

			mgr := skills.NewManager()
			skill, err := mgr.GetSkill(name)
			if err != nil {
				return err
			}

			// Confirm deletion
			if !delSkillOpts.Force {
				c.Infof("Are you sure you want to delete skill %q? [y/N]: ", name)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					c.Infof("Cancelled.\n")
					return nil
				}
			}

			if err := mgr.DeleteSkill(name); err != nil {
				return err
			}

			c.Successf("Deleted skill: %s\n", name)
			c.Infof("Removed: %s\n", skill.Dir)

			return nil
		},
	}
}

// SkillsPathCmd shows the skills directory path
func SkillsPathCmd() *gcli.Command {
	var opts = struct {
		Scope string `flag:"name=scope,s;desc=Which scope: user, project, all;default=all"`
	}{}

	return &gcli.Command{
		Name: "path",
		Desc: "Show the skills directory path",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&opts.Scope, "scope,s", "Which scope: user, project, all")
		},
		Func: func(c *gcli.Command, args []string) error {
			mgr := skills.NewManager()

			if opts.Scope == "" || opts.Scope == "all" {
				c.Infof("User skills:    %s\n", mgr.UserSkillsDir)
				c.Infof("Project skills: %s\n", mgr.ProjectSkillsDir)
			} else if opts.Scope == "user" {
				c.Infof("%s\n", mgr.UserSkillsDir)
			} else if opts.Scope == "project" {
				c.Infof("%s\n", mgr.ProjectSkillsDir)
			}

			return nil
		},
	}
}

// SkillsOpenCmd opens the skills directory in file manager
func SkillsOpenCmd() *gcli.Command {
	var opts = struct {
		Scope string `flag:"name=scope,s;desc=Which scope to open: user or project;default=user"`
	}{}

	return &gcli.Command{
		Name: "open",
		Desc: "Open the skills directory in your file manager",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&opts.Scope, "scope,s", "Which scope to open: user or project")
		},
		Func: func(c *gcli.Command, args []string) error {
			mgr := skills.NewManager()

			var dir string
			if opts.Scope == "project" {
				dir = mgr.ProjectSkillsDir
			} else {
				dir = mgr.UserSkillsDir
			}

			// Create directory if it doesn't exist
			if !fsutil.DirExist(dir) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
				c.Infof("Created directory: %s\n", dir)
			}

			// Open in file manager
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "windows":
				cmd = exec.Command("explorer", dir)
			case "darwin":
				cmd = exec.Command("open", dir)
			default: // linux
				cmd = exec.Command("xdg-open", dir)
			}

			c.Infof("Opening: %s\n", dir)
			return cmd.Start()
		},
	}
}

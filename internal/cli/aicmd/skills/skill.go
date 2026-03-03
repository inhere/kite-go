// Package skills provides skill management functionality for kite ai skills command.
// Skills are extensions that enhance AI capabilities with custom instructions.
package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"gopkg.in/yaml.v3"
)

// Skill represents a single skill with its metadata and content.
type Skill struct {
	// Name is the skill name (from directory or frontmatter)
	Name string `json:"name" yaml:"name"`
	// Description describes what the skill does
	Description string `json:"description" yaml:"description"`
	// Path is the full path to the SKILL.md file
	Path string `json:"path" yaml:"-"`
	// Scope indicates where the skill is located: "user" or "project"
	Scope string `json:"scope" yaml:"-"`
	// Content is the raw markdown content (excluding frontmatter)
	Content string `json:"content" yaml:"-"`
	// Frontmatter contains parsed YAML frontmatter fields
	Frontmatter map[string]any `json:"frontmatter,omitempty" yaml:"-"`
	// Dir is the skill directory path
	Dir string `json:"dir" yaml:"-"`
}

// SkillFile is the structure for parsing SKILL.md frontmatter
type SkillFile struct {
	Name                   string   `yaml:"name,omitempty"`
	Description            string   `yaml:"description,omitempty"`
	ArgumentHint           string   `yaml:"argument-hint,omitempty"`
	DisableModelInvocation bool     `yaml:"disable-model-invocation,omitempty"`
	UserInvocable          *bool    `yaml:"user-invocable,omitempty"`
	AllowedTools           []string `yaml:"allowed-tools,omitempty"`
	Model                  string   `yaml:"model,omitempty"`
	Context                string   `yaml:"context,omitempty"`
	Agent                  string   `yaml:"agent,omitempty"`
}

// Manager handles skill operations
type Manager struct {
	// UserSkillsDir is the user-level skills directory (~/.claude/skills/)
	UserSkillsDir string
	// ProjectSkillsDir is the project-level skills directory (./.claude/skills/)
	ProjectSkillsDir string
}

// NewManager creates a new skill manager
func NewManager() *Manager {
	homeDir, _ := os.UserHomeDir()
	userSkillsDir := filepath.Join(homeDir, ".claude", "skills")
	projectSkillsDir := filepath.Join(".", ".claude", "skills")

	return &Manager{
		UserSkillsDir:    userSkillsDir,
		ProjectSkillsDir: projectSkillsDir,
	}
}

// ScanSkills scans all skills from both user and project directories
func (m *Manager) ScanSkills(scope string) ([]*Skill, error) {
	var skills []*Skill

	// Scan user skills if scope is "user" or "all"
	if scope == "" || scope == "user" || scope == "all" {
		userSkills, err := m.scanSkillsDir(m.UserSkillsDir, "user")
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to scan user skills: %w", err)
		}
		skills = append(skills, userSkills...)
	}

	// Scan project skills if scope is "project" or "all"
	if scope == "" || scope == "project" || scope == "all" {
		projectSkills, err := m.scanSkillsDir(m.ProjectSkillsDir, "project")
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to scan project skills: %w", err)
		}
		skills = append(skills, projectSkills...)
	}

	return skills, nil
}

// scanSkillsDir scans a single directory for skills
func (m *Manager) scanSkillsDir(dir, scope string) ([]*Skill, error) {
	var skills []*Skill

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(dir, entry.Name(), "SKILL.md")
		if !fsutil.FileExists(skillPath) {
			continue
		}

		skill, err := m.loadSkill(skillPath, scope)
		if err != nil {
			// Skip invalid skills but log the error
			fmt.Printf("Warning: failed to load skill %s: %v\n", entry.Name(), err)
			continue
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// loadSkill loads a single skill from its SKILL.md file
func (m *Manager) loadSkill(path, scope string) (*Skill, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	skill := &Skill{
		Path:  path,
		Scope: scope,
		Dir:   filepath.Dir(path),
	}

	// Parse frontmatter and content
	frontmatter, body, err := parseFrontmatter(string(content))
	if err != nil {
		return nil, err
	}

	skill.Content = body
	skill.Frontmatter = frontmatter

	// Extract name from frontmatter or directory
	if name, ok := frontmatter["name"].(string); ok && name != "" {
		skill.Name = name
	} else {
		skill.Name = filepath.Base(filepath.Dir(path))
	}

	// Extract description from frontmatter
	if desc, ok := frontmatter["description"].(string); ok {
		skill.Description = desc
	} else {
		// Use first paragraph of content as description
		lines := strings.Split(strings.TrimSpace(body), "\n")
		if len(lines) > 0 {
			skill.Description = strutil.Truncate(strings.TrimSpace(lines[0]), 80, "...")
	}
	}
	return skill, nil
}

// parseFrontmatter parses YAML frontmatter from markdown content
func parseFrontmatter(content string) (map[string]any, string, error) {
	content = strings.TrimSpace(content)

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, "---") {
		return nil, content, nil
	}

	// Find the closing delimiter
	endIndex := strings.Index(content[3:], "---")
	if endIndex == -1 {
		return nil, content, nil
	}

	frontmatterStr := strings.TrimSpace(content[3 : endIndex+3])
	body := strings.TrimSpace(content[endIndex+6:])

	var frontmatter map[string]any
	if err := yaml.Unmarshal([]byte(frontmatterStr), &frontmatter); err != nil {
		return nil, content, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	return frontmatter, body, nil
}

// GetSkill retrieves a single skill by name
func (m *Manager) GetSkill(name string) (*Skill, error) {
	skills, err := m.ScanSkills("all")
	if err != nil {
		return nil, err
	}

	for _, skill := range skills {
		if skill.Name == name {
			return skill, nil
		}
	}

	return nil, fmt.Errorf("skill %q not found", name)
}

// CreateSkill creates a new skill with the given name and description
func (m *Manager) CreateSkill(name, description, scope string) error {
	// Validate name
	if name == "" {
		return fmt.Errorf("skill name is required")
	}
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("skill name contains invalid characters")
	}

	// Determine directory
	var baseDir string
	if scope == "project" {
		baseDir = m.ProjectSkillsDir
	} else {
		baseDir = m.UserSkillsDir
	}

	// Create skill directory
	skillDir := filepath.Join(baseDir, name)
	if fsutil.DirExist(skillDir) {
		return fmt.Errorf("skill %q already exists", name)
	}

	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	// Create SKILL.md
	skillContent := m.generateSkillTemplate(name, description)
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(skillContent), 0644); err != nil {
		return fmt.Errorf("failed to create SKILL.md: %w", err)
	}

	return nil
}

// generateSkillTemplate generates the default SKILL.md content
func (m *Manager) generateSkillTemplate(name, description string) string {
	if description == "" {
		description = "Description of what this skill does and when to use it."
	}

	return fmt.Sprintf(`---
name: %s
description: %s
---

# %s

Instructions for the skill go here.

## Usage

Describe how to use this skill.

## Examples

Provide examples if helpful.
`, name, description, name)
}

// DeleteSkill deletes a skill by name
func (m *Manager) DeleteSkill(name string) error {
	skill, err := m.GetSkill(name)
	if err != nil {
		return err
	}

	// Remove the skill directory
	if err := os.RemoveAll(skill.Dir); err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	return nil
}

// EditSkill opens the skill in the default editor
func (m *Manager) EditSkill(name string) error {
	skill, err := m.GetSkill(name)
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
		editor = "vim" // Could be "code", "nano", etc.
	}

	// This would typically use exec.Command to open the editor
	// For now, just return the path
	return fmt.Errorf("open %s with editor %s", skill.Path, editor)
}

// GetSkillsPath returns the path to skills directory
func (m *Manager) GetSkillsPath(scope string) string {
	if scope == "project" {
		return m.ProjectSkillsDir
	}
	return m.UserSkillsDir
}

// SkillExists checks if a skill with the given name exists
func (m *Manager) SkillExists(name string) bool {
	_, err := m.GetSkill(name)
	return err == nil
}

package manager

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

type StateTomlUpdater struct {
	newBuf byteutil.Buffer
	// raw contents of the state file
	contents []byte

	currentSection string
	// 已经处理过的键，避免重复添加
	processedKeys  map[string]bool
	processedSecs  map[string]bool
	processedPaths map[string]bool
}

var pathsRegex = regexp.MustCompile(`paths\s*=\s*\[`)

// NewTomlUpdater creates a new StateTomlUpdater
func NewTomlUpdater() *StateTomlUpdater {
	return &StateTomlUpdater{
		processedKeys:  make(map[string]bool),
		processedSecs:  make(map[string]bool),
		processedPaths: make(map[string]bool),
	}
}

// Reset the state to its initial state
func (u *StateTomlUpdater) Reset() {
	// u.contents = nil
	u.newBuf.Reset()
	u.processedKeys = make(map[string]bool)
	u.processedSecs = make(map[string]bool)
	u.processedPaths = make(map[string]bool)
}

// Update the state file with the given activity state
func (u *StateTomlUpdater) Update(state *models.ActivityState) error {
	// 读取 state.File 现在的文件内容
	data, err := os.ReadFile(state.File)
	if err != nil {
		// 如果文件不存在，则创建文件
		if os.IsNotExist(err) {
			return u.WriteNewState(state)
		}
		return fmt.Errorf("failed to read state file: %w", err)
	}

	u.contents = data
	u.Build(state)

	// 使用 newBuf 更新 state.File 内容
	return os.WriteFile(state.File, u.newBuf.Bytes(), 0644)
}

// WriteNewState write new state file
func (u *StateTomlUpdater) WriteNewState(state *models.ActivityState) error {
	contents, err := toml.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	return os.WriteFile(state.File, contents, 0644)
}

// Build the updated state file contents
func (u *StateTomlUpdater) Build(state *models.ActivityState) *StateTomlUpdater {
	if len(u.contents) == 0 {
		_ = u.WriteNewState(state)
		return u
	}

	var inPathsValues bool
	var pathsValues []string
	u.Reset()

	// 按行处理 state 数据对比更新
	scanner := bufio.NewScanner(bytes.NewReader(u.contents))

	// 遍历每一行
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			u.newBuf.WriteStr1Nl(line)
			continue
		}

		inlineComment := ""
		if idx := strings.Index(trimmed, "#"); idx > 0 {
			inlineComment = " " + trimmed[idx:]
			trimmed = strings.TrimSpace(trimmed[:idx])
		}

		// paths array start: "paths = [" 使用 regex 匹配
		if pathsRegex.MatchString(trimmed) {
			inPathsValues = true
			u.newBuf.WriteStr1Nl(line)
			continue
		}

		if inPathsValues {
			pathVal := strings.Trim(trimmed, "\"',")
			if pathVal == "" {
				continue
			}

			// paths array 结束
			if strings.HasSuffix(trimmed, "]") {
				// 对比检查是否还有未保存的path
				for _, newPath := range state.Paths {
					if !arrutil.StringsContains(pathsValues, newPath) {
						u.newBuf.Writef("  %q,%s\n", newPath)
					}
				}

				inPathsValues = false
				u.newBuf.WriteStr1Nl(line)
				u.processedSecs["paths"] = true
			} else if state.ExistsPath(pathVal) {
				// paths value
				u.newBuf.Writef("  %q,%s\n", pathVal, inlineComment)
			}
			continue
		}

		// 处理节标题 [section]: envs, tools, sdks
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			// 之前没有添加过 paths
			if u.processedSecs["paths"] == false {
				inPathsValues = false
				for _, newPath := range state.Paths {
					u.newBuf.Writef("  %q,%s\n", newPath)
				}
			}

			u.currentSection = line[1 : len(line)-1] // 去掉方括号
			u.newBuf.WriteStr1Nl(line)
			u.processedSecs[u.currentSection] = true
			continue
		}

		// 处理键值对 key = value
		if eqIndex := strings.Index(trimmed, "="); eqIndex > 0 {
			key := strings.TrimSpace(trimmed[:eqIndex])
			fullKey := strutil.OrCond(u.currentSection == "", key, u.currentSection+"."+key)
			u.processedKeys[fullKey] = true

			newVal := ""
			switch u.currentSection {
			case "envs":
				newVal = state.Envs[key]
			case "sdks":
				newVal = state.SDKs[key]
			case "tools":
				newVal = state.Tools[key]
			default:
				u.newBuf.WriteStr1Nl(line)
			}

			if newVal != "" {
				u.newBuf.Writef("%s = %q%s\n", key, newVal, inlineComment)
			}
		}
	}

	// 最后，检查添加未处理的 section
	u.addNewSections(state)
	return u
}

func (u *StateTomlUpdater) addNewSections(state *models.ActivityState) {
	// 添加未处理的 section 键值对
	sections := []string{"envs", "sdks", "tools"}
	for _, sec := range sections {
		if u.processedSecs[sec] {
			continue
		}

		switch sec {
		case "envs":
			for key, val := range state.Envs {
				u.newBuf.Writef("%s = %q\n", key, val)
			}
		case "tools":
			for key, val := range state.Tools {
				u.newBuf.Writef("%s = %q\n", key, val)
			}
		case "sdks":
			for key, val := range state.SDKs {
				u.newBuf.Writef("%s = %q\n", key, val)
			}
		}
		u.newBuf.WriteStr1Nl("")
	}
}

// NewBufBytes get new contents bytes
func (u *StateTomlUpdater) NewBufBytes() []byte {
	return u.newBuf.Bytes()
}

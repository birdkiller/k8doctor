package knowledge

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// KBEntry 知识库条目
type KBEntry struct {
	ID            string   `json:"id"`              // 唯一标识 (文件名)
	Category      string   `json:"category"`        // 分类: pod/node/network/storage/rbac/deployment
	Severity      string   `json:"severity"`        // 严重度: critical/warning/info
	Title         string   `json:"title"`           // 故障标题
	RelatedStates []string `json:"related_states"`  // 关联状态关键词
	Tags          []string `json:"tags"`            // 标签
	Content       string   `json:"content"`         // 完整正文（用于 Embedding）
	Steps         []string `json:"steps"`          // 排查步骤
	Commands      []string `json:"commands"`       // 可执行命令
}

// frontmatter 解析结果
type frontmatter struct {
	ID            string   `yaml:"id"`
	Category      string   `yaml:"category"`
	Severity      string   `yaml:"severity"`
	Title         string   `yaml:"title"`
	RelatedStates []string `yaml:"related_states"`
	Tags          []string `yaml:"tags"`
}

// KnowledgeBase 知识库
type KnowledgeBase struct {
	Path    string
	Entries []*KBEntry
}

// NewKnowledgeBase 创建知识库实例
func NewKnowledgeBase(path string) (*KnowledgeBase, error) {
	kb := &KnowledgeBase{
		Path:    path,
		Entries: make([]*KBEntry, 0),
	}

	// 遍历 kb 目录
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取知识库目录失败: %w", err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(path, f.Name())
		entry, err := parseFile(filePath)
		if err != nil {
			// 跳过解析失败的文件
			continue
		}

		kb.Entries = append(kb.Entries, entry)
	}

	return kb, nil
}

// parseFile 解析单个 Markdown 文件
func parseFile(filePath string) (*KBEntry, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	entry := &KBEntry{
		ID: filepath.Base(filePath[:len(filePath)-3]),
	}

	text := string(content)

	// 解析 frontmatter
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		// 提取 frontmatter
		endIdx := -1
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				endIdx = i
				break
			}
		}

		if endIdx > 1 {
			frontmatterContent := strings.Join(lines[1:endIdx], "\n")
			var fm frontmatter
			if err := yaml.Unmarshal([]byte(frontmatterContent), &fm); err == nil {
				entry.Category = fm.Category
				entry.Severity = fm.Severity
				entry.Title = fm.Title
				entry.RelatedStates = fm.RelatedStates
				entry.Tags = fm.Tags
			}
		}
	}

	// 移除 frontmatter，只保留正文
	contentLines := lines
	for i, line := range lines {
		if strings.TrimSpace(line) == "---" && i > 0 {
			contentLines = lines[i+1:]
			break
		}
	}

	// 解析 Markdown 内容
	contentStr := strings.Join(contentLines, "\n")

	// 提取标题（第一个 # 开头的内容）
	titleRegex := regexp.MustCompile(`(?m)^#\s+(.+)$`)
	titleMatch := titleRegex.FindStringSubmatch(contentStr)
	if len(titleMatch) > 1 && entry.Title == "" {
		entry.Title = strings.TrimSpace(titleMatch[1])
	}

	// 提取排查步骤（## 排查步骤 或 ## 修复建议 之间的内容）
	stepsRegex := regexp.MustCompile(`(?ms)## 排查步骤\s*\n(.*?)(?=##|\z)`)
	stepsMatch := stepsRegex.FindStringSubmatch(contentStr)
	if len(stepsMatch) > 1 {
		stepsText := stepsMatch[1]
		// 提取每一行作为步骤
		for _, line := range strings.Split(stepsText, "\n") {
			line = strings.TrimSpace(line)
			// 过滤掉代码块和空行
			if len(line) > 5 && !strings.HasPrefix(line, "```") && !strings.HasPrefix(line, "#") {
				// 去除列表前缀
				line = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(line, "")
				line = regexp.MustCompile(`^[-*]\s*`).ReplaceAllString(line, "")
				if len(line) > 3 {
					entry.Steps = append(entry.Steps, line)
				}
			}
		}
	}

	// 提取可执行命令（kubectl 开头的行）
	cmdRegex := regexp.MustCompile(`(?m)^\s*(kubectl\s+[^\n]+)`)
	cmdMatches := cmdRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range cmdMatches {
		if len(match) > 1 {
			cmd := strings.TrimSpace(match[1])
			if len(cmd) > 5 {
				entry.Commands = append(entry.Commands, cmd)
			}
		}
	}

	// 收集正文内容用于 Embedding
	builder := &strings.Builder{}
	builder.WriteString(entry.Title)
	if len(entry.Tags) > 0 {
		builder.WriteString(" ")
		builder.WriteString(strings.Join(entry.Tags, " "))
	}
	if len(entry.RelatedStates) > 0 {
		builder.WriteString(" ")
		builder.WriteString(strings.Join(entry.RelatedStates, " "))
	}
	// 添加正文
	for _, line := range strings.Split(contentStr, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 3 && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "```") {
			builder.WriteString(" ")
			builder.WriteString(line)
		}
	}
	entry.Content = builder.String()

	// 如果没有找到标题，使用文件名
	if entry.Title == "" {
		entry.Title = entry.ID
	}

	// 默认分类
	if entry.Category == "" {
		entry.Category = "unknown"
	}

	return entry, nil
}

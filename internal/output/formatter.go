package output

import (
	"fmt"
	"io"
	"strings"

	"k8doctor/internal/matcher"
)

// Print 输出诊断结果
func Print(w io.Writer, results []*matcher.MatchResult, rawInput string) {
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "🔍 诊断结果")
	fmt.Fprintln(w, "============")

	if len(results) == 0 {
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "⚠️  未能识别故障")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "可能原因：")
		fmt.Fprintln(w, "  1. 症状描述不够具体")
		fmt.Fprintln(w, "  2. 该故障类型暂未收录")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "💡 建议：")
		fmt.Fprintln(w, "  - 提供更多细节：Pod名称、命名空间、错误日志等")
		fmt.Fprintln(w, "  - 使用 diagnose-i 进行交互式诊断")
		return
	}

	// 原始症状
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "📝 原始症状: %s\n", rawInput)
	fmt.Fprintln(w, "")

	for i, result := range results {
		entry := result.Entry

		severity := "⚠️  Warning"
		if entry.Severity == "critical" {
			severity = "🔴 Critical"
		} else if entry.Severity == "info" {
			severity = "ℹ️  Info"
		}

		matchIcon := "🎯"
		if result.MatchType == "keyword" {
			matchIcon = "🔤"
		}

		fmt.Fprintln(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintf(w, "%s 匹配故障 [%d/%d] | 相似度: %.2f | 匹配方式: %s\n",
			matchIcon, i+1, len(results), result.Score, result.MatchType)
		fmt.Fprintf(w, "📌 故障名称: %s\n", entry.Title)
		fmt.Fprintf(w, "%s 严重程度: %s\n", severity, entry.Severity)
		fmt.Fprintf(w, "📁 分类: %s\n", entry.Category)
		if len(entry.Tags) > 0 {
			fmt.Fprintf(w, "🏷️  标签: %s\n", strings.Join(entry.Tags, ", "))
		}
		fmt.Fprintln(w, "")

		// 排查步骤
		if len(entry.Steps) > 0 {
			fmt.Fprintln(w, "## 📋 排查步骤")
			for j, step := range entry.Steps {
				// 过滤非步骤内容
				if strings.HasPrefix(step, "#") || len(step) < 5 {
					continue
				}
				fmt.Fprintf(w, "   %d. %s\n", j+1, step)
			}
			fmt.Fprintln(w, "")
		}

		// 可执行命令
		if len(entry.Commands) > 0 {
			fmt.Fprintln(w, "## 💻 可执行命令")
			for _, cmd := range entry.Commands {
				if strings.Contains(cmd, "kubectl") || strings.HasPrefix(strings.TrimSpace(cmd), "#") {
					fmt.Fprintf(w, "   %s\n", cmd)
				}
			}
			fmt.Fprintln(w, "")
		}
	}

	fmt.Fprintln(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w, "")
}

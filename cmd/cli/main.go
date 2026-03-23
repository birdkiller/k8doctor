package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8doctor/internal/cleaner"
	"k8doctor/internal/knowledge"
	"k8doctor/internal/matcher"
	"k8doctor/internal/output"
)

var (
	kbPath string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "k8doctor",
		Short: "K8s 智能诊断工具 - 基于知识库的故障排查专家",
		Long:  `通过语义理解匹配 K8s 故障知识库，给出完整可行的排查建议。`,
	}

	diagnoseCmd := &cobra.Command{
		Use:   "diagnose [症状描述]",
		Short: "诊断 K8s 故障",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			symptomText := args[0]
			return runDiagnose(symptomText)
		},
	}

	diagnoseICmd := &cobra.Command{
		Use:   "diagnose-i",
		Short: "交互式诊断",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("请输入故障症状描述：")
			var symptomText string
			fmt.Scanln(&symptomText)
			if symptomText == "" {
				fmt.Println("未输入症状描述")
				return nil
			}
			return runDiagnose(symptomText)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出所有支持的故障类型",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList()
		},
	}

	rebuildCmd := &cobra.Command{
		Use:   "rebuild-index",
		Short: "重新构建向量索引",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRebuildIndex()
		},
	}

	rootCmd.PersistentFlags().StringVar(&kbPath, "kb-path", "./kb", "知识库路径")

	rootCmd.AddCommand(diagnoseCmd)
	rootCmd.AddCommand(diagnoseICmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(rebuildCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行错误: %v\n", err)
		os.Exit(1)
	}
}

func runDiagnose(symptomText string) error {
	// 加载知识库
	kb, err := knowledge.NewKnowledgeBase(kbPath)
	if err != nil {
		return fmt.Errorf("加载知识库失败: %w", err)
	}

	// 初始化 Embedding 匹配器
	m, err := matcher.New(kb)
	if err != nil {
		return fmt.Errorf("初始化匹配器失败: %w", err)
	}

	// 清洗症状
	cleaned := cleaner.Clean(symptomText)

	// 执行匹配
	results, err := m.Match(cleaned, 3)
	if err != nil {
		return fmt.Errorf("匹配失败: %w", err)
	}

	// 输出结果
	output.Print(os.Stdout, results, cleaned.RawInput)

	return nil
}

func runList() error {
	kb, err := knowledge.NewKnowledgeBase(kbPath)
	if err != nil {
		return fmt.Errorf("加载知识库失败: %w", err)
	}

	fmt.Println("支持的故障类型：")
	fmt.Println("================")
	for _, entry := range kb.Entries {
		severity := "⚠️ "
		if entry.Severity == "critical" {
			severity = "🔴"
		}
		fmt.Printf("%s [%s] %s\n", severity, entry.Category, entry.Title)
	}
	return nil
}

func runRebuildIndex() error {
	kb, err := knowledge.NewKnowledgeBase(kbPath)
	if err != nil {
		return fmt.Errorf("加载知识库失败: %w", err)
	}

	_, err = matcher.New(kb)
	if err != nil {
		return fmt.Errorf("重建索引失败: %w", err)
	}

	fmt.Println("向量索引重建完成")
	return nil
}

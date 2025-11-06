package cmd

import (
	"github.com/phixia/phjvgen/internal/generator"
	"github.com/phixia/phjvgen/internal/utils"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "g"},
	Short:   "生成新的 Java 项目",
	Long: `生成一个新的基于 Java 25 LTS 的分层架构项目。

该命令会交互式地询问项目配置信息，然后生成完整的项目结构，包括：
  - Maven 多模块结构
  - 分层架构（common, domain, infrastructure, adapter, application, starter）
  - 基础代码（Application启动类、Result响应封装、异常处理等）
  - 配置文件（application.yml）
  - README 和 .gitignore

生成后的项目可以直接使用 Maven 构建和运行。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get project configuration
		config, err := generator.GetProjectConfig()
		if err != nil {
			return err
		}

		// Generate project
		if err := generator.GenerateProject(config); err != nil {
			utils.PrintError(err.Error())
			return err
		}

		// Print summary
		generator.PrintGenerationSummary(config)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

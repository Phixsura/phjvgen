package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "phjvgen",
	Short: "Java 25 LTS 分层架构项目生成器",
	Long: `phjvgen - Java 25 LTS Project Generator

一个基于 Java 25 LTS 的分层架构项目生成工具。
支持快速生成完整的 Spring Boot 项目结构，包括：
  - 领域驱动设计(DDD)分层架构
  - 完整的 CRUD 示例代码
  - Maven 多模块项目结构
  - MyBatis Plus 数据访问
  - Spring Boot 4.0 配置

使用示例:
  phjvgen generate         # 生成新项目
  phjvgen demo             # 生成完整 CRUD 示例
  phjvgen add payment      # 添加新业务模块
  phjvgen example          # 快速生成示例项目`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

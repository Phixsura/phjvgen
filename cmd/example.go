package cmd

import (
	"fmt"

	"github.com/phixia/phjvgen/internal/generator"
	"github.com/phixia/phjvgen/internal/utils"
	"github.com/spf13/cobra"
)

var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "快速生成示例项目",
	Long: `快速生成一个预配置的示例项目，无需交互式输入。

该命令会使用以下预设配置生成一个完整的示例项目：
  - Group ID: com.example.demo
  - Artifact ID: demo-app
  - Version: 1.0.0
  - Project Name: Demo Application
  - Description: A demo application for testing
  - Output Dir: ./demo-app

生成后可以直接进入目录构建和运行项目。

适用场景：
  - 快速测试和学习项目结构
  - CI/CD 集成测试
  - 项目模板验证`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.PrintBanner()
		fmt.Println()
		utils.PrintInfo("使用预设配置快速生成示例项目...")
		fmt.Println()

		// Create example config
		config := &generator.ProjectConfig{
			GroupID:            "com.example.demo",
			ArtifactID:         "demo-app",
			Version:            "1.0.0",
			ProjectName:        "Demo Application",
			ProjectDescription: "A demo application for testing",
			PackageName:        "com.example.demo",
			PackagePath:        "com/example/demo",
			OutputDir:          "./demo-app",
		}

		// Display configuration
		utils.PrintInfo("项目配置信息：")
		fmt.Printf("  Group ID: %s\n", config.GroupID)
		fmt.Printf("  Artifact ID: %s\n", config.ArtifactID)
		fmt.Printf("  Version: %s\n", config.Version)
		fmt.Printf("  Project Name: %s\n", config.ProjectName)
		fmt.Printf("  Package Name: %s\n", config.PackageName)
		fmt.Printf("  Output Directory: %s\n", config.OutputDir)
		fmt.Println()

		// Auto-confirm
		confirm := "y"
		fmt.Printf("确认生成项目？(y/n) [y]: %s\n", confirm)
		fmt.Println()

		// Generate project
		if err := generator.GenerateProject(config); err != nil {
			utils.PrintError(err.Error())
			return err
		}

		// Print summary
		generator.PrintGenerationSummary(config)

		fmt.Println()
		utils.PrintInfo("示例项目已生成！可以使用以下命令快速开始：")
		fmt.Println("  cd demo-app")
		fmt.Println("  mvn clean install")
		fmt.Println("  java --enable-preview -jar starter/target/starter-1.0.0.jar")
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}

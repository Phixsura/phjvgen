package cmd

import (
	"github.com/phixia/phjvgen/internal/generator"
	"github.com/phixia/phjvgen/internal/utils"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <module-name>",
	Short: "添加新的业务模块",
	Long: `在现有项目中添加一个新的 Application 业务模块。

该命令会创建一个新的 application-<module-name> 模块，包括：
  - 完整的目录结构（service, dto, assembler, executor）
  - 模块 pom.xml
  - 示例 Service 类
  - 自动更新父 pom.xml 的 modules 和 dependencyManagement

模块名称格式要求：
  - 只能包含小写字母、数字和连字符
  - 必须以小写字母开头
  - 例如：payment, order, product, user-management

使用示例：
  phjvgen add payment        # 创建 application-payment 模块
  phjvgen add order          # 创建 application-order 模块
  phjvgen add user-profile   # 创建 application-user-profile 模块

注意：必须在项目根目录（包含 pom.xml 的目录）下运行此命令。`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleName := args[0]

		if err := generator.AddApplicationModule(moduleName); err != nil {
			utils.PrintError(err.Error())
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

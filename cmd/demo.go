package cmd

import (
	"github.com/phixia/phjvgen/internal/generator"
	"github.com/phixia/phjvgen/internal/utils"
	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "生成完整的 CRUD 示例代码",
	Long: `在现有项目中生成完整的 User 模块 CRUD 示例代码。

该命令会在当前项目中生成包含以下内容的完整示例：
  - Domain层：User实体、UserRepository接口
  - Infrastructure层：UserDO、UserMapper、UserRepositoryImpl
  - Application层：UserDTO、UserService、UserAssembler、Commands
  - Adapter层：UserController、Request/Response VO、全局异常处理
  - 数据库脚本：用户表创建SQL

注意：必须在项目根目录（包含 pom.xml 的目录）下运行此命令。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generator.GenerateDemoCode(); err != nil {
			utils.PrintError(err.Error())
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
}

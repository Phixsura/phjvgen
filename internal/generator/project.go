package generator

import (
	"fmt"
	"path/filepath"

	"github.com/phixia/phjvgen/internal/templates"
	"github.com/phixia/phjvgen/internal/utils"
)

// GenerateProject generates the complete project structure
func GenerateProject(config *ProjectConfig) error {
	utils.PrintInfo("创建目录结构...")
	if err := createDirectoryStructure(config); err != nil {
		return err
	}
	utils.PrintSuccess("目录结构创建完成")

	utils.PrintInfo("生成父POM文件...")
	if err := generateParentPOM(config); err != nil {
		return err
	}
	utils.PrintSuccess("父POM文件生成完成")

	utils.PrintInfo("生成模块POM文件...")
	if err := generateModulePOMs(config); err != nil {
		return err
	}
	utils.PrintSuccess("模块POM文件生成完成")

	utils.PrintInfo("生成基础Java源文件...")
	if err := generateBasicSourceFiles(config); err != nil {
		return err
	}
	utils.PrintSuccess("基础Java源文件生成完成")

	utils.PrintInfo("生成配置文件...")
	if err := generateConfigFiles(config); err != nil {
		return err
	}
	utils.PrintSuccess("配置文件生成完成")

	utils.PrintInfo("生成README.md...")
	if err := generateREADME(config); err != nil {
		return err
	}
	utils.PrintSuccess("README.md生成完成")

	utils.PrintInfo("生成.gitignore...")
	if err := generateGitIgnore(config); err != nil {
		return err
	}
	utils.PrintSuccess(".gitignore生成完成")

	return nil
}

func createDirectoryStructure(config *ProjectConfig) error {
	baseDir := config.OutputDir
	pkgPath := config.PackagePath

	dirs := []string{
		// Common module
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/exception"),
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/response"),
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/constant"),
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/utils"),
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/annotation"),
		filepath.Join(baseDir, "common/src/main/resources"),
		filepath.Join(baseDir, "common/src/test/java", pkgPath, "common"),

		// Domain module
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/model"),
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/event"),
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/repository"),
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/service"),
		filepath.Join(baseDir, "domain/src/main/resources"),
		filepath.Join(baseDir, "domain/src/test/java", pkgPath, "domain"),

		// Infrastructure module
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/persistence/mapper"),
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/persistence/impl"),
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/cache"),
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/mq"),
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/gateway"),
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/config"),
		filepath.Join(baseDir, "infrastructure/src/main/resources/mapper"),
		filepath.Join(baseDir, "infrastructure/src/main/resources/db/migration"),
		filepath.Join(baseDir, "infrastructure/src/test/java", pkgPath, "infrastructure"),

		// Adapter-Rest module
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/controller"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/request"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/response"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/interceptor"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/filter"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/advice"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/resources"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/test/java", pkgPath, "adapter/rest"),

		// Adapter-Schedule module
		filepath.Join(baseDir, "adapter/adapter-schedule/src/main/java", pkgPath, "adapter/schedule/job"),
		filepath.Join(baseDir, "adapter/adapter-schedule/src/main/java", pkgPath, "adapter/schedule/config"),
		filepath.Join(baseDir, "adapter/adapter-schedule/src/main/resources"),
		filepath.Join(baseDir, "adapter/adapter-schedule/src/test/java", pkgPath, "adapter/schedule"),

		// Application-User module
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/service"),
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/dto"),
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/assembler"),
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/executor"),
		filepath.Join(baseDir, "application/application-user/src/main/resources"),
		filepath.Join(baseDir, "application/application-user/src/test/java", pkgPath, "application/user"),

		// Starter module
		filepath.Join(baseDir, "starter/src/main/java", pkgPath),
		filepath.Join(baseDir, "starter/src/main/resources"),
		filepath.Join(baseDir, "starter/src/test/java", pkgPath),
	}

	return utils.CreateDirs(dirs...)
}

func generateParentPOM(config *ProjectConfig) error {
	content := utils.ReplacePlaceholders(templates.ParentPOM, config.GetReplacements())
	return utils.WriteFile(filepath.Join(config.OutputDir, "pom.xml"), content)
}

func generateModulePOMs(config *ProjectConfig) error {
	replacements := config.GetReplacements()

	poms := map[string]string{
		"common/pom.xml":                           templates.CommonPOM,
		"domain/pom.xml":                           templates.DomainPOM,
		"infrastructure/pom.xml":                   templates.InfrastructurePOM,
		"adapter/adapter-rest/pom.xml":             templates.AdapterRestPOM,
		"adapter/adapter-schedule/pom.xml":         templates.AdapterSchedulePOM,
		"application/application-user/pom.xml":     templates.ApplicationUserPOM,
		"starter/pom.xml":                          templates.StarterPOM,
	}

	for path, template := range poms {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(filepath.Join(config.OutputDir, path), content); err != nil {
			return err
		}
	}

	return nil
}

func generateBasicSourceFiles(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	baseDir := config.OutputDir
	pkgPath := config.PackagePath

	files := map[string]string{
		// Starter module
		filepath.Join("starter/src/main/java", pkgPath, "Application.java"): templates.ApplicationMain,

		// Common module
		filepath.Join("common/src/main/java", pkgPath, "common/response/Result.java"):         templates.ResultClass,
		filepath.Join("common/src/main/java", pkgPath, "common/exception/BusinessException.java"): templates.BusinessExceptionClass,

		// Adapter-Rest module
		filepath.Join("adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/controller/HealthController.java"): templates.HealthControllerClass,
	}

	for path, template := range files {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(filepath.Join(baseDir, path), content); err != nil {
			return err
		}
	}

	return nil
}

func generateConfigFiles(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	baseDir := config.OutputDir

	files := map[string]string{
		"starter/src/main/resources/application.yml":     templates.ApplicationYML,
		"starter/src/main/resources/application-dev.yml": templates.ApplicationDevYML,
	}

	for path, template := range files {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(filepath.Join(baseDir, path), content); err != nil {
			return err
		}
	}

	return nil
}

func generateREADME(config *ProjectConfig) error {
	content := utils.ReplacePlaceholders(templates.README, config.GetReplacements())
	return utils.WriteFile(filepath.Join(config.OutputDir, "README.md"), content)
}

func generateGitIgnore(config *ProjectConfig) error {
	return utils.WriteFile(filepath.Join(config.OutputDir, ".gitignore"), templates.GitIgnore)
}

// PrintGenerationSummary prints a summary after project generation
func PrintGenerationSummary(config *ProjectConfig) {
	fmt.Println()
	utils.PrintSuccess("==========================================")
	utils.PrintSuccess("项目生成完成！")
	utils.PrintSuccess("==========================================")
	fmt.Println()
	utils.PrintInfo(fmt.Sprintf("项目位置: %s", config.OutputDir))
	fmt.Println()
	utils.PrintInfo("已生成的基础代码：")
	fmt.Println("  ✅ Application主类（启动类）")
	fmt.Println("  ✅ Result统一响应封装")
	fmt.Println("  ✅ BusinessException业务异常")
	fmt.Println("  ✅ HealthController健康检查接口")
	fmt.Println()
	utils.PrintInfo("后续步骤:")
	fmt.Printf("  1. cd %s\n", config.OutputDir)
	fmt.Println("  2. (可选) 运行 phjvgen demo 生成完整CRUD示例")
	fmt.Println("  3. mvn clean install")
	fmt.Println("  4. java --enable-preview -jar starter/target/starter-*.jar")
	fmt.Println("  5. 测试: curl http://localhost:8080/api/health")
	fmt.Println()
	utils.PrintInfo("添加新模块:")
	fmt.Println("  phjvgen add <模块名>")
	fmt.Println()
}

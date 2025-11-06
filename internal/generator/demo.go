package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/phixia/phjvgen/internal/templates"
	"github.com/phixia/phjvgen/internal/utils"
)

// GenerateDemoCode generates full CRUD demo code for the User entity
func GenerateDemoCode() error {
	// Check if we're in project root
	if !utils.FileExists("pom.xml") || !utils.DirExists("domain") || !utils.DirExists("application") {
		return fmt.Errorf("请在项目根目录运行此命令")
	}

	// Extract project info from pom.xml
	config, err := extractProjectInfoFromPOM()
	if err != nil {
		return err
	}

	utils.PrintInfo(fmt.Sprintf("Package: %s", config.PackageName))
	fmt.Println()
	utils.PrintWarning("此命令将生成完整的User模块CRUD示例代码，包括：")
	fmt.Println("  - Common层：Result、BusinessException、ErrorCode")
	fmt.Println("  - Domain层：User实体、UserRepository接口")
	fmt.Println("  - Infrastructure层：UserDO、UserMapper、UserRepositoryImpl")
	fmt.Println("  - Application层：UserDTO、UserService、UserAssembler")
	fmt.Println("  - Adapter层：UserController、Request/Response、ExceptionHandler")
	fmt.Println("  - Starter层：Application主类")
	fmt.Println("  - 数据库脚本：V1__create_user_table.sql")
	fmt.Println()

	confirm, err := utils.ReadInput("确认生成？(y/n): ")
	if err != nil {
		return err
	}
	if !strings.EqualFold(confirm, "y") && !strings.EqualFold(confirm, "yes") {
		utils.PrintWarning("已取消操作")
		return nil
	}

	// Generate code
	utils.PrintInfo("生成Common层代码...")
	if err := generateCommonCode(config); err != nil {
		return err
	}
	utils.PrintSuccess("Common层代码生成完成")

	utils.PrintInfo("生成Domain层代码...")
	if err := generateDomainCode(config); err != nil {
		return err
	}
	utils.PrintSuccess("Domain层代码生成完成")

	utils.PrintInfo("生成Infrastructure层代码...")
	if err := generateInfrastructureCode(config); err != nil {
		return err
	}
	utils.PrintSuccess("Infrastructure层代码生成完成")

	utils.PrintInfo("生成Application层代码...")
	if err := generateApplicationCode(config); err != nil {
		return err
	}
	utils.PrintSuccess("Application层代码生成完成")

	utils.PrintInfo("生成Adapter层代码...")
	if err := generateAdapterCode(config); err != nil {
		return err
	}
	utils.PrintSuccess("Adapter层代码生成完成")

	utils.PrintInfo("生成Starter层代码...")
	if err := generateStarterCode(config); err != nil {
		return err
	}
	utils.PrintSuccess("Starter层代码生成完成")

	printDemoSummary()
	return nil
}

func generateCommonCode(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	pkgPath := config.PackagePath

	files := map[string]string{
		filepath.Join("common/src/main/java", pkgPath, "common/response/Result.java"):         templates.ResultClass,
		filepath.Join("common/src/main/java", pkgPath, "common/exception/BusinessException.java"): templates.BusinessExceptionClass,
		filepath.Join("common/src/main/java", pkgPath, "common/constant/ErrorCode.java"):      templates.ErrorCodeClass,
	}

	for path, template := range files {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func generateDomainCode(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	pkgPath := config.PackagePath

	files := map[string]string{
		filepath.Join("domain/src/main/java", pkgPath, "domain/model/User.java"):           templates.UserEntity,
		filepath.Join("domain/src/main/java", pkgPath, "domain/repository/UserRepository.java"): templates.UserRepository,
	}

	for path, template := range files {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func generateInfrastructureCode(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	pkgPath := config.PackagePath

	files := map[string]string{
		filepath.Join("infrastructure/src/main/java", pkgPath, "infrastructure/persistence/dataobject/UserDO.java"): templates.UserDO,
		filepath.Join("infrastructure/src/main/java", pkgPath, "infrastructure/persistence/mapper/UserMapper.java"): templates.UserMapper,
		filepath.Join("infrastructure/src/main/java", pkgPath, "infrastructure/persistence/impl/UserRepositoryImpl.java"): templates.UserRepositoryImpl,
		"infrastructure/src/main/resources/db/migration/V1__create_user_table.sql": templates.UserTableSQL,
	}

	for path, template := range files {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func generateApplicationCode(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	pkgPath := config.PackagePath

	files := map[string]string{
		filepath.Join("application/application-user/src/main/java", pkgPath, "application/user/dto/UserDTO.java"):              templates.UserDTO,
		filepath.Join("application/application-user/src/main/java", pkgPath, "application/user/dto/CreateUserCommand.java"):    templates.CreateUserCommand,
		filepath.Join("application/application-user/src/main/java", pkgPath, "application/user/dto/UpdateUserCommand.java"):    templates.UpdateUserCommand,
		filepath.Join("application/application-user/src/main/java", pkgPath, "application/user/assembler/UserAssembler.java"): templates.UserAssembler,
		filepath.Join("application/application-user/src/main/java", pkgPath, "application/user/service/UserService.java"):     templates.UserService,
	}

	for path, template := range files {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func generateAdapterCode(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	pkgPath := config.PackagePath

	files := map[string]string{
		filepath.Join("adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/request/CreateUserRequest.java"):           templates.CreateUserRequest,
		filepath.Join("adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/request/UpdateUserRequest.java"):           templates.UpdateUserRequest,
		filepath.Join("adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/advice/GlobalExceptionHandler.java"):       templates.GlobalExceptionHandler,
		filepath.Join("adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/controller/UserController.java"):           templates.UserController,
		filepath.Join("adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/controller/HealthController.java"):         templates.HealthControllerClass,
	}

	for path, template := range files {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func generateStarterCode(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	pkgPath := config.PackagePath

	content := utils.ReplacePlaceholders(templates.ApplicationMain, replacements)
	return utils.WriteFile(filepath.Join("starter/src/main/java", pkgPath, "Application.java"), content)
}

func extractProjectInfoFromPOM() (*ProjectConfig, error) {
	// Read pom.xml
	content, err := os.ReadFile("pom.xml")
	if err != nil {
		return nil, fmt.Errorf("failed to read pom.xml: %w", err)
	}

	pomContent := string(content)

	// Extract groupId (simple parsing)
	groupID := extractXMLTag(pomContent, "groupId")
	if groupID == "" {
		return nil, fmt.Errorf("could not extract groupId from pom.xml")
	}

	version := extractXMLTag(pomContent, "version")
	if version == "" {
		version = "1.0.0"
	}

	artifactID := extractXMLTag(pomContent, "artifactId")
	if artifactID == "" {
		artifactID = "app"
	}

	packageName := groupID
	packagePath := strings.ReplaceAll(groupID, ".", "/")

	return &ProjectConfig{
		GroupID:     groupID,
		ArtifactID:  artifactID,
		Version:     version,
		PackageName: packageName,
		PackagePath: packagePath,
		OutputDir:   ".",
	}, nil
}

func extractXMLTag(content, tag string) string {
	startTag := "<" + tag + ">"
	endTag := "</" + tag + ">"

	startIdx := strings.Index(content, startTag)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(startTag)

	endIdx := strings.Index(content[startIdx:], endTag)
	if endIdx == -1 {
		return ""
	}

	return strings.TrimSpace(content[startIdx : startIdx+endIdx])
}

func printDemoSummary() {
	fmt.Println()
	utils.PrintSuccess("==========================================")
	utils.PrintSuccess("完整CRUD示例代码生成完成！")
	utils.PrintSuccess("==========================================")
	fmt.Println()
	utils.PrintInfo("下一步：")
	fmt.Println("  1. 检查生成的代码")
	fmt.Println("  2. 创建数据库并执行SQL脚本:")
	fmt.Println("     mysql -u root -p < infrastructure/src/main/resources/db/migration/V1__create_user_table.sql")
	fmt.Println("  3. 配置数据库连接（starter/src/main/resources/application-dev.yml）")
	fmt.Println("  4. 构建项目: mvn clean install")
	fmt.Println("  5. 运行应用: java --enable-preview -jar starter/target/starter-*.jar")
	fmt.Println("  6. 测试API:")
	fmt.Println("     - 健康检查: curl http://localhost:8080/api/health")
	fmt.Println("     - 创建用户: curl -X POST http://localhost:8080/api/users -H 'Content-Type: application/json' -d '{\"username\":\"test\",\"email\":\"test@example.com\"}'")
	fmt.Println("     - 查询用户: curl http://localhost:8080/api/users/1")
	fmt.Println()
}

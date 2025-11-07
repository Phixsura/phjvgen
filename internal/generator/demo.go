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
// It supports running from any directory by searching for the project root
func GenerateDemoCode() error {
	// Find project root by searching for pom.xml
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("未找到项目根目录: %w\n提示: 请确保在包含pom.xml的项目目录或其子目录中运行此命令", err)
	}

	utils.PrintInfo(fmt.Sprintf("项目根目录: %s", projectRoot))

	// Extract project info from pom.xml
	config, err := extractProjectInfoFromPOM(projectRoot)
	if err != nil {
		return err
	}

	// Set output dir to project root
	config.OutputDir = projectRoot

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

// findProjectRoot searches for the project root containing pom.xml
// It starts from current directory and walks up the tree
func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check current directory and parent directories
	dir := currentDir
	for {
		pomPath := filepath.Join(dir, "pom.xml")
		if utils.FileExists(pomPath) {
			// Verify it's a valid project structure
			if utils.DirExists(filepath.Join(dir, "domain")) ||
			   utils.DirExists(filepath.Join(dir, "common")) ||
			   utils.DirExists(filepath.Join(dir, "application")) {
				return dir, nil
			}
		}

		// Go up one directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// Reached root directory
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("未找到包含pom.xml和项目模块的目录")
}

func generateCommonCode(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	pkgPath := config.PackagePath
	baseDir := config.OutputDir

	files := map[string]string{
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/response/Result.java"):         templates.ResultClass,
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/exception/BusinessException.java"): templates.BusinessExceptionClass,
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/constant/ErrorCode.java"):      templates.ErrorCodeClass,
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
	baseDir := config.OutputDir

	files := map[string]string{
		// Model
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/model/User.java"):           templates.UserEntity,
		// Repository
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/repository/UserRepository.java"): templates.UserRepository,
		// Domain Service
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/service/UserDomainService.java"): templates.UserDomainService,
		// Event
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/event/UserCreatedEvent.java"): templates.UserCreatedEvent,
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
	baseDir := config.OutputDir

	files := map[string]string{
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/persistence/dataobject/UserDO.java"): templates.UserDO,
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/persistence/mapper/UserMapper.java"): templates.UserMapper,
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/persistence/impl/UserRepositoryImpl.java"): templates.UserRepositoryImpl,
		filepath.Join(baseDir, "infrastructure/src/main/resources/db/migration/V1__create_user_table.sql"): templates.UserTableSQL,
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
	baseDir := config.OutputDir

	files := map[string]string{
		// DTO
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/dto/UserDTO.java"):              templates.UserDTO,
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/dto/CreateUserCommand.java"):    templates.CreateUserCommand,
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/dto/UpdateUserCommand.java"):    templates.UpdateUserCommand,
		// Assembler
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/assembler/UserAssembler.java"): templates.UserAssembler,
		// Service
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/service/UserService.java"):     templates.UserService,
		// Executor
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/executor/RegisterUserExecutor.java"): templates.RegisterUserExecutor,
		// Listener
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/listener/UserEventListener.java"): templates.UserEventListener,
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
	baseDir := config.OutputDir

	files := map[string]string{
		// Request
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/request/CreateUserRequest.java"):           templates.CreateUserRequest,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/request/UpdateUserRequest.java"):           templates.UpdateUserRequest,
		// Response VO
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/response/UserResponseVO.java"):             templates.UserResponseVO,
		// Assembler
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/assembler/UserControllerAssembler.java"):   templates.UserControllerAssembler,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/assembler/ResponseVOAssembler.java"):       templates.ResponseVOAssembler,
		// Filter & Interceptor
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/filter/LoggingFilter.java"):                templates.LoggingFilter,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/interceptor/AuthInterceptor.java"):         templates.AuthInterceptor,
		// Config
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/config/WebMvcConfig.java"):                 templates.WebMvcConfig,
		// Advice & Controller
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/advice/GlobalExceptionHandler.java"):       templates.GlobalExceptionHandler,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/controller/UserController.java"):           templates.UserController,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/controller/HealthController.java"):         templates.HealthControllerClass,
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
	baseDir := config.OutputDir

	content := utils.ReplacePlaceholders(templates.ApplicationMain, replacements)
	return utils.WriteFile(filepath.Join(baseDir, "starter/src/main/java", pkgPath, "Application.java"), content)
}

func extractProjectInfoFromPOM(projectRoot string) (*ProjectConfig, error) {
	// Read pom.xml from project root
	pomPath := filepath.Join(projectRoot, "pom.xml")
	content, err := os.ReadFile(pomPath)
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
		OutputDir:   projectRoot,
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

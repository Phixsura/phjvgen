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

	// Generate demo code by default
	utils.PrintInfo("生成完整CRUD示例代码...")
	if err := generateDemoCodeForProject(config); err != nil {
		return err
	}
	utils.PrintSuccess("完整CRUD示例代码生成完成")

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
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/assembler"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/interceptor"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/filter"),
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/config"),
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
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/listener"),
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
	// This function is now simplified as demo code generation will handle most files
	// We only generate the starter Application class here as it's always needed
	replacements := config.GetReplacements()
	baseDir := config.OutputDir
	pkgPath := config.PackagePath

	files := map[string]string{
		// Starter module - always needed
		filepath.Join("starter/src/main/java", pkgPath, "Application.java"): templates.ApplicationMain,
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

// generateDemoCodeForProject generates demo code as part of project generation
func generateDemoCodeForProject(config *ProjectConfig) error {
	replacements := config.GetReplacements()
	baseDir := config.OutputDir
	pkgPath := config.PackagePath

	// Generate Common layer
	commonFiles := map[string]string{
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/response/Result.java"):         templates.ResultClass,
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/exception/BusinessException.java"): templates.BusinessExceptionClass,
		filepath.Join(baseDir, "common/src/main/java", pkgPath, "common/constant/ErrorCode.java"):      templates.ErrorCodeClass,
	}

	// Generate Domain layer
	domainFiles := map[string]string{
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/model/User.java"):           templates.UserEntity,
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/repository/UserRepository.java"): templates.UserRepository,
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/service/UserDomainService.java"): templates.UserDomainService,
		filepath.Join(baseDir, "domain/src/main/java", pkgPath, "domain/event/UserCreatedEvent.java"): templates.UserCreatedEvent,
	}

	// Generate Infrastructure layer
	infrastructureFiles := map[string]string{
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/persistence/dataobject/UserDO.java"): templates.UserDO,
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/persistence/mapper/UserMapper.java"): templates.UserMapper,
		filepath.Join(baseDir, "infrastructure/src/main/java", pkgPath, "infrastructure/persistence/impl/UserRepositoryImpl.java"): templates.UserRepositoryImpl,
		filepath.Join(baseDir, "infrastructure/src/main/resources/db/migration/V1__create_user_table.sql"): templates.UserTableSQL,
	}

	// Generate Application layer
	applicationFiles := map[string]string{
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/dto/UserDTO.java"):              templates.UserDTO,
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/dto/CreateUserCommand.java"):    templates.CreateUserCommand,
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/dto/UpdateUserCommand.java"):    templates.UpdateUserCommand,
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/assembler/UserAssembler.java"): templates.UserAssembler,
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/service/UserService.java"):     templates.UserService,
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/executor/RegisterUserExecutor.java"): templates.RegisterUserExecutor,
		filepath.Join(baseDir, "application/application-user/src/main/java", pkgPath, "application/user/listener/UserEventListener.java"): templates.UserEventListener,
	}

	// Generate Adapter layer
	adapterFiles := map[string]string{
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/request/CreateUserRequest.java"):           templates.CreateUserRequest,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/request/UpdateUserRequest.java"):           templates.UpdateUserRequest,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/response/UserResponseVO.java"):             templates.UserResponseVO,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/assembler/UserControllerAssembler.java"):   templates.UserControllerAssembler,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/assembler/ResponseVOAssembler.java"):       templates.ResponseVOAssembler,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/filter/LoggingFilter.java"):                templates.LoggingFilter,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/interceptor/AuthInterceptor.java"):         templates.AuthInterceptor,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/config/WebMvcConfig.java"):                 templates.WebMvcConfig,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/advice/GlobalExceptionHandler.java"):       templates.GlobalExceptionHandler,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/controller/UserController.java"):           templates.UserController,
		filepath.Join(baseDir, "adapter/adapter-rest/src/main/java", pkgPath, "adapter/rest/controller/HealthController.java"):         templates.HealthControllerClass,
	}

	// Combine all files
	allFiles := make(map[string]string)
	for k, v := range commonFiles {
		allFiles[k] = v
	}
	for k, v := range domainFiles {
		allFiles[k] = v
	}
	for k, v := range infrastructureFiles {
		allFiles[k] = v
	}
	for k, v := range applicationFiles {
		allFiles[k] = v
	}
	for k, v := range adapterFiles {
		allFiles[k] = v
	}

	// Write all files
	for path, template := range allFiles {
		content := utils.ReplacePlaceholders(template, replacements)
		if err := utils.WriteFile(path, content); err != nil {
			return err
		}
	}

	return nil
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
	utils.PrintInfo("已生成完整的示例代码：")
	fmt.Println("  ✅ Application主类（启动类）")
	fmt.Println("  ✅ Common层：Result、BusinessException、ErrorCode")
	fmt.Println("  ✅ Domain层：User实体、UserRepository接口")
	fmt.Println("  ✅ Infrastructure层：UserDO、UserMapper、UserRepositoryImpl")
	fmt.Println("  ✅ Application层：UserDTO、UserService、UserAssembler")
	fmt.Println("  ✅ Adapter层：UserController、Request/Response、ExceptionHandler")
	fmt.Println("  ✅ 数据库脚本：V1__create_user_table.sql")
	fmt.Println()
	utils.PrintInfo("后续步骤:")
	fmt.Printf("  1. cd %s\n", config.OutputDir)
	fmt.Println("  2. 创建数据库并配置连接（starter/src/main/resources/application-dev.yml）")
	fmt.Println("  3. 执行数据库脚本：infrastructure/src/main/resources/db/migration/V1__create_user_table.sql")
	fmt.Println("  4. mvn clean install")
	fmt.Println("  5. java --enable-preview -jar starter/target/starter-*.jar")
	fmt.Println("  6. 测试健康检查: curl http://localhost:8080/api/health")
	fmt.Println("  7. 测试创建用户: curl -X POST http://localhost:8080/api/users -H 'Content-Type: application/json' -d '{\"username\":\"test\",\"email\":\"test@example.com\"}'")
	fmt.Println()
	utils.PrintInfo("添加新模块:")
	fmt.Println("  phjvgen add <模块名>")
	fmt.Println()
}

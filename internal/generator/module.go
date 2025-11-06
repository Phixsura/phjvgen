package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/phixia/phjvgen/internal/templates"
	"github.com/phixia/phjvgen/internal/utils"
)

// AddApplicationModule adds a new application module
func AddApplicationModule(moduleName string) error {
	// Validate module name
	if !validateModuleName(moduleName) {
		return fmt.Errorf("模块名称格式不正确，请使用小写字母和连字符")
	}

	// Check if we're in project root
	if !utils.FileExists("pom.xml") {
		return fmt.Errorf("未找到pom.xml文件，请在项目根目录运行此命令")
	}

	if !utils.DirExists("application") {
		return fmt.Errorf("未找到application目录，请确认项目结构正确")
	}

	// Extract project info
	config, err := extractProjectInfoFromPOM()
	if err != nil {
		return err
	}

	utils.PrintInfo("项目信息：")
	fmt.Printf("  Group ID: %s\n", config.GroupID)
	fmt.Printf("  Version: %s\n", config.Version)
	fmt.Printf("  Package: %s\n", config.PackageName)

	// Check if module already exists
	moduleDir := filepath.Join("application", "application-"+moduleName)
	if utils.DirExists(moduleDir) {
		return fmt.Errorf("模块 application-%s 已存在", moduleName)
	}

	fmt.Println()
	utils.PrintInfo(fmt.Sprintf("准备创建模块: application-%s", moduleName))
	confirm, err := utils.ReadInput("确认继续？(y/n): ")
	if err != nil {
		return err
	}
	if !strings.EqualFold(confirm, "y") && !strings.EqualFold(confirm, "yes") {
		utils.PrintWarning("已取消操作")
		return nil
	}

	// Create module structure
	utils.PrintInfo("创建模块目录结构...")
	if err := createModuleStructure(config, moduleName); err != nil {
		return err
	}
	utils.PrintSuccess("目录结构创建完成")

	// Generate module POM
	utils.PrintInfo("生成模块POM文件...")
	if err := generateModulePOM(config, moduleName); err != nil {
		return err
	}
	utils.PrintSuccess("模块POM文件生成完成")

	// Update parent POM
	utils.PrintInfo("更新父POM的modules声明...")
	if err := updateParentPOMModules(moduleName); err != nil {
		return err
	}
	utils.PrintSuccess("父POM modules更新完成")

	utils.PrintInfo("更新父POM的dependencyManagement...")
	if err := updateParentPOMDependencyManagement(config, moduleName); err != nil {
		return err
	}
	utils.PrintSuccess("父POM dependencyManagement更新完成")

	// Generate sample service
	utils.PrintInfo("生成示例Service类...")
	if err := generateSampleService(config, moduleName); err != nil {
		return err
	}
	utils.PrintSuccess("示例Service类生成完成")

	printModuleSummary(config, moduleName)
	return nil
}

func validateModuleName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9-]*$`, name)
	return matched
}

func createModuleStructure(config *ProjectConfig, moduleName string) error {
	moduleDir := filepath.Join("application", "application-"+moduleName)
	pkgPath := config.PackagePath

	// Convert module-name to module_name for package path (replace - with nothing for Java package)
	modulePackage := strings.ReplaceAll(moduleName, "-", "")

	dirs := []string{
		filepath.Join(moduleDir, "src/main/java", pkgPath, "application", modulePackage, "service"),
		filepath.Join(moduleDir, "src/main/java", pkgPath, "application", modulePackage, "dto"),
		filepath.Join(moduleDir, "src/main/java", pkgPath, "application", modulePackage, "assembler"),
		filepath.Join(moduleDir, "src/main/java", pkgPath, "application", modulePackage, "executor"),
		filepath.Join(moduleDir, "src/main/resources"),
		filepath.Join(moduleDir, "src/test/java", pkgPath, "application", modulePackage),
	}

	return utils.CreateDirs(dirs...)
}

func generateModulePOM(config *ProjectConfig, moduleName string) error {
	// Convert module-name to Module Name for description
	moduleDescription := strings.ReplaceAll(moduleName, "-", " ")

	replacements := config.GetReplacements()
	replacements["{{MODULE_NAME}}"] = moduleName
	replacements["{{MODULE_DESCRIPTION}}"] = moduleDescription

	content := utils.ReplacePlaceholders(templates.ApplicationModulePOM, replacements)
	moduleDir := filepath.Join("application", "application-"+moduleName)
	return utils.WriteFile(filepath.Join(moduleDir, "pom.xml"), content)
}

func updateParentPOMModules(moduleName string) error {
	content, err := os.ReadFile("pom.xml")
	if err != nil {
		return err
	}

	pomContent := string(content)

	// Check if module already exists
	moduleEntry := fmt.Sprintf("<module>application/application-%s</module>", moduleName)
	if strings.Contains(pomContent, moduleEntry) {
		utils.PrintWarning("模块已在父POM的modules中声明")
		return nil
	}

	// Find </modules> and insert before it
	modulesEnd := strings.Index(pomContent, "</modules>")
	if modulesEnd == -1 {
		return fmt.Errorf("could not find </modules> tag in parent pom.xml")
	}

	// Insert the new module
	newModule := fmt.Sprintf("        <module>application/application-%s</module>\n    ", moduleName)
	newContent := pomContent[:modulesEnd] + newModule + pomContent[modulesEnd:]

	return os.WriteFile("pom.xml", []byte(newContent), 0644)
}

func updateParentPOMDependencyManagement(config *ProjectConfig, moduleName string) error {
	content, err := os.ReadFile("pom.xml")
	if err != nil {
		return err
	}

	pomContent := string(content)

	// Check if dependency already exists
	artifactEntry := fmt.Sprintf("<artifactId>application-%s</artifactId>", moduleName)
	if strings.Contains(pomContent, artifactEntry) {
		utils.PrintWarning("模块依赖已在父POM的dependencyManagement中声明")
		return nil
	}

	// Find application-user dependency and add after it
	userDepStart := strings.Index(pomContent, "<artifactId>application-user</artifactId>")
	if userDepStart == -1 {
		return fmt.Errorf("could not find application-user dependency in parent pom.xml")
	}

	// Find the closing </dependency> tag after application-user
	searchStart := userDepStart
	depEnd := strings.Index(pomContent[searchStart:], "</dependency>")
	if depEnd == -1 {
		return fmt.Errorf("could not find closing </dependency> tag")
	}
	depEnd += searchStart + len("</dependency>")

	// Insert new dependency
	newDep := fmt.Sprintf(`
            <dependency>
                <groupId>%s</groupId>
                <artifactId>application-%s</artifactId>
                <version>${project.version}</version>
            </dependency>`, config.GroupID, moduleName)

	newContent := pomContent[:depEnd] + newDep + pomContent[depEnd:]

	return os.WriteFile("pom.xml", []byte(newContent), 0644)
}

func generateSampleService(config *ProjectConfig, moduleName string) error {
	// Convert module-name to ModuleName (CamelCase)
	className := toCamelCase(moduleName)

	// Convert module-name to modulepackage (remove dashes for package)
	modulePackage := strings.ReplaceAll(moduleName, "-", "")

	serviceTemplate := fmt.Sprintf(`package %s.application.%s.service;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

/**
 * %s业务服务
 */
@Slf4j
@Service
public class %sService {

    /**
     * 示例方法
     */
    public String execute() {
        log.info("Executing %sService");
        return "%s service executed successfully";
    }
}
`, config.PackageName, modulePackage, className, className, className, className)

	moduleDir := filepath.Join("application", "application-"+moduleName)
	pkgPath := config.PackagePath
	servicePath := filepath.Join(moduleDir, "src/main/java", pkgPath, "application", modulePackage, "service", className+"Service.java")

	return utils.WriteFile(servicePath, serviceTemplate)
}

func toCamelCase(input string) string {
	// Split by dash and capitalize each part
	parts := strings.Split(input, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

func printModuleSummary(config *ProjectConfig, moduleName string) {
	fmt.Println()
	utils.PrintSuccess("==========================================")
	utils.PrintSuccess(fmt.Sprintf("模块 application-%s 创建完成！", moduleName))
	utils.PrintSuccess("==========================================")
	fmt.Println()
	utils.PrintInfo("下一步：")
	fmt.Printf("  1. 查看生成的文件: ls -la application/application-%s/\n", moduleName)
	fmt.Println("  2. 重新构建项目: mvn clean install")
	fmt.Println("  3. 开始开发业务逻辑")
	fmt.Println()
	utils.PrintInfo("如需在adapter-rest中使用此模块，请手动添加依赖：")
	fmt.Println("  <dependency>")
	fmt.Printf("      <groupId>%s</groupId>\n", config.GroupID)
	fmt.Printf("      <artifactId>application-%s</artifactId>\n", moduleName)
	fmt.Println("  </dependency>")
	fmt.Println()
}

package generator

import (
	"fmt"
	"strings"

	"github.com/phixia/phjvgen/internal/utils"
)

// ProjectConfig holds the project configuration
type ProjectConfig struct {
	GroupID            string
	ArtifactID         string
	Version            string
	ProjectName        string
	ProjectDescription string
	PackageName        string
	PackagePath        string
	OutputDir          string
}

// GetProjectConfig collects project configuration from user input
func GetProjectConfig() (*ProjectConfig, error) {
	utils.PrintBanner()
	fmt.Println()

	// Group ID
	groupID, err := utils.ReadValidatedInput(
		"请输入 Group ID (例如: com.mycompany): ",
		`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$`,
		"Group ID 格式不正确，请使用类似 com.mycompany 的格式",
	)
	if err != nil {
		return nil, err
	}

	// Artifact ID
	artifactID, err := utils.ReadValidatedInput(
		"请输入 Artifact ID (例如: my-app): ",
		`^[a-z][a-z0-9-]*$`,
		"Artifact ID 格式不正确，请使用小写字母和连字符",
	)
	if err != nil {
		return nil, err
	}

	// Version
	version := utils.ReadInputWithDefault(
		"请输入版本号 (默认: 1.0.0): ",
		"1.0.0",
	)

	// Project Name
	projectName := utils.ReadInputWithDefault(
		fmt.Sprintf("请输入项目名称 (默认: %s): ", artifactID),
		artifactID,
	)

	// Project Description
	projectDescription := utils.ReadInputWithDefault(
		"请输入项目描述 (默认: A Java 25 LTS Project): ",
		"A Java 25 LTS Project",
	)

	// Package name and path
	packageName := groupID
	packagePath := strings.ReplaceAll(groupID, ".", "/")

	// Output directory
	outputDir := utils.ReadInputWithDefault(
		fmt.Sprintf("请输入输出目录 (默认: ./%s): ", artifactID),
		"./"+artifactID,
	)

	config := &ProjectConfig{
		GroupID:            groupID,
		ArtifactID:         artifactID,
		Version:            version,
		ProjectName:        projectName,
		ProjectDescription: projectDescription,
		PackageName:        packageName,
		PackagePath:        packagePath,
		OutputDir:          outputDir,
	}

	// Display configuration
	fmt.Println()
	utils.PrintInfo("项目配置信息：")
	fmt.Printf("  Group ID: %s\n", config.GroupID)
	fmt.Printf("  Artifact ID: %s\n", config.ArtifactID)
	fmt.Printf("  Version: %s\n", config.Version)
	fmt.Printf("  Project Name: %s\n", config.ProjectName)
	fmt.Printf("  Package Name: %s\n", config.PackageName)
	fmt.Printf("  Output Directory: %s\n", config.OutputDir)
	fmt.Println()

	return config, nil
}

// GetReplacements returns a map for template placeholder replacement
func (c *ProjectConfig) GetReplacements() map[string]string {
	return map[string]string{
		"{{GROUP_ID}}":            c.GroupID,
		"{{ARTIFACT_ID}}":         c.ArtifactID,
		"{{VERSION}}":             c.Version,
		"{{PROJECT_NAME}}":        c.ProjectName,
		"{{PROJECT_DESCRIPTION}}": c.ProjectDescription,
		"{{PACKAGE_NAME}}":        c.PackageName,
		"{{PACKAGE_PATH}}":        c.PackagePath,
	}
}

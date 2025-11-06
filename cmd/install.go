package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/phixia/phjvgen/internal/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装 phjvgen 到系统",
	Long: `将 phjvgen 安装到系统路径，使其可以在任何位置使用。

安装位置：
  - Linux/macOS (普通用户): ~/.local/bin/phjvgen
  - Linux/macOS (root): /usr/local/bin/phjvgen
  - Windows: %USERPROFILE%\AppData\Local\phjvgen\phjvgen.exe

安装后需要确保安装目录在 PATH 环境变量中。

注意：此命令需要从 phjvgen 的构建目录或源码目录运行。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return installPhjvgen()
	},
}

func installPhjvgen() error {
	utils.PrintBanner()
	fmt.Println()

	// Determine install directory
	installDir, err := getInstallDir()
	if err != nil {
		return err
	}

	utils.PrintInfo(fmt.Sprintf("安装目标: %s", installDir))

	// Create install directory if it doesn't exist
	if err := os.MkdirAll(installDir, 0755); err != nil {
		utils.PrintError(fmt.Sprintf("无法创建目录 %s: %v", installDir, err))
		utils.PrintInfo("请尝试使用 sudo 运行此命令")
		return err
	}

	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取当前可执行文件路径: %w", err)
	}

	// Determine target executable name
	targetName := "phjvgen"
	if runtime.GOOS == "windows" {
		targetName = "phjvgen.exe"
	}
	targetPath := filepath.Join(installDir, targetName)

	// Copy executable
	utils.PrintInfo(fmt.Sprintf("复制 phjvgen 到 %s...", installDir))
	if err := copyFile(execPath, targetPath); err != nil {
		utils.PrintError(fmt.Sprintf("复制失败: %v", err))
		return err
	}

	// Make executable (Unix-like systems)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(targetPath, 0755); err != nil {
			return fmt.Errorf("设置执行权限失败: %w", err)
		}
	}

	utils.PrintSuccess("安装完成！")
	fmt.Println()

	// Check if install dir is in PATH
	checkAndPrintPathInstructions(installDir)

	// Display usage
	fmt.Println()
	utils.PrintInfo("使用方法:")
	fmt.Println("  phjvgen generate    - 生成新项目")
	fmt.Println("  phjvgen demo        - 生成CRUD示例")
	fmt.Println("  phjvgen example     - 生成示例项目")
	fmt.Println("  phjvgen add <name>  - 添加新模块")
	fmt.Println()

	// Test command
	utils.PrintInfo("测试命令...")
	if err := testInstallation(targetPath); err != nil {
		utils.PrintWarning("phjvgen 命令暂时不可用，请重新加载 shell 配置")
	} else {
		utils.PrintSuccess("phjvgen 已就绪！")
	}

	return nil
}

func getInstallDir() (string, error) {
	if os.Geteuid() == 0 {
		// Running as root
		utils.PrintWarning("检测到以 root 用户运行，将安装到 /usr/local/bin")
		return "/usr/local/bin", nil
	}

	// Regular user
	if runtime.GOOS == "windows" {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			return "", fmt.Errorf("无法获取 USERPROFILE 环境变量")
		}
		return filepath.Join(userProfile, "AppData", "Local", "phjvgen"), nil
	}

	// Unix-like systems
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("无法获取 HOME 环境变量")
	}
	return filepath.Join(home, ".local", "bin"), nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0755)
}

func checkAndPrintPathInstructions(installDir string) {
	path := os.Getenv("PATH")

	// Check if installDir is in PATH
	pathDirs := filepath.SplitList(path)
	inPath := false
	for _, dir := range pathDirs {
		if dir == installDir {
			inPath = true
			break
		}
	}

	if inPath {
		utils.PrintSuccess(fmt.Sprintf("%s 已在 PATH 中", installDir))
	} else {
		utils.PrintWarning(fmt.Sprintf("%s 不在 PATH 中", installDir))
		fmt.Println()
		utils.PrintInfo("请将以下行添加到你的 shell 配置文件中：")
		fmt.Println()

		if runtime.GOOS == "windows" {
			fmt.Printf("  将 %s 添加到系统 PATH 环境变量\n", installDir)
			fmt.Println("  或在 PowerShell 中运行:")
			fmt.Printf("  $env:PATH += \";%s\"\n", installDir)
		} else {
			home := os.Getenv("HOME")
			zshrc := filepath.Join(home, ".zshrc")
			bashrc := filepath.Join(home, ".bashrc")

			if utils.FileExists(zshrc) {
				fmt.Printf("  echo 'export PATH=\"$HOME/.local/bin:$PATH\"' >> ~/.zshrc\n")
				fmt.Println("  source ~/.zshrc")
			} else if utils.FileExists(bashrc) {
				fmt.Printf("  echo 'export PATH=\"$HOME/.local/bin:$PATH\"' >> ~/.bashrc\n")
				fmt.Println("  source ~/.bashrc")
			} else {
				fmt.Printf("  export PATH=\"%s:$PATH\"\n", installDir)
			}
		}
		fmt.Println()
	}
}

func testInstallation(targetPath string) error {
	cmd := exec.Command(targetPath, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
}

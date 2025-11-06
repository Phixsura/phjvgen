package utils

import "github.com/fatih/color"

var (
	InfoColor    = color.New(color.FgBlue).SprintFunc()
	SuccessColor = color.New(color.FgGreen).SprintFunc()
	ErrorColor   = color.New(color.FgRed).SprintFunc()
	WarningColor = color.New(color.FgYellow).SprintFunc()
)

func PrintInfo(msg string) {
	color.Blue("[INFO] %s", msg)
}

func PrintSuccess(msg string) {
	color.Green("[SUCCESS] %s", msg)
}

func PrintError(msg string) {
	color.Red("[ERROR] %s", msg)
}

func PrintWarning(msg string) {
	color.Yellow("[WARNING] %s", msg)
}

func PrintBanner() {
	color.Blue(`╔═══════════════════════════════════════════════╗
║   Java 25 LTS 分层架构项目生成器              ║
║   Project Generator based on Java 25 LTS      ║
╚═══════════════════════════════════════════════╝`)
}

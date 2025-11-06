package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 phjvgen 的版本信息。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("phjvgen version %s\n", Version)
		fmt.Println("Java 25 LTS Project Generator")
		fmt.Println("Copyright (c) 2025")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

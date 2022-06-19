package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gui bool
)

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "",
		Run: func(cmd *cobra.Command, args []string) {
			runApp()
		},
	}
	rootCmd.Flags().BoolVarP(&gui, "gui", "", false, "is gui?")

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("启动环境...")
			startEnv()
		},
	}
	// startCmd.Flags().BoolVarP(&daemon, "deamon", "d", false, "is daemon?")
	rootCmd.AddCommand(startCmd)

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("停止环境...")
			stopEnv()
		},
	}
	rootCmd.AddCommand(stopCmd)

	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "restart PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("重启环境...")
			stopEnv()
		},
	}
	rootCmd.AddCommand(restartCmd)

	var saveCmd = &cobra.Command{
		Use:   "save",
		Short: "Save Current Config",
		Run: func(cmd *cobra.Command, args []string) {
			saveConfig()
		},
	}
	rootCmd.AddCommand(saveCmd)

	var exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export Current Config",
		Run: func(cmd *cobra.Command, args []string) {
			exportNginx()
		},
	}
	rootCmd.AddCommand(exportCmd)

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check the status of NGINX and PHP",
		Run: func(cmd *cobra.Command, args []string) {
			nginxStatus, phpStatus := checkStatus()
			showStatus(nginxStatus, phpStatus)
		},
	}
	rootCmd.AddCommand(statusCmd)

	return rootCmd
}

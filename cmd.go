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
			reloadCfg()
			runApp()
		},
	}
	rootCmd.Flags().BoolVarP(&gui, "gui", "g", false, "is gui?")
	rootCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			reloadCfg()
			fmt.Println("启动环境...")
			startEnv()
		},
	}
	startCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	rootCmd.AddCommand(startCmd)

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			// 进程直接杀死就好，不需要用到配置的
			fmt.Println("终止环境...")
			stopEnv()
		},
	}
	rootCmd.AddCommand(stopCmd)

	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "restart PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			reloadCfg()
			fmt.Println("重启环境...")
			stopEnv()
		},
	}
	restartCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	rootCmd.AddCommand(restartCmd)

	var saveCmd = &cobra.Command{
		Use:   "save",
		Short: "Save Current Config",
		Run: func(cmd *cobra.Command, args []string) {
			reloadCfg()
			saveConfig()
		},
	}
	saveCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	rootCmd.AddCommand(saveCmd)

	var exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export Current Config",
		Run: func(cmd *cobra.Command, args []string) {
			reloadCfg()
			exportNginx()
		},
	}
	exportCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	rootCmd.AddCommand(exportCmd)

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check the status of NGINX and PHP",
		Run: func(cmd *cobra.Command, args []string) {
			// 进程直接查看就好，不需要调用配置信息的
			nginxStatus, phpStatus := checkStatus()
			showStatus(nginxStatus, phpStatus)
		},
	}
	rootCmd.AddCommand(statusCmd)

	return rootCmd
}

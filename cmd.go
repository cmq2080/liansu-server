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
	rootCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	if isWindows() {
		rootCmd.Flags().BoolVarP(&gui, "gui", "g", false, "run with gui?")
	}

	/************启动环境************/
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			reloadCfg()
			fmt.Println("启动环境...")
			startEnv()
		},
	}
	startCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	rootCmd.AddCommand(startCmd)

	/************终止环境************/
	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			// Windows下进程直接杀死就好，不需要用到配置的
			if isLinux() {
				reloadCfg()
			}
			fmt.Println("终止环境...")
			stopEnv()
		},
	}
	rootCmd.AddCommand(stopCmd)

	/************重启环境************/
	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart PHP Running Environment",
		Run: func(cmd *cobra.Command, args []string) {
			reloadCfg()
			fmt.Println("重启环境...")
			restartEnv()
		},
	}
	restartCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	rootCmd.AddCommand(restartCmd)

	/************保存配置************/
	// var saveCmd = &cobra.Command{
	// 	Use:   "save",
	// 	Short: "Save Current Config",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		reloadCfg()
	// 		saveConfig()
	// 	},
	// }
	// saveCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	// rootCmd.AddCommand(saveCmd)

	/************导出配置************/
	var exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Generate New Server Config Files By Current Config",
		Run: func(cmd *cobra.Command, args []string) {
			reloadCfg()
			exportToLsNginxConf()
		},
	}
	exportCmd.Flags().StringVarP(&cfgFilepath, "config", "c", "", "custom config file")
	rootCmd.AddCommand(exportCmd)

	/************查看状态************/
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

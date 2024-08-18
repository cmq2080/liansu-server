//go:build linux
// +build linux

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

/*
 * @func showConfig
 * @func startEnv
 * @func stopEnv
 * @func restartEnv
 * @func checkStatus
 */

func showConfig() {
	var tab string = "\t\t"
	fmt.Println("配置信息\n********************")
	fmt.Println("[Server]")
	fmt.Println("host" + tab + cfg.Server.Host)
	fmt.Println("port" + tab + cfg.Server.Port)
	fmt.Println("wwwroot" + tab + cfg.Server.Wwwroot)

	fmt.Printf("\n")
	fmt.Println("[Nginx]")
	// 当是Linux环境下时，用不着关心这些东西，直接用shell命令就行
	// fmt.Println("dir" + tab + cfg.Nginx.Dir)
	// fmt.Println("conf" + tab + cfg.Nginx.Conf)
	// fmt.Println("logs_dir" + "\t" + cfg.Nginx.Logs_dir)
	if isLinux() {
		fmt.Println("cmd")
		fmt.Println("  " + "start" + tab + cfg.Nginx.Cmd.Start)
		fmt.Println("  " + "stop" + tab + cfg.Nginx.Cmd.Stop)
		fmt.Println("  " + "restart" + "\t" + cfg.Nginx.Cmd.Restart)
	}

	fmt.Printf("\n")
	fmt.Println("[PHP]")
	// fmt.Println("dir" + tab + cfg.Php.Dir)
	// fmt.Println("ini" + tab + cfg.Php.Ini)
	if isLinux() {
		fmt.Println("cmd")
		fmt.Println("  " + "start" + tab + cfg.Nginx.Cmd.Start)
		fmt.Println("  " + "stop" + tab + cfg.Nginx.Cmd.Stop)
		fmt.Println("  " + "restart" + "\t" + cfg.Nginx.Cmd.Restart)
	}

	fmt.Printf("\n")
	fmt.Println("[Apps]")
	for i := 0; i < len(cfg.Apps); i++ {
		app := cfg.Apps[i]
		fmt.Println("name" + tab + app.Name)
		fmt.Println("root" + tab + app.Root)
		if app.Server_name != "" {
			fmt.Println("server_name" + "\t" + app.Server_name)
			fmt.Println("listen" + tab + app.Listen)
			fmt.Println("logs_dir" + "\t" + app.Logs_dir)
		}
		fmt.Printf("----------------\n")
	}
}

func startEnv() {
	linkToNginxConf()

	var cmdStr string
	var args []string
	// 启动php-fpm
	cmdStr, args = parseCmd(cfg.Php.Cmd.Start)

	cmd := exec.Command(cmdStr, args...)
	cmd.Output()

	// 启动nginx
	cmdStr, args = parseCmd(cfg.Nginx.Cmd.Start)

	cmd = exec.Command(cmdStr, args...)
	cmd.Output()
}

func stopEnv() {
	var cmdStr string
	var args []string
	// 终止php-fpm
	cmdStr, args = parseCmd(cfg.Php.Cmd.Stop)

	cmd := exec.Command(cmdStr, args...)
	cmd.Output()

	// 终止nginx
	cmdStr, args = parseCmd(cfg.Nginx.Cmd.Stop)

	cmd = exec.Command(cmdStr, args...)
	cmd.Output()
}

func restartEnv() {
	var cmdStr string
	var args []string

	if cfg.Php.Cmd.Restart != "" && cfg.Nginx.Cmd.Restart != "" {
		// 重新启动php-fpm
		cmdStr, args = parseCmd(cfg.Php.Cmd.Restart)

		cmd := exec.Command(cmdStr, args...)
		cmd.Output()

		// 重新启动nginx
		cmdStr, args = parseCmd(cfg.Nginx.Cmd.Restart)

		cmd = exec.Command(cmdStr, args...)
		cmd.Output()
	} else {
		stopEnv()
		startEnv()
	}
}

func parseCmd(cmdStr string) (string, []string) {
	cmdArr := strings.Split(cmdStr, " ")
	cmdStr = cmdArr[0]
	args := cmdArr[1:]

	return cmdStr, args
}

func checkStatus() (int, int) {
	var nginxStatus int
	var phpStatus int

	if isWindows() {
		cmd := exec.Command("top", []string{"-b", "-n", "1"}...)
		outBuf, err := cmd.Output()
		if err != nil {
			return 0, 0
		}

		outStr := string(outBuf)

		if strings.Contains(outStr, "nginx") {
			nginxStatus = 1
		}
		if strings.Contains(outStr, "php-fpm") {
			phpStatus = 1
		}
	}

	return nginxStatus, phpStatus
}

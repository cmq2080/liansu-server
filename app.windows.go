//go:build windows
// +build windows

package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
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
	fmt.Println("dir" + tab + cfg.Nginx.Dir)
	fmt.Println("conf" + tab + cfg.Nginx.Conf)
	fmt.Println("logs_dir" + "\t" + cfg.Nginx.Logs_dir)

	fmt.Printf("\n")
	fmt.Println("[PHP]")
	fmt.Println("dir" + tab + cfg.Php.Dir)
	fmt.Println("ini" + tab + cfg.Php.Ini)

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

	// 开启php-cgi
	subCmdStr := "" + cfg.Php.Dir + "/php-cgi.exe -c " + cfg.Php.Ini + "" // 这里已经不用再加引号
	// fmt.Println(thisDir + "/php-cgi-spawner.exe")
	// fmt.Println(subCmdStr)
	// fmt.Println(cfg.Php.Port)
	// fmt.Println("4+16")
	processNum := runtime.NumCPU()
	workerNum := processNum * 4
	cmd := exec.Command(thisDir+"/php-cgi-spawner.exe", []string{subCmdStr, cfg.Php.Port, strconv.Itoa(processNum) + "+" + strconv.Itoa(workerNum)}...)
	// fmt.Println(cmd.Args)
	cmd.Start()

	// 开启nginx
	cmd = exec.Command(cfg.Nginx.Dir+"/nginx.exe", []string{"-p", cfg.Nginx.Dir}...)
	// fmt.Println(cmd.Args)
	cmd.Start()
}

func stopEnv() {
	cmd := exec.Command("taskkill.exe", "/F", "/IM", "nginx.exe")
	cmd.Output()

	cmd = exec.Command("taskkill.exe", "/F", "/IM", "php-cgi-spawner.exe")
	cmd.Output()

	cmd = exec.Command("taskkill.exe", "/F", "/IM", "php-cgi.exe")
	cmd.Output()
}

func restartEnv() {
	stopEnv()
	startEnv()
}

func checkStatus() (int, int) {
	var nginxStatus int
	var phpStatus int

	if isWindows() {
		cmd := exec.Command("tasklist.exe")
		outBuf, err := cmd.Output()
		if err != nil {
			return 0, 0
		}

		outStr := string(outBuf)

		if strings.Contains(outStr, "nginx.exe") {
			nginxStatus = 1
		}
		if strings.Contains(outStr, "php-cgi.exe") {
			phpStatus = 1
		}
	}

	return nginxStatus, phpStatus
}

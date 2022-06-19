package main

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
)

type Config struct {
	Server ServerConfig `json:"server"`
	Nginx  NginxConfig  `json:"nginx"`
	Php    PHPConfig    `json:"php"`
	Apps   []AppConfig  `json:"apps"`
}

type ServerConfig struct {
	Ip      string `json:"ip"`
	Port    string `json:"port"`
	Wwwroot string `json:"wwwroot"`
}

type NginxConfig struct {
	Dir      string `json:"dir"`
	Conf     string `json:"conf"`
	Logs_dir string `json:"logs_dir"`
}

type PHPConfig struct {
	Dir  string `json:"dir"`
	Ini  string `json:"ini"`
	Port string `json:"port"`
}

type AppConfig struct {
	Name        string `json:"name"`
	Root        string `json:"root"`
	Server_name string `json:"server_name"`
	Listen      string `json:"listen"`
	Logs_dir    string `json:"logs_dir"`
}

func (cfg *Config) Init(filepath string) {
	buf, _ := os.ReadFile(filepath)
	var res map[string]interface{}

	json.Unmarshal(buf, &res)

	serverRes := res["server"].(map[string]interface{})
	serverCfg := ServerConfig{
		Ip:      serverRes["ip"].(string),
		Port:    serverRes["port"].(string),
		Wwwroot: getRealDir(serverRes["wwwroot"].(string)),
	}
	cfg.Server = serverCfg

	nginxRes := res["nginx"].(map[string]interface{})
	nginxCfg := NginxConfig{
		Dir:      getRealDir(nginxRes["dir"].(string)),
		Conf:     getRealDir(nginxRes["conf"].(string)),
		Logs_dir: getRealDir(nginxRes["logs_dir"].(string)),
	}
	if nginxCfg.Conf == "" {
		nginxCfg.Conf = nginxCfg.Dir + "/conf/nginx.conf"
	}
	if nginxCfg.Logs_dir == "" {
		nginxCfg.Logs_dir = nginxCfg.Dir + "/logs"
	}
	cfg.Nginx = nginxCfg

	phpRes := res["php"].(map[string]interface{})
	phpCfg := PHPConfig{
		Dir:  getRealDir(phpRes["dir"].(string)),
		Ini:  getRealDir(phpRes["ini"].(string)),
		Port: phpRes["port"].(string),
	}
	if phpCfg.Ini == "" {
		phpCfg.Ini = phpCfg.Dir + "/php.ini"
	}
	if phpCfg.Port == "" {
		phpCfg.Port = "9000"
	}
	cfg.Php = phpCfg

	appRes := res["apps"].([]interface{})
	appCfgs := []AppConfig{}
	for i := 0; i < len(appRes); i++ {
		appRes2 := appRes[i].(map[string]interface{})
		appCfg := AppConfig{
			Name:        appRes2["name"].(string),
			Root:        getRealDir(appRes2["root"].(string)),
			Server_name: appRes2["server_name"].(string),
			Listen:      appRes2["listen"].(string),
			Logs_dir:    getRealDir(appRes2["logs_dir"].(string)),
		}
		if !isAbsoluteDir(appCfg.Root) { // 不是绝对路径？那前面得加wwwroot
			appCfg.Root = serverCfg.Wwwroot + "/" + appCfg.Root
		}
		if appCfg.Listen == "" && appCfg.Server_name != "" {
			appCfg.Listen = serverCfg.Port
		}
		if appCfg.Logs_dir == "" && appCfg.Server_name != "" {
			appCfg.Logs_dir = nginxCfg.Logs_dir
		}

		appCfgs = append(appCfgs, appCfg)
	}
	cfg.Apps = appCfgs
}

func getRealDir(dir string) string {
	dir = strings.Replace(dir, "\\", "/", -1)
	dir = strings.TrimRight(dir, "/")

	return dir
}

func isAbsoluteDir(dir string) bool {
	if runtime.GOOS == "windows" {
		return strings.Contains(dir, ":")
	} else {
		return dir[0:1] == "/"
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	APP_NAME    = "liansu-server"
	APP_NAME_ZH = "liansu-server服务器控制端"
)

var (
	thisExe     string
	thisDir     string
	thisConfDir string

	cfgFilepath string
	cfg         *Config
)

func main() {
	initialize()

	getRootCmd().Execute()

	// fmt.Println(thisExe, thisDir)
}

func initialize() {
	thisExe, _ = os.Executable()
	thisDir = getRealDir(filepath.Dir(thisExe))
	thisConfDir = getRealDir(thisDir + "/conf")

	// 初始化配置
	// reloadCfg()
}

func reloadCfg() {
	if cfgFilepath == "" {
		cfgFilepath = thisDir + "/config.json"
	}

	cfg = &Config{}
	cfg.Init(cfgFilepath)
}

func runApp() {
	if gui {
		runGui()
	} else {
		showConfig()
	}
}

func showConfig() {
	var tab string = "\t\t"
	fmt.Println("配置信息\n********************")
	fmt.Println("[Server]")
	fmt.Println("ip" + tab + cfg.Server.Ip)
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
	confFilepath := cfg.Nginx.Conf

	// 获取nginx.conf文件中的内容
	buf, err := os.ReadFile(confFilepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	confText := string(buf)

	// 向nginx.conf中塞入ls-nginx.conf的文件包含
	lsConfFilepath := thisConfDir + "/nginx.conf"
	if !strings.Contains(confText, lsConfFilepath) {
		index := strings.LastIndex(confText, "}")
		confText := confText[:index]
		confText += "    include " + lsConfFilepath + ";\n}"
		// fmt.Println(confText)
		os.WriteFile(confFilepath, []byte(confText), os.ModePerm)
	}

	// 开启php-cgi
	subCmdStr := "" + cfg.Php.Dir + "/php-cgi.exe -c " + cfg.Php.Ini + "" // 这里已经不用再加引号
	// fmt.Println(thisDir + "/php-cgi-spawner.exe")
	// fmt.Println(subCmdStr)
	// fmt.Println(cfg.Php.Port)
	// fmt.Println("4+16")
	cmd := exec.Command(thisDir+"/php-cgi-spawner.exe", []string{subCmdStr, cfg.Php.Port, "4+16"}...)
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

func saveConfig() {
	buf, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(thisDir+"/config.json", buf, os.ModePerm); err != nil {
		panic(err)
	}

}

func exportNginx() {
	confFilepath := cfg.Nginx.Conf

	// 获取nginx.conf文件中的内容
	buf, err := os.ReadFile(confFilepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	confText := string(buf)

	// 预设配置
	preConf := map[string]string{"resolver": "8.8.8.8 ipv6=off"}
	if !strings.Contains(confText, "server_names_hash_bucket_size") {
		preConf["server_names_hash_bucket_size"] = "128"
	}
	if !strings.Contains(confText, "client_max_body_size") {
		preConf["client_max_body_size"] = "2048m"
	}

	lineSpliter := ";\n"
	lsConfText := "# set pre-configuration variables" + lineSpliter
	for k, v := range preConf {
		lsConfText += k + " " + v + lineSpliter
	}

	// 设置app server
	lsConfText += "\n\n"
	lsConfText += "# set common server" + lineSpliter

	lsConfText += "server {\n"
	lsConfText += "\tlisten " + cfg.Server.Port + lineSpliter
	lsConfText += "\tserver_name " + cfg.Server.Ip + lineSpliter + "\n"

	lsConfText += "\tset $www_root " + cfg.Server.Wwwroot + lineSpliter
	lsConfText += "\tset $php_port " + cfg.Server.Ip + ":" + cfg.Php.Port + lineSpliter + "\n"

	lsConfText += "\tset $app_name 0" + lineSpliter
	lsConfText += "\tset $root 0" + lineSpliter
	lsConfText += "\tset $app_default_index 0" + lineSpliter + "\n"

	lsConfText += "\tif ($uri ~* ^/app/([^\\" + "/]+)/$) {\n"
	lsConfText += "\t\tset $app_name $1" + lineSpliter
	lsConfText += "\t\trewrite ^/app/([^\\" + "/]+)/$ /index.php" + lineSpliter
	lsConfText += "\t}\n\n"

	lsConfText += "\tif ($uri ~* ^/app/([^\\" + "/]+)/(.*)$) {\n"
	lsConfText += "\t\tset $app_name $1" + lineSpliter
	lsConfText += "\t\tset $app_default_index $2" + lineSpliter
	lsConfText += "\t\trewrite ^/app/([^\\" + "/]+)/(.*)$ /$app_default_index" + lineSpliter
	lsConfText += "\t}\n\n"

	lsConfText += "\t## match app_name START" + lineSpliter
	for i := 0; i < len(cfg.Apps); i++ {
		lsConfText += "\tif ($app_name ~* ^" + cfg.Apps[i].Name + "$) {\n"
		lsConfText += "\t\tset $root " + cfg.Apps[i].Root + lineSpliter
		lsConfText += "\t}\n\n"
	}
	lsConfText += "\t## match app_name END" + lineSpliter + "\n"

	// lsConfText += "\tif ($root = 0) {\n"
	// lsConfText += "\t\tbreak" + lineSpliter
	// lsConfText += "\t}\n\n"

	lsConfText += "\tif ($root = 0) {\n"
	lsConfText += "\t}\n\n"

	lsConfText += "\tlocation / {\n"
	lsConfText += "\t\tif (!-e $request_filename) {\n"
	lsConfText += "\t\t\trewrite ^(.*)$ /index.php?s=$1 last" + lineSpliter
	lsConfText += "\t\t\tbreak" + lineSpliter
	lsConfText += "\t\t}\n"
	lsConfText += "\t}\n\n"

	lsConfText += "\tindex index.php index.html index.htm" + lineSpliter
	lsConfText += "\troot $root" + lineSpliter
	lsConfText += "\taccess_log " + cfg.Nginx.Logs_dir + "/nginx.access.log" + lineSpliter
	lsConfText += "\terror_log " + cfg.Nginx.Logs_dir + "/nginx.error.log" + lineSpliter

	lsConfText += "\tlocation ~ \\.php(.*)$ {\n"
	lsConfText += "\t\tfastcgi_pass $php_port" + lineSpliter
	lsConfText += "\t\tfastcgi_index index.php" + lineSpliter + "\n"
	lsConfText += "\t\tinclude fastcgi_params" + lineSpliter
	lsConfText += "\t\tset $real_script_name $fastcgi_script_name" + lineSpliter
	lsConfText += "\t\tset $real_path_info $fastcgi_path_info" + lineSpliter
	lsConfText += "\t\tif ($fastcgi_script_name ~ ^(.+?\\.php)(/.+)$) {\n"
	lsConfText += "\t\t\tset $real_script_name $1" + lineSpliter
	lsConfText += "\t\t\tset $real_path_info $2" + lineSpliter
	lsConfText += "\t\t}\n\n"
	lsConfText += "\t\tfastcgi_param QUERY_STRING $query_string" + lineSpliter
	lsConfText += "\t\tfastcgi_param SCRIPT_FILENAME $document_root$real_script_name" + lineSpliter
	lsConfText += "\t\tfastcgi_param APP_JT $app_name" + lineSpliter
	lsConfText += "\t\tfastcgi_param SCRIPT_NAME $real_script_name" + lineSpliter
	lsConfText += "\t\tfastcgi_param PATH_INFO $real_path_info" + lineSpliter
	lsConfText += "\t}\n\n"

	lsConfText += "\tlocation ~* .*\\.(gif|jpg|jpeg|png|bmp|swf|flv|ico)$ {\n"
	lsConfText += "\t\texpires 30d" + lineSpliter
	lsConfText += "\t\taccess_log off" + lineSpliter
	lsConfText += "\t}\n\n"

	lsConfText += "\tlocation ~* .*\\.(js|css)?$ {\n"
	lsConfText += "\t\texpires 7d" + lineSpliter
	lsConfText += "\t\taccess_log off" + lineSpliter
	lsConfText += "\t}\n\n"

	lsConfText += "\tlocation ~* \\.(eot|ttf|ttc|otf|eot|woff|woff2|svg)$ {\n"
	lsConfText += "\t\tadd_header Access-Control-Allow-Origin *" + lineSpliter
	lsConfText += "\t\taccess_log off" + lineSpliter
	lsConfText += "\t}\n\n"

	lsConfText += "}\n\n"

	// 设置自定义server
	lsConfText += "\n\n"
	lsConfText += "# set custom servers" + lineSpliter

	for i := 0; i < len(cfg.Apps); i++ {
		if cfg.Apps[i].Server_name == "" {
			continue
		}

		App := cfg.Apps[i]

		lsConfText += "server {\n"
		lsConfText += "\tlisten " + App.Listen + lineSpliter
		lsConfText += "\tserver_name " + App.Server_name + lineSpliter
		lsConfText += "\tindex index.php index.html index.htm" + lineSpliter
		lsConfText += "\troot " + App.Root + lineSpliter
		if App.Logs_dir != "" {
			lsConfText += "\taccess_log " + App.Logs_dir + "/" + App.Server_name + ".access.log" + lineSpliter
			lsConfText += "\terror_log " + App.Logs_dir + "/" + App.Server_name + ".error.log" + lineSpliter
		}
		lsConfText += "\n\n"

		lsConfText += "\tlocation / {\n"
		lsConfText += "\t\tif (!-e $request_filename) {\n"
		lsConfText += "\t\t\trewrite ^(.*)$ /index.php?s=$1 last" + lineSpliter
		lsConfText += "\t\t\tbreak" + lineSpliter
		lsConfText += "\t\t}\n"
		lsConfText += "\t}\n\n"

		lsConfText += "\tlocation ~ \\.php(.*)$ {\n"
		lsConfText += "\t\tfastcgi_pass " + cfg.Server.Ip + ":" + cfg.Php.Port + lineSpliter
		lsConfText += "\t\tfastcgi_index index.php" + lineSpliter + "\n"
		lsConfText += "\t\tinclude fastcgi_params" + lineSpliter
		lsConfText += "\t\tset $real_script_name $fastcgi_script_name" + lineSpliter
		lsConfText += "\t\tset $real_path_info $fastcgi_path_info" + lineSpliter
		lsConfText += "\t\tif ($fastcgi_script_name ~ ^(.+?\\.php)(/.+)$) {\n"
		lsConfText += "\t\t\tset $real_script_name $1" + lineSpliter
		lsConfText += "\t\t\tset $real_path_info $2" + lineSpliter
		lsConfText += "\t\t}\n\n"
		lsConfText += "\t\tfastcgi_param QUERY_STRING $query_string" + lineSpliter
		lsConfText += "\t\tfastcgi_param SCRIPT_FILENAME $document_root$real_script_name" + lineSpliter
		lsConfText += "\t\tfastcgi_param APP_JT $app_name" + lineSpliter
		lsConfText += "\t\tfastcgi_param SCRIPT_NAME $real_script_name" + lineSpliter
		lsConfText += "\t\tfastcgi_param PATH_INFO $real_path_info" + lineSpliter
		lsConfText += "\t}\n\n"

		lsConfText += "\tlocation ~* .*\\.(gif|jpg|jpeg|png|bmp|swf|flv|ico)$ {\n"
		lsConfText += "\t\texpires 30d" + lineSpliter
		lsConfText += "\t\taccess_log off" + lineSpliter
		lsConfText += "\t}\n\n"

		lsConfText += "\tlocation ~* .*\\.(js|css)?$ {\n"
		lsConfText += "\t\texpires 7d" + lineSpliter
		lsConfText += "\t\taccess_log off" + lineSpliter
		lsConfText += "\t}\n\n"

		lsConfText += "\tlocation ~* \\.(eot|ttf|ttc|otf|eot|woff|woff2|svg)$ {\n"
		lsConfText += "\t\tadd_header Access-Control-Allow-Origin *" + lineSpliter
		lsConfText += "\t\taccess_log off" + lineSpliter
		lsConfText += "\t}\n\n"

		lsConfText += "}\n\n"
	}

	confDir := thisDir + "/conf"
	os.MkdirAll(confDir, os.ModePerm)
	os.WriteFile(confDir+"/nginx.conf", []byte(lsConfText), os.ModePerm)
}

func checkStatus() (int, int) {
	var nginxStatus int
	var phpStatus int

	if isWin() {
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

func showStatus(nStat int, pStat int) {
	fmt.Println("状态\n********************")
	text := "Inactive"
	if nStat > 0 {
		text = "Active"
	}
	fmt.Println("[Nginx]..." + text)

	text = "Inactive"
	if pStat > 0 {
		text = "Active"
	}
	fmt.Println("[PHP]....." + text)
}

func isWin() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	} else if !isAbsoluteDir(cfgFilepath) { // 相对路径下自动补齐当前目录
		cfgFilepath = thisDir + cfgFilepath
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

func linkToNginxConf() {
	confFilepath := cfg.Nginx.Conf

	// 获取nginx.conf文件中的内容
	buf, err := os.ReadFile(confFilepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	confText := string(buf)

	// 向nginx.conf中塞入ls-nginx.conf的文件包含代码
	lsConfFilepath := thisConfDir + "/nginx.conf"
	// 文件不存在，则创建个空的
	if _, err := os.Stat(lsConfFilepath); err != nil {
		os.WriteFile(lsConfFilepath, []byte(""), os.ModePerm)
	}
	if !strings.Contains(confText, lsConfFilepath) {
		index := strings.LastIndex(confText, "}")
		confText := confText[:index]
		confText += "    include " + lsConfFilepath + ";\n}"
		// fmt.Println(confText)
		os.WriteFile(confFilepath, []byte(confText), os.ModePerm)
	}
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

func exportToLsNginxConf() {
	// confFilepath := cfg.Nginx.Conf

	// 获取nginx.conf文件中的内容
	// buf, err := os.ReadFile(confFilepath)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// confText := string(buf)

	// 预设配置
	// preConf := map[string]string{"resolver": "8.8.8.8 ipv6=off"}
	// if !strings.Contains(confText, "server_names_hash_bucket_size") {
	// 	preConf["server_names_hash_bucket_size"] = "128"
	// }
	// if !strings.Contains(confText, "client_max_body_size") {
	// 	preConf["client_max_body_size"] = "2048m"
	// }

	lineSplitter := ";\n"
	lsConfText := "# set pre-configuration variables" + lineSplitter
	// for k, v := range preConf {
	// 	lsConfText += k + " " + v + lineSplitter
	// }

	// 设置app server
	lsConfText += "\n\n"
	lsConfText += "# set common server" + lineSplitter

	lsConfText += "server {\n"
	lsConfText += "\tlisten " + cfg.Server.Port + lineSplitter
	lsConfText += "\tserver_name " + cfg.Server.Host + lineSplitter + "\n"

	lsConfText += "\tset $www_root " + cfg.Server.Wwwroot + lineSplitter
	lsConfText += "\tset $php_port " + cfg.Server.Host + ":" + cfg.Php.Port + lineSplitter + "\n"

	lsConfText += "\tset $app_name 0" + lineSplitter
	lsConfText += "\tset $root 0" + lineSplitter
	lsConfText += "\tset $app_default_index 0" + lineSplitter + "\n"

	lsConfText += "\tif ($uri ~* ^/app/([^\\" + "/]+)/$) {\n"
	lsConfText += "\t\tset $app_name $1" + lineSplitter
	lsConfText += "\t\trewrite ^/app/([^\\" + "/]+)/$ /index.php" + lineSplitter
	lsConfText += "\t}\n\n"

	lsConfText += "\tif ($uri ~* ^/app/([^\\" + "/]+)/(.*)$) {\n"
	lsConfText += "\t\tset $app_name $1" + lineSplitter
	lsConfText += "\t\tset $app_default_index $2" + lineSplitter
	lsConfText += "\t\trewrite ^/app/([^\\" + "/]+)/(.*)$ /$app_default_index" + lineSplitter
	lsConfText += "\t}\n\n"

	lsConfText += "\t## match app_name START" + lineSplitter
	for i := 0; i < len(cfg.Apps); i++ {
		lsConfText += "\tif ($app_name ~* ^" + cfg.Apps[i].Name + "$) {\n"
		lsConfText += "\t\tset $root " + cfg.Apps[i].Root + lineSplitter
		lsConfText += "\t}\n\n"
	}
	lsConfText += "\t## match app_name END" + lineSplitter + "\n"

	// lsConfText += "\tif ($root = 0) {\n"
	// lsConfText += "\t\tbreak" + lineSpliter
	// lsConfText += "\t}\n\n"

	lsConfText += "\tif ($root = 0) {\n"
	lsConfText += "\t}\n\n"

	lsConfText += "\tlocation / {\n"
	lsConfText += "\t\tif (!-e $request_filename) {\n"
	lsConfText += "\t\t\trewrite ^(.*)$ /index.php?s=$1 last" + lineSplitter
	lsConfText += "\t\t\tbreak" + lineSplitter
	lsConfText += "\t\t}\n"
	lsConfText += "\t}\n\n"

	lsConfText += "\tindex index.php index.html index.htm" + lineSplitter
	lsConfText += "\troot $root" + lineSplitter
	lsConfText += "\taccess_log " + cfg.Nginx.Logs_dir + "/nginx.access.log" + lineSplitter
	lsConfText += "\terror_log " + cfg.Nginx.Logs_dir + "/nginx.error.log" + lineSplitter

	lsConfText += "\tlocation ~ \\.php(.*)$ {\n"
	lsConfText += "\t\tfastcgi_pass $php_port" + lineSplitter
	lsConfText += "\t\tfastcgi_index index.php" + lineSplitter + "\n"
	lsConfText += "\t\tinclude fastcgi_params" + lineSplitter
	lsConfText += "\t\tset $real_script_name $fastcgi_script_name" + lineSplitter
	lsConfText += "\t\tset $real_path_info $fastcgi_path_info" + lineSplitter
	lsConfText += "\t\tif ($fastcgi_script_name ~ ^(.+?\\.php)(/.+)$) {\n"
	lsConfText += "\t\t\tset $real_script_name $1" + lineSplitter
	lsConfText += "\t\t\tset $real_path_info $2" + lineSplitter
	lsConfText += "\t\t}\n\n"
	lsConfText += "\t\tfastcgi_param QUERY_STRING $query_string" + lineSplitter
	lsConfText += "\t\tfastcgi_param SCRIPT_FILENAME $document_root$real_script_name" + lineSplitter
	lsConfText += "\t\tfastcgi_param APP_JT $app_name" + lineSplitter
	lsConfText += "\t\tfastcgi_param SCRIPT_NAME $real_script_name" + lineSplitter
	lsConfText += "\t\tfastcgi_param PATH_INFO $real_path_info" + lineSplitter
	lsConfText += "\t}\n\n"

	lsConfText += "\tlocation ~* .*\\.(gif|jpg|jpeg|png|bmp|swf|flv|ico)$ {\n"
	lsConfText += "\t\texpires 30d" + lineSplitter
	lsConfText += "\t\taccess_log off" + lineSplitter
	lsConfText += "\t}\n\n"

	lsConfText += "\tlocation ~* .*\\.(js|css)?$ {\n"
	lsConfText += "\t\texpires 7d" + lineSplitter
	lsConfText += "\t\taccess_log off" + lineSplitter
	lsConfText += "\t}\n\n"

	lsConfText += "\tlocation ~* \\.(eot|ttf|ttc|otf|eot|woff|woff2|svg)$ {\n"
	lsConfText += "\t\tadd_header Access-Control-Allow-Origin *" + lineSplitter
	lsConfText += "\t\taccess_log off" + lineSplitter
	lsConfText += "\t}\n\n"

	lsConfText += "}\n\n"

	// 设置自定义server
	lsConfText += "\n\n"
	lsConfText += "# set custom servers" + lineSplitter

	for i := 0; i < len(cfg.Apps); i++ {
		if cfg.Apps[i].Server_name == "" {
			continue
		}

		App := cfg.Apps[i]

		lsConfText += "server {\n"
		lsConfText += "\tlisten " + App.Listen + lineSplitter
		lsConfText += "\tserver_name " + App.Server_name + lineSplitter
		lsConfText += "\tindex index.php index.html index.htm" + lineSplitter
		lsConfText += "\troot " + App.Root + lineSplitter
		if App.Logs_dir != "" {
			lsConfText += "\taccess_log " + App.Logs_dir + "/" + App.Server_name + ".access.log" + lineSplitter
			lsConfText += "\terror_log " + App.Logs_dir + "/" + App.Server_name + ".error.log" + lineSplitter
		}
		lsConfText += "\n\n"

		lsConfText += "\tlocation / {\n"
		lsConfText += "\t\tif (!-e $request_filename) {\n"
		lsConfText += "\t\t\trewrite ^(.*)$ /index.php?s=$1 last" + lineSplitter
		lsConfText += "\t\t\tbreak" + lineSplitter
		lsConfText += "\t\t}\n"
		lsConfText += "\t}\n\n"

		lsConfText += "\tlocation ~ \\.php(.*)$ {\n"
		lsConfText += "\t\tfastcgi_pass " + cfg.Server.Host + ":" + cfg.Php.Port + lineSplitter
		lsConfText += "\t\tfastcgi_index index.php" + lineSplitter + "\n"
		lsConfText += "\t\tinclude fastcgi_params" + lineSplitter
		lsConfText += "\t\tset $real_script_name $fastcgi_script_name" + lineSplitter
		lsConfText += "\t\tset $real_path_info $fastcgi_path_info" + lineSplitter
		lsConfText += "\t\tif ($fastcgi_script_name ~ ^(.+?\\.php)(/.+)$) {\n"
		lsConfText += "\t\t\tset $real_script_name $1" + lineSplitter
		lsConfText += "\t\t\tset $real_path_info $2" + lineSplitter
		lsConfText += "\t\t}\n\n"
		lsConfText += "\t\tfastcgi_param QUERY_STRING $query_string" + lineSplitter
		lsConfText += "\t\tfastcgi_param SCRIPT_FILENAME $document_root$real_script_name" + lineSplitter
		lsConfText += "\t\tfastcgi_param APP_JT $app_name" + lineSplitter
		lsConfText += "\t\tfastcgi_param SCRIPT_NAME $real_script_name" + lineSplitter
		lsConfText += "\t\tfastcgi_param PATH_INFO $real_path_info" + lineSplitter
		lsConfText += "\t}\n\n"

		lsConfText += "\tlocation ~* .*\\.(gif|jpg|jpeg|png|bmp|swf|flv|ico)$ {\n"
		lsConfText += "\t\texpires 30d" + lineSplitter
		lsConfText += "\t\taccess_log off" + lineSplitter
		lsConfText += "\t}\n\n"

		lsConfText += "\tlocation ~* .*\\.(js|css)?$ {\n"
		lsConfText += "\t\texpires 7d" + lineSplitter
		lsConfText += "\t\taccess_log off" + lineSplitter
		lsConfText += "\t}\n\n"

		lsConfText += "\tlocation ~* \\.(eot|ttf|ttc|otf|eot|woff|woff2|svg)$ {\n"
		lsConfText += "\t\tadd_header Access-Control-Allow-Origin *" + lineSplitter
		lsConfText += "\t\taccess_log off" + lineSplitter
		lsConfText += "\t}\n\n"

		lsConfText += "}\n\n"
	}

	confDir := thisDir + "/conf"
	os.MkdirAll(confDir, os.ModePerm)
	os.WriteFile(confDir+"/nginx.conf", []byte(lsConfText), os.ModePerm)
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

func isWindows() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

func isLinux() bool {
	return strings.Contains(runtime.GOOS, "linux")
}

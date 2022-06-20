package main

import (
	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var (
	commonFont declarative.Font
	titleFont  declarative.Font
	editFont   declarative.Font

	MW_window0     *walk.MainWindow
	MW_window0_0   *walk.MainWindow
	MW_window0_0_0 *walk.MainWindow
	MW_window0_0_1 *walk.MainWindow

	LE_server_ip       *walk.LineEdit
	LE_server_port     *walk.LineEdit
	LE_server_wwwroot  *walk.LineEdit
	LE_nginx_dir       *walk.LineEdit
	LE_nginx_conf      *walk.LineEdit
	LE_nginx_logs_dir  *walk.LineEdit
	LE_php_dir         *walk.LineEdit
	LE_php_ini         *walk.LineEdit
	LE_php_port        *walk.LineEdit
	LB_appListBox      *walk.ListBox
	appList            *AppList
	LE_app_name        *walk.LineEdit
	LE_app_root        *walk.LineEdit
	LE_app_server_name *walk.LineEdit
	LE_app_listen      *walk.LineEdit
	LE_app_logs_dir    *walk.LineEdit

	// model *AppListModel
)

func runGui() {
	fmt.Println("start GUI...")

	initUi()

	declarative.MainWindow{
		Title:    APP_NAME,
		AssignTo: &MW_window0,
		Size:     declarative.Size{540, 400},
		Layout: declarative.Grid{
			Columns: 2,
		},
		Children: []declarative.Widget{

			declarative.GroupBox{
				Title:  "服务器",
				Font:   titleFont,
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.Label{
						Text: "IP",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_server_ip,
						Text:     cfg.Server.Ip,
						Font:     editFont,
					},
					declarative.Label{
						Text: "端口",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_server_port,
						Text:     cfg.Server.Port,
						Font:     editFont,
					},
					declarative.Label{
						Text: "wwwroot目录",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_server_wwwroot,
						Text:     cfg.Server.Wwwroot,
						Font:     editFont,
					},
				},
			},
			declarative.GroupBox{
				Title:  "Nginx",
				Font:   titleFont,
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.Label{
						Text: "目录",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_nginx_dir,
						Text:     cfg.Nginx.Dir,
						Font:     editFont,
					},
					declarative.Label{
						Text: "配置文件（可选）",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_nginx_conf,
						Text:     cfg.Nginx.Conf,
						Font:     editFont,
					},
					declarative.Label{
						Text: "log目录（可选）",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_nginx_logs_dir,
						Text:     cfg.Nginx.Logs_dir,
						Font:     editFont,
					},
				},
			},
			declarative.GroupBox{
				Title:  "PHP",
				Font:   titleFont,
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.Label{
						Text: "目录",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_php_dir,
						Text:     cfg.Php.Dir,
						Font:     editFont,
					},
					declarative.Label{
						Text: "配置文件（可选）",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_php_ini,
						Text:     cfg.Php.Ini,
						Font:     editFont,
					},
					declarative.Label{
						Text: "cgi端口",
						Font: commonFont,
					},
					declarative.LineEdit{
						AssignTo: &LE_php_port,
						Text:     cfg.Php.Port,
						Font:     editFont,
					},
				},
			},
			declarative.GroupBox{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						Text:      "管理应用",
						Font:      titleFont,
						OnClicked: OpenAppsWindow,
					},
					declarative.PushButton{
						Text:      "保存配置",
						Font:      titleFont,
						OnClicked: Action_saveConfig,
					},
					declarative.PushButton{
						Text:      "环境-启动",
						Font:      titleFont,
						OnClicked: Action_startEnv,
					},
					declarative.PushButton{
						Text:      "环境-停止",
						Font:      titleFont,
						OnClicked: Action_stopEnv,
					},
				},
			},
		},
	}.Create()

	// 设置 ^win.WS_MAXIMIZEBOX 禁用最大化按钮
	// 设置 ^win.WS_THICKFRAME 禁用窗口大小改变
	win.SetWindowLong(
		MW_window0.Handle(),
		win.GWL_STYLE,
		win.GetWindowLong(MW_window0.Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME,
	)

	MW_window0.Run()

}

func initUi() {
	commonFont = declarative.Font{
		Family:    "宋体",
		PointSize: 10,
	}
	titleFont = declarative.Font{
		Family:    "宋体",
		PointSize: 12,
		Bold:      true,
	}
	editFont = declarative.Font{
		Family:    "宋体",
		PointSize: 12,
	}
}

func OpenAppsWindow() {
	LB_appListBox = &walk.ListBox{}
	appList = NewAppList()

	declarative.MainWindow{
		Title:    "应用列表",
		AssignTo: &MW_window0_0,
		Size:     declarative.Size{400, 300},
		Layout: declarative.Grid{
			Columns:     4,
			MarginsZero: true,
		},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo:   &LB_appListBox,
				Model:      appList,
				Font:       editFont,
				ColumnSpan: 3,
			},
			declarative.GroupBox{
				Layout: declarative.VBox{},
				Font:   titleFont,
				Children: []declarative.Widget{
					declarative.PushButton{
						Text:      "上移",
						OnClicked: Action_moveUp,
						MaxSize:   declarative.Size{60, 35},
					},
					declarative.PushButton{
						Text:      "修改",
						OnClicked: Action_openAppEditWindow,
						MaxSize:   declarative.Size{50, 35},
					},
					declarative.PushButton{
						Text:      "下移",
						OnClicked: Action_moveDown,
						MaxSize:   declarative.Size{40, 35},
					},
					declarative.PushButton{
						Text:      "新增",
						OnClicked: Action_openAppCreateWindow,
						MaxSize:   declarative.Size{30, 35},
					},
					declarative.PushButton{
						Text:      "删除",
						OnClicked: Action_deleteFromAppList,
						MaxSize:   declarative.Size{40, 35},
					},
				},
			},
		},
	}.Create()

	MW_window0_0.Run()
}

type AppList struct {
	walk.ListModelBase

	items []AppItem
}

type AppItem struct {
	name  string
	value string
}

func NewAppList() *AppList {
	m := &AppList{
		items: make([]AppItem, len(cfg.Apps)),
	}

	for i := 0; i < len(cfg.Apps); i++ {
		name := cfg.Apps[i].Name
		if cfg.Apps[i].Server_name != "" {
			name += "[" + cfg.Apps[i].Server_name + "]"
		}
		m.items[i] = AppItem{name: name, value: ""}
	}

	return m
}

func (m *AppList) ItemCount() int {
	return len(m.items)
}

func (m *AppList) Value(index int) interface{} {
	return m.items[index].name
}

func Action_openAppCreateWindow() {
	declarative.MainWindow{
		Title:    "新增应用",
		AssignTo: &MW_window0_0_0,
		Size:     declarative.Size{300, 200},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Label{
				Text: "名称",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_name,
				Font:     editFont,
			},
			declarative.Label{
				Text: "根目录",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_root,
				Font:     editFont,
			},
			declarative.Label{
				Text: "服务器名称",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_server_name,
				Font:     editFont,
			},
			declarative.Label{
				Text: "监听端口",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_listen,
				Font:     editFont,
			},
			declarative.Label{
				Text: "logs目录(可选)",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_logs_dir,
				Font:     editFont,
			},
			declarative.PushButton{
				Text:      "确定",
				Font:      titleFont,
				OnClicked: CreateAppList,
			},
		},
	}.Create()

	MW_window0_0_0.Run()
}

func CreateAppList() {
	i := LB_appListBox.CurrentIndex()
	if i == -1 {
		i = len(cfg.Apps) - 1
	}
	appCfg := AppConfig{
		Name:        LE_app_name.Text(),
		Root:        getRealDir(LE_app_root.Text()),
		Server_name: LE_app_server_name.Text(),
		Listen:      LE_app_listen.Text(),
		Logs_dir:    getRealDir(LE_app_logs_dir.Text()),
	}
	if !isAbsoluteDir(appCfg.Root) { // 不是绝对路径？那前面得加wwwroot
		appCfg.Root = cfg.Server.Wwwroot + "/" + appCfg.Root
	}
	if appCfg.Listen == "" && appCfg.Server_name != "" {
		appCfg.Listen = cfg.Server.Port
	}
	if appCfg.Logs_dir == "" && appCfg.Server_name != "" {
		appCfg.Logs_dir = cfg.Nginx.Logs_dir
	}

	j := i + 1

	newApps := []AppConfig{}
	newApps = append(newApps, cfg.Apps[:j]...)
	// fmt.Println(newApps)
	newApps = append(newApps, appCfg)
	// fmt.Println(newApps)
	newApps = append(newApps, cfg.Apps[j:]...)
	// fmt.Println(newApps)
	cfg.Apps = newApps

	MW_window0_0_0.Close()
	reloadAppList()
}

func Action_openAppEditWindow() {
	i := LB_appListBox.CurrentIndex()
	if i == -1 {
		return
	}
	declarative.MainWindow{
		Title:    "编辑应用",
		AssignTo: &MW_window0_0_1,
		Size:     declarative.Size{300, 200},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Label{
				Text: "名称",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_name,
				Text:     cfg.Apps[i].Name,
				Font:     editFont,
			},
			declarative.Label{
				Text: "根目录",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_root,
				Text:     cfg.Apps[i].Root,
				Font:     editFont,
			},
			declarative.Label{
				Text: "服务器名称",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_server_name,
				Text:     cfg.Apps[i].Server_name,
				Font:     editFont,
			},
			declarative.Label{
				Text: "监听端口",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_listen,
				Text:     cfg.Apps[i].Listen,
				Font:     editFont,
			},
			declarative.Label{
				Text: "logs目录(可选)",
				Font: commonFont,
			},
			declarative.LineEdit{
				AssignTo: &LE_app_logs_dir,
				Text:     cfg.Apps[i].Logs_dir,
				Font:     editFont,
			},
			declarative.PushButton{
				Text:      "确定",
				Font:      titleFont,
				OnClicked: Action_saveAppList,
			},
		},
	}.Create()

	MW_window0_0_1.Run()
}

func Action_saveAppList() {
	i := LB_appListBox.CurrentIndex()
	if i == -1 {
		return
	}
	appCfg := AppConfig{
		Name:        LE_app_name.Text(),
		Root:        getRealDir(LE_app_root.Text()),
		Server_name: LE_app_server_name.Text(),
		Listen:      LE_app_listen.Text(),
		Logs_dir:    getRealDir(LE_app_logs_dir.Text()),
	}
	if !isAbsoluteDir(appCfg.Root) { // 不是绝对路径？那前面得加wwwroot
		appCfg.Root = cfg.Server.Wwwroot + "/" + appCfg.Root
	}
	if appCfg.Listen == "" && appCfg.Server_name != "" {
		appCfg.Listen = cfg.Server.Port
	}
	if appCfg.Logs_dir == "" && appCfg.Server_name != "" {
		appCfg.Logs_dir = cfg.Nginx.Logs_dir
	}
	cfg.Apps[i] = appCfg

	MW_window0_0_1.Close()
	reloadAppList()

	LB_appListBox.SetCurrentIndex(i)
}

func Action_deleteFromAppList() {
	i := LB_appListBox.CurrentIndex()
	if i == -1 {
		return
	}
	cfg.Apps = append(cfg.Apps[:i], cfg.Apps[i+1:]...)
	reloadAppList()
}

func Action_moveUp() {
	i := LB_appListBox.CurrentIndex()
	if i == -1 {
		return
	}

	if i == 0 {
		return
	}

	j := i - 1
	cfg.Apps[j], cfg.Apps[i] = cfg.Apps[i], cfg.Apps[j]
	reloadAppList()

	LB_appListBox.SetCurrentIndex(j)
}

func Action_moveDown() {
	i := LB_appListBox.CurrentIndex()
	if i == -1 {
		return
	}

	if i == len(cfg.Apps)-1 {
		return
	}

	j := i + 1
	cfg.Apps[j], cfg.Apps[i] = cfg.Apps[i], cfg.Apps[j]
	reloadAppList()

	LB_appListBox.SetCurrentIndex(j)
}

func reloadAppList() {
	appList = NewAppList()

	LB_appListBox.SetModel(appList)
}

func Action_saveConfig() {
	// 从ui界面保存配置值
	cfg.Server.Ip = LE_server_ip.Text()
	cfg.Server.Port = LE_server_port.Text()
	cfg.Server.Wwwroot = LE_server_wwwroot.Text()

	cfg.Nginx.Dir = LE_nginx_dir.Text()
	cfg.Nginx.Conf = LE_nginx_conf.Text()
	cfg.Nginx.Logs_dir = LE_nginx_logs_dir.Text()
	if cfg.Nginx.Conf == "" {
		cfg.Nginx.Conf = cfg.Nginx.Dir + "/conf/nginx.conf"
	}
	if cfg.Nginx.Logs_dir == "" {
		cfg.Nginx.Logs_dir = cfg.Nginx.Dir + "/logs"
	}

	cfg.Php.Dir = LE_php_dir.Text()
	cfg.Php.Ini = LE_php_ini.Text()
	cfg.Php.Port = LE_php_port.Text()
	if cfg.Php.Ini == "" {
		cfg.Php.Ini = cfg.Php.Dir + "/php.ini"
	}
	if cfg.Php.Port == "" {
		cfg.Php.Port = "9000"
	}
	saveConfig()

	reloadWindow0()
}

func reloadWindow0() {
	LE_server_ip.SetText(cfg.Server.Ip)
	LE_server_port.SetText(cfg.Server.Port)
	LE_server_wwwroot.SetText(cfg.Server.Wwwroot)

	LE_nginx_dir.SetText(cfg.Nginx.Dir)
	LE_nginx_conf.SetText(cfg.Nginx.Conf)
	LE_nginx_logs_dir.SetText(cfg.Nginx.Logs_dir)

	LE_php_dir.SetText(cfg.Php.Dir)
	LE_php_ini.SetText(cfg.Php.Ini)
	LE_php_port.SetText(cfg.Php.Port)
}

func Action_startEnv() {
	// 保存当前配置
	Action_saveConfig()

	// 导出nginx配置
	exportNginx()

	startEnv()
}

func Action_stopEnv() {
	stopEnv()
}

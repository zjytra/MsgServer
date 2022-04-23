// 创建时间: 2019/10/17
// 作者: zjy
// 功能介绍:
// 程序最外层 ,这里给main的入口,以及整个进程退出的控制
// 包含程序的启动,停止
package app

import (
	"flag"
	"fmt"
	"github.com/zjytra/MsgServer/app/appclient"
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/app/apploginsv"
	"github.com/zjytra/MsgServer/conf"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/timingwheel"
	"github.com/zjytra/MsgServer/engine_core/xengine"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/model"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	G_App   *App       //把app放在这里方便全局调用
)


type App struct {
	model.AppState   // app 狀態
	appLogic       xengine.ServerLogic
	appCloseTime   time.Duration       // app关闭倒计时
}

func NewApp()  {
	G_App =  new(App)
	G_App.InitApp()
	G_App.StartApp()
	G_App.RunApp()
	G_App.CloseApp()
}


//初始化
func (this *App)InitApp() {
	fmt.Println("App Init")
	conf.InitConf()
	csvdata.InitCsv()
	appdata.InitData(this)
	timingwheel.InitTimeWheel()
	dispatch.InitMainQueue()
}

// 获取命令行启动
// 1. app相关的配置文件初始化
// 2. 设置app参数
func (this *App) StartApp() {
	fmt.Println("App StartApp")
	conf.StartConf()
	csvdata.StartCsv()
	this.ParseAppArgs() //获取命令行
	this.NewLog()
}

func (this *App)NewLog() {
	logInit := &conf.LogInitModel{
		ServerName: csvdata.OutNetConf.App_name,
		LogsPath:   conf.PathModelPtr.LogsPath,
		Volatile:   conf.SvJson.LogConf,
	}
	if !xlog.NewXlog(logInit)   {
		panic("New Xlog errors ")
	}
	xlog.Debug( "NewXlog success")
}

// 程序启动获取命令行参数
func (this *App) ParseAppArgs() {
	var appid int
	flag.IntVar(&appid, "AppID", 0, "请输入app id")
	flag.Parse()
	xlog.Debug( "ParseAppArgs appid %d ",appid)
	appdata.InitAppDataByAppArgs(int32(appid))
	this.CreateAppBehavior()
	xlog.Debug( "ParseAppArgs success")
}

// 根据配置启动对应服务器
func (this *App)CreateAppBehavior() {
	//创建日志系统
	// 初始化app相关
	this.appLogic = NewAppBehavior(appdata.GetAppKind())
	if this.appLogic == nil {
		panic("appLogic == nil ")
	}
	//对应的服务器启动
	this.appLogic.OnStart()
	// 执行对应
	this.SetAppOpen()
	// //读取控制台命令 测试的时候才用
	// this.AppWG.Add(1)
	// go this.ReadConsole()
	go func() {
		ip := ":30003"
		erro := http.ListenAndServe(ip,nil)
		if erro != nil {
			fmt.Printf("%v",erro)
		}
	}()
}



func (this *App) RunApp() {
	//初始化主线程事件队列
	dispatch.MainQueue.Start()
}



//关闭程序倒计时
func (this *App)SetAppCloseTime(time time.Duration)  {
	this.appCloseTime = time
}

func (this *App)SetAppOpen()  {
	this.InitAppState()
	this.AppOpen()
}


// 关闭app 程序
func (this *App)CloseApp() {
	this.AppClose()
	this.appLogic.OnRelease() // 逻辑结束
	//timersys.OnAppClose()        // 定时器结束
	timingwheel.Stop()            //定时器结束
	csvdata.OnAppClose()         //配置文件结束
	conf.OnAppClose()
	xlog.CloseLog() // 退出日志
	fmt.Println("App Close")
}

// app 逻辑参数根据服务器启动的参数创建对应的服务器工厂
func NewAppBehavior(svKind model.AppKind)  xengine.ServerLogic {
	switch svKind {
	case model.APP_NONE:
		return nil
	case model.APP_Client:
		return new(appclient.AppClient)
	case model.APP_LoginServer: // 工厂
		return new(apploginsv.LoginServer)
	}
	return nil
}



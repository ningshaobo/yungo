package main

import (
	//	"fmt"
	"com.my/yungo/models"
	_ "com.my/yungo/routers"
	"github.com/astaxie/beego"
	beegocontext "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"
	"log"
)

// 全局 session 定义
var GlobalSessions *session.Manager

/**
* @param 初始化参数， 包括 数据库等
 */
func init() {
	// 判断是否存在 数据库， 如果不存在， 建立新库
	dbConn := models.UtilsDBConn()
	if dbConn == "" {
		log.Panicln("获取数据库连接配置失败，异常退出")
	}
	dbName := models.UtilsDBName()
	if dbName == "" {
		log.Panicln("获取数据库名称失败，异常退出")
	}

	// "nsb:nsb123@tcp(192.168.6.1:3306)/?charset=utf8"
	models.SystemDbInit(dbName, dbConn+"/?charset=utf8")
	// 开始 数据库 初始化， 开启 ORM 调试模式
	orm.Debug = false
	models.ClustersRegisterDB()

	//注册驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)                                    //注册默认数据库
	orm.RegisterDataBase("default", "mysql", dbConn+"/"+dbName+"?charset=utf8") //"root:password@/docker?charset=utf8")
	// 自动建表
	orm.RunSyncdb("default", false, true)
}

/**
* @param 初始化过滤器
 */
func filterInt() {
	var FilterSession = func(ctx *beegocontext.Context) {
		session := ctx.Input.Session("uid")
		if session == nil {
			ctx.Redirect(401, "/logins")
		} else {
			uid := session.(string)
			//			fmt.Println(uid)
			ctx.Input.SetParam("uid", uid)
		}
	}
	rtBase := models.UtilsBaseRoute()
	if rtBase == "" {
		log.Panicln("获取路由基础路径失败，异常退出")
	}
	beego.InsertFilter("/"+rtBase+"/*", beego.BeforeRouter, FilterSession)
}

/**
* @param
 */
func main() {
	// 设置session 相关配置
	// 初始化 全局session 控制器
	var cf session.ManagerConfig
	cf.CookieName = "nsbsesstion"
	cf.Gclifetime = 3600
	beego.GlobalSessions, _ = session.NewManager("memory", &cf) //("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 3600, "providerConfig": ""}`)
	go beego.GlobalSessions.GC()
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "swarmSid"

	// 初始化过滤器
	filterInt()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	//	beego.AppConfig.Set("nsb", "test")
	//	beego.AppConfig.SaveConfigFile("conf/app.conf")
	//	log.Println(beego.AppConfig.String("nsb"))
	beego.Run()
}

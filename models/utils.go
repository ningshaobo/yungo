package models

import (
	"github.com/astaxie/beego"
	"os"
)

/**
* @param 初始化过程中，一些基本信息处理
 */

/**
* @param 获取路由基础路径， 如 edge
 */
func UtilsBaseRoute() string {
	routeBase := os.Getenv("ROUTE_BASE")
	if routeBase == "" {
		routeBase = beego.AppConfig.String("appRouteBase")
	}
	return routeBase
}

/**
* @param 获取数据库连接 ：如： nsb:nsb123@tcp(192.168.6.1:3306)
 */
func UtilsDBConn() string {
	dbConn := os.Getenv("DB_CONN")
	if dbConn == "" {
		dbConn = beego.AppConfig.String("appDbConn")
	}
	return dbConn
}

/**
* @param 获取数据库名称 ：如：docker
 */
func UtilsDBName() string {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = beego.AppConfig.String("appDbName")
	}
	return dbName
}

/**
* @param 获取总台地址
 */
var appProxyUrl = ""

func UtilsGetProxyUrl() string {
	if appProxyUrl != "" {
		return appProxyUrl
	}

	proxyUrl := os.Getenv("PROXY_URL")
	if proxyUrl == "" {
		proxyUrl = beego.AppConfig.String("appProxyUrl")
	}
	if proxyUrl != "" {
		appProxyUrl = proxyUrl
	}
	return proxyUrl
}

/**
* @param 获取总台地址
 */
func UtilsSetProxyUrl(proxyUrl string) {
	appProxyUrl = proxyUrl
}

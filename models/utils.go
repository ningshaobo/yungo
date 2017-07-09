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
	var dbConn string
	dbConn = os.Getenv("DB_CONN")
	if dbConn == "" {
		dbSe,err := beego.AppConfig.GetSection("db")
		if err != nil {
			dbConn = "root:password@tcp(localhost:3306)"
		}else {
			dbConn = dbSe["dbconn"]
			if dbConn == "" {
				dbConn = "root:password@tcp(localhost:3306)"
			}
		}
		os.Setenv("DB_CONN", dbConn)
	}
	return dbConn
}

/**
* @param 获取数据库名称 ：如：docker
 */
func UtilsDBName() string {
	var dbName string
	dbName = os.Getenv("DB_NAME")
	if dbName == "" {
		dbSe,err := beego.AppConfig.GetSection("db")
		if err != nil {
			dbName = "yungo"
		}else {
			dbName = dbSe["dbname"]
			if dbName == "" {
				dbName = "yungo"
			}
		}
		os.Setenv("DB_NAME", dbName)
	}
	return dbName
}
/**
* @param 获取总台地址
 */
func UtilsMasterProxyUrl() string {
	var proxyUrl string
	proxyUrl = os.Getenv("MASTER_PROXY_URL")
	if proxyUrl == "" {
		proxySe, err	:= beego.AppConfig.GetSection("proxy")
		if err != nil {
			proxyUrl = "localhost:8080"
		}else {
			proxyUrl = proxySe["masterurl"]
			if proxyUrl == "" {
				proxyUrl = "localhost:8080"
			}
		}
		os.Setenv("MASTER_PROXY_URL", proxyUrl)
	}
	return proxyUrl
}


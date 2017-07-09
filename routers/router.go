// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"com.my/yungo/controllers"
	"com.my/yungo/models"
	"github.com/astaxie/beego"
	"log"
)

func init() {
	beego.Router("/*", &controllers.ProxyController{})
	rtBase := models.UtilsBaseRoute()
	if rtBase == "" {
		log.Panicln("获取路由基础路径失败，异常退出")
	}
	ns := beego.NewNamespace("/"+rtBase,
		beego.NSNamespace("/clusters",
			beego.NSInclude(
				&controllers.ClusterController{},
			),
		),
		beego.NSNamespace("/services",
			beego.NSInclude(
				&controllers.ServicesController{},
			),
		),
		beego.NSNamespace("/containers",
			beego.NSInclude(
				&controllers.ContainersController{},
			),
		),
	)
	beego.AddNamespace(ns)

	nsUser := beego.NewNamespace("/users", beego.NSInclude(&controllers.UsersNomalController{}))
	beego.AddNamespace(nsUser)
}

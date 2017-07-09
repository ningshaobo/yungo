package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["com.my/yungo/controllers:ClusterController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ClusterController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ClusterController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ClusterController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ContainersController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ContainersController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ContainersController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ContainersController"],
		beego.ControllerComments{
			Method: "Allotip",
			Router: `/allotip`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ContainersController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ContainersController"],
		beego.ControllerComments{
			Method: "Delip",
			Router: `/delip`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"],
		beego.ControllerComments{
			Method: "UpService",
			Router: `/upservice`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"],
		beego.ControllerComments{
			Method: "Tasks",
			Router: `/tasks`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"],
		beego.ControllerComments{
			Method: "NodeContainers",
			Router: `/nodecontainers`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:ServicesController"],
		beego.ControllerComments{
			Method: "MoreTasks",
			Router: `/moretasks`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:UsersNomalController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:UsersNomalController"],
		beego.ControllerComments{
			Method: "Zoneadmin",
			Router: `/zoneadmin`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["com.my/yungo/controllers:UsersNomalController"] = append(beego.GlobalControllerRouter["com.my/yungo/controllers:UsersNomalController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}

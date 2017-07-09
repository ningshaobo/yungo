package controllers

import (
	"fmt"
	//	"reflect"
	//	"strconv"
	"encoding/json"
	"github.com/astaxie/beego"
	beegologs "github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/session/redis"
	//	swarmtype		"github.com/docker/docker/api/types/swarm"
	"com.my/yungo/models"
	//	myUtils			"com.my/dmSvrWeb/utils"
)

// 服务操作
type ServicesController struct {
	beego.Controller
}

// @Title 创建swarm 服务
// @Description 创建 swarm 服务
// @Param	body		body 	models.CreateInBody	true		"包含 swarm 服务构建信息"
// @Success 200 {int} models.Service.Id
// @Failure 403 body is empty
// @router / [post]
func (this *ServicesController) Post() {
	var crtBody models.CreateInBody
	var err error
	clusterId := this.GetString("clusterid")
	if clusterId == "" {
		err = fmt.Errorf("clusterid 参数为空")
		beegologs.Warn("%s", err.Error())
		this.Data["json"] = err.Error()
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}
	tid := this.GetString("tid")
	if tid == "" {
		err = fmt.Errorf("tid 参数为空")
		beegologs.Warn("%s", err.Error())
		this.Data["json"] = err.Error()
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}
	//	crtBody.Mode = this.GetString("mode")
	json.Unmarshal(this.Ctx.Input.RequestBody, &crtBody)

	err = models.CreateService(tid, clusterId, crtBody)
	if err != nil {
		beegologs.Warn("创建服务失败，原因是：%s", err.Error())
		this.Data["json"] = "创建服务失败，原因是：" + err.Error()
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}

//// @Title delservice
//// @Description Logs user into the system
//// @Success 200 {string} login success
//// @Failure 403 del service fail
//// @router /delservice [post]
//func (this *ServicesController) Delservice() {
//	var svr models.ServiceTrans
//	json.Unmarshal(this.Ctx.Input.RequestBody, &svr)
//	clusterId := this.GetString("clusterid")
//	if clusterId == "" {
//		this.Ctx.Output.SetStatus(403)
//		this.Data["json"] = "clusterid 参数不能为空"
//		this.ServeJSON()
//		return
//	}
//
//	if svr.Id == "" {
//		this.Ctx.Output.SetStatus(403)
//		this.Data["json"] = "body 不能为空"
//		this.ServeJSON()
//		return
//	}
//	err := models.DeleteService(svr.Id, clusterId)
//	if err != nil {
//		this.Ctx.Output.SetStatus(403)
//		this.Data["json"] = "删除失败"
//		this.ServeJSON()
//		return
//	}
//	this.ServeJSON()
//}

//// @Title Get
//// @Description get service inspect by serviceid
//// @Success 200 {object} models.InspectService
//// @Failure 403 :serviceid is empty
//// @router /:serviceid [get]
//func (this *ServicesController) Get() {
//	var err error
//	clusterId := this.GetString("clusterid")
//	if clusterId == "" {
//		err = fmt.Errorf("clusterid 参数为空")
//		beegologs.Warn("%s", err.Error())
//		this.Data["json"] = err.Error()
//		this.Ctx.Output.SetStatus(400)
//		this.ServeJSON()
//		return
//	}
//	tid := this.GetString("tid")
//	if tid == "" {
//		err = fmt.Errorf("tid 参数为空")
//		beegologs.Warn("%s", err.Error())
//		this.Data["json"] = err.Error()
//		this.Ctx.Output.SetStatus(400)
//		this.ServeJSON()
//		return
//	}
//	serviceID := this.GetString(":serviceid")
//	if serviceID == "" {
//		err = fmt.Errorf("serviceID 为空")
//		beegologs.Warn("%s", err.Error())
//		this.Data["json"] = err.Error()
//		this.Ctx.Output.SetStatus(400)
//		this.ServeJSON()
//		return
//	}
//	services, err := models.InspectService(serviceID, clusterId)
//	if err != nil {
//		this.Data["json"] = err.Error()
//		this.Ctx.Output.SetStatus(400)
//		beegologs.Warn(err.Error())
//		this.ServeJSON()
//		return
//	}
//	session := this.GetSession("uid")
//	if session == nil {
//		err = fmt.Errorf("用户尚未登录，或者登录过期")
//		beegologs.Warn("%s", err.Error())
//		this.Data["json"] = err.Error()
//		this.ServeJSON()
//		return
//	}
//	err = fmt.Errorf("用户无租户（%s）权限", tid)
//	user := reflect.ValueOf(session).Interface().(models.User)
//	for _, ttu := range user.TenantTypeUsers {
//		if fmt.Sprintf("%d", ttu.Tenant.Id) == tid {
//			err = nil
//			break
//		}
//	}
//	if err != nil {
//		this.Data["json"] = err.Error()
//		this.Ctx.Output.SetStatus(400)
//		beegologs.Warn(err.Error())
//	} else {
//		this.Data["json"] = services
//	}
//	this.ServeJSON()
//}

// @Title GetAll
// @Description get all service
// @Success 200 {object} models.Services
// @router / [get]
func (this *ServicesController) GetAll() {
	var err error
	beegologs.Debug("GetAll services, tid = %s, uid = %s",
		this.Ctx.Input.Param("tid"), this.Ctx.Input.Param("uid"))

	tid := this.GetString("tid")
	if tid == "" {
		err = fmt.Errorf("tid 参数为空为空")
		beegologs.Warn(err.Error())
		this.Data["json"] = err.Error()
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}
	clusterId := this.GetString("clusterid")
	if clusterId == "" {
		err = fmt.Errorf("clusterId 参数为空为空")
		beegologs.Warn(err.Error())
		this.Data["json"] = err.Error()
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}

	services := models.GetAllService(tid, clusterId)
	this.Data["json"] = services
	this.ServeJSON()
}

// @Title UpService
// @Description update service
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /upservice [post]
func (this *ServicesController) UpService() {
	//	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	this.Data["json"] = models.UpdateService //"up service success"
	this.ServeJSON()
}

// @Title Tasks
// @Description get tasks from service
// @Param	serviceid
// @Success 200 [] tasks
// @Failure 400 serviceid clusterid empty
// @router /tasks [get]
func (this *ServicesController) Tasks() {
	var err error
	clusterId := this.GetString("clusterid")
	if clusterId == "" {
		this.Ctx.Output.SetStatus(400)
		err = fmt.Errorf("clusterid 参数为空")
		beegologs.Warn("%s", err.Error())
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}
	serviceId := this.GetString("serviceid")
	if serviceId == "" {
		this.Ctx.Output.SetStatus(400)
		err = fmt.Errorf("serviceId 为空")
		beegologs.Warn("%s", err.Error())
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}
	this.Data["json"] = models.GetAllSvrTasks(clusterId, serviceId)
	this.ServeJSON()
}

// @Title nodecontainers
// @Description get node info and container list by nodeId
// @Param	serviceid
// @Success 200 [] node containers
// @Failure 400 serviceid clusterid empty
// @router /nodecontainers [get]
func (this *ServicesController) NodeContainers() {
	var err error
	clusterId := this.GetString("clusterid")
	if clusterId == "" {
		this.Ctx.Output.SetStatus(400)
		err = fmt.Errorf("clusterid 参数为空")
		beegologs.Warn("%s", err.Error())
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}
	nodeId := this.GetString("nodeid")
	if nodeId == "" {
		this.Ctx.Output.SetStatus(400)
		err = fmt.Errorf("nodeId 为空")
		beegologs.Warn("%s", err.Error())
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}
	containers, nodeInfo := models.GetAllNodeContainers("", nodeId, clusterId)
	ret := make(map[string]interface{})
	ret["containers"] = containers
	ret["nodeInfo"] = nodeInfo
	this.Data["json"] = ret
	this.ServeJSON()
}

// @Title moreTasks
// @Description get node info and container list by nodeId
// @Param	serviceid
// @Success 200 [] node containers
// @Failure 400 serviceid clusterid empty
// @router /moretasks [get]
func (this *ServicesController) MoreTasks() {
	var err error
	clusterId := this.GetString("clusterid")
	if clusterId == "" {
		this.Ctx.Output.SetStatus(400)
		err = fmt.Errorf("clusterid 参数为空")
		beegologs.Warn("%s", err.Error())
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}
	serviceId := this.GetString("serviceid")
	if serviceId == "" {
		this.Ctx.Output.SetStatus(400)
		err = fmt.Errorf("serviceId 为空")
		beegologs.Warn("%s", err.Error())
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}
	this.Data["json"] = models.GetAllMoreTasks(clusterId, serviceId)
	this.ServeJSON()
}

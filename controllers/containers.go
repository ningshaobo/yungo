package controllers

import (
	"com.my/yungo/models"
	myutils "com.my/utils"
	"fmt"
	"github.com/astaxie/beego"
	beegologs "github.com/astaxie/beego/logs"
)

// Operations about Users
type ContainersController struct {
	beego.Controller
}

// @Title GetAll
// @Description 获取集群下所有容器
// @Success 200 {object} models.Containers
// @router / [get]
func (this *ContainersController) GetAll() {
	var err error
	cid := this.GetString("cid")
	if cid == "" {
		err = fmt.Errorf("获取 cid 参数失败")
		beegologs.Warn(err.Error())
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}
	tid := this.GetString("tid")
	if tid == "" {
		err = fmt.Errorf("获取 tid 参数失败")
		beegologs.Warn(err.Error())
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = err.Error()
		return
	}

	containers := models.ContainerByTidCid(tid, cid)
	this.Data["json"] = containers
	this.ServeJSON()
}

// @Title Allotip
// @Description 给容器分配 IP
// @Success 200 {string} 成功分配
// @Failure 400 failure
// @router /allotip [post]
func (this *ContainersController) Allotip() {
	beegologs.Debug("ContainersController  /Allotip")
	var body = make(map[string]string)
	myutils.CommctrBody(&this.Controller, &body)
	beegologs.Debug("body %v", body)

	err := models.ContainerAllotIp(body["host"], body["containerid"], body["ipaddr"])
	if err != nil {
		beegologs.Warn("容器（%v）分配 ip (%v) 失败，原因是：%v", body["containerid"], body["ipaddr"], err.Error())
		this.Data["json"] = fmt.Sprintf("分配ip失败")
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}
}

// @Title delip
// @Description 删除容器已经分配 IP
// @Success 200 {string} 成功删除
// @Failure 400 failure
// @router /delip [post]
func (this *ContainersController) Delip() {
	beegologs.Debug("ContainersController  /Delip")
	var body = make(map[string]string)
	myutils.CommctrBody(&this.Controller, &body)
	beegologs.Debug("body %v", body)

	err := models.ContainerDelIp(body["host"], body["containerid"])
	if err != nil {
		beegologs.Warn("容器（%v）删除 ip (%v) 失败，原因是：%v", body["containerid"], err.Error())
		this.Data["json"] = fmt.Sprintf("删除ip失败")
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}
}

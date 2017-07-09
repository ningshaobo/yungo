package controllers

import (
	"com.my/yungo/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	//	"github.com/astaxie/beego/httplib"
	dmModels "com.my/dmSvrWeb/models"
	beegologs "github.com/astaxie/beego/logs"
)

// Operations about UserNormal
type UsersNomalController struct {
	beego.Controller
}

// 定义 系统初始化 传输 body
type createAdmin struct {
	ZoneName   string
	TenantName string
	Username   string
	Password   string
}

// @Title zoneadmin
// @Description zoneadmin 初始化 zone 的admin 账号
// @Success 200 {string} zoneadmin success
// @Failure 400 failure
// @router /zoneadmin [post]
func (this *UsersNomalController) Zoneadmin() {
	beegologs.Debug("Zoneadmin")
	_, err := models.ProxyHandler(&this.Controller, nil, true)
	if err != nil {
		err = fmt.Errorf("注册Zone管理员失败")
		beegologs.Warn(err.Error())
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}

}

// @Title login
// @Description login tenant into the system
// @Success 200 {string} login success
// @Failure 400 failure
// @router /login [post]
func (this *UsersNomalController) Login() {
	beegologs.Debug("Login")
	body, err := models.ProxyHandler(&this.Controller, models.UserLoginModify, true)
	//	body, err := models.ProxyHandler(&this.Controller, nil, true)
	if err != nil {
		err = fmt.Errorf("登录失败")
		beegologs.Warn(err.Error())
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = err.Error()
		this.ServeJSON()
		return
	}
	var user dmModels.User
	json.Unmarshal(body, &user)
	beegologs.Debug("Uuid = %v", user.Id)
	this.SetSession("uid", fmt.Sprint(user.Id))
}

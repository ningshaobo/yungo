package controllers

import (
	"fmt"
	//	"reflect"
	//	"encoding/json"
	"com.my/yungo/models"
	myutils "com.my/utils"
	"github.com/astaxie/beego"
	beegologs "github.com/astaxie/beego/logs"
)

// Operations about Users
type ClusterController struct {
	beego.Controller
}

// 定义构建集群 传输 body
type createBody struct {
	Cluster models.Cluster
	Hosts   []models.Host
}

// @Title CreateCluster
// @Description 构建一个 swarm 集群，此处不关心 zone，即所有集群都在同一zone 下，且 ip addr 唯一
// @Param	body		body 	[]*models.Host	true		"body for 主机数组"
// @Success 200 {int} models.Cluster.Id
// @Failure 400 body is empty
// @router / [post]
func (this *ClusterController) Post() {
	var cluster *models.Cluster
	var hostList []models.Host
	var body createBody
	var err error
	myutils.CommctrBody(&this.Controller, &body)
	beegologs.Debug("body %v", body)
	cluster = &body.Cluster
	if cluster.Name == "" {
		err = fmt.Errorf("获取集群信息失败，原因是：集群信息为空")
		beegologs.Warn(err.Error())
		this.Data["json"] = err.Error()
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}

	hostList = body.Hosts
	if len(hostList) < 1 {
		err = fmt.Errorf("获取主机数组失败，原因是：数组为空")
		beegologs.Warn(err.Error())
		this.Data["json"] = err.Error()
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}

	// 插入集群信息到数据库中
	clusterId, err := models.ClusterInsertOne(models.ClusterSwarm, cluster)
	if err != nil {
		beegologs.Warn("插入主机集群信息失败，原因是：%v", err.Error())
		this.Data["json"] = fmt.Sprintf("插入主机集群信息失败，原因是：%v", err.Error())
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}

	// 构建集群
	var hosts []*models.Host
	for loop := 0; loop < len(hostList); loop++ {
		beegologs.Debug("get host = %v", hostList[loop])
		hosts = append(hosts, &hostList[loop])
	}
	err = models.ClusterCreate(hosts)
	if err != nil {
		models.ClusterDelOne(clusterId)
		beegologs.Warn("构建集群失败，原因是：%v", err.Error())
		this.Data["json"] = hosts //fmt.Sprintf("构建集群失败，原因是：%v", err.Error())
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}
	err = models.ClusterInsertAllHosts(cluster, hosts)
	if err != nil {
		// 此处要回退 集群数据库？？？
		models.ClusterDelOne(clusterId)
		beegologs.Warn("插入所有主机信息失败，原因是：%v", err.Error())
		this.Data["json"] = fmt.Sprintf("插入所有主机信息失败，原因是：%v", err.Error())
		this.Ctx.Output.SetStatus(400)
		this.ServeJSON()
		return
	}
	beegologs.Debug("cid = %v", clusterId)
	var retMap = make(map[string]interface{})
	retMap["Cluster"] = cluster
	retMap["Hosts"] = hosts

	this.Data["json"] = retMap
	this.ServeJSON()
}

// @Title GetAll
// @Description 获取所有集群信息
// @Success 200 {object} []*models.Cluster
// @router / [get]
func (this *ClusterController) GetAll() {
	clus := models.GetAllClusters()
	this.Data["json"] = clus
	this.ServeJSON()
}

//
//// @Title Get
//// @Description get user by uid
//// @Param	uid		path 	string	true		"The key for staticblock"
//// @Success 200 {object} models.User
//// @Failure 403 :uid is empty
//// @router /:uid [get]
//func (u *UserController) Get() {
//	uid := u.GetString(":uid")
//	if uid != "" {
//		user, err := models.GetUser(uid)
//		if err != nil {
//			u.Data["json"] = err.Error()
//		} else {
//			u.Data["json"] = user
//		}
//	}
//	u.ServeJSON()
//}
//
//// @Title Update
//// @Description update the user
//// @Param	uid		path 	string	true		"The uid you want to update"
//// @Param	body		body 	models.User	true		"body for user content"
//// @Success 200 {object} models.User
//// @Failure 403 :uid is not int
//// @router /:uid [put]
//func (u *UserController) Put() {
//	uid := u.GetString(":uid")
//	if uid != "" {
//		var user models.User
//		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
//		uu, err := models.UpdateUser(uid, &user)
//		if err != nil {
//			u.Data["json"] = err.Error()
//		} else {
//			u.Data["json"] = uu
//		}
//	}
//	u.ServeJSON()
//}
//
//// @Title Delete
//// @Description delete the user
//// @Param	uid		path 	string	true		"The uid you want to delete"
//// @Success 200 {string} delete success!
//// @Failure 403 uid is empty
//// @router /:uid [delete]
//func (u *UserController) Delete() {
//	uid := u.GetString(":uid")
//	models.DeleteUser(uid)
//	u.Data["json"] = "delete success!"
//	u.ServeJSON()
//}
//
//// @Title Login
//// @Description Logs user into the system
//// @Param	username		query 	string	true		"The username for login"
//// @Param	password		query 	string	true		"The password for login"
//// @Success 200 {string} login success
//// @Failure 403 user not exist
//// @router /login [get]
//func (u *UserController) Login() {
//	username := u.GetString("username")
//	password := u.GetString("password")
//	if models.Login(username, password) {
//		u.Data["json"] = "login success"
//	} else {
//		u.Data["json"] = "user not exist"
//	}
//	u.ServeJSON()
//}
//
//// @Title logout
//// @Description Logs out current logged in user session
//// @Success 200 {string} logout success
//// @router /logout [get]
//func (u *UserController) Logout() {
//	u.Data["json"] = "logout success"
//	u.ServeJSON()
//}

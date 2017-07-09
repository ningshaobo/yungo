package controllers

import (
	//	"net/http"
	"com.my/yungo/models"
	"github.com/astaxie/beego"
)

type ProxyController struct {
	beego.Controller
}

func (this *ProxyController) Get() {
	models.ProxyHandler(&this.Controller, nil, false)
}
func (this *ProxyController) Post() {
	models.ProxyHandler(&this.Controller, nil, false)
}

func (this *ProxyController) Put() {
	models.ProxyHandler(&this.Controller, nil, false)
}

func (this *ProxyController) Delete() {
	models.ProxyHandler(&this.Controller, nil, false)
}

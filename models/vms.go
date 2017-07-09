package models

import (
	"github.com/astaxie/beego/orm"
)

// 虚拟机定义
type Vm struct {
	Id       	int64
	Uuid     	string
	Hostaddr 	string
	Host   		*Host		`orm:"default(0);rel(fk);on_delete(set_default)"`
	Vmmacs   	[]*Vmmac	`orm:"reverse(many)"` // 一对多反向
}

// 定义 Vm 联合索引
func (u *Vm) TableIndex() [][]string {
	return [][]string{
		[]string{"Uuid", "Hostaddr"},
	}
}

/**
* @param  初始化虚拟机 数据库
 */
func VmsRegisterDB() {
	//注册 model
	orm.RegisterModel(new(Vm))
}

/**
* @param 根据参数 生成 mac 网卡字符串
xml := `<interface type='ethernet' name='ethTest0'><mac address='` + mac + `'/></interface>`
*/
func VmsMacTapXml(hostInterface string, macAddr string) string {
	baseIpXmlStr := `<interface type='direct'><mac address='` + macAddr + `'/><source dev='` + hostInterface + `'/><start mode='onboot'/></interface>`
	return baseIpXmlStr
}

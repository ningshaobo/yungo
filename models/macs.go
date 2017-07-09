package models

import (
	"fmt"
	beegologs "github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// 定义 mac 地址 对应信息
type Vmmac struct {
	Id      int64
	Macaddr string `orm:"unique"`
	Vm      *Vm    `orm:"default(0);rel(fk);on_delete(set_default)"` // 一对多
}

/**
* @param  初始化mac地址 数据库
 */
func MacsRegisterDB() {
	//注册 model
	orm.RegisterModel(new(Vmmac))
}

/**
* @param 获取某个虚拟机上所有的物理网卡信息
 */
func MacsAllotMac(host string, domuuid string) (*Vmmac, error) {
	var mac Vmmac
	var vm Vm
	var err error

	o := orm.NewOrm()

	//	_, err = o.QueryTable(new(Vmmac)).Filter("Vm__Id__exact", 0).All(&mac)
	o.Raw("select * from vmmac where vm_id = ?", 0).QueryRow(&mac)
	if err != nil {
		beegologs.Warn("获取 mac 失败，原因是: %v", err.Error())
		return &mac, err
	}

	vm.Hostaddr = host
	vm.Uuid = domuuid
	_, _, err = o.ReadOrCreate(&vm, "Hostaddr", "Uuid")
	if err != nil {
		beegologs.Warn("获取 mac 失败，原因是: %v", err.Error())
		return &mac, err
	}

	if mac.Id == 0 {
		// 数据库中找不到空闲mac， 重新生成
		var newMac string
		for loop := 0; loop < 10; loop++ {
			// 防冲突， 循环三次还冲突，异常返回
			newMac = UtilsMacGen()
			exist := o.QueryTable(new(Vmmac)).Filter("macaddr__iexact", newMac).Exist()
			if exist == false {
				mac.Macaddr = newMac
				break
			} else {
				if loop > 3 {
					err = fmt.Errorf("随机生成mac失败, 冲突次数超过 3 次")
					beegologs.Error(err.Error())
					return &mac, err
				}
			}
		}
		beegologs.Debug("new mac = %v", mac.Macaddr)
	}
	mac.Vm = &vm
	_, _, err = o.ReadOrCreate(&mac, "Macaddr", "VM")
	if err != nil {
		beegologs.Warn("获取 mac 失败，原因是: %v", err.Error())
		return &mac, err
	}

	return &mac, nil
}

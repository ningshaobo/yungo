package models

import (
	myUtils "com.my/dmSvrWeb/utils"
	"database/sql"
	"fmt"
	beegologs "github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

/**
* @param  初始化数据库
 */
func SystemDbInit(dbNAme string, dbConn string) error {
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Printf("%v \n", err.Error())
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + dbNAme + " DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci")
	if err != nil {
		fmt.Printf("%v 库已经存在， %v  \n", dbNAme, err.Error())
		return err
	} else {
		fmt.Printf("创建数据库（%s）完成", dbNAme)
	}
	return nil
}

/**
* @param  创建集群
 */
///
func systemSshCmd(host *Host) error {
	var err error
	// 删除 主机所有 yum 源配置
	err = myUtils.SshCommand("rm -rf /etc/yum.repos.d/*",
		"root", host.Password, host.Ipaddr, 22)
	if err != nil {
		beegologs.Warn("SSh(%v)删除 repo 文件失败，原因是：%v", host.Ipaddr, err.Error())
		return err
	}
	// 拷贝yum 配置文件
	err = myUtils.SshFileCopy("D:/golang_tools/gopath/src/com.my/dmSvrWeb/static/my.repo", "/etc/yum.repos.d/",
		"root", host.Password, host.Ipaddr, 22)
	if err != nil {
		beegologs.Warn("SSh(%v)打开sftp连接失败，原因是：%v", host.Ipaddr, err.Error())
		return err
	}
	// 执行 docker 安装
	err = myUtils.SshCommand("yum install docker-engine docker-engine-selinux -y ",
		"root", host.Password, host.Ipaddr, 22)
	if err != nil {
		beegologs.Warn("SSh(%v)安装 docker失败，原因是：%v", host.Ipaddr, err.Error())
		return err
	}
	// 执行 docker 启动
	err = myUtils.SshCommand("service docker restart ",
		"root", host.Password, host.Ipaddr, 22)
	if err != nil {
		beegologs.Warn("SSh(%v)启动 docker失败，原因是：%v", host.Ipaddr, err.Error())
		return err
	}
	// 拷贝 daemon.json 配置文件
	err = myUtils.SshFileCopy("D:/golang_tools/gopath/src/com.my/dmSvrWeb/static/daemon.json", "/etc/docker/",
		"root", host.Password, host.Ipaddr, 22)
	if err != nil {
		beegologs.Warn("SSh(%v)打开sftp连接失败，原因是：%v", host.Ipaddr, err.Error())
		return err
	}
	// 执行 docker 重新启动，启用新 daemon.json 配置
	err = myUtils.SshCommand("service docker restart ",
		"root", host.Password, host.Ipaddr, 22)
	if err != nil {
		beegologs.Warn("SSh(%v)重新启动 docker失败，原因是：%v", host.Ipaddr, err.Error())
		return err
	}
	myUtils.SshCommand("docker swarm leave --force ", "root", host.Password, host.Ipaddr, 22)

	return nil
}

/**
* @param  初始化 swarm 集群 leader 节点
 */
func SystemInitLeader(hosts []*Host) (int, error) {
	var leaderHostIndex = -1
	var err error

	for index := 0; index < len(hosts); index++ {
		host := hosts[index]
		if host.Type != HostManager || host.Type != HostLeader {
			continue
		}
		// 给主机安装 docker， 并启动 docker
		err = systemSshCmd(host)
		if err != nil {
			beegologs.Warn("主机（%v）安装 docker 失败，原因是： %v", host.Ipaddr, err.Error())
			host.Type = HostErr
			continue
		}
		_, err := ClusterCreateLeader(host)
		if err != nil {
			beegologs.Warn("主机（%v）初始化为 Leader 失败，原因是： %v", host.Ipaddr, err.Error())
			host.Type = HostErr
			continue
		}
		leaderHostIndex = index
		host.Type = HostLeader
		break
	}
	if leaderHostIndex < 0 {
		err = fmt.Errorf("集群有找到管理节点")
		beegologs.Warn("%v", err.Error())
		return leaderHostIndex, err
	}
	return leaderHostIndex, nil
}

package models

import (
	"fmt"
	beegologs "github.com/astaxie/beego/logs"
	orm "github.com/astaxie/beego/orm"
	swarmtype "github.com/docker/docker/api/types/swarm"
	dockercontext "golang.org/x/net/context"
	"sync"
	"time"
)

// 定义host 类型
const (
	HostErr     = -1 // -1- 异常主机
	HostNon     = 0  // 0-- 未加入集群 ，还是在集群内，只是退列了
	HostWorker  = 1  // 1-- work，
	HostManager = 2  // 2-- manager
	HostLeader  = 3  // 3-- leader
)

// 定义 host
type Host struct {
	Id        int64
	Describle string     `orm:"null;type(text)"`
	Uname     string     `orm:"null"`
	Password  string     `orm:"null"`
	Ipaddr    string     `orm:"unique"`
	Nodeuuid  string     `orm:"null"`
	Type      int        // 0-- 未加入集群 ， 1-- work，2-- manager
	Cluster   *Cluster   `orm:"rel(fk)"`                          // `orm:"rel(fk)"` //设置一对多关系
	Hostgroup *Hostgroup `orm:"null;rel(fk);on_delete(set_null)"` // `orm:"rel(fk)"` //设置一对多关系
}

// 定义 host 组， 暂时不考虑 ningshb，2017-06-19
type Hostgroup struct {
	Id        int64
	Name      string  `orm:"unique"`
	Describle string  `orm:"null;type(text)"`
	Host      []*Host `orm:"reverse(many)"`
}

// 定义集群类型
const (
	ClusterErr   = -1 // -1- 异常集群
	ClusterBare  = 0  // 0-- 裸机集群 ，
	ClusterVm    = 1  // 1-- 虚拟机集群，
	ClusterSwarm = 2  // 2-- swarm集群，
	ClusterK8s   = 3  // 2-- K8s集群
)

// 集群定义
type Cluster struct {
	Id        int64
	Name      string `orm:"unique"`
	Describle string `orm:"null;type(text)"`
	IsShare   bool   // 是否共享，共有集群 可共享， 私有集群，tenant 管理员设置后可共享
	Type      int    // -1， 异常； 0，裸机集群；1，虚拟机集群；2，swarm集群; 3,k8s集群
	Created   time.Time
	Updated   time.Time
	Hosts     []*Host `orm:"reverse(many)"` // 设置一对多的反向关系
}

/**
* @param  初始化集群 数据库
 */
func ClustersRegisterDB() {
	//注册 model
	orm.RegisterModel(new(Cluster), new(Host), new(Hostgroup))
}

/**
* @param  插入一条记录
 */
func ClusterInsertOne(clutype int, cluster *Cluster) (int64, error) {
	oDm := orm.NewOrm()
	// 插入集群 数据
	cluster.Created = time.Now()
	cluster.Updated = time.Now()
	cluster.Type = clutype
	if cluster.Describle == "" {
		cluster.Describle = "初始化集群"
	}
	cid, err := oDm.Insert(cluster)
	if err != nil {
		beegologs.Warn("集群数据插入失败，原因是：%v", err.Error())
		return cid, err
	}

	return cid, nil
}

/**
* @param  插入一条记录
 */
func ClusterDelOne(cid int64) error {
	oDm := orm.NewOrm()
	var cluster Cluster
	cluster.Id = cid
	_, err := oDm.Delete(&cluster)
	if err != nil {
		beegologs.Warn("集群数据删除失败，原因是：%v", err.Error())
		return err
	}

	return nil
}

/**
* @param  搜索所有集群
 */
func GetAllClusters() []*Cluster {
	var clus []*Cluster
	o := orm.NewOrm()
	num, err := o.QueryTable(new(Cluster)).RelatedSel().All(&clus)
	beegologs.Info("Returned Rows Num: %v, %v", num, err)
	return clus
}

/**
* @param  搜索所有 共享 集群
 */
func GetAllShareClusters() ([]*Cluster, error) {
	var clus []*Cluster
	var err error
	o := orm.NewOrm()
	num, err := o.QueryTable(new(Cluster)).RelatedSel().Filter("IsShare", true).All(&clus)
	if err != nil {
		beegologs.Warn("搜索所有共享集群失败，原因是： %v", err.Error())
	}
	beegologs.Info("搜索所有共享集群  Num: %v, %v", num, err)
	return clus, err
}

/**
* @param  根据主机类型搜索 主机
 */
func GetAllHostsByType(hostType int) []Host {
	var hosts []Host
	o := orm.NewOrm()
	var err error
	var num int64
	if hostType < 0 {
		// -1 搜索所有主机，不区分类型
		num, err = o.QueryTable("host").All(&hosts)
	} else {
		num, err = o.QueryTable("host").Filter("Type__Type", hostType).All(&hosts)
	}

	beegologs.Info("type: %v, Returned Rows Num:  %v, %v", hostType, num, err)
	return hosts
}

/**
* @param  根据主机所属集群 Id， type 类型 搜索 主机
*         搜索 >= hostType 主机数组
 */
func ClusterGetAllHostsByClusterId(clusterId int64, hostType int) ([]*Host, error) {
	var hosts []*Host
	var err error
	if clusterId < 0 {
		beegologs.Warn("Cid 不存在")
		err = fmt.Errorf("Cid 不存在")
		return hosts, err
	}

	o := orm.NewOrm()
	num, err := o.QueryTable(new(Host)).Filter("Cluster__Id", clusterId).Filter("Type__gte", hostType).All(&hosts)
	if err != nil {
		beegologs.Warn("根据主机所属集群 Id， type 类型 搜索 主机失败，原因是：%v", err.Error())
		return hosts, err
	}
	beegologs.Debug("type: %v, Returned Rows Num:  %v, %v", clusterId, num, err)

	return hosts, nil
}

/**
* @param  创建集群 leader
 */
func ClusterCreateLeader(host *Host) (string, error) {
	beegologs.Debug("createLeader host = %v", host)
	var req swarmtype.InitRequest
	var err error
	cli, err := SwarmInitClient(host.Ipaddr)
	if err != nil {
		beegologs.Warn("leader （%v） client 初始化失败，原因是：%v", host.Ipaddr, err.Error())
		return "", err
	}
	defer cli.Close()
	req.ForceNewCluster = true
	req.ListenAddr = host.Ipaddr
	req.AdvertiseAddr = host.Ipaddr
	nodeUUid, err := cli.SwarmInit(dockercontext.Background(), req)
	if err != nil {
		beegologs.Warn("Leader （%v）创建失败，原因是：%v", host.Ipaddr, err.Error())
		return "", err
	}
	return nodeUUid, nil
}

/**
* @param  加入一个主机到集群
 */
func ClusterJoinOneHost(host *Host, remoteAdd string, swarmInfo swarmtype.Swarm) error {
	var err error
	beegologs.Debug("主机（%v - %v）开始加入集群 %v", host.Ipaddr, host.Type, swarmInfo.ID)

	if !(host.Type == HostManager || host.Type == HostWorker) {
		err = fmt.Errorf("host（%v）不是 Manager或 worker", host.Ipaddr)
		beegologs.Warn(err.Error())
		return err
	}

	cli, err := SwarmInitClient(host.Ipaddr)
	defer cli.Close()
	if err != nil {
		beegologs.Warn("主机（%v）init client 失败，原因是：%v", host.Ipaddr, err.Error())
		return err
	}
	// 不管是否加入过集群， 强行退出
	err = cli.SwarmLeave(dockercontext.Background(), true)
	if err != nil {
		beegologs.Debug("host (%v) 强行退出 swarm 失败。原因是：%v", host, err.Error())
	}
	// 加入集群
	var req swarmtype.JoinRequest
	if host.Type == HostManager {
		req.JoinToken = swarmInfo.JoinTokens.Manager
	} else if host.Type == HostWorker {
		req.JoinToken = swarmInfo.JoinTokens.Worker
	} else {
		err = fmt.Errorf("host（%v）不是 Manager或 worker", host.Ipaddr)
		beegologs.Warn(err.Error())
		return err
	}
	req.ListenAddr = host.Ipaddr
	req.AdvertiseAddr = host.Ipaddr
	req.RemoteAddrs = append(req.RemoteAddrs, remoteAdd)
	err = cli.SwarmJoin(dockercontext.Background(), req)
	if err != nil {
		beegologs.Warn("host（%v）加入集群（%v）失败，原因是：%v", host.Ipaddr, swarmInfo.ID, err.Error())
		return err
	}
	return nil
}

/**
* @param  所有主机加入到集群
 */
func ClusterJoinAllHosts(hosts []*Host) error {
	var err error
	leaderIndex := -1
	for index := 0; index < len(hosts); index++ {
		host := hosts[index]
		if host.Type == HostLeader {
			leaderIndex = index
			break
		}
	}
	if leaderIndex < 0 {
		err = fmt.Errorf("集群中没有找到 leader 主机")
		beegologs.Warn("%v", err.Error())
		return err
	}
	// 集群 client
	cli, err := SwarmInitClient(hosts[leaderIndex].Ipaddr)
	defer cli.Close()
	if err != nil {
		beegologs.Warn("集群init client 失败，原因是：%v", err.Error())
		return err
	}
	// 获取集群信息，
	swarmInfo, err := cli.SwarmInspect(dockercontext.Background())
	if err != nil {
		beegologs.Warn("获取集群信息 失败，原因是：%v", err.Error())
		return err
	}
	var wg sync.WaitGroup
	for index := 0; index < len(hosts); index++ {
		if index == leaderIndex {
			continue
		}
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			host := hosts[index]
			err = ClusterJoinOneHost(host, hosts[leaderIndex].Ipaddr, swarmInfo)
			if err != nil {
				beegologs.Warn("host（%s）加入集群失败，原因是：%s", host.Ipaddr, err.Error())
				host.Type = HostErr
			}
		}(index)
	}
	wg.Wait()
	return nil
}

/**
* @param  插入所有主机信息到数据库中
 */
func ClusterInsertAllHosts(cluster *Cluster, hosts []*Host) error {
	var err error
	oDm := orm.NewOrm()
	qs, err := oDm.QueryTable(new(Host)).PrepareInsert()
	if err != nil {
		return err
	}
	defer qs.Close()
	for _, host := range hosts {
		host.Cluster = cluster

		beegologs.Debug("join host = %v", host)
		id, err := qs.Insert(host)
		if err != nil {
			beegologs.Warn("插入 (%v) 失败，原因是： %v", host.Ipaddr, err.Error())
			host.Type = HostErr
			host.Id = 0
		} else {
			host.Id = id
		}
	}
	return nil
}

/**
* @param  初始化集群
 */
func ClusterCreate(hosts []*Host) error {
	for index := 0; index < len(hosts); index++ {
		host := hosts[index]
		beegologs.Debug("host tta = %v", host)
	}
	var err error
	var leaderIndex = -1
	for loop := 0; loop < len(hosts); loop++ {
		host := hosts[loop]
		beegologs.Debug("host = %v", host)
		// 给主机安装 docker， 并启动 docker
		//		err = systemSshCmd(&host)
		//		if err != nil {
		//			beegologs.Warn("主机（%v）安装 docker 失败，原因是： %v", host.Ipaddr, err.Error())
		//			host.Type = HostErr
		//			continue
		//		}
		// 如果不是 mananger 或者 leader， 将跳过，不作为 leader 节点初始化
		if !(host.Type == HostManager || host.Type == HostLeader) {
			continue
		}
		_, err = ClusterCreateLeader(host)
		if err != nil {
			host.Type = HostErr
			continue
		}
		leaderIndex = loop
		host.Type = HostLeader
		break
	}
	if leaderIndex < 0 {
		err = fmt.Errorf("集群 leader 构建失败")
		return err
	}

	for index := 0; index < len(hosts); index++ {
		host := hosts[index]
		beegologs.Debug("host ttb = %v", host)
	}
	// 所有主机 插入集群中
	err = ClusterJoinAllHosts(hosts)
	if err != nil {
		beegologs.Warn("所有主机 插入集群失败，原因是：%v", err.Error())
		return err
	}

	return nil
}

///**
//* @param  根据 cluster id ， 搜索 集群下所有主机
// */
//func ClustersAllHosts (cid int64) []Host {
//	var hosts []Host
//	o := orm.NewOrm()
//	num, err := o.QueryTable(new(Host)).Filter("Cluster__Id", cid).All(&hosts)
//	if err != nil {
//		beegologs.Warn("get hosts failure, err : %s", err.Error())
//		return hosts
//	}
//	beegologs.Debug("num of host : %d", num)
//	return hosts
//}

package models

import (
	"fmt"
	"github.com/astaxie/beego/cache"
	beegologs "github.com/astaxie/beego/logs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	dockerclient "github.com/docker/docker/client"
	dockercontext "golang.org/x/net/context"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type swmCliStrut struct {
	Actived    bool
	HostAddr   string
	HostClient *dockerclient.Client
}

var SwarmShareClutersClients = cache.NewMemoryCache()

/**
* @param
 */
func SwarmInit(host string) (*dockerclient.Client, error) {
	var cli *dockerclient.Client
	var err error
	beegologs.Info("主机 （%s）client 链接初始化开始", host)

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err = dockerclient.NewClient("tcp://"+host+":2375", "v1.27", nil, defaultHeaders)
	if err != nil {
		beegologs.Warn("主机（%s）创建 docker client 失败 ，原因是： %s", host, err.Error())
		return cli, err
	}
	_, err = cli.Ping(dockercontext.Background())
	if err != nil {
		beegologs.Warn("主机（%s）docker ping 失败 ，原因是： %s", host, err.Error())
		return cli, err
	}
	beegologs.Info("主机 （%s）client 链接初始化成功结束", host)
	return cli, nil
}

/**
* @param 初始化client list 数组
 */
func SwarmClientListInit(hosts []*Host) ([]swmCliStrut, error) {
	var swmCliList []swmCliStrut
	var wg sync.WaitGroup
	var err error

	if hosts == nil {
		err = fmt.Errorf("初始化Swarm client list失败")
		beegologs.Warn(err.Error())
		return swmCliList, err
	}
	swmCliList = make([]swmCliStrut, len(hosts))
	cliNum := 0
	for index := 0; index < len(hosts); index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			host := hosts[index]
			swmCliList[index].Actived = false
			swmCliList[index].HostAddr = host.Ipaddr
			cli, err := SwarmInit(host.Ipaddr)
			if err == nil {
				swmCliList[index].Actived = true
				swmCliList[index].HostClient = cli
				cliNum++
			}
		}(index)
	}
	wg.Wait()
	if cliNum == 0 {
		err = fmt.Errorf("初始化Swarm client list失败")
		beegologs.Warn(err.Error())
		return swmCliList, err
	}
	return swmCliList, nil
}

/**
* @param 初始化一个 缓存，默认是 memory
 */
func checkCache(cid string) (interface{}, error) {
	var sca interface{}
	var err error
	if SwarmShareClutersClients == nil {
		SwarmShareClutersClients = cache.NewMemoryCache()
	}
	sca = SwarmShareClutersClients.Get(cid)
	if sca == nil {
		beegologs.Warn("缓存中没有找到集群 （%s）对应的数据", cid)
		// 初始化集群 manager 的client 链接
		clusterId, _ := strconv.ParseInt(cid, 10, 64)
		hosts, err := ClusterGetAllHostsByClusterId(clusterId, HostManager)
		if len(hosts) == 0 {
			err = fmt.Errorf("集群（%v）搜索管理以上主机为空", cid)
			return nil, err
		}
		swmCliList, err := SwarmClientListInit(hosts)
		if err != nil {
			beegologs.Warn("集群 （%s）初始化Swarm client list失败", cid)
			return swmCliList, err
		}

		SwarmShareClutersClients.Put(cid, swmCliList, 1*time.Hour)
		return swmCliList, nil
	}

	return sca, err
}

/**
* @param
 */
func SwarmGetHostClient(cluIdStr string) *dockerclient.Client {
	var myClient *dockerclient.Client
	var swmCliArray []swmCliStrut

	sca, err := checkCache(cluIdStr)
	if err != nil {
		beegologs.Error("集群（%v）client获取失败，原因是：%v", cluIdStr, err.Error())
		return myClient
	}
	swmCliArray = reflect.ValueOf(sca).Interface().([]swmCliStrut)

	for _, swmCli := range swmCliArray {
		if swmCli.Actived != true {
			continue
		}
		cli := swmCli.HostClient
		_, err := cli.Ping(dockercontext.Background())
		if err != nil {
			beegologs.Warn("%s 链接失败", swmCli.HostAddr)
			swmCli.Actived = false
			/* 此处通知监控 */
			continue
		}
		beegologs.Info("%s ,作为操作链接", swmCli.HostAddr)
		myClient = swmCli.HostClient
		break
	}
	return myClient
}

/**
* @param 初始化 swarm host client 连接
 */
func SwarmInitClient(host string) (*dockerclient.Client, error) {
	var cli *dockerclient.Client
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := dockerclient.NewClient("tcp://"+host+":2375", "1.24", nil, defaultHeaders)
	if err != nil {
		beegologs.Warn("创建(%s) client 连接失败，原因是： %s", host, err.Error())
		return cli, err
	}
	_, err = cli.Ping(dockercontext.Background())
	if err != nil {
		beegologs.Warn("docker ping(%s) 失败，原因是： %s", host, err.Error())
		return cli, err
	}
	return cli, nil
}

/**
* @param 获取swarm Cid 对应集群下，所有主机 IP地址
 */
func SwarmGetAllHostAddr(cid string) []string {
	var hostAddr []string

	// 获取 cid 对应的 manager client
	cli := SwarmGetHostClient(cid)

	var options types.NodeListOptions
	options.Filters = filters.NewArgs()

	nodeList, err := cli.NodeList(dockercontext.Background(), options)
	if err != nil {
		return hostAddr
	}
	for _, node := range nodeList {
		if node.Status.State == swarm.NodeStateReady {
			hostAddr = append(hostAddr, node.Status.Addr)
		}
	}
	return hostAddr
}

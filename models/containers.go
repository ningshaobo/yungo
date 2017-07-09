package models

import (
	"sync"
	//	"strings"
	beegologs "github.com/astaxie/beego/logs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	dockercontext "golang.org/x/net/context"
	"strconv"
)

type containerType struct {
	types.Container
	HostAddr string
	Name     string
}

/**
* @param  依据 cid tid 搜索 集群下 所有当前租户的容器
 */
func ContainerByTidCid(tid string, cid string) []containerType {
	var cts []containerType
	//	hosts := SwarmGetAllHostAddr(cid)
	clustId, err := strconv.Atoi(cid)
	if err != nil {
		beegologs.Warn("cid 转 int64 失败， 原因是：%s", err.Error())
		return cts
	}
	hosts, err := ClusterGetAllHostsByClusterId(int64(clustId), HostNon)
	if err != nil {
		beegologs.Warn("获取主机集失败， 原因是：%s", err.Error())
		return cts
	}

	if len(hosts) < 1 {
		beegologs.Warn("找不到主机")
	}
	var wg sync.WaitGroup
	for index := 0; index < len(hosts); index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			host := hosts[index]
			cli, err := SwarmInitClient(host.Ipaddr)
			defer cli.Close()
			if err != nil {
				beegologs.Warn("host（%v）创建client连接失败, 原因是：", host, err.Error())
				return
			}
			var options types.ContainerListOptions
			options.Filters = filters.NewArgs()
			options.Filters.Add("label", "com.my.tid="+tid)
			containers, err := cli.ContainerList(dockercontext.Background(), options)
			if err != nil {
				beegologs.Warn("获取主机(%v)上的容器失败，原因是：%v", host, err.Error())
				return
			}
			if len(containers) < 1 {
				beegologs.Debug("获取主机(%v)上的容器个数为 0 ", host)
				return
			}
			//			var hostCts containerType
			for _, ct := range containers {
				var hostC containerType
				hostC.Container = ct
				hostC.Name = ct.Names[0]
				hostC.HostAddr = host.Ipaddr

				cts = append(cts, hostC)
			}
			//			cts = append(cts, containers...)
		}(index)
	}

	wg.Wait()
	return cts
}

/**
* @param  给容器添加 ip
 */
func ContainerAllotIp(host string, cuuid string, ip string) error {
	cli, err := SwarmInitClient(host)

	defer cli.Close()
	if err != nil {
		beegologs.Warn("host（%v）创建client连接失败, 原因是：", host, err.Error())
		return err
	}
	var netoptions types.NetworkListOptions
	netoptions.Filters = filters.NewArgs()
	netoptions.Filters.Add("name", "macnet")
	macnet, err := cli.NetworkList(dockercontext.Background(), netoptions)
	if err != nil {
		beegologs.Warn("获取主机(%v)上的mac网络失败，原因是：%v", host, err.Error())
		return err
	}
	beegologs.Debug(macnet[0].ID)

	var connconfig network.EndpointSettings
	connconfig.IPAddress = ip
	connconfig.NetworkID = macnet[0].ID
	var endpoint network.EndpointIPAMConfig
	endpoint.IPv4Address = ip
	connconfig.IPAMConfig = &endpoint
	err = cli.NetworkConnect(dockercontext.Background(), macnet[0].ID, cuuid, &connconfig)
	if err != nil {
		beegologs.Warn("分配IP(%v)失败，原因是：%v", ip, err.Error())
		return err
	}
	return nil
}

/**
* @param  给容器添加 ip
 */
func ContainerDelIp(host string, cuuid string) error {
	cli, err := SwarmInitClient(host)

	defer cli.Close()
	if err != nil {
		beegologs.Warn("host（%v）创建client连接失败, 原因是：", host, err.Error())
		return err
	}
	var netoptions types.NetworkListOptions
	netoptions.Filters = filters.NewArgs()
	netoptions.Filters.Add("name", "macnet")
	macnet, err := cli.NetworkList(dockercontext.Background(), netoptions)
	if err != nil {
		beegologs.Warn("获取主机(%v)上的mac网络失败，原因是：%v", host, err.Error())
		return err
	}
	beegologs.Debug(macnet[0].ID)

	cli.NetworkDisconnect(dockercontext.Background(), macnet[0].ID, cuuid, true)
	if err != nil {
		beegologs.Warn("容器（%v）删除IP失败，原因是：%v", cuuid, err.Error())
		return err
	}
	return nil
}

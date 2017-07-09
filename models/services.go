package models

import (
	myUtils "com.my/dmSvrWeb/utils"
	"fmt"
	"github.com/astaxie/beego"
	beegologs "github.com/astaxie/beego/logs"
	dockertype "github.com/docker/docker/api/types"
	dockerfilters "github.com/docker/docker/api/types/filters"
	swarmtype "github.com/docker/docker/api/types/swarm"
	dockerclient "github.com/docker/docker/client"
	dockercontext "golang.org/x/net/context"
	"sync"
)

type SvrContainer struct {
	Name       string
	Id         string
	Service    string
	Image      string
	Status     string
	CreateTime int64
	UpdateTime int64
}

// 定义 构建 swarm 服务 传输 body
type CreateSvrBody struct {
	Name     string
	Mode     string
	Replicas uint64
	Global   swarmtype.GlobalService
}
type CreateBaseBody struct {
	Image  string
	ImgTag string
}
type CreateInBody struct {
	SvrBody  CreateSvrBody
	BaseBody CreateBaseBody
}

/**
* @param
 */
func CreateService(tid string, clusterId string, crtBody CreateInBody) error {
	cli := SwarmGetHostClient(clusterId)
	var service swarmtype.ServiceSpec
	var options dockertype.ServiceCreateOptions

	service.Name = crtBody.SvrBody.Name

	if crtBody.SvrBody.Mode == "Global" {
		service.Mode.Global = &crtBody.SvrBody.Global
	} else if crtBody.SvrBody.Mode == "Replicas" {
		//		replicas, _		:= strconv.Atoi(crtBody.SvrBody.Replicas)
		ss := crtBody.SvrBody.Replicas //uint64(replicas)
		var replicated swarmtype.ReplicatedService
		replicated.Replicas = &ss
		service.Mode.Replicated = &replicated
	} else {
		return nil
	}

	regUrl := beego.AppConfig.String("registryUrl")
	service.TaskTemplate.ContainerSpec.Image = regUrl + crtBody.BaseBody.Image

	labels := make(map[string]string)
	labels["com.my.tid"] = tid
	service.Labels = labels

	svr, err := cli.ServiceCreate(dockercontext.Background(), service, options)
	if err != nil {
		beegologs.Warn("创建 service 失败 ： %s", err.Error())
		return err
	}
	beegologs.Info("%v ", svr)
	return nil
}

/**
* @param
 */
func DeleteService(serviceID string, clusterId string) error {
	//	_, err := InspectService(serviceID, tid, clusterId)
	//	if err != nil {
	//		return err
	//	}
	//
	cli := SwarmGetHostClient(clusterId)
	err := cli.ServiceRemove(dockercontext.Background(), serviceID)

	return err
}

/**
* @param
 */
func InspectService(serviceID string, clusterId string) (swarmtype.Service, error) {
	cli := SwarmGetHostClient(clusterId)
	service, _, err := cli.ServiceInspectWithRaw(dockercontext.Background(), serviceID)
	if err != nil {
		return service, err
	}
	tidLabelKey := myUtils.TidLabelKey()
	// 查找键值是否存在
	if _, ok := service.Spec.Labels[tidLabelKey]; ok {
		//		service.Spec.Labels[tidLabelKey]		= ""
		return service, nil
	} else {
		err = fmt.Errorf("%s 无租户信息", serviceID)
		beegologs.Warn(err.Error())
		return service, err
	}
}

/**
* @param
 */
func GetAllService(tid string, clusterId string) []swarmtype.Service {
	var svrList []swarmtype.Service
	var err error

	if tid == "" {
		beegologs.Error("tid 参数不能为空")
		return svrList
	}

	cli := SwarmGetHostClient(clusterId)
	if cli == nil {
		beegologs.Error("%s", "当前 swarm 没有主机")
		return svrList
	}
	var svrOptions dockertype.ServiceListOptions
	svrFilter := dockerfilters.NewArgs()
	tidLabelKey := myUtils.TidLabelKey()
	svrFilter.Add("label", tidLabelKey+"="+tid)
	svrList, err = cli.ServiceList(dockercontext.Background(), svrOptions)
	if err != nil {
		beegologs.Warn("%s", err)
		return svrList
	}
	beegologs.Info("%d", len(svrList))
	return svrList
}

/**
* @param
 */
func GetAllNodeContainers(serviceId string, nodeId string, clusterId string) ([]dockertype.Container, swarmtype.Node) {
	var containerList []dockertype.Container
	var nodeInfo swarmtype.Node
	var err error
	cli := SwarmGetHostClient(clusterId)

	// 获取 node info
	nodeInfo, _, err = cli.NodeInspectWithRaw(dockercontext.Background(), nodeId)
	if err != nil {
		return containerList, nodeInfo
	}
	if nodeInfo.Status.State != swarmtype.NodeStateReady {
		return containerList, nodeInfo
	}
	nodeAddr := nodeInfo.Status.Addr
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	nodeCli, err := dockerclient.NewClient("tcp://"+nodeAddr+":2375", "v1.27", nil, defaultHeaders)
	if err != nil {
		beegologs.Warn("%s - %s", nodeAddr, err.Error())
		return containerList, nodeInfo
	}

	//	csFilter := dockerfilters.NewArgs()
	var csOptions dockertype.ContainerListOptions
	if serviceId != "" {
		csOptions.Filters = dockerfilters.NewArgs()
		csOptions.Filters.Add("label", "com.docker.swarm.service.id="+serviceId)
	}

	containerList, err = nodeCli.ContainerList(dockercontext.Background(), csOptions)
	if err != nil {

	}
	return containerList, nodeInfo
}

/**
* @param
 */
type NodeInfoAndContainers struct {
	NodeId     string
	NodeInfo   swarmtype.Node
	Containers []dockertype.Container
}
type MoreTask struct {
	//	Task 				swarmtype.Task
	NodeInfo      swarmtype.Node
	ContainerInfo dockertype.Container
}

/*  */
func taskMoreInfo(tasks []swarmtype.Task, clusterId string, serviceId string) []MoreTask {
	var mTasks []MoreTask
	nodes := make(map[string]NodeInfoAndContainers)
	// 首先获取所有的node
	for _, tk := range tasks {
		if tk.Status.State != swarmtype.TaskStateFailed {
			nd := tk.NodeID
			_, ok := nodes[nd]
			if !ok {
				var ndStruct NodeInfoAndContainers
				ndStruct.NodeId = nd
				nodes[nd] = ndStruct
			}
		}
	}
	if len(nodes) < 1 {
		return mTasks
	}
	var wg sync.WaitGroup
	for Id, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			node.Containers, node.NodeInfo = GetAllNodeContainers(serviceId, Id, clusterId)
			nodes[Id] = node
			//			fmt.Printf("node %v \n", node)
		}()
	}
	wg.Wait()
	for _, nd := range nodes {
		//		fmt.Printf("%v \n", nd)
		for _, c := range nd.Containers {
			var mk MoreTask
			mk.NodeInfo = nd.NodeInfo
			mk.ContainerInfo = c
			mTasks = append(mTasks, mk)
		}
	}

	return mTasks
}

/**
* @param
 */
func GetAllMoreTasks(clusterId string, serviceId string) []MoreTask {
	var mTasks []MoreTask
	tasks := GetAllSvrTasks(clusterId, serviceId)
	if len(tasks) < 1 {
		return mTasks
	}
	mTasks = taskMoreInfo(tasks, clusterId, serviceId)
	return mTasks
}

/**
* @param
 */
func GetAllSvrTasks(clusterId string, serviceId string) []swarmtype.Task {
	var tasks []swarmtype.Task
	var err error
	if clusterId == "" || serviceId == "" {
		err = fmt.Errorf("clusterId（%s） 或者 serviceId（%s） 为空", clusterId, serviceId)
		beegologs.Warn("%s", err.Error())
		return tasks
	}
	cli := SwarmGetHostClient(clusterId)

	tksfilter := dockerfilters.NewArgs()
	var tksOptions dockertype.TaskListOptions

	tksfilter.Add("service", serviceId)
	tksOptions.Filters = tksfilter

	tasks, err = cli.TaskList(dockercontext.Background(), tksOptions)
	if err != nil {
		beegologs.Warn("%s", err.Error())
		return tasks
	}
	return tasks
}

/**
* @param
 */
func UpdateService(serviceID string, tid string, clusterId int64) bool {
	//	cli := SwarmGetHostClient(clusterId)
	//	var service swarmtype.ServiceSpec
	//	var options swarmtype.ServiceUpdateOptions
	//	cli.ServiceUpdate(dockercontext.Background(), serviceID, nil, service, options)
	return true
}

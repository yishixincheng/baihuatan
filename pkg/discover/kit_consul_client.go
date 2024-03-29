package discover

import (
	"baihuatan/pkg/common"
	"fmt"
	"log"
	"strconv"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

// NewDiscoveryClientInstance 创建服务发现实例
func NewDiscoveryClientInstance(consulHost string, consulPort string) *DiscoveryClientInstance {
	port, _ := strconv.Atoi(consulPort)
	// 通过 Consul Host 和 Consul Port 创建一个 consul.Client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(port)
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil
	}

	client := consul.NewClient(apiClient)

	return &DiscoveryClientInstance{
		Host:   consulHost,
		Port:   port,
		config: consulConfig,
		client: client,
	}
}

// Register 向Consul中心注册实例
func (consulClient *DiscoveryClientInstance) Register(instanceID, svcHost, healthCheckURL, svcPort string, 
														svcName string, 
														weight int, 
														meta map[string]string, 
														tags []string, 
														logger *log.Logger) bool {
	port, _ := strconv.Atoi(svcPort)
	serviceRegistration := &api.AgentServiceRegistration{
		ID:		   instanceID,
		Name:      svcName,
		Address:   svcHost,
		Port:      port,
		Meta:      meta,
		Tags:      tags,
		Weights:   &api.AgentWeights{
			Passing: weight,
		},
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + svcHost + ":" + strconv.Itoa(port) + healthCheckURL,
			Interval:                       "15s",
		},
	}
	// 2. 发送服务注册到 Consul中
	err := consulClient.client.Register(serviceRegistration)

	if err != nil {
		if logger != nil {
			fmt.Printf("错误：%v \n", err)
			fmt.Println("http://" + svcHost + ":" + strconv.Itoa(port) + healthCheckURL)
			logger.Println("Register Service Error!")
		}
		return false
	}
	if logger != nil {
		logger.Println("Register Service Success!")
	}
	return true
}

// DeRegister 注销实例
func (consulClient *DiscoveryClientInstance) DeRegister(instanceID string, logger *log.Logger) bool {
	// 构建包含服务实例 ID 的源数结构体
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceID,
	}
	err := consulClient.client.Deregister(serviceRegistration)

	if err != nil {
		if logger != nil {
			logger.Println("DeRegister Service Error!")
		}
		return false
	}
	if logger != nil {
		logger.Println("DeRegister Service Success!")
	}
	return true
}

// DiscoverServices 发现服务
func (consulClient *DiscoveryClientInstance) DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance {
	// 该服务已监控并缓存
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}
	// 申请锁
	consulClient.mutex.Lock()
	// 再次检查是否监控
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}
	// 注册监控
	go func ()  {
		params := make(map[string]interface{})
		params["type"] = "service"
		params["service"] = serviceName
		plan, _ := watch.Parse(params)
		plan.Handler = func(u uint64, i interface{}) {
			if i == nil {
				return
			}
			v, ok := i.([]*api.ServiceEntry)
			if !ok {
				return  // 数据异常，忽略
			}
			// 没有服务实例在线
			if len(v) == 0 {
				consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
				return
			}

			var healthServices []*common.ServiceInstance
			for _, service := range v {
				if service.Checks.AggregatedStatus() == api.HealthPassing {
					healthServices = append(healthServices, newServiceInstance(service.Service))
				}
			}
			consulClient.instancesMap.Store(serviceName, healthServices)
		}
		defer plan.Stop()
		plan.Run(consulClient.config.Address)
	}()
	defer consulClient.mutex.Unlock()

	// 根据服务名请求服务实例列表
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
		if logger != nil {
			logger.Println("Discover Service Error!")
		}
		return nil
	}
	instances := make([]*common.ServiceInstance, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = newServiceInstance(entries[i].Service)
	}
	consulClient.instancesMap.Store(serviceName, instances)
	return instances
}

// newServiceInstance 创建服务实例
func newServiceInstance(service *api.AgentService) *common.ServiceInstance {
	rpcPort := service.Port -1
	if service.Meta != nil {
		if rpcPortString, ok := service.Meta["rpcPort"]; ok {
			rpcPort, _ = strconv.Atoi(rpcPortString)
		}
	}
	return &common.ServiceInstance{
		Host:    service.Address,
		Port:    service.Port,
		GrpcPort: rpcPort,
		Weight:  service.Weights.Passing,
	}
}
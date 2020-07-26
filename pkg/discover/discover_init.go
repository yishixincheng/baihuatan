package discover

import (
	"baihuatan/pkg/bootstrap"
	"baihuatan/pkg/common"
	"baihuatan/pkg/loadbalance"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	uuid "github.com/satori/go.uuid"
)

// ConsulService 服务
var ConsulService DiscoveryClient
// LoadBalance  负载均衡
var LoadBalance loadbalance.LoadBalance
// Logger 日志
var Logger *log.Logger
// ErrNoInstanceExisted 实例不存在错误
var ErrNoInstanceExisted  = errors.New("no available client")

func init() {
	// 1.实例化一个 Consul 客户端，此处实例化了原生态实现版本
	ConsulService = NewDiscoveryClientInstance(bootstrap.DiscoverConfig.Host, bootstrap.DiscoverConfig.Port)
	LoadBalance   = new(loadbalance.WeightRoundRobinLoadBalance)
	Logger        = log.New(os.Stderr, "", log.LstdFlags)
}

// CheckHealth 健康检查
func CheckHealth(writer http.ResponseWriter, reader *http.Request) {
	Logger.Println("Health check!")
	_, err := fmt.Fprintln(writer, "Server is OK!")
	if err != nil {
		Logger.Println(err)
	}
}

// DiscoveryService 发现服务
func DiscoveryService(serviceName string) (*common.ServiceInstance, error) {
	instances := ConsulService.DiscoverServices(serviceName, Logger)

	if len(instances) < 1 {
		Logger.Printf("no available client for %s.", serviceName)
		return nil, ErrNoInstanceExisted
	}
	return LoadBalance.SelectService(instances) 
}

// Register 注册服务
func Register() {
	// 实例失败，停止服务
	if ConsulService == nil {
		panic(0)
	}
	// 判空 instanceId，通过 go.uuid 获取一个服务实例ID
	instanceID := bootstrap.DiscoverConfig.InstanceID

	if instanceID == "" {
		instanceID = bootstrap.DiscoverConfig.ServiceName + uuid.NewV4().String()
	}

	if !ConsulService.Register(instanceID, bootstrap.HTTPConfig.Host, "/health",
	   bootstrap.HTTPConfig.Port, bootstrap.DiscoverConfig.ServiceName,
	   bootstrap.DiscoverConfig.Weight,
	   map[string]string{
		   "rpcPort" : bootstrap.RPCConfig.Port,
	   }, nil, Logger) {
		   Logger.Printf("register services %s failed.", bootstrap.DiscoverConfig.ServiceName)
		   // 注册失败，服务启动失败
		   panic(0)
	   }
	   Logger.Printf(bootstrap.DiscoverConfig.ServiceName + "-service for service %s success.", bootstrap.DiscoverConfig.ServiceName)
}

// Deregister 注销服务
func Deregister()  {
	if ConsulService == nil {
		panic(0)
	}
	instanceID := bootstrap.DiscoverConfig.InstanceID

	if instanceID == "" {
		instanceID = bootstrap.DiscoverConfig.ServiceName + "-" + uuid.NewV4().String()
	}
	if !ConsulService.DeRegister(instanceID, Logger) {
		Logger.Printf("deregister for service %s failed.", bootstrap.DiscoverConfig.ServiceName)
		panic(0)
	}
}
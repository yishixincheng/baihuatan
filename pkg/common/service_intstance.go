package common

// ServiceInstance 服务实例结构体定义
type ServiceInstance struct {
	Host      string   // Host
	Port      int      // Port
	Weight    int      // 权重
	CurWeight int      // 当前权重
	GrpcPort  int      // Grpc端口  
}
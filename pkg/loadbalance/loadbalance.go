package loadbalance

import (
	"baihuatan/pkg/common"
	"errors"
	"math/rand"
)

// LoadBalance 负载均衡器
type LoadBalance interface {
	SelectService(service []*common.ServiceInstance) (*common.ServiceInstance, error)
}

// RandomLoadBalance 随机负载
type RandomLoadBalance struct {	
}

// SelectService 随机负载实现
func (loadBalance *RandomLoadBalance) SelectService(services []*common.ServiceInstance) (*common.ServiceInstance, error) {
	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}
	return services[rand.Intn(len(services))], nil
}

// WeightRoundRobinLoadBalance 轮询调度负载
type WeightRoundRobinLoadBalance struct {
}

// SelectService 轮询调度负载实现
func (loadBalance *WeightRoundRobinLoadBalance) SelectService(services []*common.ServiceInstance) (*common.ServiceInstance, error) {
	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}
	total := 0
	var best *common.ServiceInstance
	for i := 0; i < len(services); i++ {
		w := services[i]
		if w ==  nil {
			continue
		}
		w.CurWeight += w.Weight
		total += w.Weight

		if best == nil || w.CurWeight > best.CurWeight {
			best = w
		}
	}

	if best == nil {
		return nil, nil
	}
	best.CurWeight -= total
	return best, nil
}

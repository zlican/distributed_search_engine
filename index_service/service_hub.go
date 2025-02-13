package indexservice

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

const (
	SERVICE_ROOT_PATH = "/engine" //区分不同的业务
	SERVICE_NAME      = "index_service"
	SERVICE_ADDR      = "127.0.0.1:8080"
	SERVICE_HEARTBEAT = 30
)

// 服务注册中心
type ServiceHub struct {
	client             *etcdv3.Client //在etcd.Client之上的实例化 （核心）
	heartbeatFrequency int64          //心跳频率
	watched            sync.Map
	loadBalancer       LoadBalancer //负载均衡策略
}

var (
	serviceHub *ServiceHub
	hubOnce    sync.Once
)
var etcdServers = []string{"localhost:2379", "localhost:22379", "localhost:32379"}

// 启动etcd.Client
func GetServiceHub(etcdServers []string, heartbeatFrequency int64) *ServiceHub {
	if serviceHub == nil {
		hubOnce.Do(func() {
			if client, err := etcdv3.New(etcdv3.Config{
				Endpoints:   etcdServers,
				DialTimeout: 5 * time.Second,
			}); err != nil {
				panic(err)
			} else {
				serviceHub = &ServiceHub{
					client:             client,
					heartbeatFrequency: heartbeatFrequency,
					loadBalancer:       &RoundRobin{},
				}
			}
		})
	}
	return serviceHub
}

// 注册自己服务（核心Client.Put）
func (hub *ServiceHub) Regist(serviceName string, endpoint string, leaseID etcdv3.LeaseID) (etcdv3.LeaseID, error) {
	ctx := context.Background()
	if leaseID <= 0 {
		//申请租约
		if lease, err := hub.client.Grant(ctx, hub.heartbeatFrequency); err != nil {
			return 0, err
		} else {
			key := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + serviceName + "/" + endpoint
			if _, err := hub.client.Put(ctx, key, endpoint, etcdv3.WithLease(lease.ID)); err != nil {
				return 0, err
			} else {
				return lease.ID, nil
			}
		}
	} else {
		//续约
		if _, err := hub.client.KeepAlive(ctx, leaseID); err == rpctypes.ErrLeaseNotFound {
			return hub.Regist(serviceName, endpoint, 0)
		} else if err != nil {
			return 0, err
		} else {
			return leaseID, nil
		}
	}
}

// 注销自己的服务(核心Client.Delete)
func (hub *ServiceHub) Unregist(serviceName string, endpoint string, leaseID etcdv3.LeaseID) error {
	ctx := context.Background()
	key := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + serviceName + "/" + endpoint
	if _, err := hub.client.Delete(ctx, key); err != nil {
		return err
	} else {
		fmt.Println("注销成功", serviceName)
		return nil
	}
}

// 获取服务列表(核心Client.Get)
func (hub *ServiceHub) GetServiceEndpoints(serviceName string) ([]string, error) {
	ctx := context.Background()
	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + serviceName + "/"
	if resp, err := hub.client.Get(ctx, prefix, etcdv3.WithPrefix()); err != nil { //前缀
		return nil, err
	} else {
		endpoints := make([]string, 0)
		for _, kv := range resp.Kvs {
			path := strings.Split(string(kv.Key), "/")
			endpoints = append(endpoints, path[len(path)-1]) //拿到服务的ip列表
		}
		return endpoints, nil
	}
}

// ServiceHub实现负载均衡
// 策略模式：根据不同的负载均衡策略，选择不同的负载均衡策略（函数）
func (hub *ServiceHub) GetServiceEndpoint(serviceName string) (string, error) {
	if endpoints, err := hub.GetServiceEndpoints(serviceName); err != nil {
		return "", err
	} else {
		return hub.loadBalancer.Take(endpoints), nil
	}
}

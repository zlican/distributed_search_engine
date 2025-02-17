package indexservice

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	etcdv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/time/rate"
)

// 定义一个接口，用于直接使用ServiceHub，或使用HubProxy
type IServiceHub interface {
	Regist(serviceName string, endpoint string, leaseID etcdv3.LeaseID) (etcdv3.LeaseID, error)
	Unregist(serviceName string, endpoint string, leaseID etcdv3.LeaseID) error
	GetServiceEndpoints(serviceName string) []string
	GetServiceEndpoint(serviceName string) string
}

type HubProxy struct {
	*ServiceHub   //继承，可以将类型当成名称调用父辈函数
	endpointCache sync.Map
	limiter       *rate.Limiter
	loadBalance   LoadBalancer
}

var (
	hubproxy  *HubProxy
	proxyOnce sync.Once
)

// 单例模式
func GetHubProxy(etcdServers []string, heartbeatFrequency int64, qps int) *HubProxy {
	if hubproxy == nil {
		proxyOnce.Do(func() {
			serviceHub := GetServiceHub(etcdServers, heartbeatFrequency)
			hubproxy = &HubProxy{
				ServiceHub:    serviceHub,
				endpointCache: sync.Map{},
				limiter:       rate.NewLimiter(rate.Every(time.Duration(1e9/qps)*time.Nanosecond), qps), //每秒产生qps个令牌
				loadBalance:   &Random{},
			}
		})
	}
	return hubproxy
}

//func (proxy *HubProxy) Register(serviceName string, endpoint string, leaseID etcdv3.LeaseID) (etcdv3.LeaseID, error) {
//return proxy.hub.Regist(serviceName, endpoint, leaseID)
//}

//func (proxy *HubProxy) Unregist(serviceName string, endpoint string, leaseID etcdv3.LeaseID) error {
//return proxy.hub.Unregist(serviceName, endpoint, leaseID)
//}

// proxy服务：限流 + 缓存	//重写
func (proxy *HubProxy) GetServiceEndpoints(serviceName string) []string {
	if !proxy.limiter.Allow() { //限流保护
		return nil
	}
	proxy.watchEndpointsOfService(serviceName) // 监听服务变化
	if endpoints, exists := proxy.endpointCache.Load(serviceName); exists {
		return endpoints.([]string)
	}

	endpoints := proxy.ServiceHub.GetServiceEndpoints(serviceName)
	if len(endpoints) > 0 {
		proxy.endpointCache.Store(serviceName, endpoints)

	}
	return endpoints
}

// proxy服务： 监听 核心函数 etcd.Watch
func (proxy *HubProxy) watchEndpointsOfService(serviceName string) {
	if _, exists := proxy.ServiceHub.watched.LoadOrStore(serviceName, true); exists {
		return
	} //当第一次监听时，往watched中添加serviceName，值为true，当第二次监听时，已存在，不进行监听
	ctx := context.Background()
	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + serviceName + "/"
	ch := proxy.ServiceHub.client.Watch(ctx, prefix, etcdv3.WithPrefix()) //在etcd上根据前缀进行监听，没一个修改都会放进ch
	fmt.Println("开始监听", serviceName)
	go func() {
		for wresp := range ch {
			for _, event := range wresp.Events {
				path := strings.Split(string(event.Kv.Key), "/")
				if len(path) > 2 {
					serviceName := path[len(path)-2] //获得变化事件的serviceName
					endpoints := proxy.ServiceHub.GetServiceEndpoints(serviceName)
					fmt.Println("监听服务变化", serviceName, event.Type, endpoints)
					if len(endpoints) > 0 {
						proxy.endpointCache.Store(serviceName, endpoints)
					} else {
						proxy.endpointCache.Delete(serviceName)
					}
				}
			}
		}
	}()
}

//func (proxy *HubProxy) GetServiceEndpoint(serviceName string) string {
//	return proxy.loadBalance.Take(proxy.GetServiceEndpoints(serviceName))
//}

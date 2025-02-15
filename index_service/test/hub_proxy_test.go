package test

import (
	indexservice "engine/index_service"
	"fmt"
	"testing"
	"time"
)

func TestHubProxy(t *testing.T) {
	const qps = 100
	hubProxy := indexservice.GetHubProxy([]string{"127.0.0.1:2379"}, 10, qps)

	leaseID1, err := hubProxy.Regist("index_service", "127.0.0.1:8080", 0)
	if err != nil {
		t.Fatalf("注册服务失败: %v", err)
	}
	defer hubProxy.Unregist("index_service", "127.0.0.1:8080", leaseID1)

	time.Sleep(100 * time.Millisecond)

	endpoints := hubProxy.GetServiceEndpoints("index_service")
	if len(endpoints) == 0 {
		t.Error("预期获取到服务节点，但是返回空")
	}
	fmt.Println("第一次注册后的节点:", endpoints)

	leaseID2, err := hubProxy.Regist("index_service", "127.0.0.2:8080", 0)
	if err != nil {
		t.Fatalf("注册服务失败: %v", err)
	}
	defer hubProxy.Unregist("index_service", "127.0.0.2:8080", leaseID2)

	time.Sleep(100 * time.Millisecond)

	endpoints = hubProxy.GetServiceEndpoints("index_service")
	if len(endpoints) != 2 {
		t.Error("预期获取到两个服务节点，但是实际节点数:", len(endpoints))
	}
	fmt.Println("第二次注册后的节点:", endpoints)

	time.Sleep(1 * time.Second)
	successCount := 0
	for i := 0; i < qps+5; i++ {
		endpoint := hubProxy.GetServiceEndpoints("index_service")
		if len(endpoint) > 0 {
			successCount++
		}
	}
	fmt.Printf("限流测试：成功请求数 %d/%d\n", successCount, qps+5)
}

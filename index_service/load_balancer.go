package indexservice

import (
	"sync/atomic"
	"time"

	"golang.org/x/exp/rand"
)

type LoadBalancer interface {
	Take([]string) string
}

// 负载均衡策略1————轮询：次数对台数求模
type RoundRobin struct {
	acc int64
}

func (b *RoundRobin) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	n := atomic.AddInt64(&b.acc, 1) //支持并发调用
	return endpoints[n%int64(len(endpoints))]
}

// 负载均衡策略2————随机
type Random struct {
}

func (b *Random) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	rand.Seed(uint64(time.Now().UnixNano()))
	return endpoints[rand.Intn(len(endpoints))]
}

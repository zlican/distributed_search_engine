package indexservice

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	types "github.com/zlican/engine/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// 分布式：调用方使用对象，哨兵
type Sentinel struct {
	hub      IServiceHub
	connPool sync.Map //连接池: grpc连接，与每个worker建立的连接，将连接缓存起来
}

func NewSentinel(etcdServers []string) *Sentinel {
	return &Sentinel{
		hub:      GetHubProxy(etcdServers, 3, 100),
		connPool: sync.Map{},
	}
}

// 获取grpc连接
func (s *Sentinel) GetGrpcConn(endpoint string) *grpc.ClientConn {
	if v, exists := s.connPool.Load(endpoint); exists {
		conn := v.(*grpc.ClientConn)
		if conn.GetState() == connectivity.TransientFailure || conn.GetState() == connectivity.Shutdown {
			conn.Close()
			s.connPool.Delete(endpoint)
		} else {
			return conn
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil
	}
	s.connPool.Store(endpoint, conn)
	return conn
}

func (sentinel *Sentinel) AddDoc(doc types.Document) (int, error) {
	endpoint := sentinel.hub.GetServiceEndpoint(INDEX_SERVICE)
	if len(endpoint) == 0 {
		return 0, errors.New("no index service available")
	}
	conn := sentinel.GetGrpcConn(endpoint)
	if conn == nil {
		return 0, errors.New("failed to get grpc connection")
	}
	client := NewIndexServiceClient(conn)
	affected, err := client.AddDoc(context.Background(), &doc)
	if err != nil {
		return 0, err
	}
	return int(affected.Count), nil
}

// 删除文档：不知道文档在哪个worker上，需要遍历所有worker，都进行删除操作
func (sentinel *Sentinel) DeleteDoc(docId string) int {
	endpoints := sentinel.hub.GetServiceEndpoints(INDEX_SERVICE)
	if len(endpoints) == 0 {
		return 0
	}
	var n int32
	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))
	for _, endpoint := range endpoints {
		go func(endpoint string) {
			defer wg.Done()
			conn := sentinel.GetGrpcConn(endpoint)
			if conn != nil {
				client := NewIndexServiceClient(conn)
				affected, err := client.DeleteDoc(context.Background(), &DocID{DocID: docId})
				if err != nil {
					fmt.Println("delete doc failed, err:", err)
				} else {
					if affected.Count > 0 {
						atomic.AddInt32(&n, affected.Count)
						fmt.Println("delete doc success, count:", affected.Count)
					}
				}
			}

		}(endpoint)
	}
	wg.Wait()
	return int(atomic.LoadInt32(&n))
}

func (sentinel *Sentinel) Search(query *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []*types.Document {
	endpoints := sentinel.hub.GetServiceEndpoints(INDEX_SERVICE)
	if len(endpoints) == 0 {
		return nil
	}
	docs := make([]*types.Document, 0, 1000)
	resultCh := make(chan *types.Document, 1000) //多个协程共同写入
	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))
	for _, endpoint := range endpoints {
		go func(endpoint string) {
			defer wg.Done()
			conn := sentinel.GetGrpcConn(endpoint)
			if conn != nil {
				client := NewIndexServiceClient(conn)
				results, err := client.Search(context.Background(), &SearchRequest{
					Query:   query,
					OnFlag:  onFlag,
					OffFlag: offFlag,
					OrFlags: orFlags,
				})
				if err != nil {
					fmt.Println("search failed, err:", err)
				} else {
					if len(results.Results) > 0 {
						for _, doc := range results.Results {
							resultCh <- doc
						}
					}
				}
			}
		}(endpoint)
	}
	receiveFinish := make(chan struct{})
	go func() {
		for {
			doc, ok := <-resultCh
			if !ok {
				break
			}
			docs = append(docs, doc)
		}
		receiveFinish <- struct{}{}
	}()
	wg.Wait()
	close(resultCh)
	<-receiveFinish
	return docs

}

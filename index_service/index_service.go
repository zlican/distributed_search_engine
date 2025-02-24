package indexservice

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	types "github.com/zlican/engine/types"

	"google.golang.org/grpc"
)

const (
	INDEX_SERVICE = "index_service"
)

// 分布式：worker方使用对象
type IndexServiceWorker struct {
	Indexer                         *Indexer
	hub                             *ServiceHub
	selfAddr                        string //grpc服务地址 endpoint
	UnimplementedIndexServiceServer        //grpc结构体
}

func (service *IndexServiceWorker) Init(DocNumEstimate int, dbtype int, DataDir string, etcdServers []string, grpcPort int) error {
	service.Indexer = new(Indexer)
	service.Indexer.Init(DocNumEstimate, dbtype, DataDir)

	//向注册中心注册自己
	if len(etcdServers) > 0 {
		if grpcPort <= 1024 {
			return fmt.Errorf("grpcPort must be greater than 1024")
		}
		//selfLocalIP, err := utils.GetLocalIP()
		//if err != nil {
		//	return err
		//}
		selfLocalIP := "127.0.0.1"
		service.selfAddr = selfLocalIP + ":" + strconv.Itoa(grpcPort)

		var heartBeat int64 = 3
		hub := GetServiceHub(etcdServers, heartBeat)
		leaseId, err := hub.Regist(INDEX_SERVICE, service.selfAddr, 0)
		if err != nil {
			return err
		}
		service.hub = hub
		go service.StartGRPCServer()
		go func() {
			for {
				hub.Regist(INDEX_SERVICE, service.selfAddr, leaseId)
				time.Sleep(time.Duration(heartBeat)*time.Second - 100*time.Millisecond) //续命
			}
		}()
	}

	return nil
}

// 从正排数据加载到倒排索引中
func (service *IndexServiceWorker) LoadFormIndexFile() int {
	return service.Indexer.LoadFormIndexFile()
}

// 关闭服务
func (service *IndexServiceWorker) Close() error {
	if service.hub != nil {
		service.hub.Unregist(INDEX_SERVICE, service.selfAddr, 0)
	}
	return service.Indexer.Close()
}

// 添加文档
func (service *IndexServiceWorker) AddDoc(ctx context.Context, doc *types.Document) (*AffectedCount, error) {
	affected, err := service.Indexer.AddDoc(*doc)
	if err != nil {
		return nil, err
	}
	return &AffectedCount{Count: int32(affected)}, nil
}

// 从索引上删除文档
func (service *IndexServiceWorker) DeleteDoc(ctx context.Context, docId *DocID) (*AffectedCount, error) {
	return &AffectedCount{Count: int32(service.Indexer.DeleteDoc(docId.DocID))}, nil
}

func (service *IndexServiceWorker) Search(ctx context.Context, request *SearchRequest) (*SearchResult, error) {
	results := service.Indexer.Search(request.Query, request.OnFlag, request.OffFlag, request.OrFlags)
	return &SearchResult{Results: results}, nil
}

// 启动 gRPC 服务器
func (service *IndexServiceWorker) StartGRPCServer() error {
	lis, err := net.Listen("tcp", service.selfAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	RegisterIndexServiceServer(grpcServer, service)

	fmt.Printf("Starting gRPC server on %s\n", service.selfAddr)
	return grpcServer.Serve(lis)
}

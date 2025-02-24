package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	indexservice "github.com/zlican/engine/index_service"
	"github.com/zlican/engine/internal/kvdb"
	"github.com/zlican/engine/types"
)

type SearchRequest struct {
	Keyword string   `json:"keyword"`
	Field   string   `json:"field"`
	OnFlag  uint64   `json:"onFlag"`
	OffFlag uint64   `json:"offFlag"`
	OrFlags []uint64 `json:"orFlags"`
}

type SearchResponse struct {
	Documents []*types.Document `json:"documents"`
	Error     string            `json:"error,omitempty"`
}

var (
	// 单机模式的索引器
	indexer *indexservice.Indexer
	// 分布式模式的哨兵
	sentinel *indexservice.Sentinel
)

func main() {
	// 初始化sentinel（分布式模式）
	sentinel = indexservice.NewSentinel([]string{"localhost:2379"})

	// 初始化索引器（单机模式）
	indexer = new(indexservice.Indexer)
	err := indexer.Init(1000, kvdb.BADGER, "./data")
	if err != nil {
		log.Fatal("初始化索引器失败:", err)
	}
	defer indexer.Close()

	// 初始化示例文档
	initializeExampleDocuments()

	// 处理静态文件
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	// API路由
	http.HandleFunc("/api/documents", handleDocuments)
	http.HandleFunc("/api/documents/", handleDocumentByID)
	http.HandleFunc("/api/search", handleSearch)

	log.Println("服务器启动在 :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	var doc types.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var affected int
	var err error

	// 尝试使用分布式模式
	if sentinel != nil {
		affected, err = sentinel.AddDoc(doc)
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"affected": affected})
			return
		}
		// 如果分布式模式失败，回退到单机模式
		log.Printf("分布式模式添加文档失败: %v, 尝试使用单机模式", err)
	}

	// 使用单机模式
	affected, err = indexer.AddDoc(doc)
	if err != nil {
		http.Error(w, fmt.Sprintf("添加文档失败: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"affected": affected})
}

func handleDocumentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 从URL中提取文档ID
	docID := strings.TrimPrefix(r.URL.Path, "/api/documents/")
	if docID == "" {
		http.Error(w, "无效的文档ID", http.StatusBadRequest)
		return
	}

	var affected int

	// 尝试使用分布式模式
	if sentinel != nil {
		affected = sentinel.DeleteDoc(docID)
	} else {
		// 使用单机模式
		affected = indexer.DeleteDoc(docID)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"affected": affected})
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 构建查询，将关键词转换为正确的格式
	keyword := &types.Keyword{
		Field: req.Field,
		Word:  req.Keyword,
	}

	query := &types.TermQuery{
		Should:  nil,
		Must:    nil,
		Keyword: keyword.ToString(),
	}

	var results []*types.Document

	// 尝试使用分布式模式
	if sentinel != nil {
		results = sentinel.Search(query, req.OnFlag, req.OffFlag, req.OrFlags)
	}

	// 如果分布式模式没有返回结果或没有启用分布式模式，使用单机模式
	if results == nil || len(results) == 0 {
		results = indexer.Search(query, req.OnFlag, req.OffFlag, req.OrFlags)
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SearchResponse{
		Documents: results,
	})
}

func initializeExampleDocuments() {
	documents := []types.Document{
		{
			Id:          "doc1",
			BitsFeature: 1,
			Keywords: []*types.Keyword{
				{Field: "title", Word: "搜索引擎"},
				{Field: "content", Word: "分布式"},
				{Field: "tag", Word: "技术"},
			},
		},
		{
			Id:          "doc2",
			BitsFeature: 2,
			Keywords: []*types.Keyword{
				{Field: "title", Word: "Go语言"},
				{Field: "content", Word: "编程"},
				{Field: "tag", Word: "技术"},
			},
		},
		{
			Id:          "doc3",
			BitsFeature: 4,
			Keywords: []*types.Keyword{
				{Field: "title", Word: "机器学习"},
				{Field: "content", Word: "人工智能"},
				{Field: "tag", Word: "AI"},
			},
		},
		{
			Id:          "doc4",
			BitsFeature: 8,
			Keywords: []*types.Keyword{
				{Field: "title", Word: "数据库"},
				{Field: "content", Word: "分布式"},
				{Field: "tag", Word: "技术"},
			},
		},
		{
			Id:          "doc5",
			BitsFeature: 16,
			Keywords: []*types.Keyword{
				{Field: "title", Word: "云计算"},
				{Field: "content", Word: "分布式"},
				{Field: "tag", Word: "技术"},
			},
		},
	}

	for _, doc := range documents {
		var affected int
		var err error

		// 尝试使用分布式模式
		if sentinel != nil {
			affected, err = sentinel.AddDoc(doc)
			if err == nil {
				log.Printf("分布式模式添加文档 %s 成功, 影响行数: %d", doc.Id, affected)
				continue
			}
			log.Printf("分布式模式添加文档失败: %v, 尝试使用单机模式", err)
		}

		// 使用单机模式
		affected, err = indexer.AddDoc(doc)
		if err != nil {
			log.Printf("添加文档 %s 失败: %v", doc.Id, err)
			continue
		}
		log.Printf("单机模式添加文档 %s 成功, 影响行数: %d", doc.Id, affected)
	}
}

# 分布式搜索引擎

这是一个基于 Go 语言实现的分布式搜索引擎，主要特点包括：

## 核心特性

- **分布式架构**: 采用 etcd 作为服务注册中心，支持多节点部署和负载均衡
- **高性能索引**: 基于跳表(Skip List)实现的倒排索引，支持高效的查询操作
- **复杂查询**: 支持 AND/OR 逻辑组合的多条件查询
- **并发控制**: 实现了自定义的并发 Map 结构，优化并发性能
- **持久化存储**: 支持 Badger 和 Bolt 两种存储引擎

## 系统架构

### 1. 索引结构

- **正排索引**: 存储完整的文档信息
- **倒排索引**: 基于跳表实现，支持高效的关键词检索
- **位图过滤**: 支持通过位图实现文档过滤

### 2. 分布式组件

- **服务注册**: 基于 etcd 实现服务注册与发现
- **负载均衡**: 支持轮询等负载均衡策略
- **RPC 通信**: 使用 gRPC 实现节点间通信

## 主要功能

### 1. 文档管理

- 添加文档
- 删除文档
- 批量操作支持

### 2. 搜索功能

- 关键词搜索
- 复杂逻辑组合查询
- 位图过滤支持

### 3. 分布式特性

- 多节点部署
- 服务自动发现
- 负载均衡
- 容错处理

## 目录结构

```bash
├── engine/
│ ├── internal/
│ │ ├── kvdb/ # 存储引擎实现
│ │ └── reverse_index/ # 倒排索引实现
│ ├── index_service/ # 分布式服务实现
│ └── utils/ # 工具函数
└── types/ # 数据类型定义
```

## 快速开始

### 安装

```bash
git clone https://github.com/zlican/engine.git
cd engine
go mod download
```

### 启动单机版本

```go
indexer := new(Indexer)
indexer.Init(1000, kvdb.BADGER, "./data")
```

### 启动分布式节点

```go
worker := new(IndexServiceWorker)
worker.Init(1000, kvdb.BADGER, "./data", []string{"localhost:2379"}, 8081)
```

## 使用示例

### 添加文档

```go
doc := types.Document{
    Id: "1",
    Keywords: []*types.Keyword{{Field: "title", Word: "golang"}},
    BitsFeature: 1,
}
indexer.AddDoc(doc)
```

### 搜索文档

```go
query := types.Be("golang").And(types.Be("database"))
results := indexer.Search(query, 0, 0, nil)
```

## 性能优化

1. 使用跳表实现倒排索引，提供 O(log n) 的查询性能
2. 自定义并发 Map 结构，优化并发访问性能
3. 批量操作支持，减少网络开销
4. 位图过滤，提供高效的文档过滤机制

## 注意事项

1. 确保 etcd 服务正常运行（分布式模式下）
2. 合理配置存储引擎参数
3. 注意并发操作的数据一致性

## 贡献指南

欢迎提交 Issue 和 Pull Request。在提交 PR 前，请确保：

1. 代码已经格式化
2. 测试用例已经补充
3. 文档已经更新

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

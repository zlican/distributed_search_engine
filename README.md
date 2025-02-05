# 分布式搜索引擎

这是一个基于跳表实现的分布式搜索引擎项目。该项目包含以下主要功能：

- **跳表交集计算**：通过跳表实现多个集合的交集计算。
- **反向索引**：使用跳表实现反向索引，支持文档的添加和删除。

## 目录结构

```
.
├── engine
│   ├── internal
│   │   └── reverse_index.go
│   ├── test
│   │   ├── skiplist_intersec_test.go
│   │   └── reverse_index_test.go
│   └── utils
│       └── skiplist_intersection.go
└── types
    └── document.go
```

## 安装

确保您已经安装了 Go 语言环境。然后克隆此仓库：

```sh
git clone https://github.com/yourusername/distributed-search-engine.git
cd distributed-search-engine
```

## 使用方法

### 跳表交集计算

在 `engine/test/skiplist_intersec_test.go` 文件中有一个测试用例 `TestSkipListIntersection`，用于测试跳表交集计算功能。您可以运行以下命令来执行测试：

```sh
go test ./engine/test -run TestSkipListIntersection
```

### 反向索引

在 `engine/internal/test/reverse_index_test.go` 文件中有两个测试用例 `TestSkipReverseIndexAdd` 和 `TestSkipReverseIndexDelete`，用于测试反向索引的添加和删除功能。您可以运行以下命令来执行测试：

```sh
go test ./engine/internal/test -run TestSkipReverseIndexAdd
go test ./engine/internal/test -run TestSkipReverseIndexDelete
```

## 运行项目

要运行整个项目，您可以使用以下命令：

```sh
go run main.go
```

请确保在运行项目之前，已经正确配置了所有依赖项。

## 贡献

欢迎提交问题和拉取请求。如果您有任何建议或改进，请随时提出。

## 许可证

此项目使用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

```

```

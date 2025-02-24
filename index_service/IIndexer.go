package indexservice

import types "github.com/zlican/engine/types"

// 单机与分布式公用的抽象接口
type IIndexer interface {
	AddDoc(doc types.Document) (int, error)
	DeleteDoc(docId string) int
	Search(query *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []*types.Document
}

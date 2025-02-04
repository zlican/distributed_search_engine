package reverseindex

import "engine/types"

type IReverseIndexer interface {
	Add(doc types.Document)
	Delete(IntId uint64, keyword *types.Keyword)
}

package test

import (
	reverseindex "engine/internal/reverse_index"
	"engine/types"
	"fmt"
	"testing"
)

func TestSkipReverseIndexAdd(t *testing.T) {
	doc := types.Document{Id: "001", IntId: uint64(1), BitsFeature: uint64(88),
		Keywords: []*types.Keyword{{Field: "title", Word: "main"}, {Field: "title", Word: "主要"}},
		Bytes:    []byte{131}}

	reverseIndex := reverseindex.NewSkipListReverseIndex(3)
	reverseIndex.Add(doc)
}

func TestSkipReverseIndexDelete(t *testing.T) {
	doc := types.Document{Id: "001", IntId: uint64(1), BitsFeature: uint64(88),
		Keywords: []*types.Keyword{{Field: "title", Word: "main"}, {Field: "title", Word: "主要"}},
		Bytes:    []byte{131}}

	reverseIndex := reverseindex.NewSkipListReverseIndex(3)
	reverseIndex.Add(doc)
	reverseIndex.Delete(doc.IntId, doc.Keywords[1])
}

func TestMine(t *testing.T) {
	var str = "陈子陵"
	fmt.Println([]byte(str))

}

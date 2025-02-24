package test

import (
	"fmt"
	"testing"

	reverseindex "github.com/zlican/engine/internal/reverse_index"
	"github.com/zlican/engine/types"
)

func TestSkipReverseIndexAdd(t *testing.T) {
	doc := types.Document{Id: "001", IntId: uint64(1), BitsFeature: uint64(1),
		Keywords: []*types.Keyword{{Field: "title", Word: "main"}, {Field: "title", Word: "主要"}},
		Bytes:    []byte{131}}

	reverseIndex := reverseindex.NewSkipListReverseIndex(3)
	reverseIndex.Add(doc)
	fmt.Println(reverseIndex.Search(&types.TermQuery{Keyword: doc.Keywords[0].ToString()}, uint64(0), uint64(0), []uint64{uint64(1)}))
}

func TestSkipReverseIndexDelete(t *testing.T) {
	doc := types.Document{Id: "001", IntId: uint64(1), BitsFeature: uint64(88),
		Keywords: []*types.Keyword{{Field: "title", Word: "main"}, {Field: "title", Word: "主要"}},
		Bytes:    []byte{131}}

	reverseIndex := reverseindex.NewSkipListReverseIndex(3)
	reverseIndex.Add(doc)
	reverseIndex.Delete(doc.IntId, doc.Keywords[1])
	fmt.Println(reverseIndex.Search(&types.TermQuery{Keyword: doc.Keywords[0].Word}, 0, 0, nil))
}

func TestMine(t *testing.T) {
	var str = "陈子陵"
	fmt.Println([]byte(str))

}

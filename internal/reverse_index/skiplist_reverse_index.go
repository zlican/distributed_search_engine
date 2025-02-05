package reverseindex

import (
	"engine/types"
	"engine/utils"
	"fmt"
	"runtime"
	"sync"

	"github.com/huandu/skiplist"
)

type SkipListReverseIndex struct { //倒排索引整体上是一个大map，里面的keyword维护一个跳表
	table *utils.MyCurrencyMap
	locks []sync.RWMutex
}

type MixValue struct {
	Id          string
	BitsFeature uint64
}

func NewSkipListReverseIndex(DocNumEstimate int) *SkipListReverseIndex { //创建一个倒排索引对象
	indexer := new(SkipListReverseIndex)
	indexer.table = utils.NewMyCurrencyMap(runtime.NumCPU(), DocNumEstimate)
	indexer.locks = make([]sync.RWMutex, 1000)
	return indexer
}

func (indexer *SkipListReverseIndex) Add(doc types.Document) {
	mixValue := MixValue{Id: doc.Id, BitsFeature: doc.BitsFeature}

	for i, key := range doc.Keywords {
		if list, exists := indexer.table.Get(key.ToString()); exists {
			list.(*skiplist.SkipList).Set(doc.IntId, mixValue) //获得对应map_key所在的跳表，加入key和mixValue
			fmt.Println("Add successfully", key.ToString(), i)
		} else {
			SkipListMap := skiplist.New(skiplist.Uint64)
			SkipListMap.Set(doc.IntId, mixValue)
			indexer.table.Set(key.ToString(), SkipListMap)
			fmt.Println("Add successfully", key.ToString(), i)
		}
	}
}

func (indexer *SkipListReverseIndex) Delete(IntId uint64, keyword *types.Keyword) {
	if list, exists := indexer.table.Get(keyword.ToString()); exists {
		result := list.(*skiplist.SkipList).Remove(IntId)
		if result != nil {
			fmt.Println("remove success", result.Key(), keyword.ToString())
		}
	} else {
		fmt.Println("remove fail", keyword.ToString)
	}
}

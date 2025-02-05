package reverseindex

import (
	"engine/types"
	"engine/utils"
	"fmt"
	"runtime"
	"sync"

	"github.com/huandu/skiplist"
	farmhash "github.com/leemcloughlin/gofarmhash"
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

func (indexer *SkipListReverseIndex) getLock(key string) *sync.RWMutex { //使用固定数量的锁
	n := int(farmhash.Hash32WithSeed([]byte(key), 0))
	return &indexer.locks[n%len(indexer.locks)]
}

func (indexer *SkipListReverseIndex) Add(doc types.Document) {
	mixValue := MixValue{Id: doc.Id, BitsFeature: doc.BitsFeature}

	for i, keyword := range doc.Keywords {
		key := keyword.ToString()
		lock := indexer.getLock(key)
		lock.Lock()
		if list, exists := indexer.table.Get(key); exists {
			list.(*skiplist.SkipList).Set(doc.IntId, mixValue) //获得对应map_key所在的跳表，加入key和mixValue
			fmt.Println("Add successfully", key, i)
		} else {
			SkipListMap := skiplist.New(skiplist.Uint64)
			SkipListMap.Set(doc.IntId, mixValue)
			indexer.table.Set(key, SkipListMap)
			fmt.Println("Add successfully", key, i)
		}
		lock.Unlock()
	}
}

func (indexer *SkipListReverseIndex) Delete(IntId uint64, keyword *types.Keyword) {
	key := keyword.ToString()
	lock := indexer.getLock(key)
	lock.Lock()
	if list, exists := indexer.table.Get(key); exists {
		result := list.(*skiplist.SkipList).Remove(IntId)
		if result != nil {
			fmt.Println("remove success", result.Key(), key)
		}
	} else {
		fmt.Println("remove fail", key)
	}
	lock.Unlock()
}

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

func NewSkipListReverseIndex(DocNumEstimate int) *SkipListReverseIndex {
	indexer := new(SkipListReverseIndex)
	indexer.table = utils.NewMyCurrencyMap(runtime.NumCPU(), DocNumEstimate)
	indexer.locks = make([]sync.RWMutex, 1000)
	return indexer
} //创建一个倒排索引对象

func (indexer *SkipListReverseIndex) getLock(key string) *sync.RWMutex {
	n := int(farmhash.Hash32WithSeed([]byte(key), 0))
	return &indexer.locks[n%len(indexer.locks)]
} //使用固定数量的锁

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

func (indexer SkipListReverseIndex) FilterByBits(bits uint64, onFlag uint64, offFlag uint64, orFlags []uint64) bool {
	//必须满足onFlag
	if bits&onFlag != onFlag {
		return false
	}
	//必须不满足offFlag
	if bits&offFlag != 0 {
		return false
	}

	//必须满足orFlags中的任意一个
	for _, orFlag := range orFlags {
		if orFlag > 0 && bits&orFlag <= 0 {
			return false
		}
	}

	return true
} //按照bits特征进行过滤

func (indexer SkipListReverseIndex) search(q *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) *skiplist.SkipList {
	if q.Keyword != "" {
		Keyword := q.Keyword
		if value, exists := indexer.table.Get(Keyword); exists {
			result := skiplist.New(skiplist.Uint64)
			list := value.(*skiplist.SkipList)
			node := list.Front()
			for node != nil {
				intId := node.Key().(uint64)
				skv, _ := node.Value.(MixValue)
				flag := skv.BitsFeature
				if intId > 0 && indexer.FilterByBits(flag, onFlag, offFlag, orFlags) { //过滤
					result.Set(intId, skv)
				}
				node = node.Next()
			}
			return result
		}
	} else if len(q.Must) > 0 {
		results := make([]*skiplist.SkipList, 0, len(q.Must))
		for _, q := range q.Must {
			results = append(results, indexer.search(q, onFlag, offFlag, orFlags)) //递归实现
		}
		return utils.SkipIntersection(results...) //&是多个跳表求交集
	} else if len(q.Should) > 0 {
		results := make([]*skiplist.SkipList, 0, len(q.Should))
		for _, q := range q.Should {
			results = append(results, indexer.search(q, onFlag, offFlag, orFlags))
		}
		return utils.SkipUnion(results...) //|是多个跳表求并集
	}
	return nil
} //内部search

func (indexer SkipListReverseIndex) Search(q *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []string {
	result := indexer.search(q, onFlag, offFlag, orFlags)
	if result == nil {
		return nil
	}
	ids := make([]string, 0, result.Len())
	node := result.Front()
	for node != nil {
		ids = append(ids, node.Value.(MixValue).Id)
		node = node.Next()
	}
	return ids
} //外部search

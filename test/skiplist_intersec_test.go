package test

import (
	"fmt"
	"testing"

	"github.com/zlican/engine/utils"

	"github.com/huandu/skiplist"
)

func TestSkipListIntersection(t *testing.T) {
	list1 := skiplist.New(skiplist.Uint64)
	list1.Set(uint64(5), uint64(0))
	list1.Set(uint64(1), uint64(0))
	list1.Set(uint64(8), uint64(0))
	list1.Set(uint64(11), uint64(0))
	list1.Set(uint64(2), uint64(0))
	list1.Set(uint64(7), uint64(0))

	list2 := skiplist.New(skiplist.Uint64)
	list2.Set(uint64(8), uint64(0))
	list2.Set(uint64(11), uint64(0))
	list2.Set(uint64(3), uint64(0))
	list2.Set(uint64(9), uint64(0))
	list2.Set(uint64(2), uint64(0))

	list3 := skiplist.New(skiplist.Uint64)
	list3.Set(uint64(2), uint64(0))
	list3.Set(uint64(11), uint64(0))
	list3.Set(uint64(8), uint64(0))
	list3.Set(uint64(9), uint64(0))

	nodes := utils.SkipIntersection(list1, list2, list3)
	node := nodes.Front()
	for {
		if node == nil || node.Key() == nil {
			return
		}
		fmt.Println(node.Key())
		node = node.Next()
	}
}

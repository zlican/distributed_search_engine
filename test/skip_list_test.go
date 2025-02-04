package test

import (
	"fmt"
	"testing"

	"github.com/huandu/skiplist"
)

func TestSkipList(t *testing.T) {
	list := skiplist.New(skiplist.Int32)
	list.Set(24, 31)
	list.Set(24, 20)
	list.Set(12, 40)
	list.Set(36, 50)
	list.Remove(36)

	if value, ok := list.GetValue(24); ok {
		fmt.Println(value)
	}
	fmt.Println("--------------------")
	node := list.Front()
	for node != nil {
		fmt.Println(node.Key(), node.Value)
		node = node.Next()
	}
}

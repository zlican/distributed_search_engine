package test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/zlican/engine/utils"
)

func TestMapIteractor(t *testing.T) {
	wg := sync.WaitGroup{}
	m := utils.NewMyCurrencyMap(1000, 100000)
	wg.Add(100)
	for i := 0; i < Currency; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				key := strconv.Itoa(j) + "zlican" + strconv.Itoa(i)
				m.Set(key, j)
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()

	// 迭代器
	iter := m.NewIterator()
	operator(iter)
}

func operator(iter utils.MapInterface) {
	i := 0
	for {
		obj := iter.Next()
		if obj == nil {
			fmt.Println(i)
			return
		}
		i++
		fmt.Println(obj)
	}
}

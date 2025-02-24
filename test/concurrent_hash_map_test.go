package test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/zlican/engine/utils"
)

const (
	Currency = 100
)

func TestMyCurrencyMap(t *testing.T) {
	wg := sync.WaitGroup{}
	m := utils.NewMyCurrencyMap(16, 1024)

	for i := 0; i < Currency; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 100000; j++ {
				key := strconv.Itoa(j) + "zlican"
				m.Set(key, j)
			}
			defer wg.Done()
		}()
	}
	for i := 0; i < Currency; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 100000; j++ {
				key := strconv.Itoa(j) + "zlican"
				value, _ := m.Get(key)
				if value == 99999 {
					fmt.Println("key:", key, "value:", value)
				}
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
}

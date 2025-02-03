package utils

import (
	"sync"

	farmhash "github.com/leemcloughlin/gofarmhash"
	"golang.org/x/exp/maps"
)

// 自制高并发锁设计思想
type MyCurrencyMap struct {
	Maps  []map[string]any
	Locks []sync.RWMutex
	Seg   int //要有几个分Map
	Seed  uint32
}

func NewMyCurrencyMap(Seg int, Cap int) *MyCurrencyMap {
	maps := make([]map[string]any, Seg)
	locks := make([]sync.RWMutex, Seg)
	for i := 0; i < Seg; i++ {
		maps[i] = make(map[string]any, Cap/Seg)
	}

	return &MyCurrencyMap{
		Maps:  maps,
		Locks: locks,
		Seg:   Seg,
		Seed:  0,
	}
}

func (m *MyCurrencyMap) GetKeyIndex(key string) int {
	hash := int(farmhash.Hash32WithSeed([]byte(key), m.Seed))
	return hash % m.Seg
}

// 写
func (m *MyCurrencyMap) Set(key string, value any) {
	index := m.GetKeyIndex(key)
	m.Locks[index].Lock()
	defer m.Locks[index].Unlock()
	m.Maps[index][key] = value
}

// 读
func (m *MyCurrencyMap) Get(key string) (value any, ok bool) {
	index := m.GetKeyIndex(key)
	m.Locks[index].RLock()
	defer m.Locks[index].RUnlock()
	value, ok = m.Maps[index][key]
	return
}

type NextMap struct {
	key   string
	value any
}

// 迭代器: 调用Next() 方法返回一个键值对
type MyCurrencyMapIterator struct {
	mcm      *MyCurrencyMap
	keys     [][]string //用于存无序的Map
	rowIndex int
	colIndex int
}

// 初始化迭代器。同时这个迭代器含有next()方法
func (m *MyCurrencyMap) NewIterator() *MyCurrencyMapIterator {
	keys := make([][]string, m.Seg)
	for _, mp := range m.Maps {
		row := maps.Keys(mp)     //每一个小map拿出key成为一个数组
		keys = append(keys, row) //二维数组
	}
	return &MyCurrencyMapIterator{
		mcm:      m,
		keys:     keys,
		rowIndex: 0,
		colIndex: 0,
	}
}

func (m *MyCurrencyMapIterator) Next() *NextMap {
	if m.rowIndex >= len(m.keys) {
		return nil
	}

	key := m.keys[m.rowIndex]

	if len(key) == 0 {
		m.rowIndex += 1
		return m.Next() //递归要使用return
	}

	if m.colIndex > len(key)-1 {
		m.rowIndex += 1
		m.colIndex = 0
		return m.Next()
	}
	value, _ := m.mcm.Get(key[m.colIndex])
	m.colIndex += 1

	return &NextMap{
		key:   key[m.colIndex-1],
		value: value,
	}
}

type MapInterface interface {
	Next() *NextMap
}

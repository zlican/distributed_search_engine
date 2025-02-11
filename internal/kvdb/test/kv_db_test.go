package test

import (
	"engine/internal/kvdb"
	"testing"
)

var (
	db       kvdb.IKeyValueDB
	setup    func()
	teardown func()
)

func init() {
	teardown = func() {
		if db != nil {
			db.Close()
		}
	}
}

// 测试数据库路径
func testGetDbPath(t *testing.T) {
	path := db.GetDbPath()
	if path == "" {
		t.Error("数据库路径不能为空")
	}
	t.Log("数据库路径:", path)
}

// 测试基本的 Get/Set/Delete/Has 操作
func testBasicOperations(t *testing.T) {
	k1 := []byte("key1")
	v1 := []byte("value1")

	// 测试 Set 和 Get
	if err := db.Set(k1, v1); err != nil {
		t.Fatalf("Set 失败: %v", err)
	}

	v, err := db.Get(k1)
	if err != nil {
		t.Fatalf("Get 失败: %v", err)
	}
	if string(v) != string(v1) {
		t.Errorf("期望值 %s, 实际值 %s", v1, v)
	}

	// 测试 Has
	if !db.Has(k1) {
		t.Error("Has 应该返回 true")
	}

	// 测试 Delete
	if err := db.Delete(k1); err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}

	if db.Has(k1) {
		t.Error("删除后 Has 应该返回 false")
	}

	// 测试获取不存在的键
	_, err = db.Get([]byte("不存在的键"))
	if err != kvdb.NoDataErr {
		t.Errorf("获取不存在的键应该返回 NoDataErr，实际返回: %v", err)
	}
}

// 测试批量操作
func testBatchOperations(t *testing.T) {
	keys := [][]byte{
		[]byte("batch1"),
		[]byte("batch2"),
		[]byte("batch3"),
	}
	values := [][]byte{
		[]byte("value1"),
		[]byte("value2"),
		[]byte("value3"),
	}

	// 测试 BatchSet
	err := db.BatchSet(keys, values)
	if err != nil {
		t.Fatalf("BatchSet 失败: %v", err)
	}

	// 测试 BatchGet
	results, err := db.BatchGet(keys)
	if err != nil {
		t.Fatalf("BatchGet 失败: %v", err)
	}
	for i, v := range results {
		if string(v) != string(values[i]) {
			t.Errorf("BatchGet: 键 %s 期望值 %s, 实际值 %s",
				string(keys[i]), string(values[i]), string(v))
		}
	}

	// 测试 BatchDelete
	err = db.BatchDelete(keys)
	if err != nil {
		t.Fatalf("BatchDelete 失败: %v", err)
	}

	// 验证删除结果
	for _, k := range keys {
		if db.Has(k) {
			t.Errorf("键 %s 应该已被删除", string(k))
		}
	}
}

// 测试迭代器
func testIterators(t *testing.T) {
	// 准备测试数据
	testData := map[string]string{
		"iter1": "value1",
		"iter2": "value2",
		"iter3": "value3",
	}

	for k, v := range testData {
		if err := db.Set([]byte(k), []byte(v)); err != nil {
			t.Fatalf("设置测试数据失败: %v", err)
		}
	}

	// 测试 IterKey
	keyCount := db.IterKey(func(k []byte) error {
		if _, exists := testData[string(k)]; !exists {
			t.Errorf("发现未知键: %s", string(k))
		}
		return nil
	})
	if keyCount != int64(len(testData)) {
		t.Errorf("IterKey: 期望计数 %d, 实际计数 %d", len(testData), keyCount)
	}

	// 测试 IterDB
	pairCount := db.IterDB(func(k, v []byte) error {
		expectedValue, exists := testData[string(k)]
		if !exists {
			t.Errorf("发现未知键值对: %s => %s", string(k), string(v))
		} else if expectedValue != string(v) {
			t.Errorf("键 %s 的值不匹配: 期望 %s, 实际 %s",
				string(k), expectedValue, string(v))
		}
		return nil
	})
	if pairCount != int64(len(testData)) {
		t.Errorf("IterDB: 期望计数 %d, 实际计数 %d", len(testData), pairCount)
	}

	// 清理测试数据
	for k := range testData {
		if err := db.Delete([]byte(k)); err != nil {
			t.Fatalf("清理测试数据失败: %v", err)
		}
	}
}

// 主测试流程
func testPipeline(t *testing.T) {
	defer teardown()
	setup()

	t.Run("测试数据库路径", func(t *testing.T) {
		testGetDbPath(t)
	})

	t.Run("测试基本操作", func(t *testing.T) {
		testBasicOperations(t)
	})

	t.Run("测试批量操作", func(t *testing.T) {
		testBatchOperations(t)
	})

	t.Run("测试迭代器", func(t *testing.T) {
		testIterators(t)
	})
}

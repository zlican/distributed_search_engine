package test

import (
	"engine/internal/kvdb"
	"os"
	"testing"
)

func TestBolt(t *testing.T) {
	dbPath := "./testdata/bolt_db"

	// 确保测试前清理旧数据
	os.RemoveAll("./testdata")

	setup = func() {
		var err error
		db, err = kvdb.GetKvDb(kvdb.BOLT, dbPath)
		if err != nil {
			t.Fatalf("初始化 Bolt 数据库失败: %v", err)
		}
	}

	t.Run("Bolt数据库测试", testPipeline)

	// 测试后清理
	os.RemoveAll("./testdata")
}

func TestBadger(t *testing.T) {
	dbPath := "./testdata/badger_db"

	// 确保测试前清理旧数据
	os.RemoveAll("./testdata")

	setup = func() {
		var err error
		db, err = kvdb.GetKvDb(kvdb.BADGER, dbPath)
		if err != nil {
			t.Fatalf("初始化 Bolt 数据库失败: %v", err)
		}
	}

	t.Run("Badger数据库测试", testPipeline)

	// 测试后清理
	os.RemoveAll("./testdata")
}

//go test -v ./internal/kvdb/test -run=^TestBolt$ -count=1
//go test -v ./internal/kvdb/test -run=^TestBadger$ -count=1

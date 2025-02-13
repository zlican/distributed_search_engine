package kvdb

import (
	"fmt"
	"os"
	"strings"
)

const (
	BOLT = iota
	BADGER
)

type IKeyValueDB interface {
	Open() error                              //初始化DB
	GetDbPath() string                        //获取储存数据的目录
	Set(k, v []byte) error                    //写入<key, value>
	BatchSet(keys, values [][]byte) error     //批量写入<key, value>
	Get(k []byte) ([]byte, error)             //读取key对应的value
	BatchGet(keys [][]byte) ([][]byte, error) //批量读取，注意不保证顺序
	Delete(k []byte) error                    //删除
	BatchDelete(keys [][]byte) error          //批量删除
	Has(k []byte) bool                        //判断key是否存在
	IterDB(fn func(k, v []byte) error) int64  //遍历DB
	IterKey(fn func(k []byte) error) int64    //遍历key
	Close() error                             //关闭DB
}

func GetKvDb(dbtype int, path string) (IKeyValueDB, error) {
	paths := strings.Split(path, "/")
	parentPath := strings.Join(paths[:len(paths)-1], "/")

	info, err := os.Stat(parentPath)
	if os.IsNotExist(err) {
		fmt.Printf("create dir %s", parentPath)
		os.MkdirAll(parentPath, 0o600)

	} else {
		if info.Mode().IsRegular() {
			fmt.Printf("file %s is not a directory", parentPath)
			os.Remove(parentPath)
		}
	}

	var db IKeyValueDB
	switch dbtype {
	case BADGER:
		db = new(Badger).WithDataPath(path)
	default:
		db = new(Bolt).WithDataPath(path).WithBucket("radic")
	}
	err = db.Open()
	return db, err

} //工厂模式初始化

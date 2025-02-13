package course

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var single *gorm.DB              //通过gorm.Open()创建连接池，只需要一个实例
var once sync.Once = sync.Once{} //once.Do() 只执行一次
var lock = &sync.Mutex{}

// 并行情况下实现单例
func GetDB1() *gorm.DB {
	if single == nil {
		lock.Lock()
		defer lock.Unlock()
		if single == nil {
			single, _ = gorm.Open(mysql.Open(""))
		}
	} else {
		fmt.Println("单例已创建")
	}
	return single
}

func init() {
	single, _ = gorm.Open(mysql.Open(""))
}
func GetDB2() *gorm.DB {
	return single
}

// 最常用的单例模式解决方案
func GetDB3() *gorm.DB {
	if single == nil {
		once.Do(func() {
			single, _ = gorm.Open(mysql.Open(""))
		})
	} else {
		fmt.Println("单例已创建")
	}
	return single
}

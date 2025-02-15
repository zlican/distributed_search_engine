package course

import (
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

var TotalQuery int32

func Handler() {
	atomic.AddInt32(&TotalQuery, 1)
	time.Sleep(50 * time.Microsecond)
}

func CallHandler() {
	limiter := rate.NewLimiter(rate.Every(100*time.Microsecond), 10) //每隔100ms产生一个令牌，令牌桶容量为10
	n := 3
	for {
		//调用方法1
		//ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*100)
		//defer cancel()
		//if err := limiter.WaitN(ctx, n); err != nil { //如果令牌不够，则等待
		//	Handler()
		//}
		//调用方法2
		//if limiter.AllowN(time.Now(), n) { //如果令牌够，则执行
		//	Handler()
		//}

		reserve := limiter.ReserveN(time.Now(), n)
		time.Sleep(reserve.Delay()) //休眠等令牌够马上调用
		Handler()
	}
}

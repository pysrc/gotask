## 定时任务

最小堆实现的定时任务


## 安装&使用

`go get github.com/pysrc/gotask`

```go
package main

import (
	"log"
	"time"

	"github.com/pysrc/gotask"
)

func main() {
	var cur = time.Now().Unix()
	log.Println("Start", cur)
	// 12s后执行
	gotask.Push(&gotask.Task{
		Start: cur + 12,
		Run: gotask.TaskRun(func() {
			log.Println(time.Now().Unix(), "12 step 0")
		}),
	})

	// 3s后执行, 且每2s执行一次，9s后不再执行
	var max = time.Now().Add(time.Second * 9).Unix()
	gotask.Push(&gotask.Task{
		Start: cur + 3,
		Run: gotask.TaskRun(func() {
			log.Println(time.Now().Unix(), "3 step 2")
		}),
		Next: gotask.TaskNext(func() int64 {
			var n = time.Now().Add(time.Second * 2).Unix()
			if n < max {
				return n
			}
			return 0
		}),
	})

	// 5s后执行，且中途panic测试
	gotask.Push(&gotask.Task{
		Start: cur + 5,
		Run: gotask.TaskRun(func() {
			log.Println(time.Now().Unix(), "5 step 0")
			panic("Task run error test")
		}),
	})
	time.Sleep(time.Second * 14)
}

```

package gotask

import (
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	var cur = time.Now().Unix()
	t.Log("Start", cur)
	Push(&Task{
		Start: cur + 12,
		Run: TaskRun(func() {
			t.Log(time.Now().Unix(), "12 step 0")
		}),
	})
	var max = time.Now().Add(time.Second * 9).Unix()
	Push(&Task{
		Start: cur + 3,
		Run: TaskRun(func() {
			t.Log(time.Now().Unix(), "3 step 2")
		}),
		Next: TaskNext(func() int64 {
			var n = time.Now().Add(time.Second * 2).Unix()
			if n < max {
				return n
			}
			panic("Task next error test")
		}),
	})
	Push(&Task{
		Start: cur + 2,
		Run: TaskRun(func() {
			t.Log(time.Now().Unix(), "2 step 0")
			panic("Task run error test")
		}),
	})
	time.Sleep(time.Second * 14)
	Stop()
	time.Sleep(time.Second)
}

package gotask

import (
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	go Run()
	var cur = time.Now().Unix()
	t.Log("Start", cur)
	Push(&Task{
		Start: cur + 12,
		Run: func() {
			t.Log(time.Now().Unix(), "12 step 0")
		},
	})
	var max = time.Now().Add(time.Second * 9).Unix()
	Push(&Task{
		Start: cur + 3,
		Run: func() {
			t.Log(time.Now().Unix(), "3 step 2")
		},
		Next: func() int64 {
			var n = time.Now().Add(time.Second * 2).Unix()
			if n < max {
				return n
			}
			return 0
		},
	})
	Push(&Task{
		Start: cur + 2,
		Run: func() {
			t.Log(time.Now().Unix(), "2 step 0")
		},
	})
	time.Sleep(time.Second * 14)
}

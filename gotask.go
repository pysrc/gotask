package gotask

import (
	"log"
	"sync"
	"time"
)

type Task struct {
	Start int64        // 开始时间
	Next  func() int64 // 下次执行时间
	Run   func()       // 任务
}

var (
	// 初始化3层堆
	tasks []*Task = make([]*Task, 0, 15)
	// 运行
	running = true
	// 堆操作锁
	mu sync.Mutex
)

func left(i int) int {
	return i*2 + 1
}

func right(i int) int {
	return i*2 + 2
}

func parent(i int) int {
	if i == 0 {
		return -1
	}
	return (i - 1) / 2
}

// 向下调整
func heapfiy(i int) {
	if i >= len(tasks) {
		return
	}
	var l, r = left(i), right(i)
	var min = i
	if l < len(tasks) && tasks[l].Start < tasks[min].Start {
		min = l
	}
	if r < len(tasks) && tasks[r].Start < tasks[min].Start {
		min = r
	}
	if min != i {
		tasks[i], tasks[min] = tasks[min], tasks[i]
		heapfiy(min)
	}
}

// 添加任务
func Push(task *Task) {
	mu.Lock()
	defer mu.Unlock()
	if task == nil {
		return
	}
	tasks = append(tasks, task)
	var p = parent(len(tasks) - 1)
	for p != -1 {
		heapfiy(p)
		p = parent(p)
	}
}

// 最早的任务
func Min() *Task {
	if len(tasks) == 0 {
		return nil
	}
	return tasks[0]
}

// 弹出最早的任务
func Pop() *Task {
	mu.Lock()
	defer mu.Unlock()
	if len(tasks) == 0 {
		return nil
	}
	var min = tasks[0]
	tasks[0] = tasks[len(tasks)-1]
	heapfiy(0)
	tasks = tasks[:len(tasks)-1]
	return min
}

// 任务执行宕机恢复
func taskRun(run func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		run()
	}
}

// 容错处理
func taskNext(run func() int64) func() int64 {
	return func() int64 {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		return run()
	}
}

func Stop() {
	running = false
}

func Run() {
	defer log.Println("Stop gotask")
	log.Println("Start gotask")
	for running {
		var task = Min()
		for task != nil && task.Start <= time.Now().Unix() {
			go taskRun(task.Run)()
			if task.Next != nil {
				task.Start = taskNext(task.Next)()
				if task.Start > time.Now().Unix() {
					mu.Lock()
					heapfiy(0)
					mu.Unlock()
				} else {
					Pop()
				}
			} else {
				Pop()
			}
			task = Min()
		}
		time.Sleep(time.Second)
	}
}

func init() {
	go Run()
}

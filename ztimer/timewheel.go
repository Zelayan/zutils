package main

import (
	"container/list"
	"sync"
	"time"
)

type JobFuc func()

type TaskId int

type TimeWheelInterface interface {
	Add(interval time.Duration, fuc JobFuc) string
	Remove(jobId string)
	Start()
	Stop()
}

// TimeWheel 时间轮结构体
type TimeWheel struct {
	ticker         *time.Ticker
	locker         sync.Mutex // 锁
	addTaskChan    chan *Task
	removeTaskChan chan TaskId
	stopChan       chan struct{}
	slots          []list.List   // 每个插槽的任务链表
	slotsNum       int64         // 插槽的数量
	interval       time.Duration // 每隔多少时间插槽后移动
	currentSlot    int64         // 当前执行的槽位
	store          map[TaskId]Task
	totalTime      time.Duration // 转一圈需要的总时间
}

func (t *TimeWheel) Start() {
	defer func() {

	}()
	for {
		select {
		case <-t.ticker.C:
			t.handerTasks()
		case taskId := <-t.removeTaskChan:
			t.handlerRemoveTask(taskId)
		case task := <-t.addTaskChan:
			t.handlerAddTask(task)
		case <-t.stopChan:
			t.ticker.Stop()
			return
		}
	}
}

func (t *TimeWheel) Stop() {
	//TODO implement me
	panic("implement me")
}

func (t *TimeWheel) Add(interval time.Duration, fuc JobFuc) string {
	// 填充槽，circle
	total := interval.Milliseconds() / t.interval.Milliseconds()
	circle := total / t.slotsNum
	slot := total % t.slotsNum
	// circle
	taskid := 1
	task := &Task{
		circle:     circle,
		job:        fuc,
		taskId:     TaskId(taskid),
		initCircle: circle,
		panicFunc:  nil,
		slot:       slot,
		slotNum:    total,
	}
	t.addTaskChan <- task
	return "11"
}

func (t *TimeWheel) Remove(jobId string) {
	//TODO implement me
	panic("implement me")
}

func (t *TimeWheel) handlerAddTask(task *Task) {
	t.locker.Lock()
	t.slots[task.slot].PushBack(task)
	t.locker.Unlock()
}

func (t *TimeWheel) handlerRemoveTask(id TaskId) {

}

func (t *TimeWheel) handerTasks() {
	l := t.slots[t.currentSlot]
	if e := l.Front(); e != nil {
		task := e.Value.(*Task)
		if task.circle == 0 {
			go task.job()
			l.Remove(e)
			// 对于已经运行过的 job 需要重新计算下次运行的槽的位置, 并插入
			t.reAdd(task)
		} else {
			// 圈数减一
			task.circle--
		}
	}
}

func (t *TimeWheel) reAdd(task *Task) {
	// task 的 interval 小于 t.num
	if t.totalTime > task.delay {
		// 圈数不够
		if t.slotsNum-t.currentSlot < task.slotNum {
			task.circle++
			task.
		} else {

		}
	} else {

	}
	// task 的 interval 大于 t.num

}

func NewTimeWheel() TimeWheelInterface {
	var slotsNum int64 = 10
	var interval = time.Duration(1000) * time.Millisecond
	return &TimeWheel{
		ticker:         time.NewTicker(time.Duration(1) * time.Second),
		locker:         sync.Mutex{},
		addTaskChan:    make(chan *Task, 1),
		removeTaskChan: make(chan TaskId, 1),
		stopChan:       make(chan struct{}),
		slots:          nil,
		slotsNum:       slotsNum,
		interval:       interval,
		currentSlot:    1,
		store:          nil,
		totalTime:      time.Duration(slotsNum * interval.Milliseconds()),
	}
}

// Task 每个具体的任务
type Task struct {
	circle     int64  // 具体需要转几圈
	job        JobFuc // 需要执行的任务
	taskId     TaskId
	initCircle int64         //初始circle
	delay      time.Duration // task的运行周期
	panicFunc  func()
	slot       int64 // 在哪个槽
	slotNum    int64 // 总共需要占用多少个槽
}

func (t *Task) run() {
	defer func() {
		if err := recover(); err != nil {
			if t.panicFunc != nil {
				t.panicFunc()
			}
		}
	}()
	t.job()
	t.circle = t.initCircle
}

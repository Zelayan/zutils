package main

import (
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
	slots          [][]*Task     // 每个插槽的任务链表
	slotsNum       int           // 插槽的数量
	interval       time.Duration // 每隔多少时间插槽后移动
	currentSlot    uint64        // 当前执行的槽位
	store          map[TaskId]Task
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
	//
	i := int(interval.Seconds()) % t.slotsNum
	// circle
	circle := int(interval.Seconds()) / t.slotsNum
	taskid := 1
	task := &Task{
		circle:     circle,
		job:        fuc,
		taskId:     TaskId(taskid),
		initCircle: circle,
		panicFunc:  nil,
		slot:       i,
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
	t.slots[task.slot] = append(t.slots[task.slot], task)
	t.locker.Unlock()
}

func (t *TimeWheel) handlerRemoveTask(id TaskId) {

}

func (t *TimeWheel) handerTasks() {
	l := t.slots[t.currentSlot]
	for i := range l {
		if l[i].circle == 0 {
			go l[i].run()
		} else {
			l[i].circle--
		}
	}
}

func NewTimeWheel() TimeWheelInterface {
	return &TimeWheel{
		ticker:         time.NewTicker(time.Duration(1) * time.Second),
		locker:         sync.Mutex{},
		addTaskChan:    make(chan *Task, 1),
		removeTaskChan: make(chan TaskId, 1),
		stopChan:       make(chan struct{}),
		slots:          make([][]*Task, 10),
		slotsNum:       10,
		interval:       time.Duration(1) * time.Second,
		currentSlot:    1,
		store:          nil,
	}
}

// Task 每个具体的任务
type Task struct {
	circle     int    // 具体需要转几圈
	job        JobFuc // 需要执行的任务
	taskId     TaskId
	initCircle int //初始circle
	panicFunc  func()
	slot       int // 在哪个槽
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

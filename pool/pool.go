package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var taskPool sync.Pool
var workPool sync.Pool
var defaultPool Pool

func init() {
	taskPool.New = newTask
	workPool.New = newWorker
	defaultPool = NewPool("default pool", 10000)
}

func newTask() interface{} {
	return &task{}
}
func newWorker() interface{} {
	return &worker{}
}

type Pool interface { // 池子的名称

	// SetCap 设置池子内Goroutine的容量
	SetCap(cap int32)

	// Go 执行 f 函数
	Go(f func())

	// CtxGo 带 ctx，执行 f 函数
	CtxGo(ctx context.Context, f func())

	// SetPanicHandler 设置发生panic时调用的函数
	SetPanicHandler(f func(context.Context, interface{}))
}

type pool struct {
	// name
	name string
	// pool池的容量
	cap int32
	// 任务锁
	taskLock sync.Mutex
	// 任务头
	taskHead *task
	// 新添加的任务
	taskTail *task
	// 当前的任务数量
	taskCount int32

	// 当前运行的worker数量
	workerCount int32
	// 发生异常时调用的方法
	panicHandler func(context.Context, interface{})
}

func (p *pool) SetCap(cap int32) {
	p.cap = cap
}

func (p *pool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f
	p.taskLock.Lock()
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()
	atomic.AddInt32(&p.taskCount, 1)
	if p.workerCount == 0 || atomic.LoadInt32(&p.taskCount) <= atomic.LoadInt32(&p.cap) {
		p.incWorkerCount()
		w := workPool.Get().(*worker)
		w.pool = p
		w.run()
	}

}

func (p *pool) SetPanicHandler(f func(context.Context, interface{})) {
	p.panicHandler = f
}

func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}

func Go(f func()) {
	CtxGo(context.Background(), f)
}

func CtxGo(ctx context.Context, f func()) {
	defaultPool.CtxGo(ctx, f)
}

func NewPool(name string, cap int32) Pool {
	p := &pool{
		name: name,
		cap:  cap,
	}
	return p
}

type task struct {
	ctx  context.Context
	f    func()
	next *task
}

type worker struct {
	pool *pool
}

// 用于消费pool中的数据
func (w *worker) run() {
	go func() {
		for {
			var t *task
			w.pool.taskLock.Lock()
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCount, -1)
			}

			if t == nil {
				w.pool.taskLock.Unlock()
				return
			}
			w.pool.taskLock.Unlock()

			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("pool: %s panic\n", w.pool.name)
						if w.pool.panicHandler != nil {
							w.pool.panicHandler(t.ctx, r)
						}
					}
				}()
				t.f()
			}()
		}
	}()
}

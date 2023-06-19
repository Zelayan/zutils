package main

import (
	"fmt"
	"time"
)

type Timer struct {
	interval time.Duration
	callback func()
}

type TimeWheel struct {
	tickInterval time.Duration
	maxTimeout   time.Duration
	slots        [][]*Timer
	currentPos   int
	stopChan     chan struct{}
}

func NewTimeWheel(tickInterval, maxTimeout time.Duration, numSlots int) *TimeWheel {
	tw := &TimeWheel{
		tickInterval: tickInterval,
		maxTimeout:   maxTimeout,
		slots:        make([][]*Timer, numSlots),
		currentPos:   0,
		stopChan:     make(chan struct{}),
	}

	for i := 0; i < numSlots; i++ {
		tw.slots[i] = make([]*Timer, 0)
	}

	return tw
}

func (tw *TimeWheel) Start() {
	ticker := time.NewTicker(tw.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tw.rotate()
		case <-tw.stopChan:
			return
		}
	}
}

func (tw *TimeWheel) Stop() {
	close(tw.stopChan)
}

func (tw *TimeWheel) rotate() {
	tw.currentPos = (tw.currentPos + 1) % len(tw.slots)
	timers := tw.slots[tw.currentPos]
	tw.slots[tw.currentPos] = make([]*Timer, 0)

	for _, timer := range timers {
		timer.callback()
	}
}

func (tw *TimeWheel) AddTimer(timer *Timer) {
	pos := (tw.currentPos + int(timer.interval/tw.tickInterval)) % len(tw.slots)
	tw.slots[pos] = append(tw.slots[pos], timer)
}

func main() {
	tw := NewTimeWheel(time.Second, 10*time.Second, 10)

	timer1 := &Timer{
		interval: 3 * time.Second,
		callback: func() {
			fmt.Println("Timer 1 fired")
		},
	}

	timer2 := &Timer{
		interval: 5 * time.Second,
		callback: func() {
			fmt.Println("Timer 2 fired")
		},
	}

	tw.AddTimer(timer1)
	tw.AddTimer(timer2)

	go tw.Start()

	time.Sleep(20 * time.Second)

	tw.Stop()

	fmt.Println("TimeWheel stopped")
}

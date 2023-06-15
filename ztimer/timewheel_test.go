package main

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {
	wheel := NewTimeWheel()
	go wheel.Start()

	wheel.Add(time.Duration(1)*time.Second, func() {
		fmt.Println("xxx")
	})

	time.Sleep(time.Second * 10)

}

package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestNewPool(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	pool := NewPool("xx", 10)
	pool.Go(func() {
		t.Log("working!")
		wg.Done()
	})
	wg.Wait()
}

func TestPool_SetPanicHandler(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	pool := NewPool("panic pool", 10)
	pool.SetPanicHandler(func(ctx context.Context, i interface{}) {
		fmt.Printf("panic %s\n", i)
		wg.Done()
	})
	pool.Go(func() {
		t.Log("worker")
		panic("wow panic !!")
	})
	wg.Wait()
}

func TestGo(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		t.Log("hah")
		wg.Done()
	})
	wg.Wait()
}

package zmap

import (
	"strconv"
	"sync"
)

func HelpBench(n int, zmap Zmap) {
	wg := sync.WaitGroup{}
	wg.Add(n * 2)
	for i := 0; i < n; i++ {
		go func(a int) {
			defer wg.Done()
			zmap.Store(strconv.Itoa(a), a)
		}(i)
	}

	for i := 0; i < 1000; i++ {
		go func(a int) {
			defer wg.Done()
			zmap.Load(strconv.Itoa(a))
		}(i)
	}
	wg.Wait()
}

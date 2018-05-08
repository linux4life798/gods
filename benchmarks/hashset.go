package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/linux4life798/safetyfast"
	"github.com/linux4life798/testutils"
)

func TestHashSetRW(numgoroutines int, reads, writes, updates bool, readvals, writevals []int32, c safetyfast.AtomicContext) time.Duration {
	var multiplier int = 0
	if reads {
		multiplier++
	}
	if writes {
		multiplier++
	}
	if updates {
		multiplier++
	}

	fmt.Printf("# Starting: GoRoutines=%d | reads=%v | writes=%v | updates=%v\n", numgoroutines*multiplier, reads, writes, updates)

	m := make(map[int32]int32)
	for i := range readvals {
		m[readvals[i]] = 0
	}
	var wg sync.WaitGroup

	writerroutine := func() {
		for _, v := range writevals {
			c.Atomic(func() {
				m[v] = 0
			})
		}
		wg.Done()
	}

	readerroutine := func() {
		for _, v := range readvals {
			c.Atomic(func() {
				if _, ok := m[v]; !ok {
					panic("Element doesn't exist")
				}
			})
		}
		wg.Done()
	}

	updaterrroutine := func() {
		for _, v := range readvals {
			c.Atomic(func() {
				val, ok := m[v]
				if !ok {
					panic("Element doesn't exist")
				}
				m[v] = val + 1
			})
		}
		wg.Done()
	}

	runtime.KeepAlive(writerroutine)
	runtime.KeepAlive(readerroutine)
	runtime.KeepAlive(updaterrroutine)

	wg.Add(numgoroutines * multiplier)
	fmt.Printf("# Starting %d goroutines\n", numgoroutines*multiplier)
	start := time.Now()
	for tid := 0; tid < numgoroutines; tid++ {
		if reads {
			go readerroutine()
		}
		if writes {
			go writerroutine()
		}
		if updates {
			go updaterrroutine()
		}
	}
	wg.Wait()
	dur := time.Since(start)

	fmt.Printf("# Finished %d goroutines\n", numgoroutines*multiplier)
	fmt.Printf("# Finished in %v\n", dur)
	fmt.Println("")

	return dur
}

func TestHashMap(fileprefix string) {

	maxgoroutines := runtime.GOMAXPROCS(-1)

	var pl *testutils.PerfPlot
	var filename string
	var series string
	var title string

	const nreads = 3000000
	const nwrites = 3000000
	readvals := testutils.NewRandValues().AddSparseInt32(nreads).GetAllInt32()
	writevals := testutils.NewRandValues().AddSparseInt32(nwrites).GetAllInt32()

	pl = testutils.NewPerfPlot()
	runtime.GC()
	series = "NoMutex"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines; count++ {
		dur := TestHashSetRW(count, true, false, false, readvals, writevals, safetyfast.NewLockedContext(new(safetyfast.NoMutex)))
		pl.AddMetric(series, int64(count), dur/time.Duration(nreads))
	}
	runtime.GC()
	series = "SystemLock"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines; count++ {
		dur := TestHashSetRW(count, true, false, false, readvals, writevals, safetyfast.NewLockedContext(new(sync.Mutex)))
		pl.AddMetric(series, int64(count), dur/time.Duration(nreads))
	}
	runtime.GC()
	series = "HLESpinLock"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines; count++ {
		dur := TestHashSetRW(count, true, false, false, readvals, writevals, safetyfast.NewLockedContext(new(safetyfast.SpinHLEMutex)))
		pl.AddMetric(series, int64(count), dur/time.Duration(nreads))
	}
	runtime.GC()
	series = "RTM"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines; count++ {
		dur := TestHashSetRW(count, true, false, false, readvals, writevals, safetyfast.NewRTMContex(new(sync.Mutex)))
		pl.AddMetric(series, int64(count), dur/time.Duration(nreads))
	}

	filename = fileprefix + "maps-reads.svg"
	title = fmt.Sprintf("Map %d-Read Performance", nreads)
	pl.Plot(filename, "Number of Goroutines", "Duration (ns)", title, false, true)
	testutils.OpenPlot(filename)

	//////////////

	pl = testutils.NewPerfPlot()
	runtime.GC()
	series = "SystemLock"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines/2; count++ {
		dur := TestHashSetRW(count, true, false, true, readvals, writevals, safetyfast.NewLockedContext(new(sync.Mutex)))
		pl.AddMetric(series, int64(count*2), dur/time.Duration(nreads))
	}
	runtime.GC()
	series = "HLESpinLock"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines/2; count++ {
		dur := TestHashSetRW(count, true, false, true, readvals, writevals, safetyfast.NewLockedContext(new(safetyfast.SpinHLEMutex)))
		pl.AddMetric(series, int64(count*2), dur/time.Duration(nreads))
	}
	runtime.GC()
	series = "RTM"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines/2; count++ {
		dur := TestHashSetRW(count, true, false, true, readvals, writevals, safetyfast.NewRTMContex(new(sync.Mutex)))
		pl.AddMetric(series, int64(count*2), dur/time.Duration(nreads))
	}

	filename = fileprefix + "maps-reads-updates.svg"
	title = fmt.Sprintf("Map %d-Read/%d-Update Performance", nreads, nreads)
	pl.Plot(filename, "Number of Goroutines", "Duration (ns)", title, false, true)
	testutils.OpenPlot(filename)

	//////////////

	pl = testutils.NewPerfPlot()
	runtime.GC()
	series = "SystemLock"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines/2; count++ {
		dur := TestHashSetRW(count, true, true, false, readvals, writevals, safetyfast.NewLockedContext(new(sync.Mutex)))
		pl.AddMetric(series, int64(count*2), dur/time.Duration(nreads))
	}
	runtime.GC()
	series = "HLESpinLock"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines/2; count++ {
		dur := TestHashSetRW(count, true, true, false, readvals, writevals, safetyfast.NewLockedContext(new(safetyfast.SpinHLEMutex)))
		pl.AddMetric(series, int64(count*2), dur/time.Duration(nreads))
	}
	runtime.GC()
	series = "RTM"
	fmt.Println("# Running experiment", series)
	for count := 1; count <= maxgoroutines/2; count++ {
		dur := TestHashSetRW(count, true, true, false, readvals, writevals, safetyfast.NewRTMContex(new(sync.Mutex)))
		pl.AddMetric(series, int64(count*2), dur/time.Duration(nreads))
	}

	filename = fileprefix + "maps-reads-puts.svg"
	title = fmt.Sprintf("Map %d-Read/%d-Put Performance", nreads, nwrites)
	pl.Plot(filename, "Number of Goroutines", "Duration (ns)", title, false, true)
	testutils.OpenPlot(filename)
}

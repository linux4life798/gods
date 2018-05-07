package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/linux4life798/safetyfast"
	"github.com/linux4life798/testutils"
)

func TestHashSet() {
	const numbasevalues = 2000
	const numops = 20000

	pl := testutils.NewPerfPlot()

	for ratio := 0; ratio <= 100; ratio += 25 {
		// const ratio = 50

		for numgoroutines := 1; numgoroutines <= runtime.GOMAXPROCS(-1); numgoroutines++ {
			h := hashset.New()

			fmt.Printf("Setting up for %d Goroutines\n", numgoroutines)
			basevalues := testutils.NewRandValues().AddSparseInt64(numbasevalues).GetAll()
			h.Add(basevalues...)

			actions := make([][]bool, numgoroutines)
			putvalues := make([][]interface{}, numgoroutines)

			for i := range actions {
				a := testutils.NewRandValues()
				a.AddIdenticalBool(true, numops*ratio)        // get value
				a.AddIdenticalBool(false, numops*(100-ratio)) // put value
				a.Shuffle()
				actions[i] = a.GetAllBool()

				putvalues[i] = testutils.NewRandValues().AddSparseInt64(numops * 100).GetAll()
			}

			var lock sync.Mutex
			var wg sync.WaitGroup

			routinelocks := func(tid int) {
				putvals := putvalues[tid]
				for i, isget := range actions[tid] {
					if isget {
						lock.Lock()
						_ = h.Contains(basevalues[i%len(basevalues)])
						lock.Unlock()
					} else {
						lock.Lock()
						h.Add(putvals[i])
						lock.Unlock()
					}
				}
				wg.Done()
			}

			fmt.Printf("Starting %d Goroutines\n", numgoroutines)
			wg.Add(numgoroutines)
			start := time.Now()
			for tid := 0; tid < numgoroutines; tid++ {
				go routinelocks(tid)
			}
			wg.Wait()
			dur := time.Since(start)

			fmt.Printf("%d Goroutines finished in %v\n", numgoroutines, dur)
			pl.AddMetric(fmt.Sprintf("HS-%d%%Gets", ratio), int64(numgoroutines), dur/time.Duration(numops))

			runtime.GC()
			runtime.Gosched()

		}

		for numgoroutines := 1; numgoroutines <= runtime.GOMAXPROCS(-1); numgoroutines++ {
			h := hashset.New()

			fmt.Printf("Setting up for %d Goroutines\n", numgoroutines)
			basevalues := testutils.NewRandValues().AddSparseInt64(numbasevalues).GetAll()
			h.Add(basevalues...)

			actions := make([][]bool, numgoroutines)
			putvalues := make([][]interface{}, numgoroutines)

			for i := range actions {
				a := testutils.NewRandValues()
				a.AddIdenticalBool(true, numops*ratio)        // get value
				a.AddIdenticalBool(false, numops*(100-ratio)) // put value
				a.Shuffle()
				actions[i] = a.GetAllBool()

				putvalues[i] = testutils.NewRandValues().AddSparseInt64(numops * 100).GetAll()
			}

			// var lock sync.Mutex
			rtm := safetyfast.NewRTMContexDefault()
			var wg sync.WaitGroup

			routine := func(tid int) {
				putvals := putvalues[tid]
				for i, isget := range actions[tid] {
					if isget {
						rtm.Commit(func() {
							_ = h.Contains(basevalues[i%len(basevalues)])
						})
					} else {
						rtm.Commit(func() {
							h.Add(putvals[i])
						})
					}
				}
				wg.Done()
			}

			fmt.Printf("Starting %d Goroutines\n", numgoroutines)
			wg.Add(numgoroutines)
			start := time.Now()
			for tid := 0; tid < numgoroutines; tid++ {
				go routine(tid)
			}
			wg.Wait()
			dur := time.Since(start)

			fmt.Printf("%d Goroutines finished in %v\n", numgoroutines, dur)
			pl.AddMetric(fmt.Sprintf("HSRTM-%d%%Gets", ratio), int64(numgoroutines), dur/time.Duration(numops))

			runtime.GC()
			runtime.Gosched()

		}

	}

	fmt.Println("# Plotting")
	pl.Plot("hashset.svg", "NumGoroutines", "Duration Per Operation (ns)", "Hash Set", false)
	fmt.Println("# Done Plotting")
}

func main() {
	// TestAVL("avl")
	TestHashSet()
}

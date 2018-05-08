package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/linux4life798/safetyfast"
	"github.com/linux4life798/testutils"
)

func TestHashSet() {
	const numbasevalues = 100
	const numops = 7000

	pl := testutils.NewPerfPlot()

	// for ratio := 0; ratio <= 100; ratio += 25 {
	const ratio = 100

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

		var lock safetyfast.SpinHLEMutex
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
		pl.AddMetric(fmt.Sprintf("HSHLE-%d%%Gets", ratio), int64(numgoroutines), dur/time.Duration(numops))

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

		rtm := safetyfast.NewRTMContex(new(safetyfast.SpinHLEMutex))
		var wg sync.WaitGroup

		routine := func(tid int) {
			putvals := putvalues[tid]
			for i, isget := range actions[tid] {
				if isget {
					rtm.Atomic(func() {
						_ = h.Contains(basevalues[i%len(basevalues)])
					})
				} else {
					rtm.Atomic(func() {
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

	// }

	fmt.Println("# Plotting")
	pl.Plot("hashset.svg", "NumGoroutines", "Duration Per Operation (ns)", "Hash Set", false, true)
	fmt.Println("# Done Plotting")
}

/////////////////////////

// func TestAVLTree() {
// 	const numbasevalues = 20
// 	const numops = 20000

// 	pl := testutils.NewPerfPlot()

// 	// for ratio := 0; ratio <= 100; ratio += 25 {
// 	for ratio := 50; ratio <= 50; ratio += 25 {
// 		// const ratio = 50

// 		for numgoroutines := 4; numgoroutines <= runtime.GOMAXPROCS(-1); numgoroutines++ {
// 			h := avltree.NewWith(utils.Int64Comparator)

// 			fmt.Printf("Setting up for %d Goroutines\n", numgoroutines)
// 			basevalues := testutils.NewRandValues().AddSparseInt64(numbasevalues).GetAll()

// 			for _, v := range basevalues {
// 				h.Put(v, v)
// 			}

// 			actions := make([][]bool, numgoroutines)
// 			putvalues := make([][]interface{}, numgoroutines)

// 			for i := range actions {
// 				a := testutils.NewRandValues()
// 				a.AddIdenticalBool(true, numops*ratio)        // get value
// 				a.AddIdenticalBool(false, numops*(100-ratio)) // put value
// 				a.Shuffle()
// 				actions[i] = a.GetAllBool()

// 				putvalues[i] = testutils.NewRandValues().AddSparseInt64(numops * 100).GetAll()
// 			}

// 			var lock sync.Mutex
// 			var wg sync.WaitGroup

// 			routinelocks := func(tid int) {
// 				putvals := putvalues[tid]
// 				for i, isget := range actions[tid] {
// 					if isget {
// 						lock.Lock()
// 						_, _ = h.Get(basevalues[i%len(basevalues)])
// 						lock.Unlock()
// 					} else {
// 						lock.Lock()
// 						h.Put(putvals[i], putvals[i])
// 						lock.Unlock()
// 					}
// 				}
// 				wg.Done()
// 			}

// 			fmt.Printf("Starting %d Goroutines\n", numgoroutines)
// 			wg.Add(numgoroutines)
// 			start := time.Now()
// 			for tid := 0; tid < numgoroutines; tid++ {
// 				go routinelocks(tid)
// 			}
// 			wg.Wait()
// 			dur := time.Since(start)

// 			fmt.Printf("%d Goroutines finished in %v\n", numgoroutines, dur)
// 			pl.AddMetric(fmt.Sprintf("AVL-%d%%Gets", ratio), int64(numgoroutines), dur/time.Duration(numops))

// 			runtime.GC()
// 			runtime.Gosched()
// 		}

// 		for numgoroutines := 4; numgoroutines <= runtime.GOMAXPROCS(-1); numgoroutines++ {
// 			h := avltree.NewWith(utils.Int64Comparator)

// 			fmt.Printf("Setting up for %d Goroutines\n", numgoroutines)
// 			basevalues := testutils.NewRandValues().AddSparseInt64(numbasevalues).GetAll()

// 			for _, v := range basevalues {
// 				h.Put(v, v)
// 			}
// 			actions := make([][]bool, numgoroutines)
// 			putvalues := make([][]interface{}, numgoroutines)

// 			for i := range actions {
// 				a := testutils.NewRandValues()
// 				a.AddIdenticalBool(true, numops*ratio)        // get value
// 				a.AddIdenticalBool(false, numops*(100-ratio)) // put value
// 				a.Shuffle()
// 				actions[i] = a.GetAllBool()

// 				putvalues[i] = testutils.NewRandValues().AddSparseInt64(numops * 100).GetAll()
// 			}

// 			// var lock sync.Mutex
// 			rtm := safetyfast.NewRTMContexDefault()
// 			var wg sync.WaitGroup

// 			routine := func(tid int) {
// 				putvals := putvalues[tid]
// 				for i, isget := range actions[tid] {
// 					if isget {
// 						rtm.Atomic(func() {
// 							_, _ = h.Get(basevalues[i%len(basevalues)])
// 						})
// 					} else {
// 						rtm.Atomic(func() {
// 							h.Put(putvals[i], putvals[i])
// 						})
// 					}
// 				}
// 				wg.Done()
// 			}

// 			fmt.Printf("Starting %d Goroutines\n", numgoroutines)
// 			wg.Add(numgoroutines)
// 			start := time.Now()
// 			for tid := 0; tid < numgoroutines; tid++ {
// 				go routine(tid)
// 			}
// 			wg.Wait()
// 			dur := time.Since(start)

// 			fmt.Println("RTM Capacity =", rtm.CapacityAborts())

// 			fmt.Printf("%d Goroutines finished in %v\n", numgoroutines, dur)
// 			pl.AddMetric(fmt.Sprintf("AVLRTM-%d%%Gets", ratio), int64(numgoroutines), dur/time.Duration(numops))

// 			runtime.GC()
// 			runtime.Gosched()

// 		}

// 	}

// 	fmt.Println("# Plotting")
// 	pl.Plot("avltree.svg", "NumGoroutines", "Duration Per Operation (ns)", "AVL Tree", false, true)
// 	fmt.Println("# Done Plotting")
// }

var FlagFilePrefix string

func init() {
	flag.StringVar(&FlagFilePrefix, "prefix", "", "This is the prefix to all files generated")
}

func main() {
	flag.Parse()
	// TestAVL("avl")
	// TestHashSet()

	TestHashMap(FlagFilePrefix)

	TestAVLTree(FlagFilePrefix)

	// TestAVLTree()
}

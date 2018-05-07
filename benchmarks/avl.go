package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/emirpasic/gods/trees/avltree"
	"github.com/emirpasic/gods/utils"
	"github.com/linux4life798/testutils"
)

func TestAVLGet() *testutils.PerfPlot {
	const ngets = int64(100000)

	p := testutils.NewPerfPlot()
	t := avltree.NewWith(utils.Int64Comparator)

	for n := int64(1); n < 513; n++ {
		t.Clear()

		vals := testutils.NewRandValues().AddConsecutiveInt64(0, int(n))
		vals.Shuffle()

		for _, v := range vals.GetAll() {
			t.Put(v.(int64), v.(int64))
		}

		getvals := vals.Shuffle().GetAllInt64()

		var ret interface{}
		start := time.Now()
		for i := int64(0); i < ngets; i++ {
			var ok bool
			ret, ok = t.Get(getvals[i%n])
			if !ok {
				panic("failed to find my value")
			}
		}
		dur := time.Since(start)
		runtime.KeepAlive(ret)

		perop := dur / time.Duration(ngets)
		p.AddMetric("Get", n, perop)
		fmt.Println("Get: n =", n)
		runtime.GC()
		runtime.Gosched()
	}

	return p
}

func TestAVLPut() *testutils.PerfPlot {
	// const nputs = int64(100000)
	const nputs = int64(100)
	const repeatcount = 2000

	p := testutils.NewPerfPlot()
	var trees [repeatcount]*avltree.Tree

	// t := avltree.NewWith(utils.Int64Comparator)
	for r := 0; r < repeatcount; r++ {
		trees[r] = avltree.NewWith(utils.Int64Comparator)
	}

	for n := int64(1); n < 513; n++ {
		// t.Clear()
		for r := 0; r < repeatcount; r++ {
			trees[r].Clear()
		}

		vals := testutils.NewRandValues().AddConsecutiveInt64(0, int(n+nputs))
		valsarr := vals.Shuffle().GetAllInt64()

		// place n values in the tree
		for i := int64(0); i < nputs; i++ {
			for r := int64(0); r < int64(repeatcount); r++ {
				trees[r].Put(valsarr[i], valsarr[i])
			}
		}

		// start := time.Now()
		// // time the placement of nput items
		// for i := int64(0); i < nputs; i++ {
		// 	j := i % n
		// 	t.Put(valsarr[j], valsarr[j])
		// }
		// dur := time.Since(start)

		start := time.Now()
		// var start time.Time
		// var dur time.Duration
		imax := n + nputs
		// time the placement of nput items
		for i := int64(n); i < imax; i++ {
			// start = time.Now()
			for r := int64(0); r < int64(repeatcount); r++ {
				// j := (i + r) % imax
				// trees[r].Put(valsarr[j], valsarr[j])
				// j := i % n
				trees[r].Put(valsarr[i], valsarr[i])
			}
			// dur += time.Since(start)
			// t.Remove(valsarr[i])
		}
		dur := time.Since(start)

		perop := dur / time.Duration(nputs*int64(repeatcount))
		p.AddMetric("Put", n, perop)
		fmt.Println("Put: n =", n)
		runtime.GC()
		runtime.Gosched()
	}

	return p
}

func TestAVL(fileprefix string) {
	// getfile := fileprefix + "-get" + ".svg"
	// fmt.Println("Starting Get Test")
	// getpl := TestAVLGet()
	// fmt.Println("Finished Get Test")
	// getpl.Plot(getfile, "N", "Operation Time (ns)", "AVL Get Performance", false)
	// testutils.OpenPlot(getfile)

	putfile := fileprefix + "-put" + ".svg"
	fmt.Println("Starting Put Test")
	putpl := TestAVLPut()
	fmt.Println("Finished Put Test")
	// putpl.LimitMax(getpl.GetMax())
	// putpl.LimitMax(time.Nanosecond * 300)
	putpl.Plot(putfile, "N", "Operation Time (ns)", "AVL Put Performance", false)
	testutils.OpenPlot(putfile)
}

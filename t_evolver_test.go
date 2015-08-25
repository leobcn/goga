// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_evo01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo01. organise sequence of ints")
	io.Pf("\n")

	// initialise random numbers generator
	rnd.Init(0) // 0 => use current time as seed

	// parameters
	C := NewConfParams()
	C.Nisl = 1
	C.Ninds = 20
	C.RegTol = 0
	//C.GAtype = "crowd"
	C.CrowdSize = 2
	C.Tf = 50
	C.Verbose = chk.Verbose
	C.CalcDerived()

	// mutation function
	C.MtIntFunc = func(A []int, time, nchanges int, pm float64, extra interface{}) {
		size := len(A)
		if !rnd.FlipCoin(pm) || size < 1 {
			return
		}
		pos := rnd.IntGetUniqueN(0, size, nchanges)
		for _, i := range pos {
			if A[i] == 1 {
				A[i] = 0
			}
			if A[i] == 0 {
				A[i] = 1
			}
		}
	}

	// generation function
	nvals := 20
	C.PopIntGen = func(ninds, nova, noor, nbases int, noise float64, args interface{}, irange [][]int) Population {
		o := make([]*Individual, ninds)
		genes := make([]int, nvals)
		for i := 0; i < ninds; i++ {
			for j := 0; j < nvals; j++ {
				genes[j] = rand.Intn(2)
			}
			o[i] = NewIndividual(nova, noor, nbases, genes)
		}
		return o
	}

	// objective function
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		score := 0.0
		count := 0
		for _, val := range ind.Ints {
			if val == 0 && count%2 == 0 {
				score += 1.0
			}
			if val == 1 && count%2 != 0 {
				score += 1.0
			}
			count++
		}
		ind.Ovas[0] = 1.0 / (1.0 + score)
		return
	}

	// run optimisation
	nova := 1
	noor := 0
	evo := NewEvolver(nova, noor, C)
	evo.Run()

	// results
	ideal := 1.0 / (1.0 + float64(nvals))
	io.PfGreen("\nBest = %v\nBestOV = %v  (ideal=%v)\n", evo.Best.Ints, evo.Best.Ovas[0], ideal)
}

func Test_evo02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo02")

	// initialise random numbers generator
	//rnd.Init(0) // 0 => use current time as seed
	rnd.Init(1111) // 0 => use current time as seed

	// parameters
	C := NewConfParams()
	C.Pll = false
	C.Nisl = 1
	C.Ninds = 20
	C.RegTol = 0
	C.GAtype = "crowd"
	C.ParetoPhi = 0.01
	C.Elite = false
	C.Verbose = false
	C.RangeFlt = [][]float64{
		{-2, 2}, // gene # 0: min and max
		{-2, 2}, // gene # 1: min and max
	}
	C.PopFltGen = PopFltGen
	if chk.Verbose {
		C.FnKey = "test_evo02"
		C.DoPlot = true
	}
	C.CalcDerived()

	f := func(x []float64) float64 { return x[0]*x[0]/2.0 + x[1]*x[1] - x[0]*x[1] - 2.0*x[0] - 6.0*x[1] }
	c1 := func(x []float64) float64 { return x[0] + x[1] - 2.0 }      // ≤ 0
	c2 := func(x []float64) float64 { return -x[0] + 2.0*x[1] - 2.0 } // ≤ 0
	c3 := func(x []float64) float64 { return 2.0*x[0] + x[1] - 3.0 }  // ≤ 0
	c4 := func(x []float64) float64 { return -x[0] }                  // ≤ 0
	c5 := func(x []float64) float64 { return -x[1] }                  // ≤ 0

	// objective function
	p := 1.0
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		x := ind.GetFloats()
		ind.Ovas[0] = f(x)
		ind.Oors[0] = utl.GtePenalty(0, c1(x), p)
		ind.Oors[1] = utl.GtePenalty(0, c2(x), p)
		ind.Oors[2] = utl.GtePenalty(0, c3(x), p)
		ind.Oors[3] = utl.GtePenalty(0, c4(x), p)
		ind.Oors[4] = utl.GtePenalty(0, c5(x), p)
		return
	}

	// evolver
	nova := 1
	noor := 5
	evo := NewEvolver(nova, noor, C)
	pop0 := evo.Islands[0].Pop.GetCopy()
	evo.Run()

	// results
	io.PfGreen("\nx=%g (%g)\n", evo.Best.GetFloat(0), 2.0/3.0)
	io.PfGreen("y=%g (%g)\n", evo.Best.GetFloat(1), 4.0/3.0)
	io.PfGreen("BestOV=%g (%g)\n", evo.Best.Ovas[0], f([]float64{2.0 / 3.0, 4.0 / 3.0}))

	// plot contour
	if C.DoPlot {
		xmin := []float64{-2, -2}
		xmax := []float64{2, 2}
		PlotTwoVarsContour("/tmp/goga", "contour_evo02", pop0, evo.Islands[0].Pop, evo.Best, 41, nil, true,
			xmin, xmax, false, false, nil, nil, f, c1, c2, c3, c4, c5)
	}
}

func Test_evo03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo03")

	rnd.Init(0)

	// parameters
	C := NewConfParams()
	C.Pll = false
	C.Nisl = 1
	C.Ninds = 20
	C.GAtype = "crowd"
	C.Elite = true
	C.RangeFlt = [][]float64{
		{-1, 3}, // gene # 0: min and max
		{-1, 3}, // gene # 1: min and max
	}
	C.PopFltGen = PopFltGen
	if chk.Verbose {
		C.FnKey = "test_evo03"
		C.DoPlot = true
	}
	C.CalcDerived()

	// geometry
	xe := 1.0                      // centre of circle
	le := -0.4                     // selected level of f(x)
	ys := xe - (1.0+le)/math.Sqrt2 // coordinates of minimum point with level=le
	y0 := 2.0*ys + xe              // vertical axis intersect of straight line defined by c(x)
	xc := []float64{xe, xe}        // centre
	nx := len(xc)
	f := func(x []float64) (res float64) {
		for i := 0; i < nx; i++ {
			res += (x[i] - xc[i]) * (x[i] - xc[i])
		}
		return math.Sqrt(res) - 1
	}
	c := func(x []float64) (res float64) {
		return x[0] + x[1] + xe - y0
	}

	// objective function
	p := 1.0
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		x := ind.GetFloats()
		fp := utl.GtePenalty(1e-2, math.Abs(c(x)), p)
		ind.Ovas[0] = f(x) + fp
		ind.Oors[0] = fp
		return
	}

	// evolver
	nova := 1
	noor := 1
	evo := NewEvolver(nova, noor, C)
	pop0 := evo.Islands[0].Pop.GetCopy()
	evo.Run()

	// results
	xbest := []float64{evo.Best.GetFloat(0), evo.Best.GetFloat(1)}
	io.PfGreen("\nx=%g (%g)\n", xbest[0], ys)
	io.PfGreen("y=%g (%g)\n", xbest[1], ys)
	io.PfGreen("BestOV=%g (%g)\n\n", evo.Best.Ovas[0], le)

	// plot contour
	if C.DoPlot {
		extra := func() {
			plt.PlotOne(ys, ys, "'o', markeredgecolor='yellow', markerfacecolor='none', markersize=10")
		}
		xmin := []float64{-1, -1}
		xmax := []float64{3, 3}
		PlotTwoVarsContour("/tmp/goga", "contour_evo03", pop0, evo.Islands[0].Pop, evo.Best, 41, extra, true,
			xmin, xmax, false, false, nil, nil, f, c)
	}
}

func Test_evo04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo04. TSP")

	// location / coordinates of stations
	locations := [][]float64{
		{60, 200}, {180, 200}, {80, 180}, {140, 180}, {20, 160}, {100, 160}, {200, 160},
		{140, 140}, {40, 120}, {100, 120}, {180, 100}, {60, 80}, {120, 80}, {180, 60},
		{20, 40}, {100, 40}, {200, 40}, {20, 20}, {60, 20}, {160, 20},
	}
	nstations := len(locations)

	// parameters
	C := NewConfParams()
	C.Nisl = 1
	C.Ninds = 20
	C.RegTol = 0.3
	C.RegPct = 0.2
	//C.Dtmig = 30
	C.IntOrd = true
	C.GAtype = "crowd"
	C.ParetoPhi = 0.1
	C.Elite = false
	C.DoPlot = false //chk.Verbose
	C.PopOrdGen = PopOrdGen
	C.OrdNints = nstations
	//C.Rws = true
	C.CalcDerived()

	// initialise random numbers generator
	rnd.Init(0)

	// objective value function
	C.OvaOor = func(ind *Individual, idIsland, t int, report *bytes.Buffer) {
		L := locations
		ids := ind.Ints
		dist := 0.0
		for i := 1; i < nstations; i++ {
			a, b := ids[i-1], ids[i]
			dist += math.Sqrt(math.Pow(L[b][0]-L[a][0], 2.0) + math.Pow(L[b][1]-L[a][1], 2.0))
		}
		a, b := ids[nstations-1], ids[0]
		dist += math.Sqrt(math.Pow(L[b][0]-L[a][0], 2.0) + math.Pow(L[b][1]-L[a][1], 2.0))
		ind.Ovas[0] = dist
		return
	}

	// evolver
	nova := 1
	noor := 0
	evo := NewEvolver(nova, noor, C)

	// print initial population
	pop := evo.Islands[0].Pop
	//io.Pf("\n%v\n", pop.Output(nil, false))

	// 0,4,8,11,14,17,18,15,12,19,13,16,10,6,1,3,7,9,5,2 894.363
	if false {
		for i, x := range []int{0, 4, 8, 11, 14, 17, 18, 15, 12, 19, 13, 16, 10, 6, 1, 3, 7, 9, 5, 2} {
			pop[0].Ints[i] = x
		}
		evo.Islands[0].CalcOvs(pop, 0)
		evo.Islands[0].CalcDemeritsAndSort(pop)
	}

	// check initial population
	ints := make([]int, nstations)
	if false {
		for i := 0; i < C.Ninds; i++ {
			for j := 0; j < nstations; j++ {
				ints[j] = pop[i].Ints[j]
			}
			sort.Ints(ints)
			chk.Ints(tst, "ints", ints, utl.IntRange(nstations))
		}
	}

	// run
	evo.Run()
	//io.Pf("%v\n", pop.Output(nil, false))
	io.Pfgreen("best = %v\n", evo.Best.Ints)
	io.Pfgreen("best OVA = %v  (871.117353844847)\n\n", evo.Best.Ovas[0])

	// best = [18 17 14 11 8 4 0 2 5 9 12 7 6 1 3 10 16 13 19 15]
	// best OVA = 953.4643474956656

	// best = [8 11 14 17 18 15 12 19 16 13 10 6 1 3 7 9 5 2 0 4]
	// best OVA = 871.117353844847

	// best = [5 2 0 4 8 11 14 17 18 15 12 19 16 13 10 6 1 3 7 9]
	// best OVA = 871.1173538448469

	// best = [6 10 13 16 19 15 18 17 14 11 8 4 0 2 5 9 12 7 3 1]
	// best OVA = 880.7760751923065

	// check final population
	if false {
		for i := 0; i < C.Ninds; i++ {
			for j := 0; j < nstations; j++ {
				ints[j] = pop[i].Ints[j]
			}
			sort.Ints(ints)
			chk.Ints(tst, "ints", ints, utl.IntRange(nstations))
		}
	}

	// plot travelling salesman path
	if C.DoPlot {
		plt.SetForEps(1, 300)
		X, Y := make([]float64, nstations), make([]float64, nstations)
		for k, id := range evo.Best.Ints {
			X[k], Y[k] = locations[id][0], locations[id][1]
			plt.PlotOne(X[k], Y[k], "'r.', ms=5, clip_on=0, zorder=20")
			plt.Text(X[k], Y[k], io.Sf("%d", id), "fontsize=7, clip_on=0, zorder=30")
		}
		plt.Plot(X, Y, "'b-', clip_on=0, zorder=10")
		plt.Plot([]float64{X[0], X[nstations-1]}, []float64{Y[0], Y[nstations-1]}, "'b-', clip_on=0, zorder=10")
		plt.Equal()
		plt.AxisRange(10, 210, 10, 210)
		plt.Gll("$x$", "$y$", "")
		plt.SaveD("/tmp/goga", "test_evo04.eps")
	}
}

func Test_evo05(tst *testing.T) {

	verbose()
	chk.PrintTitle("evo04. sin⁶(5 π x)")

	// configuration
	C := NewConfParams()
	C.Nisl = 1
	C.Ninds = 12
	C.GAtype = "crowd"
	C.CrowdSize = 3
	C.ParetoPhi = 0
	C.Noise = 0.05
	C.DoPlot = false
	C.RegTol = 0
	C.Pc = 0.8
	C.Pm = 0.01
	C.MtExtra = map[string]interface{}{"flt": 1.1}
	C.Tf = 100
	C.Dtmig = 101
	C.RangeFlt = [][]float64{{0, 1}}
	C.PopFltGen = PopFltGen
	C.CalcDerived()

	// initialise random numbers generator
	rnd.Init(0)

	// function
	yfcn := func(x float64) float64 {
		return math.Pow(math.Sin(5.0*math.Pi*x), 6.0)
	}

	// objective value function
	C.OvaOor = func(ind *Individual, idIsland, t int, report *bytes.Buffer) {
		x := ind.GetFloat(0)
		ind.Ovas[0] = -yfcn(x)
		ind.Oors[0] = utl.GtePenalty(x, 0, 1)
		ind.Oors[1] = utl.GtePenalty(1, x, 1)
	}

	// post-processing function
	values := utl.Deep3alloc(C.Tf/10, C.Nisl, C.Ninds)
	C.PostProc = func(idIsland, time int, pop Population) {
		if time%10 == 0 && false {
			k := time / 10
			for i, ind := range pop {
				values[k][idIsland][i] = ind.GetFloat(0)
			}
		}
	}

	// run
	nova := 1
	noor := 2
	evo := NewEvolver(nova, noor, C)
	evo.Run()

	// print population
	for _, isl := range evo.Islands {
		io.Pf("%v", isl.Pop.Output(nil, true, false, -1))
	}

	// write histograms and plot
	if chk.Verbose {

		// write histograms
		var buf bytes.Buffer
		hist := rnd.Histogram{Stations: utl.LinSpace(0, 1, 26)}
		for k := 0; k < C.Tf/10; k++ {
			for i := 0; i < C.Nisl; i++ {
				clear := false
				if i == 0 {
					clear = true
				}
				hist.Count(values[k][i], clear)
			}
			io.Ff(&buf, "\ntime=%d\n%v", k*10, rnd.TextHist(hist.GenLabels("%4.2f"), hist.Counts, 60))
		}
		io.WriteFileVD("/tmp/goga", "test_evo05_hist.txt", &buf)

		// plot
		plt.SetForEps(0.8, 300)
		xmin := evo.Islands[0].Pop[0].GetFloat(0)
		xmax := xmin
		for k := 0; k < C.Nisl; k++ {
			for _, ind := range evo.Islands[k].Pop {
				x := ind.GetFloat(0)
				y := yfcn(x)
				xmin = utl.Min(xmin, x)
				xmax = utl.Max(xmax, x)
				plt.PlotOne(x, y, "'r.',clip_on=0,zorder=20")
			}
		}
		np := 401
		//X := utl.LinSpace(xmin, xmax, np)
		X := utl.LinSpace(0, 1, np)
		Y := make([]float64, np)
		for i := 0; i < np; i++ {
			Y[i] = yfcn(X[i])
		}
		plt.Plot(X, Y, "'b-',clip_on=0,zorder=10")
		plt.Gll("$x$", "$y$", "")
		//plt.AxisXrange(0, 1)
		plt.SaveD("/tmp/goga", "test_evo05_func.eps")
	}
}

// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	PI = math.Pi
)

func DTLZ2mGenerator(opt *goga.Optimiser, nf int) (ng int, fcn goga.MinProb_t) {
	nx := nf + 10
	ng = nx * 2
	opt.FltMin = make([]float64, nx)
	opt.FltMax = make([]float64, nx)
	for i := 0; i < nx; i++ {
		opt.FltMin[i], opt.FltMax[i] = -0.01, 1.01
	}
	fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
		var failed bool
		for i := 0; i < nx; i++ {
			g[0+i*2] = x[i]
			g[1+i*2] = 1.0 - x[i]
			if g[0+i*2] < 0 {
				failed = true
			}
			if g[1+i*2] < 0 {
				failed = true
			}
		}
		if failed {
			return
		}
		var c float64
		for i := nf - 1; i < nx; i++ {
			c += math.Pow((x[i] - 0.5), 2.0)
		}
		for i := 0; i < nf; i++ {
			f[i] = (1.0 + c)
			for j := 0; j < nf-1-i; j++ {
				f[i] *= math.Cos(x[j] * PI / 2.0)
			}
			if i > 0 {
				j := nf - 1 - i
				f[i] *= math.Sin(x[j] * PI / 2.0)
			}
		}
	}
	opt.Multi_fcnErr = func(f []float64) float64 {
		var sum float64
		for i := 0; i < nf; i++ {
			sum += f[i] * f[i]
		}
		return sum - 1.0
	}
	opt.RptFmin = make([]float64, nf)
	opt.RptFmax = make([]float64, nf)
	for i := 0; i < nf; i++ {
		opt.RptFmax[i] = 1
	}
	return
}

func solve_problem(problem int) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 300
	opt.Ncpu = 6
	opt.Tf = 500
	opt.Nsamples = 2
	opt.DEC = 0.01

	// options for report
	opt.HistNsta = 6
	opt.HistLen = 13
	opt.RptFmtE = "%.4e"
	opt.RptFmtL = "%.4e"
	opt.RptFmtEdev = "%.3e"
	opt.RptFmtLdev = "%.3e"

	// problem variables
	var nf, ng, nh int     // number of functions
	var fcn goga.MinProb_t // functions

	// problems
	switch problem {
	case 1:
		nf = 5
		opt.RptName = io.Sf("DTLZ2m%d", nf)
		ng, fcn = DTLZ2mGenerator(opt, nf)

	case 2:
		nf = 7
		opt.RptName = io.Sf("DTLZ2m%d", nf)
		ng, fcn = DTLZ2mGenerator(opt, nf)

	case 3:
		nf = 10
		opt.RptName = io.Sf("DTLZ2m%d", nf)
		ng, fcn = DTLZ2mGenerator(opt, nf)

	case 4:
		nf = 13
		opt.RptName = io.Sf("DTLZ2m%d", nf)
		ng, fcn = DTLZ2mGenerator(opt, nf)

	case 5:
		nf = 15
		opt.RptName = io.Sf("DTLZ2m%d", nf)
		ng, fcn = DTLZ2mGenerator(opt, nf)

	case 6:
		nf = 20
		opt.RptName = io.Sf("DTLZ2m%d", nf)
		ng, fcn = DTLZ2mGenerator(opt, nf)

	default:
		chk.Panic("problem %d is not available", problem)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.RunMany("", "")
	goga.StatMulti(opt, true)

	// check
	goga.CheckFront0(opt, true)

	// star plot
	if true {
		plt.SetForEps(1, 300)
		goga.PlotStar(opt)
		plt.SaveD("/tmp/goga", io.Sf("starplot_%s.eps", opt.RptName))
	}
	return
}

func main() {
	P := utl.IntRange2(1, 7)
	//P := []int{4}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem)
	}
	io.Pf("\n-------------------------- generating report --------------------------\nn")
	rpt := goga.NewTexReport(opts)
	rpt.NRowPerTab = 9
	rpt.Type = 1
	rpt.Title = "Unconstrained many objective problems"
	rpt.Fnkey = "many-obj"
	rpt.Generate()
}

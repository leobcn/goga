// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/utl"

type Group struct {
	Ncur    int         // number of current solutions == len(All) / 2
	All     []*Solution // current and future solutions => view to Solutions and FutureSols
	Indices []int       // indices of current solutions
	Pairs   [][]int     // randomly selected pairs from Indices
	Metrics *Metrics    // metrics
}

func (o *Group) Init(cpu, ncpu int, solutions, futuresols []*Solution) {
	nsol := len(solutions)
	start, endp1 := (cpu*nsol)/ncpu, ((cpu+1)*nsol)/ncpu
	o.Ncur = endp1 - start
	o.All = make([]*Solution, o.Ncur*2)
	o.Indices = make([]int, o.Ncur)
	o.Pairs = utl.IntsAlloc(o.Ncur/2, 2)
	for i := 0; i < o.Ncur; i++ {
		o.All[i] = solutions[start+i]
		o.All[o.Ncur+i] = futuresols[start+i]
		o.Indices[i] = i
	}
	nova := len(solutions[0].Ova)
	nflt := len(solutions[0].Flt)
	nint := len(solutions[0].Int)
	o.Metrics = new(Metrics)
	o.Metrics.Init(nova, nflt, nint, len(o.All))
}

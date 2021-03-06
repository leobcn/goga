// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/utl"

// Group holds a group of solutions
type Group struct {
	Ncur    int         // number of current solutions == len(All) / 2
	All     []*Solution // current and future solutions. half part is a view to Solutions
	Indices []int       // indices of current solutions
	Pairs   [][]int     // randomly selected pairs from Indices
	Metrics *Metrics    // metrics
}

// Init initialises group
func (o *Group) Init(cpu, ncpu int, solutions []*Solution, prms *Parameters) {
	nsol := len(solutions)
	start, endp1 := (cpu*nsol)/ncpu, ((cpu+1)*nsol)/ncpu
	o.Ncur = endp1 - start
	o.All = make([]*Solution, o.Ncur*2)
	o.Indices = make([]int, o.Ncur)
	o.Pairs = utl.IntsAlloc(o.Ncur/2, 2)
	for i := 0; i < o.Ncur; i++ {
		o.All[i] = solutions[start+i]
		o.All[o.Ncur+i] = NewSolution(-i, nsol, prms)
		o.Indices[i] = i
	}
	o.Metrics = new(Metrics)
	o.Metrics.Init(len(o.All), prms)
}

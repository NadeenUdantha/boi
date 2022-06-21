// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boigp

type BoiInterface interface {
	Init(n int)
	Update(i int)
	Score(i int) float64
	Swap(i, j int)
	Mutate(i int)
	Crossover(i, a, b int)
}

// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boigp

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type SimpleInstances []*SimpleInstance

type SimpleInstance struct {
	v     uint64
	Score float64
}

func (bois SimpleInstances) Init(n int) {
	for i := 0; i < n; i++ {
		boi := new(SimpleInstance)
		boi.v = rand.Uint64() & 0xff
		bois[i] = boi
	}
}

func (bois SimpleInstances) Update(_boi int) {
	boi := bois[_boi]
	v := boi.v
	c := 0
	for v > 0 {
		v &= (v - 1)
		c++
	}
	boi.Score = float64(c)
}

func (bois SimpleInstances) Mutate(_boi int) {
	boi := bois[_boi]
	p := rand.Int() % 64
	boi.v ^= (1 << p)
}

func (bois SimpleInstances) Crossover(_boi, _a, _b int) {
	boi := bois[_boi]
	a := bois[_a]
	b := bois[_b]
	p := rand.Int() % 64
	for i := 0; i < 64; i++ {
		var v uint64
		if i < p {
			v = a.v
		} else {
			v = b.v
		}
		boi.v = (boi.v & ^(1 << i)) | (v & (1 << i))
	}
}

func (bois SimpleInstances) Score(_boi int) float64 {
	return bois[_boi].Score
}

func (bois SimpleInstances) Swap(i, j int) {
	bois[i], bois[j] = bois[j], bois[i]
}

func (boi SimpleInstance) String() string {
	return fmt.Sprintf("simple[%d]", int(boi.Score))
}

func TestSimple(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	bgp := BoiGP{
		nbois:     6,
		p_nbois:   3,
		c_nbois:   1,
		max_gen:   1e9,
		log_gen:   1e0,
		max_score: 64,
	}
	bgp.Run(make(SimpleInstances, bgp.nbois))
}

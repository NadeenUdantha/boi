package boigp

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

type BoiGP struct {
	nbois     int
	p_nbois   int
	c_nbois   int
	max_gen   int
	log_gen   int
	mp        bool
	max_score float64
	bois      BoiInterface
}

func (bgp *BoiGP) Len() int {
	return bgp.nbois
}

func (bgp *BoiGP) Swap(i, j int) {
	bgp.bois.Swap(i, j)
}

func (bgp *BoiGP) Less(i, j int) bool {
	return bgp.bois.Score(i) < bgp.bois.Score(j)
}

func (bgp *BoiGP) UpdateAll() {
	if bgp.mp {
		var wg sync.WaitGroup
		wg.Add(bgp.nbois)
		for i := 0; i < bgp.nbois; i++ {
			go func(i int) {
				defer wg.Done()
				bgp.bois.Update(i)
			}(i)
		}
		wg.Wait()
	} else {
		for i := 0; i < bgp.nbois; i++ {
			bgp.bois.Update(i)
		}
	}
}

func (bgp *BoiGP) Run(bois BoiInterface) {
	bgp.bois = bois
	nb := bgp.nbois
	pnb := bgp.p_nbois
	cnb := bgp.p_nbois
	gen := 1
	bois.Init(nb)
	bgp.UpdateAll()
	for {
		sort.Sort(bgp)
		if gen%bgp.log_gen == 0 {
			fmt.Println(gen, bois)
		}
		if gen >= bgp.max_gen {
			break
		}
		if bois.Score(nb-1) >= bgp.max_score {
			break
		}
		for i := 0; i < cnb; i++ {
			ai := nb - 1 - (rand.Int() % pnb)
			bi := nb - 1 - (rand.Int() % pnb)
			bois.Crossover(i, ai, bi)
			if rand.Float32() < .2 {
				bois.Mutate(i)
			}
			bois.Update(i)
		}
		gen += 1
	}
	if gen%bgp.log_gen != 0 {
		fmt.Println(gen, bois)
	}
}

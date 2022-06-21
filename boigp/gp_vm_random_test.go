// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boigp

import (
	"boi/boilang"
	"boi/boivm"
	"bytes"
	srand "crypto/rand"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

type VMRandomInstances []*VMRandomInstance

type VMRandomInstance struct {
	vm         *boivm.BoiVM
	Score      float64
	ni, ne, no int
}

func (bois VMRandomInstances) Init(n int) {
	for i := 0; i < n; i++ {
		boi := new(VMRandomInstance)
		vm := new(boivm.BoiVM)
		vm.Init()
		srand.Read(vm.Mem.X)
		boi.vm = vm
		bois[i] = boi
	}
}

func (bois VMRandomInstances) Update(_boi int) {
	boi := bois[_boi]
	vm := boi.vm
	vm.Reset()
	rand.Seed(time.Now().Unix())
	i := 128
	if boi.ni >= i {
		i = boi.ni * 2
	}
	z := 0
	vm.Sys = func(u uint16) {
		//vm.Pop()
		z += 10
	}
	ni, no, errs := vm.Run(i, false, false)
	boi.ni = ni
	boi.no = no
	boi.ne = len(errs)
	boi.Score = float64(ni + (len(errs) * -1) + (no * 5) + z)
}

func (bois VMRandomInstances) Mutate(_boi int) {
	vm := bois[_boi].vm
	z := rand.Int() % (1024 * 8)
	for i := 0; i < z; i++ {
		vm.Mem.X[rand.Int()%0x8000] ^= (1 << (rand.Int() % 8))
	}
}

func (bois VMRandomInstances) Crossover(_boi, _a, _b int) {
	boi := bois[_boi].vm.Mem.X
	/*for i := 0; i < 0x8000; i++ {
		boi[i] = 0
	}*/
	a := bois[_a].vm.Mem.X
	b := bois[_b].vm.Mem.X
	z := 0x0010
	for i := 0; i < 0x8000; i += z {
		p := rand.Int() % z
		var v []uint8
		if i < p {
			v = a
		} else {
			v = b
		}
		zs := i
		ze := i + z
		copy(boi[zs:ze], v[zs:ze])
	}
}

func (bois VMRandomInstances) Score(_boi int) float64 {
	return bois[_boi].Score
}

func (bois VMRandomInstances) Swap(i, j int) {
	bois[i], bois[j] = bois[j], bois[i]
}

func (boi VMRandomInstance) String() string {
	return fmt.Sprintf("vm[%d,%d,%d,%d]", boi.ni, boi.ne, boi.no, int(boi.Score))
}

func TestVMRandom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	bgp := BoiGP{
		nbois:     10,
		p_nbois:   6,
		c_nbois:   3,
		max_gen:   1e9,
		log_gen:   1e0,
		max_score: 1e6,
		mp:        true,
	}
	vms := make(VMRandomInstances, bgp.nbois)
	bgp.Run(vms)
	mem := vms[len(vms)-1].vm.Mem.X
	//ioutil.WriteFile("R:/boi.mem", mem, 0)
	s, nerrs := boilang.Disassemble(bytes.NewBuffer(mem))
	os.Remove("R:/boi_random.boi")
	ioutil.WriteFile("R:/boi_random.boi", []byte(s), 0)
	fmt.Println(nerrs)
}

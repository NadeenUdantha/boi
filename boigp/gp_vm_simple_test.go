// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boigp

import (
	"boi/boilang"
	"boi/boivm"
	"bytes"
	srand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

const VMSimpleMaxSteps = 16
const VMSimpleAddr = 0x0c00

type VMSimpleInstances []*VMSimpleInstance

type VMSimpleInstance struct {
	vm    *boivm.BoiVM
	v     int
	Score float64
}

func (bois VMSimpleInstances) Init(n int) {
	for i := 0; i < n; i++ {
		boi := new(VMSimpleInstance)
		vm := new(boivm.BoiVM)
		vm.Init()
		srand.Read(vm.Mem.X[:VMSimpleMaxSteps*8])
		binary.BigEndian.PutUint32(vm.Mem.X[VMSimpleAddr:], 0)
		boi.vm = vm
		bois[i] = boi
	}
}

func (bois VMSimpleInstances) Update(_boi int) {
	boi := bois[_boi]
	vm := boi.vm
	vm.Reset()
	rand.Seed(time.Now().Unix())
	vm.Sys = func(u uint16) {}
	x := make([]uint8, len(vm.Mem.X))
	copy(x, vm.Mem.X)
	_, _, errs := vm.Run(VMSimpleMaxSteps, true, false)
	boi.v = int(binary.BigEndian.Uint32(vm.Mem.X[VMSimpleAddr:]))
	if len(errs) != 0 {
		boi.v--
	}
	boi.Score = float64(boi.v)
	copy(vm.Mem.X, x)
}

func (bois VMSimpleInstances) Mutate(_boi int) {
	vm := bois[_boi].vm
	z := rand.Int() % (1024 * 8)
	for i := 0; i < z; i++ {
		vm.Mem.X[rand.Int()%0x8000] ^= (1 << (rand.Int() % 8))
	}
}

func (bois VMSimpleInstances) Crossover(_boi, _a, _b int) {
	boi := bois[_boi].vm.Mem.X
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

func (bois VMSimpleInstances) Score(_boi int) float64 {
	return bois[_boi].Score
}

func (bois VMSimpleInstances) Swap(i, j int) {
	bois[i], bois[j] = bois[j], bois[i]
}

func (boi VMSimpleInstance) String() string {
	return fmt.Sprintf("vm[v=%d]", boi.v)
}

func TestVMSimple(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	bgp := BoiGP{
		nbois:     10,
		p_nbois:   6,
		c_nbois:   3,
		max_gen:   1e9,
		log_gen:   1e2,
		max_score: 0xffffff,
		mp:        true,
	}
	vms := make(VMSimpleInstances, bgp.nbois)
	bgp.Run(vms)
	boi := vms[len(vms)-1]
	mem := boi.vm.Mem.X
	os.Remove("R:/boi.mem")
	ioutil.WriteFile("R:/boi.mem", mem, 0)
	s, errs := boilang.Disassemble(bytes.NewBuffer(mem))
	fmt.Println(len(errs))
	os.Remove("R:/boi.boi")
	ioutil.WriteFile("R:/boi.boi", []byte(s), 0)
	boi.vm.Reset()
	boi.vm.Run(VMSimpleMaxSteps, true, true)
}

// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boivm

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	vm := new(BoiVM)
	vm.Init()
	rand.Read(vm.Mem.X)
	/*for i := 0; i < len(vm.Mem.X); i++ {
		vm.Mem.X[i] = uint8(i & 0xff)
	}*/
	for i := 0; i < 1e2; i++ {
		nsteps, nerrors, nops := vm.Run(10, false, false)
		fmt.Printf("nsteps=%d nerrors=%d nops=%d\n", nsteps, nerrors, nops)
	}
}

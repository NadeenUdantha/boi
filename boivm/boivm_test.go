// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boivm

import (
	"boi/boilang"
	"testing"
)

func Test(t *testing.T) {
	vm := new(BoiVM)
	vm.Init()
	/*f, err := elf.Open(`D:\nadeen\boi\tcc\boi\test.o`)
	assert.NoError(t, err)
	d, err := f.Section(".text").Data()
	assert.NoError(t, err)
	copy(vm.mem.x, d)*/
	copy(vm.Mem.X, boilang.Assemble("nop()").Bytes())
	vm.Run(1, false, false)
}

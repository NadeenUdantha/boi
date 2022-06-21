// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boivm

import (
	"boi/boilang"
	. "boi/boilang"
	"bytes"
	"fmt"
	"strings"
)

type BoiVM struct {
	Mem *BoiMem
	ip  uint16
	sp  uint16
	Sys func(uint16)
}

func (vm *BoiVM) Init() {
	vm.Mem = &BoiMem{X: make([]byte, 0x8000)}
	vm.Reset()
}

func (vm *BoiVM) Reset() {
	vm.ip = 0
	vm.sp = 0x8000
}

func (vm *BoiVM) Push(v uint16) {
	if vm.sp == 0 {
		panic("stackoverflow")
	}
	vm.Mem.w16(vm.sp, v)
	vm.sp -= 2
}

func (vm *BoiVM) Pop() uint16 {
	if vm.sp == 0x8000 {
		panic("stackunderflow")
	}
	v := vm.Mem.r16(vm.sp)
	vm.sp += 2
	return v
}

func (vm *BoiVM) Step(debug bool) (uint8, error) {
	defer func() { recover() }()
	m := vm.Mem
	av := m.av
	i := m.r8(vm.ip)
	//fmt.Printf("ip=%#04x i=%#02x\n", vm.ip, i)
	if debug {
		code, _ := boilang.Disassemble(bytes.NewBuffer(vm.Mem.X[vm.ip:][:8]))
		z := strings.Split(code, "\n")
		fmt.Printf("ip=%#04x i=%#02x %s\n", vm.ip, i, z[0])
	}
	vm.ip += 1
	if i == I_nop {
	} else if i == I_mov {
		m.w16(m.r16(vm.ip), av(m.r16(vm.ip+1)))
		vm.ip += 4
	} else if i == I_push {
		vm.Push(vm.Mem.av(vm.Mem.r16(vm.ip)))
		vm.ip += 2
	} else if i == I_pop {
		m.w16(m.r16(vm.ip), vm.Pop())
		vm.ip += 2
	} else if i == I_not {
		m.w16(m.r16(vm.ip), ^av(m.r16(vm.ip+1)))
		vm.ip += 4
	} else if i == I_call {
		f := av(m.r16(vm.ip))
		vm.ip += 2
		m.w16(vm.sp, vm.ip)
		vm.sp -= 2
		vm.ip = f
	} else if i == I_ret {
		vm.ip = m.r16(vm.sp)
		vm.sp += 2
	} else if i == I_sys {
		id := av(m.r16(vm.ip))
		vm.ip += 2
		vm.Sys(id)
		//Throw("syscall%d", id)
	} else if i == I_jmp {
		vm.ip = av(m.r16(vm.ip))
	} else if i == I_jz {
		if av(m.r16(vm.ip+1)) == 0 {
			vm.ip = av(m.r16(vm.ip))
		} else {
			vm.ip += 4
		}
	} else {
		if i == I_add || i == I_sub || i == I_mul || i == I_div || i == I_and || i == I_or || i == I_xor || i == I_shl || i == I_shr {
			dst := m.r16(vm.ip + 1)
			vm.ip += 2
			src1 := av(m.r16(vm.ip + 1))
			vm.ip += 2
			src2 := av(m.r16(vm.ip + 1))
			vm.ip += 2
			var dstv uint16
			if i == I_add {
				dstv = src1 + src2
			} else if i == I_sub {
				dstv = src1 - src2
			} else if i == I_mul {
				dstv = src1 * src2
			} else if i == I_div {
				if src2 == 0 {
					dstv = 0
				} else {
					dstv = src1 / src2
				}
			} else if i == I_and {
				dstv = src1 & src2
			} else if i == I_or {
				dstv = src1 | src2
			} else if i == I_xor {
				dstv = src1 ^ src2
			} else if i == I_shl {
				dstv = src1 << src2
			} else if i == I_shr {
				dstv = src1 >> src2
			}
			m.w16(dst, dstv)
		} else if i == I_je || i == I_jg || i == I_jl {
			dsta := av(m.r16(vm.ip))
			vm.ip += 2
			a := av(m.r16(vm.ip + 1))
			vm.ip += 2
			b := av(m.r16(vm.ip + 1))
			vm.ip += 2
			var cmp bool
			if i == I_je {
				cmp = a == b
			} else if i == I_jg {
				cmp = a > b
			} else if i == I_jl {
				cmp = a < b
			}
			if cmp {
				vm.ip = dsta
			}
		} else {
			return 0, fmt.Errorf("illegal instruction: %d ip=%x", i, vm.ip)
		}
	}
	return i, nil
}

func (vm *BoiVM) Run(max_steps int, return_on_err bool, debug bool) (nsteps, nops int, errs []error) {
	for nsteps < max_steps {
		i, err := vm.Step(debug)
		nsteps += 1
		if i > I_nop && i <= I_jz {
			nops++
		}
		if err != nil {
			errs = append(errs, err)
			if return_on_err {
				return
			}
		}
	}
	return
}

// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boilang

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var r_instr = regexp.MustCompile(`(\w+)\((.*)\)`)
var r_args = regexp.MustCompile(`([\[\w\]]+),?`)

func assemble_instr(bb *BoiBuf, line string) {
	mx := r_instr.FindStringSubmatch(line)
	//log.Println(len(mx), mx)
	if len(mx) == 0 || len(mx) != 3 || mx[0] != line {
		panic("wtf: " + line)
	}
	n := mx[1]
	z := r_args.FindAllString(mx[2], -1)
	for x, v := range z {
		if v[len(v)-1] == ',' {
			z[x] = v[:len(v)-1]
		}
	}
	fmt.Printf("[asm] %s.%d(%s)\n", n, len(z), strings.Join(z, ","))
	if len(z) == 1 && len(z[0]) == 0 {
		z = z[1:]
	}
	zl := func(i int) {
		if len(z) != i {
			panic(fmt.Sprintf("%s %d %s%d", n, i, z, len(z)))
		}
	}
	if n == "nop" {
		zl(0)
		bb.u8(I_nop)
	} else if n == "mov" {
		zl(2)
		bb.u8(I_mov)
		bb.u16(asm_a(z[0]))
		bb.u16(asm_av(z[1]))
	} else if n == "push" {
		zl(1)
		bb.u8(I_push)
		bb.u16(asm_av(z[0]))
	} else if n == "pop" {
		zl(1)
		bb.u8(I_pop)
		bb.u16(asm_a(z[0]))
	} else if n == "not" {
		zl(2)
		bb.u8(I_not)
		bb.u16(asm_a(z[0]))
		bb.u16(asm_av(z[1]))
	} else if n == "call" {
		zl(1)
		bb.u8(I_call)
		bb.u16(asm_av(z[0]))
	} else if n == "ret" {
		zl(0)
		bb.u8(I_ret)
	} else if n == "sys" {
		zl(1)
		bb.u8(I_sys)
		bb.u16(asm_av(z[0]))
	} else if n == "jmp" {
		zl(1)
		bb.u8(I_jmp)
		bb.u16(asm_av(z[0]))
	} else if n == "jz" {
		zl(2)
		bb.u8(I_jz)
		bb.u16(asm_av(z[0]))
		bb.u16(asm_av(z[1]))
	} else {
		if n == "je" || n == "jg" || n == "jl" {
			zl(3)
			if n == "je" {
				bb.u8(I_je)
			} else if n == "jg" {
				bb.u8(I_jg)
			} else if n == "jl" {
				bb.u8(I_jl)
			}
			bb.u16(asm_av(z[0]))
			bb.u16(asm_av(z[1]))
			bb.u16(asm_av(z[2]))
		} else if n == "add" || n == "sub" || n == "mul" || n == "div" || n == "and" || n == "or" || n == "xor" || n == "shl" || n == "shr" {
			zl(3)
			if n == "add" {
				bb.u8(I_add)
			} else if n == "sub" {
				bb.u8(I_sub)
			} else if n == "mul" {
				bb.u8(I_mul)
			} else if n == "div" {
				bb.u8(I_div)
			} else if n == "and" {
				bb.u8(I_and)
			} else if n == "or" {
				bb.u8(I_or)
			} else if n == "xor" {
				bb.u8(I_xor)
			} else if n == "shl" {
				bb.u8(I_shl)
			} else if n == "shr" {
				bb.u8(I_shr)
			}
			bb.u16(asm_a(z[0]))
			bb.u16(asm_av(z[1]))
			bb.u16(asm_av(z[2]))

		} else {
			Throw(`wtf is "%s"?`, n)
		}
	}
}

func asm_u15(x string) uint16 {
	//log.Println("u15", x)
	z, err := strconv.ParseUint(x, 10, 15)
	if err != nil {
		panic(err)
	}
	return uint16(z) & 0x7fff
}

const (
	asm_v_mask      = 0x7fff
	asm_isaddr_mask = 0x8000
)

func asm_a(x string) uint16 {
	//log.Println("z_a", x)
	if x[0] == '[' && x[len(x)-1] == ']' {
		return asm_u15(x[1:len(x)-1]) | asm_isaddr_mask
	}
	panic(x)
}

func asm_av(x string) uint16 {
	//log.Println("z_av", x)
	if x[0] == '[' && x[len(x)-1] == ']' {
		return asm_u15(x[1:len(x)-1]) | asm_isaddr_mask
	}
	return asm_u15(x) & asm_v_mask
}

var r_line = regexp.MustCompile(`[\r\n]+`)

func Assemble(source string) *bytes.Buffer {
	buf := &BoiBuf{buf: new(bytes.Buffer)}
	for _, line := range r_line.Split(source, -1) {
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] != '#' {
			assemble_instr(buf, line)
		}
	}
	return buf.buf
}

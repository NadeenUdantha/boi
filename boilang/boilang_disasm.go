package boilang

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func disassemble_instr(bb *BoiBuf) (string, bool) {
	defer func() { recover() }()
	i := bb.ru8()
	sb := strings.Builder{}
	zi := func(n string) {
		sb.WriteString(n)
		sb.WriteByte('(')
	}
	dasm_avz := func(a bool, c bool) {
		x := bb.ru16()
		isaddr := (x & asm_isaddr_mask) == asm_isaddr_mask
		if a && !isaddr {
			Throw("!addr: %d %d", x, x&asm_v_mask)
		}
		x = x & 0x7fff
		if isaddr {
			sb.WriteByte('[')
		}
		sb.WriteString(strconv.FormatUint(uint64(x), 10))
		if isaddr {
			sb.WriteByte(']')
		}
		if c {
			sb.WriteByte(',')
		}
	}
	dasm_av := func(c bool) {
		dasm_avz(false, c)
	}
	dasm_a := func(c bool) {
		dasm_avz(true, c)
	}
	if i == I_nop {
		zi("nop")
	} else if i == I_mov {
		zi("mov")
		dasm_a(true)
		dasm_av(false)
	} else if i == I_push {
		zi("push")
		dasm_av(false)
	} else if i == I_pop {
		zi("pop")
		dasm_a(false)
	} else if i == I_not {
		zi("not")
		dasm_a(true)
		dasm_av(false)
	} else if i == I_call {
		zi("call")
		dasm_av(false)
	} else if i == I_ret {
		zi("ret")
	} else if i == I_sys {
		zi("sys")
		dasm_av(false)
	} else if i == I_jmp {
		zi("jmp")
		dasm_av(false)
	} else if i == I_jz {
		zi("jz")
		dasm_av(true)
		dasm_av(false)
	} else {
		if i == I_add || i == I_sub || i == I_mul || i == I_div || i == I_and || i == I_or || i == I_xor || i == I_shl || i == I_shr {
			if i == I_add {
				zi("add")
			} else if i == I_sub {
				zi("sub")
			} else if i == I_mul {
				zi("mul")
			} else if i == I_div {
				zi("div")
			} else if i == I_and {
				zi("and")
			} else if i == I_or {
				zi("or")
			} else if i == I_xor {
				zi("xor")
			} else if i == I_shl {
				zi("shl")
			} else if i == I_shr {
				zi("shr")
			}
			dasm_a(true)
			dasm_av(true)
			dasm_av(false)
		} else if i == I_je || i == I_jg || i == I_jl {
			if i == I_je {
				zi("je")

			} else if i == I_jg {
				zi("jg")

			} else if i == I_jl {
				zi("jl")
			}
			dasm_av(true)
			dasm_av(true)
			dasm_av(false)
		} else {
			e := fmt.Sprintf("# i=%d?", i)
			//fmt.Println(e)
			return e, true
		}
	}
	sb.WriteByte(')')
	fmt.Printf("[dasm] %s\n", sb.String())
	return sb.String(), false
}

func Disassemble(x *bytes.Buffer) (string, int) {
	buf := &BoiBuf{buf: x}
	sb := &strings.Builder{}
	nerrs := 0
	for {
		s, err := disassemble_instr(buf)
		s = strings.TrimSpace(s)
		if err {
			nerrs += 1
		}
		if len(s) > 0 {
			if err {
				//sb.WriteString(s)
				//sb.WriteByte('\n')
			} else {
				sb.WriteString(s)
				sb.WriteByte('\n')
			}
		}
		if x.Len() == 0 {
			break
		}
	}
	return sb.String(), nerrs
}

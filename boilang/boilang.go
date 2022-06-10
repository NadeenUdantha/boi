package boilang

import "fmt"

const (
	I_nop = uint8(iota)
	I_mov
	I_push
	I_pop
	I_add
	I_sub
	I_mul
	I_div
	I_and
	I_or
	I_xor
	I_not
	I_shl
	I_shr
	I_call
	I_ret
	I_sys
	I_jmp
	I_jz
	I_je
	I_jg
	I_jl
)

func Throw(x string, z ...interface{}) {
	panic(fmt.Errorf(x, z...))
}

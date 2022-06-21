// Copyright (c) 2022 Nadeen Udantha <me@nadeen.lk>. All rights reserved.

package boilang

import (
	"bytes"
	"debug/elf"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCC(t *testing.T) {
	f, err := elf.Open(`D:\nadeen\boi\tcc\boi\test.o`)
	assert.NoError(t, err)
	d, err := f.Section(".text").Data()
	assert.NoError(t, err)
	Disassemble(bytes.NewBuffer(d))
}

func TestASM(t *testing.T) {
	xx, _ := ioutil.ReadFile("instr.boi")
	for _, x := range strings.Split(string(xx), "\n") {
		x = strings.TrimSpace(x)
		if len(x) == 0 {
			continue
		}
		fmt.Println("input: ", x)
		y := Assemble(x)
		fmt.Println("output: ", y.Bytes())
		yy, _ := Disassemble(y)
		yy = strings.TrimSpace(yy)
		fmt.Println("input2: ", yy)
		assert.Equal(t, x, yy)
	}
}

func TestBoiBuf(t *testing.T) {
	bb := BoiBuf{buf: new(bytes.Buffer)}
	bb.u8(123)
	assert.Equal(t, uint8(123), bb.ru8())
	bb.u16(123)
	assert.Equal(t, uint16(123), bb.ru16())
	bb.u16(1234)
	assert.Equal(t, uint16(1234), bb.ru16())
	bb.u16(12345)
	assert.Equal(t, uint16(12345), bb.ru16())
}

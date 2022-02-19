package registers

import (
	"fmt"
)

type Registers struct {
	A     byte
	F     byte
	B     byte
	C     byte
	D     byte
	E     byte
	H     byte
	L     byte
	PC    uint16 // It's not really a 64 bit integer but to read bytes from the game file we need it like this
	SP    uint16
	Stack []byte
	ZF    bool
	NF    bool
	HF    bool
	CF    bool
}

// GetInitializedRegisters initializes a new set of registers to their zero values (for the GB, ofc)
// after the checks that the GB is supposed to perform.
func GetInitializedRegisters() *Registers {
	r := new(Registers)
	r.A = 0x01
	r.F = 0xB0
	r.B = 0x00
	r.C = 0x13
	r.D = 0x00
	r.E = 0xD8
	r.H = 0x01
	r.L = 0x4D
	r.PC = 0x100
	r.SP = 0xE000
	r.Stack = make([]byte, 0)
	r.ZF = false
	r.NF = false
	r.HF = false
	r.CF = false
	return r
}

func (r *Registers) String() string {
	return fmt.Sprintf(
		"A: %X\nF: %X\nB: %X\nC: %X\nD: %X\nE: %X\nH: %X\nL: %X\nPC: %X\nSP: %X\n",
		r.A, r.F, r.B, r.C, r.D, r.E, r.H, r.L, r.PC, r.SP,
	)
}

func (r *Registers) AF() uint16 {
	return uint16(r.A)<<8 + uint16(r.F)
}

func (r *Registers) BC() uint16 {
	return uint16(r.B)<<8 + uint16(r.C)
}

func (r *Registers) DE() uint16 {
	return uint16(r.D)<<8 + uint16(r.E)
}

func (r *Registers) HL() uint16 {
	return uint16(r.H)<<8 + uint16(r.L)
}

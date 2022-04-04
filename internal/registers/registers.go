package registers

import (
	"fmt"
)

// Registers in a GB CPU. Flags are represented as booleans too so that it's easier to use them as such.
type Registers struct {
	A      byte
	F      byte
	B      byte
	C      byte
	D      byte
	E      byte
	H      byte
	L      byte
	PC     uint16
	SP     uint16
	ZF     bool
	NF     bool
	HF     bool
	CF     bool
	Halted bool
}

// GetInitializedRegisters initializes a new set of registers to their zero values (for the GB, ofc)
// after the checks that the GB is supposed to perform.
func GetInitializedRegisters() *Registers {
	r := Registers{
		A:      0x01,
		F:      0xB0,
		B:      0x00,
		C:      0x13,
		D:      0x00,
		E:      0xD8,
		H:      0x01,
		L:      0x4D,
		PC:     0x100,
		SP:     0xE000,
		ZF:     false,
		NF:     false,
		HF:     false,
		CF:     false,
		Halted: false,
	}

	return &r
}

func (r *Registers) String() string {
	return fmt.Sprintf(
		"A: %X\nF: %X\nB: %X\nC: %X\nD: %X\nE: %X\nH: %X\nL: %X\nPC: %X\nSP: %X\nZF: %t\nNF: %t\nHF: %t\nCF: %t\n",
		r.A, r.F, r.B, r.C, r.D, r.E, r.H, r.L, r.PC, r.SP, r.ZF, r.NF, r.HF, r.CF,
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

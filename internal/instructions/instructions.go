package instructions

import (
	"fmt"
	"go-boy/internal/memory"
	"go-boy/internal/registers"
	"go-boy/internal/utils"
)

func unimplemented(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	return fmt.Errorf("unimplemented instruction reached at PC=%04X: %02X %02X", r.PC, args[0], args[1]), 0
}

// 0x00
// Doesn't do anything.
func nop(_ *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return nil, 1
}

// 0x01
// Loads two bytes immediate into BC
func ldBCnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.B = args[2]
	r.C = args[1]
	return nil, 3
}

// 0x03
// Increments the value in BC
func incBC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	bc := r.BC() + 1
	r.B = byte(bc >> 8)
	r.C = byte(bc)
	return nil, 1
}

// 0x04
// Increments the value in B
func incB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.B&0x0F == 0x0F
	r.B++
	r.ZF = r.B == 0
	return nil, 1
}

// 0x05
// Decrements the value in B.
func decB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.HF = r.B&0x0F == 0
	r.B--
	r.ZF = r.B == 0
	return nil, 1
}

// 0x06
// Loads an 8 bit int into B
func ldBn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.B = args[1]
	return nil, 2
}

// 0x07
// Rotates A left, bit 7 to carry.
func rlcA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = false
	carry := r.A >> 7
	r.CF = carry == 1
	r.A <<= 1
	r.A += carry
	r.ZF = r.A == 0
	return nil, 1
}

// 0x08
// Loads SP into memory address nn.
func ldnnSP(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	m.Store(uint16(args[2])<<8+uint16(args[1]), byte(r.SP))
	m.Store(uint16(args[2])<<8+uint16(args[1])+1, byte(r.SP>>8))
	return nil, 3
}

// 0x09
// Adds HL to BC, result to HL.
func addHLBC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	bc := r.BC()
	hl := r.HL()
	r.HF = bc&0xFFF+hl&0xFFF > 0xFFF
	newHL := uint32(bc) + uint32(hl)
	r.CF = newHL > 0xFFFF
	r.H = byte(newHL >> 8)
	r.L = byte(newHL)
	return nil, 1
}

// 0x0A
// Loads A from address pointed to by BC
func ldABC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.A = m.Read(r.BC())
	return nil, 1
}

// 0x0B
// Decrements BC
func decBC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	bc := r.BC() - 1
	r.B = byte(bc >> 8)
	r.C = byte(bc)
	return nil, 1
}

// 0x0C
// Increments the value in C
func incC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.C&0x0F == 0x0F
	r.C++
	r.ZF = r.C == 0
	return nil, 1
}

// 0x0D
// Decrements the value in C
func decC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.HF = r.C&0x0F == 0
	r.C--
	r.ZF = r.C == 0
	return nil, 1
}

// 0x0E
// Loads an 8 bit int into C
func ldCn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.C = args[1]
	return nil, 2
}

// 0x11
// Loads a 16 bit int into DE
func ldDEnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.D = args[2]
	r.E = args[1]
	return nil, 3
}

// 0x12
// Copy a to memory address pointed by DE
func ldDEA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.DE(), r.A)
	return nil, 1
}

// 0x13
// Increments the value in DE
func incDE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	de := r.DE() + 1
	r.D = byte(de >> 8)
	r.E = byte(de)
	return nil, 1
}

// 0x14
// Increments the value in D
func incD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.D&0x0F == 0x0F
	r.D++
	r.ZF = r.D == 0
	return nil, 1
}

// 0x15
// Decrements the value in D
func decD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.HF = r.D&0x0F == 0
	r.D--
	r.ZF = r.D == 0
	return nil, 1
}

// 0x16
// Load an 8 bit int into D
func ldDn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.D = args[1]
	return nil, 2
}

// 0x17
// Rotates A left through carry flag
func rlA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = false
	carry := r.A >> 7
	r.A <<= 1
	if r.CF {
		r.A += 1
	}
	r.CF = carry == 1
	r.ZF = r.A == 0
	return nil, 1
}

// 0x18
// Add n to address and jump to it
func jrn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.PC = uint16(int16(r.PC) + int16(int8(args[1])))
	return nil, 2
}

// 0x19
// Adds DE to HL. Stores the result in HL.
func addHLDE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	de := r.DE()
	hl := r.HL()
	r.HF = de&0xFFF+hl&0xFFF > 0xFFF
	newHL := uint32(de) + uint32(hl)
	r.CF = newHL > 0xFFFF
	r.H = byte(newHL >> 8)
	r.L = byte(newHL)
	return nil, 1
}

// 0x1A
// Saves byte in memory address pointed by DE in A.
func ldADE(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.A = m.Read(r.DE())
	return nil, 1
}

// 0x1C
// Increments the value in E
func incE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.E&0x0F == 0x0F
	r.E++
	r.ZF = r.E == 0
	return nil, 1
}

// 0x1D
// Decrements the value in E
func decE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.HF = r.E&0x0F == 0
	r.E--
	r.ZF = r.E == 0
	return nil, 1
}

// 0x1E
// Load an 8 bit int into E
func ldEn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.E = args[1]
	return nil, 2
}

// 0x1F
// Rotate A right
func rrA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = false
	var carry byte
	if r.CF {
		carry = 0x80
	} else {
		carry = 0
	}
	r.CF = r.A&0b1 == 1
	r.A = (r.A >> 1) | carry
	r.ZF = r.A == 0
	return nil, 1
}

// 0x20
// Adds a specific signed amount to PC if Z flag is unset
func jrNZn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if !r.ZF {
		r.PC = uint16(int16(r.PC) + int16(int8(args[1])))
	}
	return nil, 2
}

// 0x21
// Loads a 16 bit int into HL
func ldHLnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.L = args[1]
	r.H = args[2]
	return nil, 3
}

// 0x22
// Saves A in memory address pointed at by HL, then increments HL
func ldiHLA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), r.A)
	_, _ = incHL(r, nil, nil)
	return nil, 1
}

// 0x23
// Increments the value in HL
func incHL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	hl := r.HL() + 1
	r.H = byte(hl >> 8)
	r.L = byte(hl)
	return nil, 1
}

// 0x25
// Decrements the value in H
func decH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.HF = r.H&0x0F == 0
	r.H--
	r.ZF = r.H == 0
	return nil, 1
}

// 0x28
// If ZF is set, add to PC and jump
func jrZn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if r.ZF {
		r.PC = uint16(int16(r.PC) + int16(int8(args[1])))
	}
	return nil, 2
}

// 0x29
// Adds HL to HL. Basically HL*2.
func addHL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	hl := r.HL()
	r.HF = (hl&0xFFF)<<1 > 0xFFF
	r.CF = uint32(hl)<<1 > 0xFFFF
	r.H = byte((hl << 1) >> 8)
	r.L = byte(hl << 1)
	return nil, 1
}

// 0x2A
// Stores the contents of the memory address HL into A, then increments HL
func ldiAHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.A = m.Read(r.HL())
	_, _ = incHL(r, nil, nil)
	return nil, 1
}

// 0x2B
// Decrements the value in HL.
func decHL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	hl := r.HL()
	hl--
	r.H = byte(hl >> 8)
	r.L = byte(hl)
	return nil, 1
}

// 0x2C
// Increments the value in L.
func incL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.L&0x0F == 0x0F
	r.L++
	r.ZF = r.L == 0
	return nil, 1
}

// 0x2D
// Decrements the value in L.
func decL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.HF = r.L&0x0F == 0
	r.L--
	r.ZF = r.L == 0
	return nil, 1
}

// 0x2F
// Complements A (flip all bits / not A)
func cpl(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = ^r.A
	return nil, 1
}

// 0x30
// Adds a specific signed amount to PC if C flag is unset
func jrNCn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if !r.CF {
		r.PC = uint16(int16(r.PC) + int16(int8(args[1])))
	}
	return nil, 2
}

// 0x31
// Loads nn into SP.
func ldSPnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.SP = uint16(args[2])<<8 + uint16(args[1])
	return nil, 3
}

// 0x32
// Stores the contents of A into the memory address HL, then decrements HL.
func lddHLA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), r.A)
	_, _ = decHL(r, nil, nil)
	return nil, 1
}

// 0x34
// Increments the contents in the memory address HL.
func incPHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	data := m.Read(r.HL())
	r.NF = false
	r.HF = data&0x0F == 0x0F
	data++
	r.ZF = data == 0
	m.Store(r.HL(), data)
	return nil, 1
}

// 0x35
// Decrements the contents in the memory address HL.
func decPHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	data := m.Read(r.HL())
	r.NF = true
	r.HF = data&0x0F == 0
	data--
	r.ZF = data == 0
	m.Store(r.HL(), data)
	return nil, 1
}

// 0x36
// Stores an immediate byte into the memory address HL.
func ldHLn(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	m.Store(r.HL(), args[1])
	return nil, 2
}

// 0x37
// Sets carry flag.
func scf(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.CF = true
	r.NF = false
	r.HF = false
	return nil, 1
}

// 0x38
// If CF is set, add to PC and jump
func jrCn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if r.CF {
		r.PC = uint16(int16(r.PC) + int16(int8(args[1])))
	}
	return nil, 2
}

// 0x3A
// Stores the contents of memory address HL into A, then decrements HL.
func lddAHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.A = m.Read(r.HL())
	_, _ = decHL(r, nil, nil)
	return nil, 1
}

// 0x3C
// Increments A.
func incA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.A&0x0F == 0x0F
	r.A++
	r.ZF = r.A == 0
	return nil, 1
}

// 0x3D
// Decrements A.
func decA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.HF = r.A&0x0F == 0
	r.A--
	r.ZF = r.A == 0
	return nil, 1
}

// 0x3E
// Stores the immediate byte into A.
func ldAn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.A = args[1]
	return nil, 2
}

// 0x40
// Copies B to B.
func ldBB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return nil, 1
}

// 0x41
// Copies C to B.
func ldBC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.B = r.C
	return nil, 1
}

// 0x46
// Copies contents of memory address HL into B
func ldBHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.B = m.Read(r.HL())
	return nil, 1
}

// 0x47
// Copies A to B
func ldBA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.B = r.A
	return nil, 1
}

// 0x4E
// Copies contents of memory address HL into C
func ldCHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.C = m.Read(r.HL())
	return nil, 1
}

// 0x4F
// Copies A to C
func ldCA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.C = r.A
	return nil, 1
}

// 0x54
// Copies H to D
func ldDH(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.D = r.H
	return nil, 1
}

// 0x56
// Copies the contents in memory address HL to D
func ldDHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.D = m.Read(r.HL())
	return nil, 1
}

// 0x57
// Copies A to D
func ldDA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.D = r.A
	return nil, 1
}

// 0x5D
// Copies L to E
func ldEL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.E = r.L
	return nil, 1
}

// 0x5E
// Copies the contents in memory address HL to E
func ldEHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.E = m.Read(r.HL())
	return nil, 1
}

// 0x5F
// Copies A to E
func ldEA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.E = r.A
	return nil, 1
}

// 0x60
// Copies B to H
func ldHB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.H = r.B
	return nil, 1
}

// 0x62
// Copies D to H
func ldHD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.H = r.D
	return nil, 1
}

// 0x66
// Copies value in memory address HL to H
func ldHHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.H = m.Read(r.HL())
	return nil, 1
}

// 0x67
// Copies A to H
func ldHA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.H = r.A
	return nil, 1
}

// 0x68
// Copies B to L
func ldLB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.L = r.B
	return nil, 1
}

// 0x69
// Copies C to L
func ldLC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.L = r.C
	return nil, 1
}

// 0x6B
// Copies E to L
func ldLE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.L = r.E
	return nil, 1
}

// 0x6F
// Copies A to L
func ldLA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.L = r.A
	return nil, 1
}

// 0x71
// Stores the contents of C into the memory address HL.
func ldHLC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), r.C)
	return nil, 1
}

// 0x72
// Stores the contents of D into the memory address HL.
func ldHLD(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), r.D)
	return nil, 1
}

// 0x73
// Stores the contents of E into the memory address HL.
func ldHLE(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), r.E)
	return nil, 1
}

// 0x77
// Stores the contents of A into the memory address HL.
func ldHLA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), r.A)
	return nil, 1
}

// 0x78
// Copies B to A
func ldAB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = r.B
	return nil, 1
}

// 0x79
// Copies C to A
func ldAC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = r.C
	return nil, 1
}

// 0x7A
// Copies D to A
func ldAD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = r.D
	return nil, 1
}

// 0x7B
// Stores the contents of E into A.
func ldAE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = r.E
	return nil, 1
}

// 0x7C
// Copies H into A
func ldAH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = r.H
	return nil, 1
}

// 0x7D
// Copies L into A
func ldAL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = r.L
	return nil, 1
}

// 0x7E
// Copies contents of memory address HL into A
func ldAHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.A = m.Read(r.HL())
	return nil, 1
}

// 0x80
// Adds A + B, result to A.
func addAB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.A&0x0F+r.B&0x0F > 0x0F
	r.CF = uint16(r.A)+uint16(r.B) > 0x00FF
	r.A = r.A + r.B
	r.ZF = r.A == 0
	return nil, 1
}

// 0x82
// Adds A + D, result to A.
func addAD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.A&0x0F+r.D&0x0F > 0x0F
	r.CF = uint16(r.A)+uint16(r.D) > 0x00FF
	r.A = r.A + r.D
	r.ZF = r.A == 0
	return nil, 1
}

// 0x85
// Adds A + L, result to A.
func addAL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.A&0x0F+r.L&0x0F > 0x0F
	r.CF = uint16(r.A)+uint16(r.L) > 0x00FF
	r.A = r.A + r.L
	r.ZF = r.A == 0
	return nil, 1
}

// 0x87
// Adds A + A.
func addAA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = r.A&0x0F+r.A&0x0F > 0x0F
	r.CF = uint16(r.A)+uint16(r.A) > 0x00FF
	r.A = r.A + r.A
	r.ZF = r.A == 0
	return nil, 1
}

// 0x8A
// Adds D + carry to A.
func adcAD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	var carry byte
	if r.CF {
		carry = 1
	} else {
		carry = 0
	}
	r.HF = r.A&0x0F+r.D&0x0F+carry > 0x0F
	r.CF = uint16(r.A)+uint16(r.D)+uint16(carry) > 0x00FF
	r.A = r.A + r.D + carry
	r.ZF = r.A == 0
	return nil, 1
}

// 0x93
// Subtracts E from A, result to A.
func subAE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.CF = r.E > r.A
	r.HF = r.E&0x0F > r.A&0x0F
	r.A -= r.E
	r.ZF = r.A == 0
	return nil, 1
}

// 0xA1
// Performs AND of A against C, result to A.
func andC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A &= r.C
	r.ZF = r.A == 0
	r.NF = false
	r.HF = true
	r.CF = false
	return nil, 1
}

// 0xA7
// Performs AND of A against itself.
func andA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A &= r.A
	r.ZF = r.A == 0
	r.NF = false
	r.HF = true
	r.CF = false
	return nil, 1
}

// 0xA9
// Performs an XOR of the register C against A, result to A.
func xorC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A ^= r.C
	r.ZF = r.A == 0
	r.NF = false
	r.HF = false
	r.CF = false
	return nil, 1
}

// 0xAF
// Performs an XOR of the register A against itself. To simplify, I just set A as 0.
func xorA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = 0
	r.ZF = true
	r.NF = false
	r.HF = false
	r.CF = false
	return nil, 1
}

// 0xB0
// Performs an OR of B against A, stores result in A.
func orB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A |= r.B
	r.ZF = r.A == 0
	r.NF = false
	r.HF = false
	r.CF = false
	return nil, 1
}

// 0xB1
// Performs an OR of C against A, stores result in A.
func orC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A |= r.C
	r.ZF = r.A == 0
	r.NF = false
	r.HF = false
	r.CF = false
	return nil, 1
}

// 0xB6
// Performs an OR of memory address HL against A, stores result in A.
func orHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.A |= m.Read(r.HL())
	r.ZF = r.A == 0
	r.NF = false
	r.HF = false
	r.CF = false
	return nil, 1
}

// 0xB9
// Compares C against A.
func cpC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = true
	r.CF = r.C > r.A
	r.HF = r.C&0x0F > r.A&0x0F
	r.ZF = r.A-r.C == 0
	return nil, 1
}

// 0xBF
// Compares A against A. Basically sets some flags.
func cpA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.ZF = true
	r.NF = true
	r.HF = false
	r.CF = false
	return nil, 1
}

// 0xC0
// Returns if ZF is reset.
func retNZ(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	if !r.ZF {
		r.PC = utils.PopStackShort(r, m)
		return nil, 0
	}
	return nil, 1
}

// 0xC1
// Pops BC.
func popBC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	bc := utils.PopStackShort(r, m)
	r.B = byte(bc >> 8)
	r.C = byte(bc)
	return nil, 1
}

// 0xC2
// Sets PC as specified in the arguments if Z flag is reset
func jpNZnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if !r.ZF {
		r.PC = uint16(args[1]) + uint16(args[2])<<8
		return nil, 0
	}
	return nil, 3
}

// 0xC3
// Sets PC to the specified address in the arguments.
func jpnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.PC = uint16(args[1]) + uint16(args[2])<<8
	return nil, 0 // Return 0 because it's a jump
}

// 0xC5
// Pushes BC into the stack.
func pushBC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.BC())
	return nil, 1
}

// 0xC6
// Adds A + immediate byte, result to A.
func addAn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.NF = false
	r.HF = r.A&0x0F+args[1]&0x0F > 0x0F
	r.CF = uint16(r.A)+uint16(args[1]) > 0x00FF
	r.A = r.A + args[1]
	r.ZF = r.A == 0
	return nil, 2
}

// 0xC8
// Returns if ZF is set.
func retZ(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	if r.ZF {
		r.PC = utils.PopStackShort(r, m)
		return nil, 0
	}
	return nil, 1
}

// 0xC9
// Pops two bytes from the stack, then sets PC as those.
func ret(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.PC = utils.PopStackShort(r, m)
	return nil, 0
}

// 0xCA
// Sets PC as specified in the arguments if Z flag is set
func jpZnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if r.ZF {
		r.PC = uint16(args[1]) + uint16(args[2])<<8
		return nil, 0
	}
	return nil, 3
}

// 0xCB
// Execute a 2 byte instruction
func execCB(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	instruction := CBTable[args[1]]
	return instruction(r, m, args)
}

// 0xCD
// Calls a function (pushes the next instruction address to the stack and goes to nn)
func callnn(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.PC+3)
	r.PC = uint16(args[1]) + uint16(args[2])<<8
	return nil, 0
}

// 0xCE
// Adds immediate byte + carry to A.
func adcAn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.NF = false
	var carry byte
	if r.CF {
		carry = 1
	} else {
		carry = 0
	}
	r.HF = r.A&0x0F+args[1]&0x0F+carry > 0x0F
	r.CF = uint16(r.A)+uint16(args[1])+uint16(carry) > 0x00FF
	r.A = r.A + args[1] + carry
	r.ZF = r.A == 0
	return nil, 2
}

// 0xD1
// Pops DE.
func popDE(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	de := utils.PopStackShort(r, m)
	r.D = byte(de >> 8)
	r.E = byte(de)
	return nil, 1
}

// 0xD2
// Sets PC as specified in the arguments if C flag is unset
func jpNCnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if !r.CF {
		r.PC = uint16(args[1]) + uint16(args[2])<<8
		return nil, 0
	}
	return nil, 3
}

// 0xD5
// Pushes DE into the stack.
func pushDE(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.DE())
	return nil, 1
}

// 0xD9
// Pops two bytes from the stack and assigns them to PC, then enables interrupts.
func reti(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.PC = utils.PopStackShort(r, m)
	m.IME = true
	return nil, 0
}

// 0xE0
// Put A into memory address FF00 + n
func ldhnA(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	m.Store(0xFF00+uint16(args[1]), r.A)
	return nil, 2
}

// 0xE1
// Pops HL
func popHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	hl := utils.PopStackShort(r, m)
	r.H = byte(hl >> 8)
	r.L = byte(hl)
	return nil, 1
}

// 0xE2
// Put A into memory address FF00 + register C
func ldhCA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(0xFF00+uint16(r.C), r.A)
	return nil, 1
}

// 0xE5
// Pushes HL into the stack.
func pushHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.HL())
	return nil, 1
}

// 0xE6
// Perform AND between a number and A, store result in A.
func andn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.A &= args[1]
	r.ZF = r.A == 0
	r.NF = false
	r.HF = true
	r.CF = false
	return nil, 2
}

// 0xE9
// Jumps to the address stored in HL
func jpHL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.PC = r.HL()
	return nil, 0
}

// 0xEA
// Put A into immediate memory address
func ldnnA(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	m.Store(uint16(args[2])<<8+uint16(args[1]), r.A)
	return nil, 3
}

// 0xEF
// Calls routine at 0x0028
func rst28(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.PC+1)
	r.PC = 0x0028
	return nil, 0
}

// 0xF0
// Load the contents of FF00 + n into A
func ldhAn(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	r.A = m.Read(0xFF00 + uint16(args[1]))
	return nil, 2
}

// 0xF1
// Pops AF.
func popAF(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	af := utils.PopStackShort(r, m)
	r.A = byte(af >> 8)
	r.F = byte(af)
	return nil, 1
}

// 0xF3
// Disable interrupts
func di(_ *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.IMESteps = 1
	m.IMEReqType = false
	return nil, 1
}

// 0xF5
// Push AF into the stack
func pushAF(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.AF())
	return nil, 1
}

// 0xF6
// Performs an OR of an immediate byte against A, stores result in A.
func orn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.A |= args[1]
	r.ZF = r.A == 0
	r.NF = false
	r.HF = false
	r.CF = false
	return nil, 1
}

// 0xFA
// Puts in A the value in memory address nn
func ldAnn(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	r.A = m.Read(uint16(args[2])<<8 + uint16(args[1]))
	return nil, 3
}

// 0xFB
// Enable interrupts
func ei(_ *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.IMESteps = 1
	m.IMEReqType = true
	return nil, 1
}

// 0xFE
// Compare A with a byte
func cpn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	n := args[1]
	r.NF = true
	r.CF = n > r.A
	r.HF = n&0x0F > r.A&0x0F
	r.ZF = r.A-n == 0
	return nil, 2
}

// 0xFF
// Calls routine at 0x0038
func rst38(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.PC+1)
	r.PC = 0x0038
	return nil, 0
}

var InstructionTable = [256]func(_ *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16){
	nop, // 0x00
	ldBCnn,
	unimplemented,
	incBC,
	incB,
	decB,
	ldBn,
	rlcA,
	ldnnSP,
	addHLBC,
	ldABC,
	decBC,
	incC,
	decC,
	ldCn,
	unimplemented,
	unimplemented, // 0x10
	ldDEnn,
	ldDEA,
	incDE,
	incD,
	decD,
	ldDn,
	rlA,
	jrn,
	addHLDE,
	ldADE,
	unimplemented,
	incE,
	decE,
	ldEn,
	rrA,
	jrNZn, // 0x20
	ldHLnn,
	ldiHLA,
	incHL,
	unimplemented,
	decH,
	unimplemented,
	unimplemented,
	jrZn,
	addHL,
	ldiAHL,
	decHL,
	incL,
	decL,
	unimplemented,
	cpl,
	jrNCn, // 0x30
	ldSPnn,
	lddHLA,
	unimplemented,
	incPHL,
	decPHL,
	ldHLn,
	scf,
	jrCn,
	unimplemented,
	lddAHL,
	unimplemented,
	incA,
	decA,
	ldAn,
	unimplemented,
	ldBB, // 0x40
	ldBC,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	ldBHL,
	ldBA,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	ldCHL,
	ldCA,
	unimplemented, // 0x50
	unimplemented,
	unimplemented,
	unimplemented,
	ldDH,
	unimplemented,
	ldDHL,
	ldDA,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	ldEL,
	ldEHL,
	ldEA,
	ldHB, // 0x60
	unimplemented,
	ldHD,
	unimplemented,
	unimplemented,
	unimplemented,
	ldHHL,
	ldHA,
	ldLB,
	ldLC,
	unimplemented,
	ldLE,
	unimplemented,
	unimplemented,
	unimplemented,
	ldLA,
	unimplemented, // 0x70
	ldHLC,
	ldHLD,
	ldHLE,
	unimplemented,
	unimplemented,
	unimplemented,
	ldHLA,
	ldAB,
	ldAC,
	ldAD,
	ldAE,
	ldAH,
	ldAL,
	ldAHL,
	unimplemented,
	addAB, // 0x80
	unimplemented,
	addAD,
	unimplemented,
	unimplemented,
	addAL,
	unimplemented,
	addAA,
	unimplemented,
	unimplemented,
	adcAD,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x90
	unimplemented,
	unimplemented,
	subAE,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0xA0
	andC,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	andA,
	unimplemented,
	xorC,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	xorA,
	orB, // 0xB0
	orC,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	orHL,
	unimplemented,
	unimplemented,
	cpC,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	cpA,
	retNZ, // 0xC0
	popBC,
	jpNZnn,
	jpnn,
	unimplemented,
	pushBC,
	addAn,
	unimplemented,
	retZ,
	ret,
	jpZnn,
	execCB,
	unimplemented,
	callnn,
	adcAn,
	unimplemented,
	unimplemented, // 0xD0
	popDE,
	jpNCnn,
	unimplemented,
	unimplemented,
	pushDE,
	unimplemented,
	unimplemented,
	unimplemented,
	reti,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	ldhnA, // 0xE0
	popHL,
	ldhCA,
	unimplemented,
	unimplemented,
	pushHL,
	andn,
	unimplemented,
	unimplemented,
	jpHL,
	ldnnA,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	rst28,
	ldhAn, // 0xF0
	popAF,
	unimplemented,
	di,
	unimplemented,
	pushAF,
	orn,
	unimplemented,
	unimplemented,
	unimplemented,
	ldAnn,
	ei,
	unimplemented,
	unimplemented,
	cpn,
	rst38}

var cyclesTable = [256]int{
	2, 6, 4, 4, 2, 2, 4, 4, 10, 4, 4, 4, 2, 2, 4, 4, // 0x0_
	2, 6, 4, 4, 2, 2, 4, 4, 4, 4, 4, 4, 2, 2, 4, 4, // 0x1_
	0, 6, 4, 4, 2, 2, 4, 2, 0, 4, 4, 4, 2, 2, 4, 2, // 0x2_
	4, 6, 4, 4, 6, 6, 6, 2, 0, 4, 4, 4, 2, 2, 4, 2, // 0x3_
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 0x4_
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 0x5_
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 0x6_
	4, 4, 4, 4, 4, 4, 2, 4, 2, 2, 2, 2, 2, 2, 4, 2, // 0x7_
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 0x8_
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 0x9_
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 0xa_
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 0xb_
	0, 6, 0, 6, 0, 8, 4, 8, 0, 2, 0, 0, 0, 6, 4, 8, // 0xc_
	0, 6, 0, 0, 0, 8, 4, 8, 0, 8, 0, 0, 0, 0, 4, 8, // 0xd_
	6, 6, 4, 0, 0, 8, 4, 8, 8, 2, 8, 0, 0, 0, 4, 8, // 0xe_
	6, 6, 4, 2, 0, 8, 4, 8, 6, 4, 8, 2, 0, 0, 4, 8, // 0xf_
}

// 0xCB19
// Rotate C right through carry flag.
func rrC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = false
	var carry byte = 0x00
	if r.CF {
		carry = 0x80
	}
	r.CF = r.C&0x01 == 1
	r.C = (r.C >> 1) | carry
	r.ZF = r.C == 0
	return nil, 2
}

// 0xCB27
// Shift left A into carry, LSB = 0.
func slaA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = false
	r.CF = r.A>>7 == 1
	r.A <<= 1
	r.ZF = r.A == 0
	return nil, 2
}

// 0xCB37
// Swap nibbles of A.
func swapA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = r.A<<4 + r.A>>4
	r.ZF = r.A == 0
	r.CF = false
	r.HF = false
	r.NF = false
	return nil, 2
}

// 0xCB38
// Shift right B into carry, MSB = 0.
func srlB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = false
	r.CF = r.B&0x01 == 1
	r.B >>= 1
	r.ZF = r.B == 0
	return nil, 2
}

// 0xCB50
// Test bit 2 of B.
func test2B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.ZF = (r.B>>2)&0x01 == 0
	r.HF = true
	r.NF = false
	return nil, 2
}

// 0xCB58
// Test bit 3 of B.
func test3B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.ZF = (r.B>>3)&0x01 == 0
	r.HF = true
	r.NF = false
	return nil, 2
}

// 0xCB60
// Test bit 4 of B.
func test4B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.ZF = (r.B>>4)&0x01 == 0
	r.HF = true
	r.NF = false
	return nil, 2
}

// 0xCB68
// Test bit 5 of B.
func test5B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.ZF = (r.B>>5)&0x01 == 0
	r.HF = true
	r.NF = false
	return nil, 2
}

// 0xCB7E
// Test bit 7 of value in memory address HL.
func test7HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.ZF = (m.Read(r.HL())>>7)&0x01 == 0
	r.HF = true
	r.NF = false
	return nil, 2
}

// 0xCB7F
// Test bit 7 of A.
func test7A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.ZF = (r.A>>7)&0x01 == 0
	r.HF = true
	r.NF = false
	return nil, 2
}

// 0xCB86
// Reset bit 0 of value in memory address HL.
func res0HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())&0xFE)
	return nil, 2
}

// 0xCB87
// Reset bit 0 of A.
func res0A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A &= 0xFE
	return nil, 2
}

var CBTable = [256]func(_ *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16){
	unimplemented, // 0x00
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x10
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	rrC,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x20
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	slaA,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x30
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	swapA,
	srlB,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x40
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	test2B, // 0x50
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	test3B,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	test4B, // 0x60
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	test5B,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x70
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	test7HL,
	test7A,
	unimplemented, // 0x80
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	res0HL,
	res0A,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x90
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0xA0
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0xB0
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0xC0
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0xD0
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0xE0
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0xF0
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
}

var cbCyclesTable = [256]int{
	8, 8, 8, 8, 8, 8, 16, 8, 8, 8, 8, 8, 8, 8, 16, 8, // 0x0_
	8, 8, 8, 8, 8, 8, 16, 8, 8, 8, 8, 8, 8, 8, 16, 8, // 0x1_
	8, 8, 8, 8, 8, 8, 16, 8, 8, 8, 8, 8, 8, 8, 16, 8, // 0x2_
	8, 8, 8, 8, 8, 8, 16, 8, 8, 8, 8, 8, 8, 8, 16, 8, // 0x3_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0x4_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0x5_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0x6_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0x7_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0x8_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0x9_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0xa_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0xb_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0xc_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0xd_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0xe_
	8, 8, 8, 8, 8, 8, 12, 8, 8, 8, 8, 8, 8, 8, 12, 8, // 0xf_
}

func Execute(r *registers.Registers, m *memory.Memory, instructionArray []byte) (error, uint16, int) {
	// The instruction is the first byte in the array
	opCode := instructionArray[0]
	// Look up the instructions table and obtain the function that executes the instruction
	operation := InstructionTable[opCode]
	cycles := cyclesTable[opCode]
	// Execute the operation
	err, jump := operation(r, m, instructionArray)
	var f8 byte = 0x00
	var f7 byte = 0x00
	var f6 byte = 0x00
	var f5 byte = 0x00
	if r.ZF {
		f8 = 0x80
	}
	if r.NF {
		f7 = 0x40
	}
	if r.HF {
		f6 = 0x20
	}
	if r.CF {
		f5 = 0x10
	}
	r.F = f8 + f7 + f6 + f5
	return err, jump, cycles
}

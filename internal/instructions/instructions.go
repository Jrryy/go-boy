package instructions

import (
	"fmt"
	"go-boy/internal/memory"
	"go-boy/internal/registers"
	"go-boy/internal/utils"
)

var flagsRecovered bool = false

// Generic 8 bit add function.
func add(r *byte, op byte, nf, zf, hf, cf *bool) (error, uint16) {
	*nf = false
	*hf = *r&0x0F+op&0x0F > 0x0F
	*cf = uint16(*r)+uint16(op) > 0x00FF
	*r += op
	*zf = *r == 0
	return nil, 1
}

// Generic 16 bit add function.
func add16(op1, op2 uint16, rh, rl *byte, nf, zf, hf, cf *bool) (error, uint16) {
	*nf = false
	*hf = op1&0xFFF+op2&0xFFF > 0xFFF
	newHL := uint32(op1) + uint32(op2)
	*cf = newHL > 0xFFFF
	*rh = byte(newHL >> 8)
	*rl = byte(newHL)
	return nil, 1
}

// Generic 8 bit add with carry function.
func adc(r *byte, op byte, nf, zf, hf, cf *bool) (error, uint16) {
	*nf = false
	var carry byte
	if *cf {
		carry = 1
	} else {
		carry = 0
	}
	*hf = *r&0x0F+op&0x0F+carry > 0x0F
	*cf = uint16(*r)+uint16(op)+uint16(carry) > 0x00FF
	*r += op + carry
	*zf = *r == 0
	return nil, 1
}

// Generic subtract function.
func sub(r *byte, op byte, nf, zf, hf, cf *bool, immediate bool) (error, uint16) {
	*nf = true
	*cf = op > *r
	*hf = op&0x0F > *r&0x0F
	*r -= op
	*zf = *r == 0
	if immediate {
		return nil, 2
	} else {
		return nil, 1
	}
}

// Generic 8 bit load function.
func ld(r *byte, n byte, immediate bool) (error, uint16) {
	*r = n
	if immediate {
		return nil, 2
	} else {
		return nil, 1
	}
}

// Generic 16 bit load function.
func ld16(rh, rl *byte, n uint16, immediate bool) (error, uint16) {
	*rh = byte(n >> 8)
	*rl = byte(n)
	if !immediate {
		return nil, 1
	} else {
		return nil, 3
	}
}

// Generic or function.
func or(r *byte, op byte, nf, zf, hf, cf *bool) (error, uint16) {
	*r |= op
	*zf = *r == 0
	*nf = false
	*hf = false
	*cf = false
	return nil, 1
}

// Generic and function.
func and(r *byte, op byte, nf, zf, hf, cf *bool) (error, uint16) {
	*r &= op
	*zf = *r == 0
	*nf = false
	*hf = true
	*cf = false
	return nil, 1
}

// Generic xor function.
func xor(r *byte, op byte, nf, zf, hf, cf *bool, immediate bool) (error, uint16) {
	*r ^= op
	*zf = *r == 0
	*nf = false
	*hf = false
	*cf = false
	if immediate {
		return nil, 2
	} else {
		return nil, 1
	}
}

// Generic compare function.
func cp(op1, op2 byte, nf, zf, hf, cf *bool) (error, uint16) {
	*nf = true
	*cf = op2 > op1
	*hf = op2&0x0F > op1&0x0F
	*zf = op1-op2 == 0
	return nil, 1
}

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
	return ld16(&r.B, &r.C, uint16(args[2])<<8+uint16(args[1]), true)
}

// 0x02
// Copy A to memory address pointed by BC
func ldBCA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.BC(), r.A)
	return nil, 1
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
	return ld(&r.B, args[1], true)
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
	return add16(r.BC(), r.HL(), &r.H, &r.L, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x0A
// Loads A from address pointed to by BC
func ldABC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.A, m.Read(r.BC()), false)
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
	return ld(&r.C, args[1], true)
}

// 0x11
// Loads a 16 bit int into DE
func ldDEnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.D = args[2]
	r.E = args[1]
	return nil, 3
}

// 0x12
// Copy A to memory address pointed by DE
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
	return ld(&r.D, args[1], true)
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
	return add16(r.DE(), r.HL(), &r.H, &r.L, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x1A
// Saves byte in memory address pointed by DE in A.
func ldADE(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.A = m.Read(r.DE())
	return nil, 1
}

// 0x1B
// Decrements the value in DE
func decDE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	de := r.DE() - 1
	r.D = byte(de >> 8)
	r.E = byte(de)
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
	return ld(&r.E, args[1], true)
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

// 0x26
// Load an 8 bit int into H
func ldHn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	return ld(&r.H, args[1], true)
}

// 0x27
// Decimal adjust after addition.
func daa(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	a := uint16(r.A)
	if !r.NF {
		if a&0x0F > 9 || r.HF {
			a += 6
		}
		if a > 0x9F || r.CF {
			a += 0x60
		}
	} else {
		if r.HF {
			a -= 6
			a &= 0xFF
		}
		if r.CF {
			a -= 0x60
		}
	}
	r.HF = false
	r.ZF = a == 0
	r.CF = a >= 0x100
	r.A = byte(a)
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
// Adds HL to HL.
func addHLHL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return add16(r.HL(), r.HL(), &r.H, &r.L, &r.NF, &r.ZF, &r.HF, &r.CF)
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

// 0x2E
// Load an 8 bit int into L
func ldLn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	return ld(&r.L, args[1], true)
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
	return ld(&r.B, r.B, false)
}

// 0x41
// Copies C to B.
func ldBC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.B, r.C, false)
}

// 0x42
// Copies D to B.
func ldBD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.B, r.D, false)
}

// 0x43
// Copies D to B.
func ldBE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.B, r.E, false)
}

// 0x44
// Copies H to B.
func ldBH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.B, r.H, false)
}

// 0x45
// Copies D to B.
func ldBL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.B, r.L, false)
}

// 0x46
// Copies contents of memory address HL into B.
func ldBHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.B, m.Read(r.HL()), false)
}

// 0x47
// Copies A to B.
func ldBA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.B, r.A, false)
}

// 0x48
// Copies B to C.
func ldCB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.C, r.B, false)
}

// 0x49
// Copies C to C.
func ldCC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.C, r.C, false)
}

// 0x4A
// Copies D to C.
func ldCD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.C, r.D, false)
}

// 0x4B
// Copies D to C.
func ldCE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.C, r.E, false)
}

// 0x4C
// Copies H to C.
func ldCH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.C, r.H, false)
}

// 0x4D
// Copies D to C.
func ldCL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.C, r.L, false)
}

// 0x4E
// Copies contents of memory address HL into C
func ldCHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.C, m.Read(r.HL()), false)
}

// 0x4F
// Copies A to C
func ldCA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.C, r.A, false)
}

// 0x50
// Copies B to D
func ldDB(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.D, r.B, false)
}

// 0x51
// Copies C to D
func ldDC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.D, r.C, false)
}

// 0x52
// Copies D to D
func ldDD(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.D, r.D, false)
}

// 0x53
// Copies E to D
func ldDE(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.D, r.E, false)
}

// 0x54
// Copies H to D
func ldDH(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.D, r.H, false)
}

// 0x55
// Copies L to D
func ldDL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.D, r.L, false)
}

// 0x56
// Copies the contents in memory address HL to D
func ldDHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.D, m.Read(r.HL()), false)
}

// 0x57
// Copies A to D
func ldDA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.D, r.A, false)
}

// 0x58
// Copies B to E
func ldEB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.E, r.B, false)
}

// 0x59
// Copies C to E
func ldEC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.E, r.C, false)
}

// 0x5A
// Copies D to E
func ldED(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.E, r.D, false)
}

// 0x5B
// Copies E to E
func ldEE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.E, r.E, false)
}

// 0x5C
// Copies H to E
func ldEH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.E, r.H, false)
}

// 0x5D
// Copies L to E
func ldEL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.E, r.L, false)
}

// 0x5E
// Copies the contents in memory address HL to E
func ldEHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.E, m.Read(r.HL()), false)
}

// 0x5F
// Copies A to E
func ldEA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.E, r.A, false)
}

// 0x60
// Copies B to H
func ldHB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.H, r.B, false)
}

// 0x61
// Copies C to H
func ldHC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.H, r.C, false)
}

// 0x62
// Copies D to H
func ldHD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.H, r.D, false)
}

// 0x63
// Copies E to H
func ldHE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.H, r.E, false)
}

// 0x64
// Copies H to H
func ldHH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.H, r.H, false)
}

// 0x65
// Copies L to H
func ldHL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.H, r.L, false)
}

// 0x66
// Copies value in memory address HL to H
func ldHHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.H, m.Read(r.HL()), false)
}

// 0x67
// Copies A to H
func ldHA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.H, r.A, false)
}

// 0x68
// Copies B to L
func ldLB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.L, r.B, false)
}

// 0x69
// Copies C to L
func ldLC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.L, r.C, false)
}

// 0x6A
// Copies D to L
func ldLD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.L, r.D, false)
}

// 0x6B
// Copies E to L
func ldLE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.L, r.E, false)
}

// 0x6C
// Copies H to L
func ldLH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.L, r.H, false)
}

// 0x6D
// Copies L to L
func ldLL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.L, r.L, false)
}

// 0x6E
// Copies value in memory address HL to L
func ldLHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.L, m.Read(r.HL()), false)
}

// 0x6F
// Copies A to L
func ldLA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return ld(&r.L, r.A, false)
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

// 0x76
// Stop CPU until interruption occurs
func halt(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.Halted = true
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
	return add(&r.A, r.B, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x81
// Adds A + C, result to A.
func addAC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return add(&r.A, r.C, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x82
// Adds A + D, result to A.
func addAD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return add(&r.A, r.D, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x83
// Adds A + E, result to A.
func addAE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return add(&r.A, r.E, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x84
// Adds A + H, result to A.
func addAH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return add(&r.A, r.H, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x85
// Adds A + L, result to A.
func addAL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return add(&r.A, r.L, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x86
// Adds A + value pointed at by HL, result to A.
func addAHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return add(&r.A, m.Read(r.HL()), &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x87
// Adds A + A.
func addAA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return add(&r.A, r.A, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x88
// Adds B + carry to A.
func adcAB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return adc(&r.A, r.B, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x89
// Adds C + carry to A.
func adcAC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return adc(&r.A, r.C, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x8A
// Adds D + carry to A.
func adcAD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return adc(&r.A, r.D, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x8B
// Adds E + carry to A.
func adcAE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return adc(&r.A, r.E, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x8C
// Adds H + carry to A.
func adcAH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return adc(&r.A, r.H, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x8D
// Adds L + carry to A.
func adcAL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return adc(&r.A, r.L, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x8E
// Adds value in memory pointed at by HL + carry to A.
func adcAHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return adc(&r.A, m.Read(r.HL()), &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x8F
// Adds A + carry to A.
func adcAA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return adc(&r.A, r.A, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0x90
// Subtracts B from A, result to A.
func subAB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sub(&r.A, r.B, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0x91
// Subtracts C from A, result to A.
func subAC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sub(&r.A, r.C, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0x92
// Subtracts D from A, result to A.
func subAD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sub(&r.A, r.D, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0x93
// Subtracts E from A, result to A.
func subAE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sub(&r.A, r.E, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0x94
// Subtracts H from A, result to A.
func subAH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sub(&r.A, r.H, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0x95
// Subtracts L from A, result to A.
func subAL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sub(&r.A, r.L, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0x96
// Subtracts value in memory address HL from A, result to A.
func subAHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return sub(&r.A, m.Read(r.HL()), &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0x97
// Subtracts A from A, result to A.
func subAA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sub(&r.A, r.A, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xA0
// Performs AND of A against B, result to A.
func andB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return and(&r.A, r.B, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xA1
// Performs AND of A against C, result to A.
func andC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return and(&r.A, r.C, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xA2
// Performs AND of A against D, result to A.
func andD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return and(&r.A, r.D, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xA3
// Performs AND of A against E, result to A.
func andE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return and(&r.A, r.E, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xA4
// Performs AND of A against H, result to A.
func andH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return and(&r.A, r.H, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xA5
// Performs AND of A against L, result to A.
func andL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return and(&r.A, r.L, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xA6
// Performs AND of A against value in memory address HL, result to A.
func andHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return and(&r.A, m.Read(r.HL()), &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xA7
// Performs AND of A against itself.
func andA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return and(&r.A, r.A, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xA8
// Performs an XOR of the register B against A, result to A.
func xorB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return xor(&r.A, r.B, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xA9
// Performs an XOR of the register C against A, result to A.
func xorC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return xor(&r.A, r.C, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xAA
// Performs an XOR of the register D against A, result to A.
func xorD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return xor(&r.A, r.D, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xAB
// Performs an XOR of the register E against A, result to A.
func xorE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return xor(&r.A, r.E, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xAC
// Performs an XOR of the register H against A, result to A.
func xorH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return xor(&r.A, r.H, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xAD
// Performs an XOR of the register L against A, result to A.
func xorL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return xor(&r.A, r.L, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xAE
// Performs an XOR of the value in memory address HL against A, result to A.
func xorHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return xor(&r.A, m.Read(r.HL()), &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xAF
// Performs an XOR of the register A against itself.
func xorA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return xor(&r.A, r.A, &r.NF, &r.ZF, &r.HF, &r.CF, false)
}

// 0xB0
// Performs an OR of B against A, stores result in A.
func orB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return or(&r.A, r.B, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB1
// Performs an OR of C against A, stores result in A.
func orC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return or(&r.A, r.C, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB2
// Performs an OR of D against A, stores result in A.
func orD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return or(&r.A, r.D, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB3
// Performs an OR of E against A, stores result in A.
func orE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return or(&r.A, r.E, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB4
// Performs an OR of H against A, stores result in A.
func orH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return or(&r.A, r.H, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB5
// Performs an OR of L against A, stores result in A.
func orL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return or(&r.A, r.L, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB6
// Performs an OR of memory address HL against A, stores result in A.
func orHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return or(&r.A, m.Read(r.HL()), &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB7
// Performs an OR of A against A, stores result in A.
func orA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return or(&r.A, r.A, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB8
// Compares B against A.
func cpB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return cp(r.A, r.B, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xB9
// Compares C against A.
func cpC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return cp(r.A, r.C, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xBA
// Compares D against A.
func cpD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return cp(r.A, r.D, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xBB
// Compares E against A.
func cpE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return cp(r.A, r.E, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xBC
// Compares H against A.
func cpH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return cp(r.A, r.H, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xBD
// Compares L against A.
func cpL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return cp(r.A, r.L, &r.NF, &r.ZF, &r.HF, &r.CF)
}

// 0xBE
// Compares value pointed at by HL against A.
func cpHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return cp(r.A, m.Read(r.HL()), &r.NF, &r.ZF, &r.HF, &r.CF)
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

// 0xC4
// Calls a function if ZF is reset.
func callNZnn(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	if !r.ZF {
		utils.PushStackShort(r, m, r.PC+3)
		r.PC = uint16(args[1]) + uint16(args[2])<<8
		return nil, 0
	} else {
		return nil, 3
	}
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

// 0xCC
// Calls a function if ZF is set.
func callZnn(r *registers.Registers, m *memory.Memory, args []byte) (error, uint16) {
	if r.ZF {
		utils.PushStackShort(r, m, r.PC+3)
		r.PC = uint16(args[1]) + uint16(args[2])<<8
		return nil, 0
	} else {
		return nil, 3
	}
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

// 0xCF
// Calls routine at 0x0008
func rst08(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.PC+1)
	r.PC = 0x0008
	return nil, 0
}

// 0xD0
// Returns if CF is reset.
func retNC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	if !r.CF {
		r.PC = utils.PopStackShort(r, m)
		return nil, 0
	}
	return nil, 1
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

// 0xD6
// Subtracts 8 bit immediate from A.
func subAn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	return sub(&r.A, args[1], &r.NF, &r.ZF, &r.HF, &r.CF, true)
}

// 0xD7
// Calls routine at 0x0010
func rst10(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.PC+1)
	r.PC = 0x0010
	return nil, 0
}

// 0xD8
// Returns if CF is set.
func retC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	if r.CF {
		r.PC = utils.PopStackShort(r, m)
		return nil, 0
	}
	return nil, 1
}

// 0xD9
// Pops two bytes from the stack and assigns them to PC, then enables interrupts.
func reti(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.PC = utils.PopStackShort(r, m)
	m.IME = true
	return nil, 0
}

// 0xDA
// Sets PC as specified in the arguments if C flag is set
func jpCnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if r.CF {
		r.PC = uint16(args[1]) + uint16(args[2])<<8
		return nil, 0
	}
	return nil, 3
}

// 0xDF
// Calls routine at 0x0018
func rst18(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.PC+1)
	r.PC = 0x0018
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

// 0xE7
// Calls routine at 0x0020
func rst20(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.PC+1)
	r.PC = 0x0020
	return nil, 0
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

// 0xEE
// xor of an 8 bit immediate against A.
func xorn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	return xor(&r.A, args[1], &r.NF, &r.ZF, &r.HF, &r.CF, true)
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
	flagsRecovered = true
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

// 0xF7
// Calls routine at 0x0030
func rst30(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	utils.PushStackShort(r, m, r.PC+1)
	r.PC = 0x0030
	return nil, 0
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
	ldBCA,
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
	decDE,
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
	ldHn,
	daa,
	jrZn,
	addHLHL,
	ldiAHL,
	decHL,
	incL,
	decL,
	ldLn,
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
	ldBD,
	ldBE,
	ldBH,
	ldBL,
	ldBHL,
	ldBA,
	ldCB,
	ldCC,
	ldCD,
	ldCE,
	ldCH,
	ldCL,
	ldCHL,
	ldCA,
	ldDB, // 0x50
	ldDC,
	ldDD,
	ldDE,
	ldDH,
	ldDL,
	ldDHL,
	ldDA,
	ldEB,
	ldEC,
	ldED,
	ldEE,
	ldEH,
	ldEL,
	ldEHL,
	ldEA,
	ldHB, // 0x60
	ldHC,
	ldHD,
	ldHE,
	ldHH,
	ldHL,
	ldHHL,
	ldHA,
	ldLB,
	ldLC,
	ldLD,
	ldLE,
	ldLH,
	ldLL,
	ldLHL,
	ldLA,
	unimplemented, // 0x70
	ldHLC,
	ldHLD,
	ldHLE,
	unimplemented,
	unimplemented,
	halt,
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
	addAC,
	addAD,
	addAE,
	addAH,
	addAL,
	addAHL,
	addAA,
	adcAB,
	adcAC,
	adcAD,
	adcAE,
	adcAH,
	adcAL,
	adcAHL,
	adcAA,
	subAB, // 0x90
	subAC,
	subAD,
	subAE,
	subAH,
	subAL,
	subAHL,
	subAA,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	andB, // 0xA0
	andC,
	andD,
	andE,
	andH,
	andL,
	andHL,
	andA,
	xorB,
	xorC,
	xorD,
	xorE,
	xorH,
	xorL,
	xorHL,
	xorA,
	orB, // 0xB0
	orC,
	orD,
	orE,
	orH,
	orL,
	orHL,
	orA,
	cpB,
	cpC,
	cpD,
	cpE,
	cpH,
	cpL,
	cpHL,
	cpA,
	retNZ, // 0xC0
	popBC,
	jpNZnn,
	jpnn,
	callNZnn,
	pushBC,
	addAn,
	unimplemented,
	retZ,
	ret,
	jpZnn,
	execCB,
	callZnn,
	callnn,
	adcAn,
	rst08,
	retNC, // 0xD0
	popDE,
	jpNCnn,
	unimplemented,
	unimplemented,
	pushDE,
	subAn,
	rst10,
	retC,
	reti,
	jpCnn,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	rst18,
	ldhnA, // 0xE0
	popHL,
	ldhCA,
	unimplemented,
	unimplemented,
	pushHL,
	andn,
	rst20,
	unimplemented,
	jpHL,
	ldnnA,
	unimplemented,
	unimplemented,
	unimplemented,
	xorn,
	rst28,
	ldhAn, // 0xF0
	popAF,
	unimplemented,
	di,
	unimplemented,
	pushAF,
	orn,
	rst30,
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

// Generic shift left function.
func sla(r *byte, zf, nf, hf, cf *bool) (error, uint16) {
	*nf = false
	*hf = false
	*cf = *r>>7 == 1
	*r <<= 1
	*zf = *r == 0
	return nil, 2
}

// Generic function for swapping nibbles.
func swap(r *byte, zf, nf, hf, cf *bool) (error, uint16) {
	*r = *r<<4 + *r>>4
	*zf = *r == 0
	*cf = false
	*hf = false
	*nf = false
	return nil, 2
}

// Generic function for testing the nth bit of a register.
func test(bit int, r byte, zf, hf, nf *bool) (error, uint16) {
	*zf = (r>>bit)&0x01 == 0
	*hf = true
	*nf = false
	return nil, 2
}

// Generic reset bit function
func reset(r *byte, bit byte) (error, uint16) {
	*r &= (0xFF ^ (0x01 << bit))
	return nil, 2
}

// Generic set bit function
func set(r *byte, bit byte) (error, uint16) {
	*r |= (0x01 << bit)
	return nil, 2
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

// 0xCB20
// Shift left B into carry, LSB = 0.
func slaB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sla(&r.B, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB21
// Shift left C into carry, LSB = 0.
func slaC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sla(&r.C, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB22
// Shift left D into carry, LSB = 0.
func slaD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sla(&r.D, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB23
// Shift left E into carry, LSB = 0.
func slaE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sla(&r.E, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB24
// Shift left H into carry, LSB = 0.
func slaH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sla(&r.H, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB25
// Shift left L into carry, LSB = 0.
func slaL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sla(&r.L, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB26
// Shift left value in address HL into carry, LSB = 0.
func slaHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	data := m.Read(r.HL())
	r.NF = false
	r.HF = false
	r.CF = data>>7 == 1
	data <<= 1
	r.ZF = data == 0
	m.Store(r.HL(), data)
	return nil, 2
}

// 0xCB27
// Shift left A into carry, LSB = 0.
func slaA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return sla(&r.A, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB30
// Swap nibbles of B.
func swapB(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return swap(&r.B, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB31
// Swap nibbles of C.
func swapC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return swap(&r.C, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB32
// Swap nibbles of D.
func swapD(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return swap(&r.D, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB33
// Swap nibbles of E.
func swapE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return swap(&r.E, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB34
// Swap nibbles of H.
func swapH(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return swap(&r.H, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB35
// Swap nibbles of L.
func swapL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return swap(&r.L, &r.ZF, &r.NF, &r.HF, &r.CF)
}

// 0xCB36
// Swap nibbles of memory pointed at by HL.
// Can't be generalized because of the use of memory.
func swapHL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	data := m.Read(r.HL())
	data = data<<4 + data>>4
	r.ZF = data == 0
	r.CF = false
	r.HF = false
	r.NF = false
	m.Store(r.HL(), data)
	return nil, 2
}

// 0xCB37
// Swap nibbles of A.
func swapA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return swap(&r.A, &r.ZF, &r.NF, &r.HF, &r.CF)
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

// 0xCB3F
// Shift right A into carry, MSB = 0.
func srlA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.NF = false
	r.HF = false
	r.CF = r.A&0x01 == 1
	r.A >>= 1
	r.ZF = r.A == 0
	return nil, 2
}

// 0xCB40
// Test bit 0 of B.
func test0B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(0, r.B, &r.ZF, &r.HF, &r.NF)
}

// 0xCB41
// Test bit 0 of C.
func test0C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(0, r.C, &r.ZF, &r.HF, &r.NF)
}

// 0xCB42
// Test bit 0 of D.
func test0D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(0, r.D, &r.ZF, &r.HF, &r.NF)
}

// 0xCB43
// Test bit 0 of E.
func test0E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(0, r.E, &r.ZF, &r.HF, &r.NF)
}

// 0xCB44
// Test bit 0 of H.
func test0H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(0, r.H, &r.ZF, &r.HF, &r.NF)
}

// 0xCB45
// Test bit 0 of L.
func test0L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(0, r.L, &r.ZF, &r.HF, &r.NF)
}

// 0xCB46
// Test bit 0 of value in memory address HL.
func test0HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return test(0, m.Read(r.HL()), &r.ZF, &r.HF, &r.NF)
}

// 0xCB47
// Test bit 0 of A.
func test0A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(0, r.A, &r.ZF, &r.HF, &r.NF)
}

// 0xCB48
// Test bit 1 of B.
func test1B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(1, r.B, &r.ZF, &r.HF, &r.NF)
}

// 0xCB41
// Test bit 1 of C.
func test1C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(1, r.C, &r.ZF, &r.HF, &r.NF)
}

// 0xCB42
// Test bit 1 of D.
func test1D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(1, r.D, &r.ZF, &r.HF, &r.NF)
}

// 0xCB43
// Test bit 1 of E.
func test1E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(1, r.E, &r.ZF, &r.HF, &r.NF)
}

// 0xCB44
// Test bit 1 of H.
func test1H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(1, r.H, &r.ZF, &r.HF, &r.NF)
}

// 0xCB45
// Test bit 1 of L.
func test1L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(1, r.L, &r.ZF, &r.HF, &r.NF)
}

// 0xCB46
// Test bit 1 of value in memory address HL.
func test1HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return test(1, m.Read(r.HL()), &r.ZF, &r.HF, &r.NF)
}

// 0xCB47
// Test bit 1 of A.
func test1A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(1, r.A, &r.ZF, &r.HF, &r.NF)
}

// 0xCB50
// Test bit 2 of B.
func test2B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(2, r.B, &r.ZF, &r.HF, &r.NF)
}

// 0xCB51
// Test bit 2 of C.
func test2C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(2, r.C, &r.ZF, &r.HF, &r.NF)
}

// 0xCB52
// Test bit 2 of D.
func test2D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(2, r.D, &r.ZF, &r.HF, &r.NF)
}

// 0xCB53
// Test bit 2 of E.
func test2E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(2, r.E, &r.ZF, &r.HF, &r.NF)
}

// 0xCB54
// Test bit 2 of H.
func test2H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(2, r.H, &r.ZF, &r.HF, &r.NF)
}

// 0xCB55
// Test bit 2 of L.
func test2L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(2, r.L, &r.ZF, &r.HF, &r.NF)
}

// 0xCB56
// Test bit 2 of value in memory address HL.
func test2HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return test(2, m.Read(r.HL()), &r.ZF, &r.HF, &r.NF)
}

// 0xCB57
// Test bit 2 of A.
func test2A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(2, r.A, &r.ZF, &r.HF, &r.NF)
}

// 0xCB58
// Test bit 3 of B.
func test3B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(3, r.B, &r.ZF, &r.HF, &r.NF)
}

// 0xCB59
// Test bit 3 of C.
func test3C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(3, r.C, &r.ZF, &r.HF, &r.NF)
}

// 0xCB5A
// Test bit 3 of D.
func test3D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(3, r.D, &r.ZF, &r.HF, &r.NF)
}

// 0xCB5B
// Test bit 3 of E.
func test3E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(3, r.E, &r.ZF, &r.HF, &r.NF)
}

// 0xCB5C
// Test bit 3 of H.
func test3H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(3, r.H, &r.ZF, &r.HF, &r.NF)
}

// 0xCB5D
// Test bit 3 of L.
func test3L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(3, r.L, &r.ZF, &r.HF, &r.NF)
}

// 0xCB5E
// Test bit 3 of value in memory address HL.
func test3HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return test(3, m.Read(r.HL()), &r.ZF, &r.HF, &r.NF)
}

// 0xCB5F
// Test bit 3 of A.
func test3A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(3, r.A, &r.ZF, &r.HF, &r.NF)
}

// 0xCB60
// Test bit 4 of B.
func test4B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(4, r.B, &r.ZF, &r.HF, &r.NF)
}

// 0xCB61
// Test bit 4 of C.
func test4C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(4, r.C, &r.ZF, &r.HF, &r.NF)
}

// 0xCB62
// Test bit 4 of D.
func test4D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(4, r.D, &r.ZF, &r.HF, &r.NF)
}

// 0xCB63
// Test bit 4 of E.
func test4E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(4, r.E, &r.ZF, &r.HF, &r.NF)
}

// 0xCB64
// Test bit 4 of H.
func test4H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(4, r.H, &r.ZF, &r.HF, &r.NF)
}

// 0xCB65
// Test bit 4 of L.
func test4L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(4, r.L, &r.ZF, &r.HF, &r.NF)
}

// 0xCB66
// Test bit 4 of value in memory address HL.
func test4HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return test(4, m.Read(r.HL()), &r.ZF, &r.HF, &r.NF)
}

// 0xCB67
// Test bit 4 of A.
func test4A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(4, r.A, &r.ZF, &r.HF, &r.NF)
}

// 0xCB68
// Test bit 5 of B.
func test5B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(5, r.B, &r.ZF, &r.HF, &r.NF)
}

// 0xCB69
// Test bit 5 of C.
func test5C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(5, r.C, &r.ZF, &r.HF, &r.NF)
}

// 0xCB6A
// Test bit 5 of D.
func test5D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(5, r.D, &r.ZF, &r.HF, &r.NF)
}

// 0xCB6B
// Test bit 5 of E.
func test5E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(5, r.E, &r.ZF, &r.HF, &r.NF)
}

// 0xCB6C
// Test bit 5 of H.
func test5H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(5, r.H, &r.ZF, &r.HF, &r.NF)
}

// 0xCB6D
// Test bit 5 of L.
func test5L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(5, r.L, &r.ZF, &r.HF, &r.NF)
}

// 0xCB6E
// Test bit 5 of value in memory address HL.
func test5HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return test(5, m.Read(r.HL()), &r.ZF, &r.HF, &r.NF)
}

// 0xCB6F
// Test bit 5 of A.
func test5A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(5, r.A, &r.ZF, &r.HF, &r.NF)
}

// 0xCB70
// Test bit 6 of B.
func test6B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(6, r.B, &r.ZF, &r.HF, &r.NF)
}

// 0xCB71
// Test bit 6 of C.
func test6C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(6, r.C, &r.ZF, &r.HF, &r.NF)
}

// 0xCB72
// Test bit 6 of D.
func test6D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(6, r.D, &r.ZF, &r.HF, &r.NF)
}

// 0xCB73
// Test bit 6 of E.
func test6E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(6, r.E, &r.ZF, &r.HF, &r.NF)
}

// 0xCB74
// Test bit 6 of H.
func test6H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(6, r.H, &r.ZF, &r.HF, &r.NF)
}

// 0xCB75
// Test bit 6 of L.
func test6L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(6, r.L, &r.ZF, &r.HF, &r.NF)
}

// 0xCB76
// Test bit 6 of value in memory address HL.
func test6HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return test(6, m.Read(r.HL()), &r.ZF, &r.HF, &r.NF)
}

// 0xCB77
// Test bit 6 of A.
func test6A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(6, r.A, &r.ZF, &r.HF, &r.NF)
}

// 0xCB78
// Test bit 7 of B.
func test7B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(7, r.B, &r.ZF, &r.HF, &r.NF)
}

// 0xCB79
// Test bit 7 of C.
func test7C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(7, r.C, &r.ZF, &r.HF, &r.NF)
}

// 0xCB7A
// Test bit 7 of D.
func test7D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(7, r.D, &r.ZF, &r.HF, &r.NF)
}

// 0xCB7B
// Test bit 7 of E.
func test7E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(7, r.E, &r.ZF, &r.HF, &r.NF)
}

// 0xCB7C
// Test bit 7 of H.
func test7H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(7, r.H, &r.ZF, &r.HF, &r.NF)
}

// 0xCB7D
// Test bit 7 of L.
func test7L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(7, r.L, &r.ZF, &r.HF, &r.NF)
}

// 0xCB7E
// Test bit 7 of value in memory address HL.
func test7HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	return test(3, m.Read(r.HL()), &r.ZF, &r.HF, &r.NF)
}

// 0xCB7F
// Test bit 7 of A.
func test7A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return test(7, r.A, &r.ZF, &r.HF, &r.NF)
}

// 0xCB80
// Reset bit 0 of B.
func res0B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.B, 0)
}

// 0xCB81
// Reset bit 0 of C.
func res0C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.C, 0)
}

// 0xCB82
// Reset bit 0 of D.
func res0D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.D, 0)
}

// 0xCB83
// Reset bit 0 of E.
func res0E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.E, 0)
}

// 0xCB84
// Reset bit 0 of H.
func res0H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.H, 0)
}

// 0xCB85
// Reset bit 0 of L.
func res0L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.L, 0)
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
	return reset(&r.A, 0)
}

// 0xCB88
// Reset bit 1 of B.
func res1B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.B, 1)
}

// 0xCB89
// Reset bit 1 of C.
func res1C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.C, 1)
}

// 0xCB8A
// Reset bit 1 of D.
func res1D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.D, 1)
}

// 0xCB8B
// Reset bit 1 of E.
func res1E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.E, 1)
}

// 0xCB8C
// Reset bit 1 of H.
func res1H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.H, 1)
}

// 0xCB8D
// Reset bit 1 of L.
func res1L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.L, 1)
}

// 0xCB8E
// Reset bit 1 of value in memory address HL.
func res1HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())&0xFD)
	return nil, 2
}

// 0xCB8F
// Reset bit 1 of A.
func res1A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.A, 1)
}

// 0xCB90
// Reset bit 2 of B.
func res2B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.B, 2)
}

// 0xCB91
// Reset bit 2 of C.
func res2C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.C, 2)
}

// 0xCB92
// Reset bit 2 of D.
func res2D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.D, 2)
}

// 0xCB93
// Reset bit 2 of E.
func res2E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.E, 2)
}

// 0xCB94
// Reset bit 2 of H.
func res2H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.H, 2)
}

// 0xCB95
// Reset bit 2 of L.
func res2L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.L, 2)
}

// 0xCB96
// Reset bit 2 of value in memory address HL.
func res2HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())&0xFB)
	return nil, 2
}

// 0xCB97
// Reset bit 2 of A.
func res2A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.A, 2)
}

// 0xCB98
// Reset bit 3 of B.
func res3B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.B, 3)
}

// 0xCB99
// Reset bit 3 of C.
func res3C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.C, 3)
}

// 0xCB9A
// Reset bit 3 of D.
func res3D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.D, 3)
}

// 0xCB9B
// Reset bit 3 of E.
func res3E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.E, 3)
}

// 0xCB9C
// Reset bit 3 of H.
func res3H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.H, 3)
}

// 0xCB9D
// Reset bit 3 of L.
func res3L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.L, 3)
}

// 0xCB9E
// Reset bit 3 of value in memory address HL.
func res3HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())&0xF7)
	return nil, 2
}

// 0xCB9F
// Reset bit 3 of A.
func res3A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.A, 3)
}

// 0xCBA0
// Reset bit 4 of B.
func res4B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.B, 4)
}

// 0xCBA1
// Reset bit 4 of C.
func res4C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.C, 4)
}

// 0xCBA2
// Reset bit 4 of D.
func res4D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.D, 4)
}

// 0xCBA3
// Reset bit 4 of E.
func res4E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.E, 4)
}

// 0xCBA4
// Reset bit 4 of H.
func res4H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.H, 4)
}

// 0xCBA5
// Reset bit 4 of L.
func res4L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.L, 4)
}

// 0xCBA6
// Reset bit 4 of value in memory address HL.
func res4HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())&0xEF)
	return nil, 2
}

// 0xCBA7
// Reset bit 4 of A.
func res4A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.A, 4)
}

// 0xCBA8
// Reset bit 5 of B.
func res5B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.B, 5)
}

// 0xCBA9
// Reset bit 5 of C.
func res5C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.C, 5)
}

// 0xCBAA
// Reset bit 5 of D.
func res5D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.D, 5)
}

// 0xCBAB
// Reset bit 5 of E.
func res5E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.E, 5)
}

// 0xCBAC
// Reset bit 5 of H.
func res5H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.H, 5)
}

// 0xCBAD
// Reset bit 5 of L.
func res5L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.L, 5)
}

// 0xCBAE
// Reset bit 5 of value in memory address HL.
func res5HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())&0xDF)
	return nil, 2
}

// 0xCBAF
// Reset bit 5 of A.
func res5A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.A, 5)
}

// 0xCBB0
// Reset bit 6 of B.
func res6B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.B, 6)
}

// 0xCBB1
// Reset bit 6 of C.
func res6C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.C, 6)
}

// 0xCBB2
// Reset bit 6 of D.
func res6D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.D, 6)
}

// 0xCBB3
// Reset bit 6 of E.
func res6E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.E, 6)
}

// 0xCBB4
// Reset bit 6 of H.
func res6H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.H, 6)
}

// 0xCBB5
// Reset bit 6 of L.
func res6L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.L, 6)
}

// 0xCBB6
// Reset bit 6 of value in memory address HL.
func res6HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())&0xBF)
	return nil, 2
}

// 0xCBB7
// Reset bit 6 of A.
func res6A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.A, 6)
}

// 0xCBB8
// Reset bit 7 of B.
func res7B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.B, 7)
}

// 0xCBB9
// Reset bit 7 of C.
func res7C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.C, 7)
}

// 0xCBBA
// Reset bit 7 of D.
func res7D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.D, 7)
}

// 0xCBBB
// Reset bit 7 of E.
func res7E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.E, 7)
}

// 0xCBBC
// Reset bit 7 of H.
func res7H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.H, 7)
}

// 0xCBBD
// Reset bit 7of L.
func res7L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.L, 7)
}

// 0xCBBE
// Reset bit 7 of value in memory address HL.
func res7HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())&0x7F)
	return nil, 2
}

// 0xCBBF
// Reset bit 7 of A.
func res7A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return reset(&r.A, 7)
}

// 0xCBC0
// Set bit 0 of B.
func set0B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.B, 0)
}

// 0xCBC1
// Set bit 0 of C.
func set0C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.C, 0)
}

// 0xCBC2
// Set bit 0 of D.
func set0D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.D, 0)
}

// 0xCBC3
// Set bit 0 of E.
func set0E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.E, 0)
}

// 0xCBC4
// Set bit 0 of H.
func set0H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.H, 0)
}

// 0xCBC5
// Set bit 0 of L.
func set0L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.L, 0)
}

// 0xCBC6
// Set bit 0 of value in memory address HL.
func set0HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())|0x01)
	return nil, 2
}

// 0xCBC7
// Set bit 0 of A.
func set0A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.A, 0)
}

// 0xCBC8
// Set bit 1 of B.
func set1B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.B, 1)
}

// 0xCBC9
// Set bit 1 of C.
func set1C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.C, 1)
}

// 0xCBCA
// Set bit 1 of D.
func set1D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.D, 1)
}

// 0xCBCB
// Set bit 1 of E.
func set1E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.E, 1)
}

// 0xCBCC
// Set bit 1 of H.
func set1H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.H, 1)
}

// 0xCBCD
// Set bit 1 of L.
func set1L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.L, 1)
}

// 0xCBCE
// Set bit 1 of value in memory address HL.
func set1HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())|0x02)
	return nil, 2
}

// 0xCBCF
// Set bit 1 of A.
func set1A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.A, 1)
}

// 0xCBD0
// Set bit 2 of B.
func set2B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.B, 2)
}

// 0xCBD1
// Set bit 2 of C.
func set2C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.C, 2)
}

// 0xCBD2
// Set bit 2 of D.
func set2D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.D, 2)
}

// 0xCBD3
// Set bit 2 of E.
func set2E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.E, 2)
}

// 0xCBD4
// Set bit 2 of H.
func set2H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.H, 2)
}

// 0xCBD5
// Set bit 2 of L.
func set2L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.L, 2)
}

// 0xCBD6
// Set bit 2 of value in memory address HL.
func set2HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())|0x04)
	return nil, 2
}

// 0xCBD7
// Set bit 2 of A.
func set2A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.A, 2)
}

// 0xCBD8
// Set bit 3 of B.
func set3B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.B, 3)
}

// 0xCBD9
// Set bit 3 of C.
func set3C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.C, 3)
}

// 0xCBDA
// Set bit 3 of D.
func set3D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.D, 3)
}

// 0xCBDB
// Set bit 3 of E.
func set3E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.E, 3)
}

// 0xCBDC
// Set bit 3 of H.
func set3H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.H, 3)
}

// 0xCBDD
// Set bit 3 of L.
func set3L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.L, 3)
}

// 0xCBDE
// Set bit 3 of value in memory address HL.
func set3HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())|0x08)
	return nil, 2
}

// 0xCBDF
// Set bit 3 of A.
func set3A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.A, 3)
}

// 0xCBE0
// Set bit 4 of B.
func set4B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.B, 4)
}

// 0xCBE1
// Set bit 4 of C.
func set4C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.C, 4)
}

// 0xCBE2
// Set bit 4 of D.
func set4D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.D, 4)
}

// 0xCBE3
// Set bit 4 of E.
func set4E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.E, 4)
}

// 0xCBE4
// Set bit 4 of H.
func set4H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.H, 4)
}

// 0xCBE5
// Set bit 4 of L.
func set4L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.L, 4)
}

// 0xCBE6
// Set bit 4 of value in memory address HL.
func set4HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())|0x10)
	return nil, 2
}

// 0xCBE7
// Set bit 4 of A.
func set4A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.A, 4)
}

// 0xCBE8
// Set bit 5 of B.
func set5B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.B, 5)
}

// 0xCBE9
// Set bit 5 of C.
func set5C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.C, 5)
}

// 0xCBEA
// Set bit 5 of D.
func set5D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.D, 5)
}

// 0xCBEB
// Set bit 5 of E.
func set5E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.E, 5)
}

// 0xCBEC
// Set bit 5 of H.
func set5H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.H, 5)
}

// 0xCBED
// Set bit 5 of L.
func set5L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.L, 5)
}

// 0xCBEE
// Set bit 5 of value in memory address HL.
func set5HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())|0x20)
	return nil, 2
}

// 0xCBEF
// Set bit 5 of A.
func set5A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.A, 5)
}

// 0xCBF0
// Set bit 6 of B.
func set6B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.B, 6)
}

// 0xCBF1
// Set bit 6 of C.
func set6C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.C, 6)
}

// 0xCBF2
// Set bit 6 of D.
func set6D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.D, 6)
}

// 0xCBF3
// Set bit 6 of E.
func set6E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.E, 6)
}

// 0xCBF4
// Set bit 6 of H.
func set6H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.H, 6)
}

// 0xCBF5
// Set bit 6 of L.
func set6L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.L, 6)
}

// 0xCBF6
// Set bit 6 of value in memory address HL.
func set6HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())|0x40)
	return nil, 2
}

// 0xCBF7
// Set bit 6 of A.
func set6A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.A, 6)
}

// 0xCBF8
// Set bit 7 of B.
func set7B(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.B, 7)
}

// 0xCBF9
// Set bit 7 of C.
func set7C(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.C, 7)
}

// 0xCBFA
// Set bit 7 of D.
func set7D(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.D, 7)
}

// 0xCBFB
// Set bit 7 of E.
func set7E(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.E, 7)
}

// 0xCBFC
// Set bit 7 of H.
func set7H(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.H, 7)
}

// 0xCBFD
// Set bit 7 of L.
func set7L(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.L, 7)
}

// 0xCBFE
// Set bit 7 of value in memory address HL.
func set7HL(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	m.Store(r.HL(), m.Read(r.HL())|0x80)
	return nil, 2
}

// 0xCBFF
// Set bit 7 of A.
func set7A(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	return set(&r.A, 7)
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
	slaB, // 0x20
	slaC,
	slaD,
	slaE,
	slaH,
	slaL,
	slaHL,
	slaA,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	swapB, // 0x30
	swapC,
	swapD,
	swapE,
	swapH,
	swapL,
	swapHL,
	swapA,
	srlB,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	srlA,
	test0B, // 0x40
	test0C,
	test0D,
	test0E,
	test0H,
	test0L,
	test0HL,
	test0A,
	test1B,
	test1C,
	test1D,
	test1E,
	test1H,
	test1L,
	test1HL,
	test1A,
	test2B, // 0x50
	test2C,
	test2D,
	test2E,
	test2H,
	test2L,
	test2HL,
	test2A,
	test3B,
	test3C,
	test3D,
	test3E,
	test3H,
	test3L,
	test3HL,
	test3A,
	test4B, // 0x60
	test4C,
	test4D,
	test4E,
	test4H,
	test4L,
	test4HL,
	test4A,
	test5B,
	test5C,
	test5D,
	test5E,
	test5H,
	test5L,
	test5HL,
	test5A,
	test6B, // 0x70
	test6C,
	test6D,
	test6E,
	test6H,
	test6L,
	test6HL,
	test6A,
	test7B,
	test7C,
	test7D,
	test7E,
	test7H,
	test7L,
	test7HL,
	test7A,
	res0B, // 0x80
	res0C,
	res0D,
	res0E,
	res0H,
	res0L,
	res0HL,
	res0A,
	res1B,
	res1C,
	res1D,
	res1E,
	res1H,
	res1L,
	res1HL,
	res1A,
	res2B, // 0x90
	res2C,
	res2D,
	res2E,
	res2H,
	res2L,
	res2HL,
	res2A,
	res3B,
	res3C,
	res3D,
	res3E,
	res3H,
	res3L,
	res3HL,
	res3A,
	res4B, // 0xA0
	res4C,
	res4D,
	res4E,
	res4H,
	res4L,
	res4HL,
	res4A,
	res5B,
	res5C,
	res5D,
	res5E,
	res5H,
	res5L,
	res5HL,
	res5A,
	res6B, // 0xB0
	res6C,
	res6D,
	res6E,
	res6H,
	res6L,
	res6HL,
	res6A,
	res7B,
	res7C,
	res7D,
	res7E,
	res7H,
	res7L,
	res7HL,
	res7A,
	set0B, // 0xC0
	set0C,
	set0D,
	set0E,
	set0H,
	set0L,
	set0HL,
	set0A,
	set1B,
	set1C,
	set1D,
	set1E,
	set1H,
	set1L,
	set1HL,
	set1A,
	set2B, // 0xD0
	set2C,
	set2D,
	set2E,
	set2H,
	set2L,
	set2HL,
	set2A,
	set3B,
	set3C,
	set3D,
	set3E,
	set3H,
	set3L,
	set3HL,
	set3A,
	set4B, // 0xE0
	set4C,
	set4D,
	set4E,
	set4H,
	set4L,
	set4HL,
	set4A,
	set5B,
	set5C,
	set5D,
	set5E,
	set5H,
	set5L,
	set5HL,
	set5A,
	set6B, // 0xF0
	set6C,
	set6D,
	set6E,
	set6H,
	set6L,
	set6HL,
	set6A,
	set7B,
	set7C,
	set7D,
	set7E,
	set7H,
	set7L,
	set7HL,
	set7A,
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
	// Handle changes in flags.
	// This shouldn't be done this way and all flags info should be in the F register.
	// I didn't foresee this mess. I'll fix it some time later.
	if !flagsRecovered {
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
	} else {
		r.ZF = r.F&0x80 != 0
		r.NF = r.F&0x40 != 0
		r.HF = r.F&0x20 != 0
		r.CF = r.F&0x10 != 0
		flagsRecovered = false
	}
	return err, jump, cycles
}

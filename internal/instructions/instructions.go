package instructions

import (
	"fmt"
	"go-boy/internal/memory"
	"go-boy/internal/registers"
	"reflect"
	"runtime"
)

func unimplemented(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	return fmt.Errorf("unimplemented instruction reached at PC=%X: %X", r.PC, args[0]), 0
}

// 0x00
// Doesn't do anything.
func nop(_ *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
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

// 0x0A
// Loads A from address pointed to by BC
func ldABC(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	r.A = m.Read(r.BC())
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

// 0x19
// Adds DE to HL. Stores the result in HL.
func addDE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
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
// Adds a specific amount to PC if Z flag is unset
func jrNZn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if !r.ZF {
		r.PC = r.PC + uint16(args[1])
		return nil, 0
	} else {
		return nil, 1
	}
}

// 0x21
// Loads a 16 bit int into HL
func ldHLnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.L = args[1]
	r.H = args[2]
	return nil, 3
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

// 0x2B
// Decrements the value in HL.
func decHL(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	hl := r.HL()
	hl--
	r.H = byte(hl >> 8)
	r.L = byte(hl)
	return nil, 1
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
	hl := r.HL()
	m.Store(hl, r.A)
	_, _ = decHL(r, nil, nil)
	return nil, 1
}

// 0x41
// Copies C to B.
func ldBC(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.B = r.C
	return nil, 1
}

// 0x77
// Stores the contents of A into the memory address HL.
func ldHLA(r *registers.Registers, m *memory.Memory, _ []byte) (error, uint16) {
	hl := r.HL()
	m.Store(hl, r.A)
	return nil, 1
}

// 0x7B
// Stores the contents of E into A.
func ldAE(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.A = r.E
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

// 0xBF
// Compares A with A. Basically sets some flags.
func cpA(r *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16) {
	r.ZF = true
	r.NF = true
	r.HF = false
	r.CF = false
	return nil, 1
}

// 0xC3
// Sets PC to the specified address in the arguments.
func jpnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	r.PC = uint16(args[1]) + uint16(args[2])<<8
	return nil, 0 // Return 0 because it's a jump
}

// 0xD2
// Sets PC as specified in the arguments if C flag is unset
func jpNCnn(r *registers.Registers, _ *memory.Memory, args []byte) (error, uint16) {
	if !r.CF {
		r.PC = uint16(args[1]) + uint16(args[2])<<8
		return nil, 0
	}
	return nil, 1
}

var InstructionTable = [256]func(_ *registers.Registers, _ *memory.Memory, _ []byte) (error, uint16){
	nop, // 0x00
	unimplemented,
	unimplemented,
	incBC,
	unimplemented,
	decB,
	ldBn,
	rlcA,
	ldnnSP,
	unimplemented,
	ldABC,
	unimplemented,
	incC,
	decC,
	ldCn,
	unimplemented,
	unimplemented, // 0x10
	ldDEnn,
	unimplemented,
	unimplemented,
	incD,
	decD,
	ldDn,
	unimplemented,
	unimplemented,
	addDE,
	unimplemented,
	unimplemented,
	unimplemented,
	decE,
	ldEn,
	rrA,
	jrNZn, // 0x20
	ldHLnn,
	unimplemented,
	incHL,
	unimplemented,
	decH,
	unimplemented,
	unimplemented,
	unimplemented,
	addHL,
	unimplemented,
	decHL,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x30
	ldSPnn,
	lddHLA,
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
	unimplemented, // 0x40
	ldBC,
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
	unimplemented, // 0x50
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
	unimplemented, // 0x60
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
	unimplemented, // 0x70
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	ldHLA,
	unimplemented,
	unimplemented,
	unimplemented,
	ldAE,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented, // 0x80
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
	unimplemented,
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
	xorA,
	orB, // 0xB0
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
	cpA,
	unimplemented, // 0xC0
	unimplemented,
	unimplemented,
	jpnn,
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
	jpNCnn,
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
	unimplemented}

func Execute(r *registers.Registers, m *memory.Memory, instructionArray []byte) (error, uint16) {
	// The instruction is the first byte in the array
	opCode := instructionArray[0]
	// Look up the instructions table and obtain the function that executes the instruction
	operation := InstructionTable[opCode]
	fmt.Println(runtime.FuncForPC(reflect.ValueOf(operation).Pointer()).Name())
	// Execute the operation
	return operation(r, m, instructionArray)
}

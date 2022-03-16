package utils

import (
	"go-boy/internal/memory"
	"go-boy/internal/registers"
)

func PushStack(r *registers.Registers, m *memory.Memory, data byte) {
	m.Store(r.SP, data)
	r.SP--
}

func PushStackShort(r *registers.Registers, m *memory.Memory, data uint16) {
	m.Store(r.SP, byte(data))
	r.SP--
	m.Store(r.SP, byte(data>>8))
	r.SP--
}

func PopStack(r *registers.Registers, m *memory.Memory) byte {
	r.SP++
	data := m.Read(r.SP)
	return data
}

func PopStackShort(r *registers.Registers, m *memory.Memory) uint16 {
	data := uint16(m.Read(r.SP+1))<<8 + uint16(m.Read(r.SP+2))
	r.SP += 2
	return data
}

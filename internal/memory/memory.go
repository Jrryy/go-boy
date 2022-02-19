package memory

import (
	"fmt"
	"os"
)

// Memory represents the different parts of the GB memory.
// It's been split in different parts only to help understand it better.
type Memory struct {
	EchoRAM   []byte
	RAM       []byte
	Cartridge []byte
}

func GetInitializedMemory(gameFile *os.File) *Memory {
	m := new(Memory)
	m.RAM = make([]byte, 0x2000)
	m.EchoRAM = make([]byte, 0x2000)
	m.Cartridge = make([]byte, 0x8000)
	_, err := gameFile.ReadAt(m.Cartridge, 0)
	if err != nil {
		panic(err)
	}
	return m
}

func (m *Memory) getMemoryPart(address uint16) (*[]byte, uint16) {
	if address < 0x8000 {
		return &m.Cartridge, 0
	}
	if address >= 0xC000 && address < 0xE000 {
		return &m.RAM, address - 0xC000
	}
	return nil, 0
}

// Store stores a byte in an address of the memory.
func (m *Memory) Store(address uint16, n byte) {
	memoryPart, offset := m.getMemoryPart(address)
	if memoryPart == nil {
		panic(fmt.Sprintf("Memory part not implemented: %X", address))
	}
	(*memoryPart)[offset] = n
}

func (m *Memory) Read(address uint16) byte {
	memoryPart, offset := m.getMemoryPart(address)
	if memoryPart == nil {
		panic(fmt.Sprintf("Memory part not implemented: %X", address))
	}
	return (*memoryPart)[offset]
}

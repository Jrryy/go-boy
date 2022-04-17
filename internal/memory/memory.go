package memory

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// These represent the three different states of input handling:
// P14 = arrows, P15 = Start, Select, A, B, P00 when input unavailable.
const (
	P00 = iota
	P14
	P15
)

// Memory represents the different parts of the GB memory.
// It's been split in different parts only to help understand it better.
type Memory struct {
	InputMode   int
	IME         bool
	IMEReqType  bool
	IMESteps    byte
	IER         []byte // FFFF
	InternalRAM []byte // FF80 - FFFE
	UnusableIO2 []byte // FF4C - FF7F
	IOPorts     []byte // FF00 - FF4B
	UnusableIO1 []byte // FEA0 - FEFF
	OAM         []byte // FE00 - FE9F
	EchoRAM     []byte // E000 - FDFF
	RAM         []byte // C000 - DFFF
	SRAM        []byte // A000 - BFFF
	VRAM        []byte // 8000 - 9FFF
	Cartridge   []byte // 0000 - 7FFF
}

func GetInitializedMemory(gameFile *os.File) *Memory {
	m := new(Memory)
	m.IER = make([]byte, 1)
	m.InternalRAM = make([]byte, 0x7F)
	m.UnusableIO2 = make([]byte, 0x34)
	m.IOPorts = make([]byte, 0x4C)
	m.UnusableIO1 = make([]byte, 0x60)
	m.OAM = make([]byte, 0xA0)
	m.EchoRAM = make([]byte, 0x1E00)
	m.RAM = make([]byte, 0x2000)
	m.SRAM = make([]byte, 0x2000)
	m.VRAM = make([]byte, 0x2000)
	m.Cartridge = make([]byte, 0x8000)
	_, err := gameFile.ReadAt(m.Cartridge, 0)
	if err != nil {
		panic(err)
	}
	m.Store(0xFF05, 0x00)
	m.Store(0xFF06, 0x00)
	m.Store(0xFF07, 0x00)
	m.Store(0xFF10, 0x80)
	m.Store(0xFF11, 0xBF)
	m.Store(0xFF12, 0xF3)
	m.Store(0xFF14, 0xBF)
	m.Store(0xFF16, 0x3F)
	m.Store(0xFF17, 0x00)
	m.Store(0xFF19, 0xBF)
	m.Store(0xFF1A, 0x7F)
	m.Store(0xFF1B, 0xFF)
	m.Store(0xFF1C, 0x9F)
	m.Store(0xFF1E, 0xBF)
	m.Store(0xFF20, 0xFF)
	m.Store(0xFF21, 0x00)
	m.Store(0xFF22, 0x00)
	m.Store(0xFF23, 0xBF)
	m.Store(0xFF24, 0x77)
	m.Store(0xFF25, 0xF3)
	m.Store(0xFF26, 0xF1)
	m.Store(0xFF40, 0x91)
	m.Store(0xFF42, 0x00)
	m.Store(0xFF43, 0x00)
	m.Store(0xFF45, 0x00)
	m.Store(0xFF47, 0xFC)
	m.Store(0xFF48, 0xFF)
	m.Store(0xFF49, 0xFF)
	m.Store(0xFF4A, 0x00)
	m.Store(0xFF4B, 0x00)
	m.Store(0xFFFF, 0x00)
	return m
}

func (m *Memory) getMemoryPart(address uint16) (*[]byte, uint16) {
	if address < 0x8000 {
		return &m.Cartridge, address
	}
	if address < 0xA000 {
		return &m.VRAM, address - 0x8000
	}
	if address < 0xC000 {
		return &m.SRAM, address - 0xA000
	}
	if address < 0xE000 {
		return &m.RAM, address - 0xC000
	}
	if address < 0xFE00 {
		// Echo RAM. Basically the same as RAM so we return it and avoid redundance
		return &m.RAM, address - 0xE000
	}
	if address < 0xFEA0 {
		return &m.OAM, address - 0xFE00
	}
	if address < 0xFF00 {
		return &m.UnusableIO1, address - 0xFEA0
	}
	if address < 0xFF4C {
		return &m.IOPorts, address - 0xFF00
	}
	if address < 0xFF80 {
		return &m.UnusableIO2, address - 0xFF4C
	}
	if address < 0xFFFF {
		return &m.InternalRAM, address - 0xFF80
	}
	if address == 0xFFFF {
		return &m.IER, 0
	}
	return nil, 0
}

// Store stores a byte in an address of the memory.
func (m *Memory) Store(address uint16, n byte) {
	if address < 0x8000 {
		// Bank switching would go here, but we ain't doing this yet.
	} else if address == 0xFF00 {
		// Handling input.
		if n == 0x10 {
			m.InputMode = P14
		} else if n == 0x20 {
			m.InputMode = P15
		} else {
			m.InputMode = P00
		}
	} else {
		memoryPart, offset := m.getMemoryPart(address)
		if memoryPart == nil {
			panic(fmt.Sprintf("Memory part not implemented: %X", address))
		}
		(*memoryPart)[offset] = n
	}
}

func (m *Memory) Read(address uint16) byte {
	if address == 0xFF00 {
		return m.getUserInput()
	} else {
		memoryPart, offset := m.getMemoryPart(address)
		if memoryPart == nil {
			panic(fmt.Sprintf("Memory part not implemented: %X", address))
		}
		return (*memoryPart)[offset]
	}
}

func (m *Memory) ReadInstruction(address uint16) []byte {
	return []byte{m.Read(address), m.Read(address + 1), m.Read(address + 2)}
}

// Processes the user's input
func (m *Memory) getUserInput() byte {
	capturedInput := byte(0x00)
	if m.InputMode == P14 {
		if ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonRightBottom) {
			// A
			capturedInput |= 0x01
		}
		if ebiten.IsKeyPressed(ebiten.KeyX) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonRightRight) {
			// B
			capturedInput |= 0x02
		}
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonCenterLeft) {
			// Select
			capturedInput |= 0x04
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonCenterRight) {
			// Start
			capturedInput |= 0x08
		}
		return ^capturedInput
	} else if m.InputMode == P15 {
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonLeftRight) {
			// Right
			capturedInput |= 0x01
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonLeftLeft) {
			// Left
			capturedInput |= 0x02
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonLeftTop) {
			// Up
			capturedInput |= 0x04
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonLeftBottom) {
			// Down
			capturedInput |= 0x08
		}
		return ^capturedInput
	} else {
		return 0xCF
	}
}

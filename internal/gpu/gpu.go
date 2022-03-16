package gpu

import "go-boy/internal/memory"

type GPU struct {
	ticks           int
	scanLine        uint8
	Mode            uint8
	VBlankInterrupt bool
}

const (
	HBLANK = iota
	VBLANK
	OAM
	VRAM
)

func InitGPU() *GPU {
	newGPU := new(GPU)
	newGPU.ticks = 0
	newGPU.scanLine = 0
	newGPU.Mode = HBLANK
	return newGPU
}

func (gpu *GPU) Step(cycles int, m *memory.Memory) {
	gpu.ticks += cycles
	switch gpu.Mode {
	case HBLANK:
		if gpu.ticks >= 204 {
			gpu.scanLine++

			if gpu.scanLine == 143 {
				gpu.Mode = VBLANK
				m.Store(0xFF41, (m.Read(0xFF41)&0xFC)|VBLANK)
				if m.IER[0]&0x01 == 1 {
					m.Store(0xFF0F, m.Read(0xFF0F)|0x01)
				}
			} else {
				gpu.Mode = OAM
				m.Store(0xFF41, (m.Read(0xFF41)&0xFC)|OAM)
			}
			gpu.ticks -= 204
		}
	case VBLANK:
		if gpu.ticks >= 456 {
			gpu.ticks -= 456
			gpu.scanLine++

			if gpu.scanLine > 153 {
				gpu.Mode = OAM
				gpu.scanLine = 0
				m.Store(0xFF41, (m.Read(0xFF41)&0xFC)|OAM)
			}
		}
	case OAM:
		if gpu.ticks >= 80 {
			gpu.ticks -= 80
			gpu.Mode = VRAM
			m.Store(0xFF41, (m.Read(0xFF41)&0xFC)|VRAM)
		}
	case VRAM:
		if gpu.ticks >= 172 {
			gpu.ticks -= 172
			gpu.Mode = HBLANK
			m.Store(0xFF41, (m.Read(0xFF41)&0xFC)|HBLANK)
		}
	}
	m.Store(0xFF44, gpu.scanLine)
	if m.Read(0xFF44) == m.Read(0xFF45) {
		m.Store(0xFF41, (m.Read(0xFF41)&0xBF)|0x40)
	} else {
		m.Store(0xFF41, m.Read(0xFF41)&0xBF)
	}
}

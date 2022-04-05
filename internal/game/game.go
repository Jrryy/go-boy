package game

import (
	"fmt"
	"go-boy/internal/gpu"
	"go-boy/internal/instructions"
	"go-boy/internal/memory"
	"go-boy/internal/registers"
	"go-boy/internal/utils"
	"image/color"
	"os"
	"reflect"
	"runtime"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var frequency = 4194304
var cyclesPerFrame = frequency / 60
var cyclesPerDivUpdate = cyclesPerFrame / (16384 / 60)
var cyclesPerTimaUpdate = [4]int{
	cyclesPerFrame / (4096 / 60),
	cyclesPerFrame / (262144 / 60),
	cyclesPerFrame / (65536 / 60),
	cyclesPerFrame / (16384 / 60),
}
var gameFont font.Face

var color00 = color.RGBA{0xE0, 0xF8, 0xCF, 0xFF}
var color01 = color.RGBA{0x86, 0xC0, 0x6C, 0xFF}
var color10 = color.RGBA{0x30, 0x68, 0x50, 0xFF}
var color11 = color.RGBA{0x07, 0x18, 0x21, 0xFF}
var colors = [4]color.RGBA{color00, color01, color10, color11}

type Game struct {
	R     *registers.Registers
	M     *memory.Memory
	GPU   *gpu.GPU
	Pause bool
	Debug bool
}

func init() {
	fontFileBytes, err := os.ReadFile("assets/Hack-Regular.ttf")
	if err != nil {
		panic(err)
	}
	tt, _ := opentype.Parse(fontFileBytes)
	gameFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  36,
	})
}

func (g *Game) Update() error {
	currentCycles := 0
	divCycles := 0
	timaCycles := 0
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Debug = true
	}
	for currentCycles < cyclesPerFrame {
		// Joypad stuff. For now let's just consider the use is not pressing a single button.
		g.M.Store(0xFF00, 0xCF)
		var err error
		var bytes uint16
		var cycles int
		if !g.R.Halted {
			// Read always 3 bytes: op code and 2 possible arguments
			instructionArray := g.M.ReadInstruction(g.R.PC)
			if g.Debug {
				fmt.Println(runtime.FuncForPC(reflect.ValueOf(instructions.InstructionTable[instructionArray[0]]).Pointer()).Name())
				fmt.Println(g.R)
			}
			// fmt.Printf("%X %X %X\n", instructionArray[0], instructionArray[1], instructionArray[2])
			err, bytes, cycles = instructions.Execute(g.R, g.M, instructionArray)
			if err != nil {
				panic(err)
			}
		} else {
			bytes = 0
			cycles = 1
		}
		// Add cycles executed to the current cycles of the frame
		currentCycles += cycles
		// Update DIV
		divCycles += cycles
		if divCycles >= cyclesPerDivUpdate {
			divCycles = 0
			g.M.Store(0xFF04, g.M.Read(0xFF04)+1)
		}
		// Update TIMA
		tac := g.M.Read(0xFF07)
		if tac&0x04 == 0x04 {
			timaCycles += cycles
			if timaCycles >= cyclesPerTimaUpdate[tac&0x03] {
				timaCycles = 0
				currentTIMA := uint16(g.M.Read(0xFF05)) + 1
				if currentTIMA > 0xFF {
					if g.M.IER[0]&0x04 == 0x04 {
						g.M.Store(0xFF0F, g.M.Read(0xFF0F)|0x04)
					}
					g.M.Store(0xFF05, g.M.Read(0xFF06))
				} else {
					g.M.Store(0xFF05, byte(currentTIMA))
				}
			}
		}
		// Augment the PC as much as the amount of bytes the instruction has used
		g.R.PC += bytes
		// Check if there are pending interruptions
		g.CheckInterruptRequests()
		// Run a gpu step
		g.GPU.Step(cycles, g.M)
		// Run an interruptions step
		g.InterruptStep()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.Gray{0x77})

	lcdc := g.M.Read(0xFF40)
	// scx := g.M.Read(0xFF42)
	// scy := g.M.Read(0xFF43)

	g.transferOAM()

	if lcdc&0x01 == 1 {
		g.drawBackground(screen, lcdc)

		if lcdc&0x02 != 0 {
			g.drawSprites(screen, lcdc)
		}
	}

	g.debugMemory(screen)
}

func (g *Game) transferOAM() {
	address := uint16(g.M.Read(0xFF46)) << 8

	for i := 0; i < 0xA0; i++ {
		g.M.OAM[i] = g.M.Read(address + uint16(i))
	}
}

func (g *Game) drawBackground(screen *ebiten.Image, lcdc byte) {
	var tileMapAddr, tileDataAddr uint16
	var signed bool

	bgp := g.M.Read(0xFF47)
	bgpColors := [4]color.RGBA{
		colors[bgp&0x03],
		colors[bgp&0x0C>>2],
		colors[bgp&0x30>>4],
		colors[bgp>>6],
	}

	if lcdc&0x08 == 0 {
		tileMapAddr = 0x9800
	} else {
		tileMapAddr = 0x9C00
	}

	if lcdc&0x10 == 0 {
		signed = true
		tileDataAddr = 0x9000
	} else {
		signed = false
		tileDataAddr = 0x8000
	}

	for y := 0; y < 18; y++ {
		for x := 0; x < 20; x++ {

			tileNumber := g.M.Read(tileMapAddr)
			tileAddr := tileDataAddr + uint16(tileNumber)*16
			if signed && tileNumber >= 0x7F {
				tileAddr = (0x0800 + uint16(tileNumber)) << 4
			}
			// fmt.Printf("%04X ", tileAddr)

			for i := 0; i < 8; i++ {
				tileLineUp := g.M.Read(tileAddr)
				tileAddr++
				tileLineDown := g.M.Read(tileAddr)
				tileAddr++

				binaryTileLineUp := fmt.Sprintf("%08b", tileLineUp)
				binaryTileLineDown := fmt.Sprintf("%08b", tileLineDown)

				for j := 0; j < 8; j++ {
					pair := string(binaryTileLineDown[j]) + string(binaryTileLineUp[j])
					bgColor, _ := strconv.ParseInt(pair, 2, 8)
					pixelColor := bgpColors[bgColor]
					screen.Set(200+8*x+j, 8*y+i, pixelColor)
				}
			}
			tileMapAddr++
		}
		tileMapAddr += 12
	}
}

func (g *Game) drawWindow(screen *ebiten.Image, lcdc byte) {

}

func (g *Game) drawSprites(screen *ebiten.Image, lcdc byte) {
	var height int
	obp0 := g.M.Read(0xFF48)
	obp1 := g.M.Read(0xFF49)
	obps := [2][4]color.RGBA{
		{
			colors[obp0&0x03],
			colors[obp0&0x0C>>2],
			colors[obp0&0x30>>4],
			colors[obp0>>6],
		},
		{
			colors[obp1&0x03],
			colors[obp1&0x0C>>2],
			colors[obp1&0x30>>4],
			colors[obp1>>6],
		},
	}

	if lcdc&0x04 == 0 {
		height = 1
	} else {
		height = 2
	}

	var finalSpriteAddr uint16 = 0xFE9C

	for nSprite := uint16(0); nSprite < 40; nSprite++ {
		yPosition := int(g.M.Read(finalSpriteAddr - 4*nSprite))
		xPosition := int(g.M.Read(finalSpriteAddr - 4*nSprite + 1))
		patternNumber := g.M.Read(finalSpriteAddr - 4*nSprite + 2)
		if height == 2 {
			patternNumber &= 0xFE
		}
		flags := g.M.Read(finalSpriteAddr - 4*nSprite + 3)
		priority := flags&0x80 == 0
		tileAddr := 0x8000 + uint16(patternNumber)*16
		obp := obps[flags&0x10>>4]

		for i := 0; i < 8*height; i++ {
			tileLineUp := g.M.Read(tileAddr)
			tileAddr++
			tileLineDown := g.M.Read(tileAddr)
			tileAddr++

			binaryTileLineUp := fmt.Sprintf("%08b", tileLineUp)
			binaryTileLineDown := fmt.Sprintf("%08b", tileLineDown)

			for j := 0; j < 8; j++ {
				// Check the pixel is on screen
				if xPosition+j >= 8 && xPosition+j < 168 && yPosition+i >= 16 && yPosition+i < 160 {
					// Check that the sprite either has priority or is on a 00 pixel.
					if priority || screen.At(200+xPosition+j-8, yPosition+i-16) == color00 {
						pair := string(binaryTileLineDown[j]) + string(binaryTileLineUp[j])
						obColor, _ := strconv.ParseInt(pair, 2, 8)
						if obColor != 0 {
							pixelColor := obp[obColor]
							screen.Set(200+xPosition+j-8, yPosition+i-16, pixelColor)
						}
					}
				}
			}
		}
	}
}

func (g *Game) debugMemory(screen *ebiten.Image) {
	bytesToWrite := ""
	current := 0xFE00
	for current <= 0xFEA0 {
		endLine := current + 0x07
		for current <= endLine {
			bytesToWrite += fmt.Sprintf("%02X ", g.M.Read(uint16(current)))
			current++
		}
		bytesToWrite += "\n"
	}

	ebitenutil.DebugPrint(screen, bytesToWrite)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 576
}

func (g *Game) CheckInterruptRequests() {
	// This is so dirty but simple to check if we've let the next instruction run
	if g.M.IMESteps == 1 {
		g.M.IMESteps++
	} else if g.M.IMESteps == 2 {
		g.M.IME = g.M.IMEReqType
		g.M.IMESteps = 0
	}
}

func (g *Game) InterruptStep() {
	iFlags := g.M.Read(0xFF0F)
	if g.M.IME && g.M.IER[0] != 0 && iFlags != 0 {
		g.M.IME = false
		if g.M.IER[0]&0x01 == 1 && iFlags&0x01 == 1 {
			// VBlank
			g.M.Store(0xFF0F, iFlags&0xFE)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0040
		} else if g.M.IER[0]&0x02 == 0x02 && iFlags&0x02 == 0x02 {
			// LCDC
			g.M.Store(0xFF0F, iFlags&0xFD)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0048
		} else if g.M.IER[0]&0x04 == 0x04 && iFlags&0x04 == 0x04 {
			// TIMA overflow
			g.M.Store(0xFF0F, iFlags&0xFB)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0050
		} else if g.M.IER[0]&0x08 == 0x08 && iFlags&0x08 == 0x08 {
			// Serial I/O transfer complete
			g.M.Store(0xFF0F, iFlags&0xF7)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0058
		} else if g.M.IER[0]&0x10 == 0x10 && iFlags&0x10 == 0x10 {
			// High to low P10-P13
			g.M.Store(0xFF0F, iFlags&0xEF)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0060
		}
		g.R.Halted = false
	}
}

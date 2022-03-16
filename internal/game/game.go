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

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var cyclesPerFrame = 4194304 / 4 / 60
var cyclesPerDivUpdate = cyclesPerFrame / (16384 / 60)
var gameFont font.Face

type Game struct {
	R     *registers.Registers
	M     *memory.Memory
	GPU   *gpu.GPU
	Pause bool
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
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Pause = true
	}
	if !g.Pause {
		for currentCycles < cyclesPerFrame {
			// Read always 3 bytes: op code and 2 possible arguments
			instructionArray := g.M.ReadInstruction(g.R.PC)
			fmt.Printf("%X %X %X\n", instructionArray[0], instructionArray[1], instructionArray[2])
			err, bytes, cycles := instructions.Execute(g.R, g.M, instructionArray)
			if err != nil {
				panic(err)
			}
			// Add cycles executed to the current cycles of the frame
			currentCycles += cycles
			// Check if we have to update DIV
			divCycles += cycles
			if divCycles >= cyclesPerDivUpdate {
				divCycles = 0
				g.M.Store(0xFF04, g.M.Read(0xFF04)+1)
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
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.White)

	// lcdc := g.M.Read(0xFF40)
	// scx := g.M.Read(0xFF42)
	// scy := g.M.Read(0xFF43)

	tileMapAddr := uint16(0x9800)
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {

			tileNumber := g.M.Read(tileMapAddr)
			tileAddr := 0x8000 + uint16(tileNumber)*16
			fmt.Printf("%04X ", tileAddr)

			for i := 0; i < 8; i++ {
				tileLineUp := g.M.Read(tileAddr)
				tileAddr++
				tileLineDown := g.M.Read(tileAddr)
				tileAddr++

				binaryTileLineUp := fmt.Sprintf("%08b", tileLineUp)
				binaryTileLineDown := fmt.Sprintf("%08b", tileLineDown)

				for j := 0; j < 8; j++ {
					pair := string(binaryTileLineDown[j]) + string(binaryTileLineUp[j])
					var pixelColor color.Color
					switch pair {
					case "00":
						pixelColor = color.RGBA{0xE0, 0xF8, 0xCF, 0xFF}
					case "01":
						pixelColor = color.RGBA{0x86, 0xC0, 0x6C, 0xFF}
					case "10":
						pixelColor = color.RGBA{0x30, 0x68, 0x50, 0xFF}
					case "11":
						pixelColor = color.RGBA{0x07, 0x18, 0x21, 0xFF}
					}
					screen.Set(8*x+j, 8*y+i, pixelColor)
				}
			}
			tileMapAddr++
		}
		tileMapAddr += 12
	}
	/*
		tileMapAddr := uint16(0x8000)
		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				for i := 0; i < 8; i++ {
					tileLineUp := g.M.Read(tileMapAddr)
					tileMapAddr++
					tileLineDown := g.M.Read(tileMapAddr)
					tileMapAddr++

					binaryTileLineUp := fmt.Sprintf("%08b", tileLineUp)
					binaryTileLineDown := fmt.Sprintf("%08b", tileLineDown)

					for j := 0; j < 8; j++ {
						pair := string(binaryTileLineDown[j]) + string(binaryTileLineUp[j])
						var pixelColor color.Color
						switch pair {
						case "00":
							pixelColor = color.RGBA{0xE0, 0xF8, 0xCF, 0xFF}
						case "01":
							pixelColor = color.RGBA{0x86, 0xC0, 0x6C, 0xFF}
						case "10":
							pixelColor = color.RGBA{0x30, 0x68, 0x50, 0xFF}
						case "11":
							pixelColor = color.RGBA{0x07, 0x18, 0x21, 0xFF}
						}
						screen.Set(8*y+j, 8*x+i, pixelColor)
					}
				}
			}
		}
	*/
	bytesToWrite := ""
	current := 0x8000
	for current <= 0x8100 {
		endLine := current + 0x10
		for current <= endLine {
			bytesToWrite += fmt.Sprintf("%02X ", g.M.Read(uint16(current)))
			current++
		}
		bytesToWrite += "\n"
	}
	// ebitenutil.DebugPrint(screen, bytesToWrite)
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
			g.M.Store(0xFF0F, iFlags&0xFE)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0040
		}
	}
}

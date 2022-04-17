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

// Frequency of the Game Boy (cycles per second)
var frequency = 4194304

// Divide by 60 and we get how many cycles can be ran per Update.
// Necessary to limit the speed of the emulator to that of an actual GB.
var cyclesPerFrame = frequency / 60

// Cycles required to update value of DIV register
var cyclesPerDivUpdate = cyclesPerFrame / (16384 / 60)

// Cycles required to update value of TIMA register (one for each possible speed)
var cyclesPerTimaUpdate = [4]int{
	cyclesPerFrame / (4096 / 60),
	cyclesPerFrame / (262144 / 60),
	cyclesPerFrame / (65536 / 60),
	cyclesPerFrame / (16384 / 60),
}
var gameFont font.Face

// The 4 different colors in the Game Boy, from lighter to darker.
var color00 = color.RGBA{0xE0, 0xF8, 0xCF, 0xFF}
var color01 = color.RGBA{0x86, 0xC0, 0x6C, 0xFF}
var color10 = color.RGBA{0x30, 0x68, 0x50, 0xFF}
var color11 = color.RGBA{0x07, 0x18, 0x21, 0xFF}
var colors = [4]color.RGBA{color00, color01, color10, color11}

// This struct represents the console and its components.
// It's where the main instructions disassembling and execution routine happens.
type Game struct {
	R     *registers.Registers
	M     *memory.Memory
	GPU   *gpu.GPU
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

// Update function. Here the instructions are executed, the GPU states updated and the interruptions processed.
func (g *Game) Update() error {
	// Initialize all cycle counters at 0
	currentCycles := 0
	divCycles := 0
	timaCycles := 0
	// Run instructions until we reach the maximum an actual GB would have ran in the same time.
	for currentCycles < cyclesPerFrame {
		var err error
		var bytes uint16
		var cycles int
		if !g.R.Halted {
			// Read always 3 bytes: op code and 2 possible arguments
			instructionArray := g.M.ReadInstruction(g.R.PC)
			if g.Debug {
				fmt.Printf("%X %X %X\n", instructionArray[0], instructionArray[1], instructionArray[2])
				fmt.Println(runtime.FuncForPC(reflect.ValueOf(instructions.InstructionTable[instructionArray[0]]).Pointer()).Name())
				fmt.Println(g.R)
			}
			// Execute the next instruction.
			err, bytes, cycles = instructions.Execute(g.R, g.M, instructionArray)
			if err != nil {
				panic(err)
			}
		} else {
			// The CPU is halted. The clock ticks, but no instructions are executed until a new interruption happens.
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
					// TIMA overflow! Set the TIMA interruption flag and set TIMA as TAM
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

// Draw function. Prints the tiles and sprites, but does not execute instructions.
func (g *Game) Draw(screen *ebiten.Image) {

	// Fill the whole screen with gray, so that looking at it doesn't hurt our eyes.
	screen.Fill(color.Gray{0x77})

	// Take the LCD controller data
	lcdc := g.M.Read(0xFF40)
	// scx := g.M.Read(0xFF42)
	// scy := g.M.Read(0xFF43)

	// Transfer sprites data to OAM
	g.transferOAM()

	// Display background and window?
	if lcdc&0x01 == 1 {
		// Draw the background
		g.drawBackground(screen, lcdc)

		// Display sprites?
		if lcdc&0x02 != 0 {
			//Draw sprites
			g.drawSprites(screen, lcdc)
		}
	}
	// g.debugMemory(screen)
}

// Transfer sprites data to OAM.
// Basically copy 0xA0 bytes of data to OAM starting at the address in 0xFF46 followed by two 0s.
func (g *Game) transferOAM() {
	address := uint16(g.M.Read(0xFF46)) << 8

	for i := 0; i < 0xA0; i++ {
		g.M.OAM[i] = g.M.Read(address + uint16(i))
	}
}

// Method to draw the background of the game.
func (g *Game) drawBackground(screen *ebiten.Image, lcdc byte) {
	var tileMapAddr, tileDataAddr uint16
	var signed bool

	// The colour palette for the background. Each 2 bits indicate one of the 4 colours defined above.
	bgp := g.M.Read(0xFF47)
	bgpColors := [4]color.RGBA{
		colors[bgp&0x03],
		colors[bgp&0x0C>>2],
		colors[bgp&0x30>>4],
		colors[bgp>>6],
	}

	// 4th bit of LCDC indicates whether the tile map for the background starts at 0x9800 or 0x9C00.
	if lcdc&0x08 == 0 {
		tileMapAddr = 0x9800
	} else {
		tileMapAddr = 0x9C00
	}

	// 5th bit of LCDC indicates whether the 0 address of the tiles data (the actual graphic data,
	// not their disposition on the screen) is 0x8000 or 0x9000.
	if lcdc&0x10 == 0 {
		signed = true
		tileDataAddr = 0x9000
	} else {
		signed = false
		tileDataAddr = 0x8000
	}

	// GB screen is 20x18 tiles.
	for y := 0; y < 18; y++ {
		for x := 0; x < 20; x++ {
			// Get the number of the current tile to be drawn.
			tileNumber := g.M.Read(tileMapAddr)
			// Adjust with the data address
			tileAddr := tileDataAddr + uint16(tileNumber)*16
			// Haha funny. So apparently if the 0 address of the data is 0x8000, the tile number is unsigned,
			// but if it's 0x9000, then it is signed and ranges from -126 to 127 and we need to adjust that too.
			if signed && tileNumber >= 0x7F {
				tileAddr = (0x0800 + uint16(tileNumber)) << 4
			}

			for i := 0; i < 8; i++ {
				// To print the tile, we need to read 16 bytes in groups of 2.
				// Every 2 bytes represent one line of 8 pixels in the tile. The way this works is:
				// The first tile represents the LSB of the 8 pixels, the second represents the MSB.
				// Group each LSB with its respective MSB and you'll get a list of numbers from 0 to 3,
				// representing one of the 4 colours to be printed in that spot. Now to the code:

				// Read the two lines.
				tileLineLSB := g.M.Read(tileAddr)
				tileAddr++
				tileLineMSB := g.M.Read(tileAddr)
				tileAddr++

				// Transform them to binary.
				binaryTileLineLSB := fmt.Sprintf("%08b", tileLineLSB)
				binaryTileLineMSB := fmt.Sprintf("%08b", tileLineMSB)

				for j := 0; j < 8; j++ {
					// Pair the LSB with the MSB, parse as binary, select the color and print the pixel.
					pair := string(binaryTileLineMSB[j]) + string(binaryTileLineLSB[j])
					bgColor, _ := strconv.ParseInt(pair, 2, 8)
					pixelColor := bgpColors[bgColor]
					screen.Set(8*x+j, 8*y+i, pixelColor)
				}
			}
			// Next tile.
			tileMapAddr++
		}
		// We've ended a line. Jump the next 12 tiles, since they're outside the visible screen.
		tileMapAddr += 12
	}
}

func (g *Game) drawWindow(screen *ebiten.Image, lcdc byte) {

}

// Method to draw the sprites.
func (g *Game) drawSprites(screen *ebiten.Image, lcdc byte) {
	var height int

	// Sprites can be coloured with two different palettes, OBP0 and OBP1.
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

	// Sprites can also be 8x8 (height == 1), or 8x16 (height == 2).
	if lcdc&0x04 == 0 {
		height = 1
	} else {
		height = 2
	}

	// Draw sprites from end to start, because the ones in the start have more priority and should be drawn above the others.
	var finalSpriteAddr uint16 = 0xFE9C

	for nSprite := uint16(0); nSprite < 40; nSprite++ {
		// Sprites have 4 bytes of data:
		// Byte 0: Y position on the screen.
		// Byte 1: X position on the screen.
		// Byte 2: number of pattern in the tile map (sprites always start at 0x8000)
		// Byte 3: priority, flip, and palette flags.
		yPosition := int(g.M.Read(finalSpriteAddr - 4*nSprite))
		xPosition := int(g.M.Read(finalSpriteAddr - 4*nSprite + 1))
		patternNumber := g.M.Read(finalSpriteAddr - 4*nSprite + 2)
		if height == 2 {
			patternNumber &= 0xFE
		}
		tileAddr := 0x8000 + uint16(patternNumber)*16
		flags := g.M.Read(finalSpriteAddr - 4*nSprite + 3)
		priority := flags&0x80 == 0
		yFlip := flags&0x40 != 0
		xFlip := flags&0x20 != 0
		obp := obps[flags&0x10>>4]

		// The way to draw a sprite is almost the same as to draw the background.
		// However, we need to do a few more checks.
		for i := 0; i < 8*height; i++ {
			tileLineLSB := g.M.Read(tileAddr)
			tileAddr++
			tileLineMSB := g.M.Read(tileAddr)
			tileAddr++

			binaryTileLineLSB := fmt.Sprintf("%08b", tileLineLSB)
			binaryTileLineMSB := fmt.Sprintf("%08b", tileLineMSB)

			for j := 0; j < 8; j++ {
				var screenXPosition, screenYPosition int

				// If flip flags are set, the sprites need to be drawn starting in the opposite side of the axis.
				if !xFlip {
					screenXPosition = xPosition + j - 8
				} else {
					screenXPosition = xPosition - j
				}
				if !yFlip {
					screenYPosition = yPosition + i - 16
				} else {
					screenYPosition = yPosition - i
				}

				// Check that the pixel is on screen and should be drawn.
				if screenXPosition >= 0 && screenXPosition < 160 && screenYPosition >= 0 && screenYPosition < 144 {
					adjustedXPosition := screenXPosition
					adjustedYPosition := screenYPosition
					// Check that the sprite either has priority or is on a 00 pixel.
					// Sprites should not be drawn on top of the background unless one of these two conditions is true.
					if priority || screen.At(adjustedXPosition, adjustedYPosition) == color00 {
						pair := string(binaryTileLineMSB[j]) + string(binaryTileLineLSB[j])
						obColor, _ := strconv.ParseInt(pair, 2, 8)
						// if the color of the pixel is 0, it is considered transparent and should not replace the background.
						if obColor != 0 {
							pixelColor := obp[obColor]
							screen.Set(adjustedXPosition, adjustedYPosition, pixelColor)
						}
					}
				}
			}
		}
	}
}

// Method to print the contents of a part of the memory. Only for debugging
func (g *Game) debugMemory(screen *ebiten.Image) {
	bytesToWrite := ""
	// First address to print.
	current := 0xFF00
	// Print until we reach the specified address.
	for current <= 0xFF00 {
		endLine := current + 0x0F
		for current <= endLine {
			bytesToWrite += fmt.Sprintf("%02X ", g.M.Read(uint16(current)))
			current++
		}
		bytesToWrite += "\n"
	}

	ebitenutil.DebugPrint(screen, bytesToWrite)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 4, outsideHeight / 4
}

// Check if there is an interrupt requested, and let one more instruction run.
// This looks dirty to me but simple. It just works. Whatever.
func (g *Game) CheckInterruptRequests() {
	if g.M.IMESteps == 1 {
		g.M.IMESteps++
	} else if g.M.IMESteps == 2 {
		g.M.IME = g.M.IMEReqType
		g.M.IMESteps = 0
	}
}

// Interrupts step. Executed after every instruction.
// Checks if the IME is set, if there are any interrupts requested, and if that interrupt flag is set in the IER.
// If so, reset IME, immediately call the corresponding interrupt subroutine, and resume CPU activity.
func (g *Game) InterruptStep() {
	iFlags := g.M.Read(0xFF0F)
	if g.M.IME && g.M.IER[0] != 0 && iFlags != 0 {
		g.M.IME = false
		if g.M.IER[0]&0x01 == 1 && iFlags&0x01 == 1 {
			// VBlank
			g.M.Store(0xFF0F, iFlags&0xFE)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0040
			g.R.Halted = false
		} else if g.M.IER[0]&0x02 == 0x02 && iFlags&0x02 == 0x02 {
			// LCDC
			g.M.Store(0xFF0F, iFlags&0xFD)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0048
			g.R.Halted = false
		} else if g.M.IER[0]&0x04 == 0x04 && iFlags&0x04 == 0x04 {
			// TIMA overflow
			g.M.Store(0xFF0F, iFlags&0xFB)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0050
			g.R.Halted = false
		} else if g.M.IER[0]&0x08 == 0x08 && iFlags&0x08 == 0x08 {
			// Serial I/O transfer complete
			g.M.Store(0xFF0F, iFlags&0xF7)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0058
			g.R.Halted = false
		} else if g.M.IER[0]&0x10 == 0x10 && iFlags&0x10 == 0x10 {
			// High to low P10-P13
			g.M.Store(0xFF0F, iFlags&0xEF)
			utils.PushStackShort(g.R, g.M, g.R.PC)
			g.R.PC = 0x0060
			g.R.Halted = false
		}
	}
}

package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"go-boy/internal/instructions"
	"go-boy/internal/memory"
	"go-boy/internal/registers"
)

type Game struct {
	R *registers.Registers
	M *memory.Memory
}

func (g *Game) Update() error {
	// Read always 3 bytes: op code and 2 possible arguments
	instructionArray := g.M.Cartridge[g.R.PC : g.R.PC+3]
	fmt.Printf("%X %X %X\n", instructionArray[0], instructionArray[1], instructionArray[2])
	err, bytes := instructions.Execute(g.R, g.M, instructionArray)
	if err != nil {
		panic(err)
	}
	// Augment the PC as much as the amount of bytes the instruction has used
	g.R.PC += bytes
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 576
}

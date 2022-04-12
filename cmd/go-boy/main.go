package main

import (
	"fmt"
	game2 "go-boy/internal/game"
	"go-boy/internal/gpu"
	"go-boy/internal/memory"
	"go-boy/internal/registers"
	"log"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	// First of all, check that the user passed a file as game. Panic otherwise.
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	// Hardcoded bytes for the Nintendo scrolling graphic.
	scrollingNintendoGraphic := [48]byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
		0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
		0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
		0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
		0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
		0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E}

	// Array to get the scrolling graphic from the cartridge.
	scrollingNintendoGraphicB := make([]byte, 48)
	// Array for the bytes with the title of the game.
	titleB := make([]byte, 16)

	// Read the scrolling graphic bytes from the cartridge.
	// If they don't exactly match the hardcoded scrolling graphic data, panic.
	_, err = file.ReadAt(scrollingNintendoGraphicB, 0x104)
	for i := range scrollingNintendoGraphicB {
		if scrollingNintendoGraphicB[i] != scrollingNintendoGraphic[i] {
			panic(fmt.Errorf("the scrolling graphic isn't correct. Game not valid"))
		}
	}
	// Read the title
	_, err = file.ReadAt(titleB, 0x134)
	if err != nil {
		panic(err)
	}

	// Initialize a Game struct
	game := &game2.Game{
		R:     registers.GetInitializedRegisters(),
		M:     memory.GetInitializedMemory(file),
		GPU:   gpu.InitGPU(),
		Debug: false,
	}

	if err = file.Close(); err != nil {
		panic(err)
	}

	// Set the window's size and name.
	ebiten.SetWindowSize(640, 576)
	ebiten.SetWindowTitle(string(titleB))
	// Run the emulator's main loop.
	if err = ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

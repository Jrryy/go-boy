package main

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
	"gitlab.com/jrryy/go-boy/internal/instructions"
	"gitlab.com/jrryy/go-boy/internal/registers"
	"os"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())

	glfw.WindowHint(glfw.Resizable, glfw.False) // Make window non resizable

	file, err := os.Open("/home/gerard/Downloads/Tetris.gb")
	if err != nil {
		panic(err)
	}

	scrollingNintendoGraphic := []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
		0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
		0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
		0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
		0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
		0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E}

	scrollingNintendoGraphicB := make([]byte, 48) // The scrolling graphic thing. Only for integrity purposes.
	titleB := make([]byte, 16) // Obtain the bytes with the title of the game.
	instructionArray := make([]byte, 3) // Read always 3 bytes: op code and 2 possible arguments

	_, err = file.ReadAt(scrollingNintendoGraphicB, 0x104)
	for i := range scrollingNintendoGraphicB {
		if scrollingNintendoGraphicB[i] != scrollingNintendoGraphic[i] {
			panic(fmt.Errorf("the scrolling graphic isn't correct. Game not valid"))
		}
	}
	_, err = file.ReadAt(titleB, 0x134)
	if err != nil {
		panic(err)
	}

	window, err := glfw.CreateWindow(640, 576, vk.ToString(titleB), nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	r := registers.InitializeRegisters()

	for !window.ShouldClose() {
		_, err = file.ReadAt(instructionArray, r.PC)
		if err != nil {
			panic(err)
		}

		err, bytes := instructions.Execute(r, instructionArray)
		if err != nil {
			 panic(err)
		}

		// Augment the PC as much as the amount of bytes of the instruction
		r.PC += int64(bytes)

		window.SwapBuffers()
		glfw.PollEvents()

	}
}

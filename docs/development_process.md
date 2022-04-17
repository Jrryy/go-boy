# Development process

## First thoughts
At first all I knew was that I wanted to make this using Go, but given my background I had no idea how to build any GUI or anything of the sort, nor what to use. I came across an article about ebiten, a library to build 2D games with Go. It was stable, under active development, looked easy to use and handled the main game update and graphics draw loops. 

With all this decided, I could begin with the actual development.

## How does a Game Boy work?
No idea, but that's what the Internet is for. My main source of information was the [unofficial Game Boy CPU manual](http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf), which although contains a few errors here and there, is complete and clear enough to create a functioning emulator.

## First steps
What I did first of all was explore the basic information data in the cartridge. There's a lot of data in them, but since I was going to start with basic ROM only games (Tetris, Dr. Mario, etc...) I could skip most of them and keep the two most basic ones to me: the game title and the scrolling Nintendo graphic. I use the scrolling graphic just to check that the file input is actually a GB game, but I skip directly to the first instruction. The title is simply what I use to name the window I create.

When this worked and a window was correctly created, I could start actually working on the CPU.

## Registers
The most basic part of the CPU are the registers. There are 8 that contain a byte of data (A, B, C, D, E, F, H and L), and 2 that contained 2 bytes (PC, program counter, and SP, stack pointer). I created a struct with them in it, and added the flags as bools as well to handle them more easily (this will prove to have been a bad decision later on, since the register F should have been more than enough). I also created some functions to retrieve the 16 bit combinations of 2 registers (AF, BC, DE and HL), the function to initialize them as the guide indicates, and the function to format a sprint with information on all the registers for debugging purposes. With that I was good to start coding the instructions.

## Instructions
My first goal was to code enough instructions to reach a loop representing the [very first screen](https://i.imgur.com/LvApfib.png) on Tetris, at which point I planned on starting to work on the actual graphics to actually see that screen. At this point, the speed at which my emulator worked was not important, so I could leave the handling of the cycles for later when it would become necessary. 

The way I wanted the execution of instructions to work had to be generic and simple. In the main update loop of ebiten, I took three bytes from the cartridge where PC indicates (1 for the instruction, 2 for immediate arguments), and then executed a function to which I would only need to send those 3 bytes of data and it would return both an error and an amount to add to PC. Also, for debugging purposes, I printed out the instruction being executed and the contents of all the registers and flags.

To store and execute the instructions, I built an array in which every function was stored in the position corresponding to that instruction's opcode. I might be wrong, but this is probably the easiest way to store and handle them. Then, I created an unimplemented function that, when the instruction to be executed was not yet coded, would panic and print the opcode and args of that instruction. With this, I could code the instructions in 4 steps:

1. Build and execute the emulator.
2. Wait for it to panic and check the opcode.
3. Look it up in the [GB opcode table](http://imrannazar.com/Gameboy-Z80-Opcode-Map).
4. Code it according to the guide and go back to step 1.

Very soon, though, instructions that handled memory would appear. Time to build the memory.

## Memory
For the memory, I built a struct with arrays of bytes for every defined part of it (cartridge, VRAM, etc...). Has this been useful? So far, not really. In fact, I would say it has made development and my code slightly harder and more complex, and a single array with all data would have worked just well. Maybe it will be useful to handle cartridges with MBC... I don't know.

Anyway, I initialized the memory according to the guide, I made the functions to read from and write to the memory, and... for ROM only cartridges, this was really it, I was good to go.

## Timings, interrupts and GPU
All these kinda came along together. Reading through the guide, I realised that I would never see any graphics if there weren't VBlank interrupts in my emulator, there would never be VBlank interrupts if it did never go into VBlank mode, and it would never go into VBlank mode if it didn't have a properly functioning GPU.

### Timings

First of all, I had to simulate the timings. My way of doing this took advantage of the fact that ebiten's `Update` and `Draw` functions are executed an average of 60 times a second. All I had to do was divide the GB's frequency / 60, and I would already know how many ticks I'd have to emulate every call to `Update`. I also added the information on ticks to the instructions table. With this, I could start working on the interrupts.

### Interrupts
The interrupts were not very complex either. I was planning on building a struct for them, but since all I really had to make was the master enable (IME), I decided to store that in the memory where the enable register (IER) already was.

Due to the way interrupts are enabled and disabled on the GB, I also added two more pieces of data: `IMEReqType` and `IMESteps`. `IMEReqType` was a boolean set to `true` or `false` if the interrupts are to be enabled or disabled respectively, and `IMESteps` exists because this change of state doesn't happen until one more instruction has been executed after requesting it, so it's just set to 1 when the change has been requested, then to 2 the next instruction, then to 0 and the change of state happens.

The way to execute them is very straight forward: if IME is enabled, a certain interrupt is enabled in IER, and in the interrupt flags register (`0xFF0F` in memory) that interrput's flag is set, call a routine in a fixed address.

All of this would be in a function executed right after every instruction.

### GPU
It was a little difficult to figure out how to make the GPU. At this point I didn't quite understand what the 4 different modes were exactly for, but I figured the most important parts were the ticks that each one of them lasted for, and the fact that as soon as it entered VBlank mode, the interrupt that runs the routine to draw the screen (or sets the sprites for background and windows, really) was run. So far, this was enough for the interrupt routines to run. I figured this because now the emulator reached a `RETI` (return from interrupt) instruction. 

## Drawing the background
To begin with, I didn't even try to draw the background, just the tiles in order to check that I was drawing them properly. Putting them in their place would come later. I was using Tetris, so I knew in which memory address the tiles were stored, and that every sprite was of size 8x8. I developed a small loop to print the tiles pixel by pixel and after a few tries, [this appeared](https://i.imgur.com/G5TtAUS.png). Well that looks quite alright, you can see the numbers, letters, and even the building in the player selection screen. I then mapped the sprites to their right place using the map starting at `0x9800`, and saw [the copyright screen](https://i.imgur.com/LvApfib.png). Eventually, after coding a few more instructions, the game reached the [player selection screen](https://i.imgur.com/0SxqTPc.png) too. The background is working all right now!

## Drawing the sprites
The sprites are quite the same thing as a tile map: their data contains the tile to print, and the x and y coordinates of its top left corner, and a few extra flags to flip the sprite or to print it above winows (I hadn't done these at this point yet). However, it took me a while to realise that the process to copy the sprites to OAM isn't in the game, but in the hardware. Apparently, the byte at `0xFF46` followed by two 0s indicates the address of the first piece of data to be copied to OAM and then, in the memory, there's a bit of code that loops for long enough to let the whole sprites map be copied and placed where it has to be. Since the emulator can do all of that instantly, I'm simplifying this so that whenever the contents of `0xFF46` change, I immediately copy the sprites map to OAM. This is not a problem since the loop happens immediately after that memory value is set.

## Handling user input
For this, all I really needed to do was map the different controls to the different 4 least significant bits of `0xFF00` depending on the input mode active. Quite easy, although for some reason the guide has the modes swapped. Not a big problem, but it boggled my mind for a few minutes. Since it was all quite easy to implement, I went ahead and also made a controller usable. I tested it and it worked alright. Nice stuff.
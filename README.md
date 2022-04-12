# go-boy
The go-boy is a Game Boy emulator made entirely with Golang. The graphics and main loop have been done using [ebiten](https://ebiten.org/).

To build and use it:

```
cd go-boy
go build ./cmd/go-boy
./go-boy ~/path/to/the/game.gb
```

Also, for reference on my tought process while building this, check out [my development process](docs/development_process.md).

## Current state and next steps
The demo gameplay on Tetris and Dr. Mario are working and looping correctly.

Next steps: handle user input.

## Aren't there enough emulators already?
Yes, but I made this one myself :)

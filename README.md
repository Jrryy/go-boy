# go-boy
The go-boy is a Game Boy emulator made entirely with Golang. The graphics and main loop have been done using [ebiten](https://ebiten.org/).

To build and use it:

```
cd go-boy
go build ./cmd/go-boy
./go-boy ~/path/to/the/game.gb
```

Button mapping:
```
A -> Z
B -> X
Start -> Enter
Select -> Backspace
Down -> Arrow down
Up -> Arrow up
Left -> Arrow left
Right -> Arrow right
```
A controller can also be used.

Also, for reference on my tought process while building this, check out [my development process](docs/development_process.md).

## Current state and next steps
ROM only games playable with both keyboard and controller.

Next steps: Add all the instructions.

## Aren't there enough emulators already?
Yes, but I made this one myself :)

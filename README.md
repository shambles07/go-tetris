A fork of a simple terminal tetris implementation in Go, originally by [cespare](https://github.com/cespare).

Features listed in bold are additions.

## Installation:
    $ go get github.com/shambles07/go-tetris

## Usage:
    $ go-tetris

## Controls
* Move piece down: `↓`, `j`
* Move piece left: `←`, `h`
* Move piece right: `→`, `l`
* Rotate CW: `↑`, `k`, **`x`**
* **Rotate CCW: `z`**
* **Rotate 180: `a`**
* Quick drop: `space`
* Quit: `q`, `ctrl-c`

## Implemented features
* Simplified 7-bag piece order
* Quick drop
* Rotation (CW/**CCW**)
* Next box
* Speeding up
* Pausing

## Feature Roadmap
* High scores
* Hold piece
* Ghost piece
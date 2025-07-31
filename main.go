package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Grid struct {
	w, h int
	data []byte
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func get_data() (Grid, int) {
	filename := os.Args[1]
	iterations, err := strconv.Atoi(os.Args[2])
	check(err)
	dat, err := os.ReadFile(filename)
	check(err)

	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
	h := len(lines)
	w := len(lines[0])
	data := make([]byte, h*w)

	for y, line := range lines {
		for x := 0; x < w; x++ {
			if line[x] == 'X' {
				data[y*w+x] = 1
			} else {
				data[y*w+x] = 0
			}
		}
	}
	return Grid{w, h, data}, iterations
}

func printState(grid Grid) {
	for y := 0; y < grid.h; y++ {
		rowStart := y * grid.w
		for x := 0; x < grid.w; x++ {
			if grid.data[rowStart+x] == 1 {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func Step(current, next Grid) {
	w, h := current.w, current.h
	offsets1D := []int{-w - 1, -w, -w + 1, -1, 1, w - 1, w, w + 1}
	offsets2D := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			idx := y*w + x
			liveNeighbors := 0
			for _, offset := range offsets1D {
				liveNeighbors += int(current.data[idx+offset])
			}
			next.data[idx] = updateCell(current.data[idx], liveNeighbors)
		}
	}

	// Process edges with bounds checks
	// top and bottom rows
	for x := 0; x < w; x++ {
		for _, y := range []int{0, h - 1} {
			idx := y*w + x
			liveNeighbors := 0
			for _, offset := range offsets2D {
				nx, ny := x+offset[0], y+offset[1]
				if nx >= 0 && nx < w && ny >= 0 && ny < h {
					liveNeighbors += int(current.data[ny*w+nx])
				}
			}
			next.data[idx] = updateCell(current.data[idx], liveNeighbors)
		}
	}
	for y := 1; y < h-1; y++ {
		for _, x := range []int{0, w - 1} {
			idx := y*w + x
			liveNeighbors := 0
			for _, offset := range offsets2D {
				nx, ny := x+offset[0], y+offset[1]
				if nx >= 0 && nx < w && ny >= 0 && ny < h {
					liveNeighbors += int(current.data[ny*w+nx])
				}
			}
			next.data[idx] = updateCell(current.data[idx], liveNeighbors)
		}
	}
}

func updateCell(cell byte, liveNeighbors int) byte {
	if cell == 1 {
		if liveNeighbors < 2 || liveNeighbors > 3 {
			return 0
		}
		return 1
	} else {
		if liveNeighbors == 3 {
			return 1
		}
		return 0
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./life <filename> <iterations>")
		os.Exit(1)
	}

	current, iterations := get_data()
	next := Grid{current.w, current.h, make([]byte, current.h*current.w)}

	var i int
	_ = i
	for i = range iterations {
		Step(current, next)
		current, next = next, current
	}
}

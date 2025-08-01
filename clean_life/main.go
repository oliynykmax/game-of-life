package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	// "time"
)

type Grid struct {
	w, h int
	data []byte
}

var rules = [2][9]byte{
	{0, 0, 0, 1, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 0, 0, 0, 0, 0},
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
		rowOffset := y * w
		for x := 0; x < w; x++ {
			if line[x] == 'X' {
				data[rowOffset+x] = 1
			}
		}
	}
	return Grid{w, h, data}, iterations
}

func printState(grid Grid) {
	var b strings.Builder
	for y := 0; y < grid.h; y++ {
		rowStart := y * grid.w
		for x := 0; x < grid.w; x++ {
			if grid.data[rowStart+x] == 1 {
				b.WriteByte('X')
			} else {
				b.WriteByte('.')
			}
		}
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	fmt.Print(b.String())
}

func Step(current, next Grid) {
	w, h := current.w, current.h
	data := current.data
	nextData := next.data

	// Process interior cells (no bounds checking needed)
	for y := 1; y < h-1; y++ {
		rowAbove := (y - 1) * w
		rowCurrent := y * w
		rowBelow := (y + 1) * w

		for x := 1; x < w-1; x++ {
			// Unrolled neighbor sum - direct array access
			neighbors := int(data[rowAbove+x-1]) + int(data[rowAbove+x]) + int(data[rowAbove+x+1]) +
				int(data[rowCurrent+x-1]) + int(data[rowCurrent+x+1]) +
				int(data[rowBelow+x-1]) + int(data[rowBelow+x]) + int(data[rowBelow+x+1])

			idx := rowCurrent + x
			nextData[idx] = rules[data[idx]][neighbors]
		}
	}

	// Process edges and corners efficiently (each cell processed exactly once)

	// Top and bottom rows
	for x := 0; x < w; x++ {
		// Top row
		neighbors := 0
		for dy := 0; dy <= 1; dy++ {
			ny := dy
			for dx := -1; dx <= 1; dx++ {
				nx := x + dx
				if nx >= 0 && nx < w && !(dx == 0 && dy == 0) {
					neighbors += int(data[ny*w+nx])
				}
			}
		}
		nextData[x] = rules[data[x]][neighbors]

		// Bottom row (skip if h == 1)
		if h > 1 {
			neighbors = 0
			bottomIdx := (h-1)*w + x
			for dy := -1; dy <= 0; dy++ {
				ny := h - 1 + dy
				for dx := -1; dx <= 1; dx++ {
					nx := x + dx
					if nx >= 0 && nx < w && !(dx == 0 && dy == 0) {
						neighbors += int(data[ny*w+nx])
					}
				}
			}
			nextData[bottomIdx] = rules[data[bottomIdx]][neighbors]
		}
	}

	// Left and right columns (excluding corners already processed)
	for y := 1; y < h-1; y++ {
		// Left column
		neighbors := 0
		leftIdx := y * w
		for dy := -1; dy <= 1; dy++ {
			ny := y + dy
			for dx := 0; dx <= 1; dx++ {
				nx := dx
				if !(dx == 0 && dy == 0) {
					neighbors += int(data[ny*w+nx])
				}
			}
		}
		nextData[leftIdx] = rules[data[leftIdx]][neighbors]

		// Right column (skip if w == 1)
		if w > 1 {
			neighbors = 0
			rightIdx := y*w + (w - 1)
			for dy := -1; dy <= 1; dy++ {
				ny := y + dy
				for dx := -1; dx <= 0; dx++ {
					nx := w - 1 + dx
					if !(dx == 0 && dy == 0) {
						neighbors += int(data[ny*w+nx])
					}
				}
			}
			nextData[rightIdx] = rules[data[rightIdx]][neighbors]
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./life <filename> <iterations>")
		os.Exit(1)
	}

	current, iterations := get_data()
	next := Grid{current.w, current.h, make([]byte, current.h*current.w)}
	// fmt.Print("\033[H\033[2J")
	// printState(current)
	for i := 0; i < iterations; i++ {
		Step(current, next)
		current, next = next, current
		// fmt.Print("\033[H\033[2J")
		// printState(current)
		// time.Sleep(200 * time.Millisecond) // or time.Second
	}

}

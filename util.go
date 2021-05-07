package main

import (
	"fmt"
	// "flag"
	"bufio"
    "os"
	"strconv"
	"strings"
	// "time"
	// "sync"
	// "math"
)

func assert(b bool) {
	if !b {
		panic(b)
	}
}

var INF = 1 << 27
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func read_input(input_path string) (int, int, []map[int]int) {
	file, err := os.Open(input_path)
	assert(err == nil)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	first := true
	var n, m int
	var adj []map[int]int
	line_ct := 0
    for scanner.Scan() {
		scanner2 := bufio.NewScanner(strings.NewReader(scanner.Text()))
		scanner2.Split(bufio.ScanWords)
		var line []int
		for scanner2.Scan() {
			x, err := strconv.Atoi(scanner2.Text())
			assert(err == nil)
			line = append(line, x)
		}
		if first {
			assert(len(line) == 2)
			n = line[0]
			m = line[1]
			adj = make([]map[int]int, n)
			for i := 0; i < n; i++ {
				adj[i] = make(map[int]int)
			}
			first = false
		} else {
			assert(len(line) == 3)
			a, b, d := line[0], line[1], line[2]
			if val, ok := adj[a][b]; ok {
				adj[a][b] = min(val, d)
			} else {
				adj[a][b] = d
			}
			if val, ok := adj[b][a]; ok {
				adj[b][a] = min(val, d)
			} else {
				adj[b][a] = d
			}
		}
		line_ct += 1
	}
	assert(line_ct == m + 1)
	return n, m, adj
}

func write_output(output_path string, mat [][]int) {
	file, err := os.Create(output_path)
	assert(err == nil)

	for i := range mat {
		fmt.Fprintln(file, mat[i])
	}
	file.Close()
}
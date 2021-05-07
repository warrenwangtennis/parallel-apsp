package util

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

func Assert(b bool) {
	if !b {
		panic(b)
	}
}

var INF = 1 << 27
func Min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func Make_mat(n, m int) [][]int {
	var ret [][]int
	for i := 0; i < n; i++ {
		var row []int
		for j := 0; j < m; j++ {
			row = append(row, 0)
		}
		ret = append(ret, row)
	}
	return ret
}

func Copy_mat(src, dst [][]int) {
	Assert(len(src) == len(dst))
	for i := range(src) {
		Assert(len(src[i]) == len(dst[i]))
		for j := range(src[i]) {
			dst[i][j] = src[i][j]
		}
	}
}

func Read_input(input_path string) (int, int, []map[int]int) {
	file, err := os.Open(input_path)
	Assert(err == nil)
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
			Assert(err == nil)
			line = append(line, x)
		}
		if first {
			Assert(len(line) == 2)
			n = line[0]
			m = line[1]
			adj = make([]map[int]int, n)
			for i := 0; i < n; i++ {
				adj[i] = make(map[int]int)
			}
			first = false
		} else {
			Assert(len(line) == 3)
			a, b, d := line[0], line[1], line[2]
			if val, ok := adj[a][b]; ok {
				adj[a][b] = Min(val, d)
			} else {
				adj[a][b] = d
			}
			if val, ok := adj[b][a]; ok {
				adj[b][a] = Min(val, d)
			} else {
				adj[b][a] = d
			}
		}
		line_ct += 1
	}
	Assert(line_ct == m + 1)
	return n, m, adj
}

func Write_output(output_path string, mat [][]int) {
	file, err := os.Create(output_path)
	Assert(err == nil)

	for i := range mat {
		for j := range mat[i] {
			if mat[i][j] == INF {
				fmt.Fprint(file, "X ")
			} else {
				fmt.Fprint(file, mat[i][j], " ")
			}
		}
		fmt.Fprintln(file)
	}
	file.Close()
}
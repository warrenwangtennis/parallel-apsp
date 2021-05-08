package main

import (
	"fmt"
	"flag"
	"os/exec"
	"../util"
	"time"
	"bytes"
)

var n, m int
var adj []map[int]int
var adj_mat, apsp [][]int

func load_adj_mat() [][]int {
	adj_mat = make([][]int, n)
	for i := 0; i < n; i++ {
		adj_mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				adj_mat[i][j] = 0
			} else {
				adj_mat[i][j] = util.INF
			}
		}
	}
	
	for u := 0; u < n; u++ {
		for v, d := range adj[u] {
			adj_mat[u][v] = d
		}
	}

	return adj_mat
}

func solve() {
	apsp = make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			apsp[i] = append(apsp[i], adj_mat[i][j])
		}
		util.Assert(len(apsp[i]) == n)
	}

	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				apsp[i][j] = util.Min(apsp[i][j], apsp[i][k] + apsp[k][j])				
			}
		}
	}
}

func main() {
	total_start := time.Now()

	var input_path, output_path string
	flag.StringVar(&input_path, "i", "", "string-valued path to an input file")
	flag.StringVar(&output_path, "o", "./floyd_seq_out.txt", "string-valued path to an output file")
	flag.Parse()
	n, m, adj = util.Read_input(input_path)

	load_adj_mat()

	solve_start := time.Now()
	solve()
	fmt.Println("solve time", time.Since(solve_start).Seconds())

	util.Write_output(output_path, apsp)

	fmt.Println("e2e time", time.Since(total_start).Seconds())

	cmd := exec.Command("diff", "-sq", output_path, input_path[:len(input_path) - 4] + "_key.txt")
	var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
	util.Assert(err == nil)
	fmt.Print(out.String())
}

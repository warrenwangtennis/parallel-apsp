package main

import (
	// "fmt"
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
				adj_mat[i][j] = INF
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
		assert(len(apsp[i]) == n)
	}

	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				apsp[i][j] = min(apsp[i][j], apsp[i][k] + apsp[k][j])				
			}
		}
	}
}

func main() {
	n, m, adj = read_input("t2.txt")

	load_adj_mat()

	solve()

	write_output("./floyd_seq_out.txt", apsp)
}
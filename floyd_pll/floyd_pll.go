package main

import (
	"fmt"
	"../util"
	"sync"
	"flag"
	"time"
	"os/exec"
	"bytes"
)

var n, m, nthreads int
var adj []map[int]int
var adj_mat, apsp, apsp_tmp [][]int

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

func worker_update(id int, k int, wg *sync.WaitGroup) {
	for ij := id; ij < n * n; ij += nthreads {
		i, j := ij / n, ij % n
		apsp_tmp[i][j] = util.Min(apsp[i][j], apsp[i][k] + apsp[k][j])
	}
	wg.Done()
}

func worker_copy(id int, wg *sync.WaitGroup) {
	for ij := id; ij < n * n; ij += nthreads {
		i, j := ij / n, ij % n
		apsp[i][j] = apsp_tmp[i][j]
	}
	wg.Done()
}

func solve() {
	apsp = util.Make_mat(n, n)
	apsp_tmp = util.Make_mat(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			apsp[i][j] = adj_mat[i][j]
		}
	}

	for k := 0; k < n; k++ {
		var wg_update, wg_copy sync.WaitGroup

		for i := 0; i < nthreads; i++ {
			wg_update.Add(1)
			go worker_update(i, k, &wg_update)
		}
		wg_update.Wait()
		for i := 0; i < nthreads; i++ {
			wg_copy.Add(1)
			go worker_copy(i, &wg_copy)
		}
		wg_copy.Wait()
	}
}

func main() {
	total_start := time.Now()

	var input_path, output_path string
	flag.IntVar(&nthreads, "t", 1, "num threads")
	flag.StringVar(&input_path, "i", "", "string-valued path to an input file")
	flag.StringVar(&output_path, "o", "./floyd_pll_out.txt", "string-valued path to an output file")
	flag.Parse()

	fmt.Println("nthreads", nthreads)

	n, m, adj = util.Read_input(input_path)

	load_adj_mat()

	solve_start := time.Now()
	solve()
	fmt.Println("solve time", time.Since(solve_start).Seconds())

	util.Write_output(output_path, apsp)

	fmt.Println("e2e time", time.Since(total_start).Seconds())

	cmd := exec.Command("diff", "-sq", output_path, input_path[:len(input_path) - 4] + "_key.txt")
	// fmt.Println(output_path)
	// fmt.Println(input_path[:len(input_path) - 4] + "_key.txt")
	var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
	// fmt.Println(err)
	util.Assert(err == nil)
	fmt.Print(out.String())
}
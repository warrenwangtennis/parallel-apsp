package main

import (
	"fmt"
	"sync"
	"flag"
	// "math"
	"time"
	"../util"
	"os/exec"
	"bytes"
)

type Info struct {
	i int
	j int
	d int
}

var n, m, T, sqrtT, subsize int
var adj []map[int]int
var adj_mat, apsp [][]int
var chans [][]chan Info

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

func idx_range(id int) (int, int) {
	return util.Min(id * subsize, n), util.Min((id + 1) * subsize, n)
}

// func idx_range2(idi, idj int) (int, int, int, int) {
// 	return idi * sz, (idi + 1) * sz, idj * sz, (idj + 1) * sz
// }


func worker(idi, idj, k int, wg *sync.WaitGroup) {
	li, ri := idx_range(idi)
	lj, rj := idx_range(idj)
	// fmt.Println("here", li, ri, lj, rj)
	if li <= k && k < ri {
		for z := 0; z < sqrtT; z++ {
			for j := lj; j < rj; j++ {
				chans[z][idj] <- Info{k, j, apsp[k][j]}
			}
		}
	}
	if lj <= k && k < rj {
		for z := 0; z < sqrtT; z++ {
			for i := li; i < ri; i++ {
				chans[idi][z] <- Info{i, k, apsp[i][k]}
			}
		}
	}

	// fmt.Println(lj, rj)
	hori_bar := make([]int, rj - lj)
	vert_bar := make([]int, ri - li)
	for lv := 0; lv < (rj - lj) + (ri - li); lv++ {
		info := <- chans[idi][idj]
		i, j, d := info.i, info.j, info.d
		// fmt.Println("received", i,  j, d)
		util.Assert(i == k || j == k)
		// if (i < li || i >= ri) && j == k {
		// 	fmt.Println(k, li, ri, lj, rj, i, j, d)
		// }
		// if (j < lj || j >= rj) && i == k {
		// 	fmt.Println(k, li, ri, lj, rj, i, j, d)
		// }
		if i == k && (lj <= j && j < rj) {
			// fmt.Println(j - lj, rj - lj)
			hori_bar[j - lj] = d 
		}
		if j == k && (li <= i && i < ri) {
			// fmt.Println(i - li, ri - li)
			vert_bar[i - li] = d
		}
	}

	for i := li; i < ri; i++ {
		for j := lj; j < rj; j++ {
			apsp[i][j] = util.Min(apsp[i][j], hori_bar[j - lj] + vert_bar[i - li])
		}
	}

	wg.Done()
}

func solve() {
	apsp = util.Make_mat(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			apsp[i][j] = adj_mat[i][j]
		}
	}

	for idi := 0; idi < sqrtT; idi++ {
		chans = append(chans, make([]chan Info, sqrtT))
		for idj := 0; idj < sqrtT; idj++ {
			chans[idi][idj] = make(chan Info, subsize * 5)
		}
	}
	

	for k := 0; k < n; k++ {
		var wg sync.WaitGroup
		for idi := 0; idi < sqrtT; idi++ {
			for idj := 0; idj < sqrtT; idj++ {
				wg.Add(1)
				go worker(idi, idj, k, &wg)
			}
		}
		wg.Wait()
	}
}

func main() {
	total_start := time.Now()

	var input_path, output_path string
	flag.IntVar(&T, "t", 1, "num threads")
	flag.StringVar(&input_path, "i", "", "string-valued path to an input file")
	flag.StringVar(&output_path, "o", "./floyd_pll2_out.txt", "string-valued path to an output file")
	flag.Parse()

	n, m, adj = util.Read_input(input_path)

	if T > n*n {
		T = n*n
	}
	sqrtT = 1
	for ; (sqrtT+1)*(sqrtT+1) <= T; sqrtT++ {
		
	}
	T = sqrtT * sqrtT
	subsize = (n + sqrtT - 1) / sqrtT
	
	fmt.Println(input_path)
	fmt.Println("nthreads", T)

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
package main

import (
	"fmt"
	"flag"
	"../util"
	"sync"
)

var n, m, nthreads int
var adj []map[int]int
var adj_mat, apsp, apsp_tmp [][]int

var dist []int

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

func bellman_ford_plus(source int) bool{
	dist = make([]int, n+1)
	for i := 0; i < n+1; i++ {
		if i != source {
			dist[i] = util.INF
		} else {
			dist[i] = 0
		}
	}
	//Calculate shortest distances
	for index := 0; index < n-1 +1; index++ {
		for u := 0; u < n+1; u++ {
			if u != source {
				for v, d := range adj[u] {
					if dist[v] > dist[u] + d {
						dist[v] = dist[u] + d
					}
				}
				//Source node is the last edge for each node
				// if dist[source] > dist[u] + 0 {
				// 	dist[source] = dist[u] + 0
				// }
			} else {
				//For source node:
				for v := 0; v < n; v++ {
					if dist[v] > dist[u] + 0 {
						dist[v] = dist[u] + 0
					}
				}
			}
		}
		
	}
	//Report negative weight cycle
	for u := 0; u < n+1; u++ {
		if u != source {
			for v, d := range adj[u] {
				if dist[v] > dist[u] + d {
					return false
				}
			}
			if dist[source] > dist[u] + 0 {
				return false
			}
		} else {
			//For source node:
			for v := 0; v < n; v++ {
				if dist[v] > dist[source] + 0 {
					return false
				}
			}
		}
	}
	return true
}

func dijkstra_worker(neighbor int, d int, minNode int, distances *[]int, wg *sync.WaitGroup) {
	path := (*distances)[minNode] + d
	if path < (*distances)[neighbor] {
		(*distances)[neighbor] = path
	}
	wg.Done()
}

func dijkstra(id int, wg *sync.WaitGroup) {
	for index := id; index < n; index += nthreads {
		distances := make([]int, n)
		queue := make(map[int]bool)
		for i := 0; i < n; i++ {
			if i != index {
				distances[i] = util.INF
			} else {
				distances[i] = 0
			}
			queue[i] = true;
		}

		for len(queue) > 0 {
			min := util.INF
			minNode := util.INF
			for key, _ := range queue {
				if distances[key] <= min{
					min = distances[key]
					minNode = key
				}
			}
			delete(queue, minNode)

			var wg_dijkstra sync.WaitGroup
			for neighbor, d := range adj[minNode] {
				if _, ok := queue[neighbor]; ok {
					wg_dijkstra.Add(1)
					go dijkstra_worker(neighbor, d, minNode, &distances, &wg_dijkstra)
				}
			}
			wg_dijkstra.Wait()
		}

		for j := 0; j < n; j++ {
			apsp[index] = append(apsp[index], distances[j])				
		}
	}
	wg.Done()
}

func worker_reweigh(id int, wg *sync.WaitGroup) {
	for u := id; u < n; u += nthreads {
		for v := range adj[u] {
			adj[u][v] += dist[u] - dist[v]
		}
	}
	wg.Done()
}

func solve() {	
	if !bellman_ford_plus(n) {
		fmt.Println("Negative weight cycle")
		var wg_reweigh sync.WaitGroup
		//Negative weight cycle, reweigh edges
		for i := 0; i < nthreads; i++ {
			wg_reweigh.Add(1)
			go worker_reweigh(i, &wg_reweigh)
		}
		wg_reweigh.Wait();
	}
	
	fmt.Println(dist)

	apsp = make([][]int, n)

	var wg_dij sync.WaitGroup
	for i := 0; i < nthreads; i++ {
		//Run Dijkstra's on each vertex
		wg_dij.Add(1)
		go dijkstra(i, &wg_dij)
	}
	wg_dij.Wait()
}

func main() {
	var input_path string
	flag.IntVar(&nthreads, "t", 1, "num threads")
	flag.StringVar(&input_path, "i", "", "string-valued path to an input file")
	flag.Parse()
	n, m, adj = util.Read_input(input_path)

	load_adj_mat()

	solve()

	util.Write_output("./johnson_pll_out.txt", apsp)
}
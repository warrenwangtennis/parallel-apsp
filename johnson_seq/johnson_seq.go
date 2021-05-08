package main

import (
	"fmt"
	"flag"
	"../util"
)

var n, m int
var adj []map[int]int
var adj_mat, apsp [][]int

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

// func insert_node() {
// 	for u := 0; u < n; u++ {
// 		adj[u][n] = 0
// 		adj[n][u] = 0
// 	}
// }

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
		for u := 0; u < n; u++ {
			for v, d := range adj[u] {
				if dist[v] > dist[u] + d {
					dist[v] = dist[u] + d
				}
			}
			//Source node is the last edge for each node
			if dist[source] > dist[u] + 0 {
				dist[source] = dist[u] + 0
			}
		}
		//For source node:
		//for v, d := range adj[n]
		//this array is all 0
	}
	//Report negative weight cycle
	for u := 0; u < n; u++ {
		for v, d := range adj[u] {
			if dist[v] > dist[u] + d {
				return false
			}
		}
		if dist[source] > dist[u] + 0 {
			return false
		}
	}
	//For source node:
	for v := 0; v < n; v++ {
		if dist[v] > dist[source] + 0 {
			return false
		}
	}

	return true
}

func dijkstra(source int) []int{
	distances := make([]int, n)
	queue := make(map[int]bool)
	for i := 0; i < n; i++ {
		if i != source {
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

		for neighbor, d := range adj[minNode] {
			if _, ok := queue[neighbor]; ok {
				path := distances[minNode] + d
				if path < distances[neighbor] {
					distances[neighbor] = path
				}
			}
		}
	}

	return distances
}


func solve() {
	//insert_node()

	if !bellman_ford_plus(n) {
		//Negative weight cycle, reweigh edges
		for u := 0; u < n; u++ {
			for v := range adj[u] {
				adj[u][v] += dist[u] - dist[v]
			}
		}
	}

	fmt.Println(dist)

	apsp = make([][]int, n)
	for i := 0; i < n; i++ {
		//Run Dijkstra's on each vertex
		ans :=  dijkstra(i)
		for j := 0; j < n; j++ {
			apsp[i] = append(apsp[i], ans[j])				
		}
	}
}

func main() {
	var input_path string
	flag.StringVar(&input_path, "i", "", "string-valued path to an input file")
	flag.Parse()
	n, m, adj = util.Read_input(input_path)

	load_adj_mat()

	solve()

	util.Write_output("./johnson_seq_out.txt", apsp)
}
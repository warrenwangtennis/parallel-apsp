package main

import (
	"fmt"
)

func main() {
	n, m, adj := read_input("t1.txt")

	fmt.Println(n, m, adj)

	var mat = [][]int{
		{0, 1, 2, 3},
		{4, 5, 6, 7},
	}

	write_output("./floyd_seq_out.txt", mat)
}
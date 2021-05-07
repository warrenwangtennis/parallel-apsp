floyd_seq : floyd_seq.go util.go
	go build -o floyd_seq

johnson_seq : johnson_seq.go
	go build -o johnson_seq util/util.go
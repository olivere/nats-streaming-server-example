.PHONY: pub sub

default: pub sub

pub:
	go build -o pub pub.go
	go build -o sub sub.go

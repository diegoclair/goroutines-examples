package main

import (
	"fmt"
	"runtime"
)

func main() {

	//go run *.go
	ExampleMutex()

}

func exampleCPU() {
	//how to see how many CPU we have
	fmt.Println(runtime.NumCPU())
}

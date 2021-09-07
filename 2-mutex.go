package main

import (
	"fmt"
	"runtime"
	"sync"
)

var contador int

func ExampleMutex() {
	//how to solve race condition problem with mutex
	mu := sync.Mutex{}

	fmt.Println("CPUs:", runtime.NumCPU())
	fmt.Println("Goroutines:", runtime.NumGoroutine())

	totalGoRoutines := 15

	wg := sync.WaitGroup{}
	wg.Add(totalGoRoutines)

	for i := 0; i < totalGoRoutines; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			v := contador
			runtime.Gosched() //yield
			v++
			contador = v
			mu.Unlock()
		}()
		fmt.Println("Goroutines:", runtime.NumGoroutine())
	}

	fmt.Println("Goroutines:", runtime.NumGoroutine())
	wg.Wait()
	fmt.Println(contador)
}

func raceCondition() {
	//race condition problem, we can execute go run -race main,go
	fmt.Println("CPUs:", runtime.NumCPU())
	fmt.Println("Goroutines:", runtime.NumGoroutine())

	totalGoRoutines := 1000

	wg := sync.WaitGroup{}
	wg.Add(totalGoRoutines)

	for i := 0; i < totalGoRoutines; i++ {
		go func() {
			defer wg.Done()
			v := contador
			runtime.Gosched() //yield
			v++
			contador = v
		}()
		fmt.Println("Goroutines:", runtime.NumGoroutine())
	}

	fmt.Println("Goroutines:", runtime.NumGoroutine())
	wg.Wait()
	fmt.Println(contador)
}

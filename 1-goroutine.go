package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
)

func GoRoutines() {
	Example3()
}

//=========================================================
func Example1() {
	//running two functions concurrently using the sync package
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		func1()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		func2()
	}()

	fmt.Println("number of goroutines right now: ", runtime.NumGoroutine()) //how to see how many goroutines is running

	wg.Wait()
}

func func1() {
	for i := 0; i < 10; i++ {
		fmt.Println("func1:", i)
	}
}

func func2() {
	for i := 0; i < 10; i++ {
		fmt.Println("func2:", i)
	}
}

//=========================================================
type employee struct {
	Name   string
	Age    int
	Salary float64
}

func Example2() {
	//this example do not cancel other go routines with some routine has an error. See example 3 with cancel approach
	t := []employee{}

	errChan := make(chan error)
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int, errChan chan error) {
			defer wg.Done()
			tt, err := getEmployee(i)
			if err != nil {
				errChan <- err
				return
			}
			t = append(t, tt)
			fmt.Println(tt)

		}(i, errChan)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for e := range errChan {
		if e != nil {
			log.Print("Deu erro: ", e)
		}
	}

	fmt.Println(t)
}

func Example3() {
	//cancel all go routines if we had some error in one of goroutines
	t := []employee{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error)

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				fmt.Println("OPA")
				return // Error somewhere, terminate
			default: // Default is must to avoid blocking
			}
			tt, err := getEmployee(i)
			if err != nil {
				cancel()
				errChan <- err
				return
			}
			t = append(t, tt)
			fmt.Println("Log externo", i, tt)

		}(i)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for e := range errChan {
		if e != nil {
			log.Fatal("Deu erro", e.Error())
		}
	}

	fmt.Println(t)
}
func getEmployee(i int) (e employee, err error) {

	e = employee{
		Name:   "Diego" + strconv.Itoa(i),
		Age:    27,
		Salary: 10000.00,
	}

	fmt.Println("Log interno", i, e)
	if i == 1 {
		return e, fmt.Errorf("Erro no processo " + strconv.Itoa(i))
	}

	return e, err
}

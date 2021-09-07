package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Channels() {
	// channels()
	// channelsWithRange()
	//channelsWithSelect()
	//channelsWithSelectAndQuit()
	convergeChannels()
}

//=========================================================
func channels() {
	channel := make(chan int)
	go func() {
		channel <- 42
	}()
	fmt.Println("channel1: ", <-channel)

	/*
		we can't do this. We need go routines to add or receive channel values, not in the same thread
		canal := make(chan int)
		canal <- 42
		fmt.Println(<-canal)
	*/

	//if we put the size of the channel, we create a buffer, so we can add values in the same thread
	//try never use this method
	channel2 := make(chan int, 1)
	channel2 <- 21
	fmt.Println("channel2: ", <-channel2)

	/*send and receive channels*/

	channel3 := make(chan int)
	go send(channel3)
	receive(channel3)
}
func send(s chan<- int) {
	s <- 55
}
func receive(r <-chan int) {
	fmt.Println("Value removed from channel:", <-r)
}

//=========================================================
func channelsWithRange() {

	c := make(chan int)
	go loop(8, c)
	for v := range c {
		fmt.Println("Final range value: ", v)
	}
}
func loop(t int, s chan<- int) {
	defer close(s)
	for i := 0; i < t; i++ {
		s <- i
	}
}

//=========================================================
func channelsWithSelect() {

	a := make(chan int)
	b := make(chan int)

	x := 100

	go func(x int) {
		for i := 0; i < x; i++ {
			a <- i
		}
	}(x / 2)
	go func(x int) {
		for i := 0; i < x; i++ {
			b <- i
		}
	}(x / 2)

	for i := 0; i < x; i++ {
		//we can use the select to receive values from any channels
		select {
		case v := <-a:
			fmt.Println("Received from A: ", v)
		case v := <-b:
			fmt.Println("Received from B: ", v)
		}
	}
}

//=========================================================
func channelsWithSelectAndQuit() {

	even := make(chan int)
	odd := make(chan int)
	quit := make(chan bool)

	total := 5

	go sendNum(total, even, odd, quit)
	receiveSelect(even, odd, quit)
}
func sendNum(total int, even, odd chan int, quit chan bool) {
	for i := 0; i < total; i++ {
		if i%2 == 0 {
			even <- i
		} else {
			odd <- i
		}
	}
	close(even)
	close(odd)
	quit <- true
}
func receiveSelect(even, odd chan int, quit chan bool) {
	for {
		select {
		case v, ok := <-even:
			if ok {
				fmt.Println("The number", v, "is even.")
			}
		case v, ok := <-odd:
			if ok {
				fmt.Println("The number", v, "is odd")
			}
		case <-quit:
			return
		}
	}
}

//=========================================================
func convergeChannels() {

	channel := converge(work("first "), work("second"))
	for x := 0; x < 16; x++ {
		fmt.Println(<-channel)
	}
}
func work(s string) chan string {
	channel := make(chan string)

	go func(s string, c chan string) {
		for i := 1; ; i++ {
			c <- fmt.Sprintf("Function %v say: %v", s, i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1e3))) //1e3 -> 1000
		}
	}(s, channel)
	return channel
}
func converge(x, y chan string) chan string {
	newChannel := make(chan string)
	go func() {
		for {
			newChannel <- <-x
		}
	}()
	go func() {
		for {
			newChannel <- <-y
		}
	}()
	return newChannel
}

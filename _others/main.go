package main

import (
	"fmt"
)

func simpleNumberPrinter() {
	numCh := make(chan int)
	oddEnder := make(chan bool)
	evenEnder := make(chan bool)
	// evenChan := make(chan int)

	go func(ch chan int) {
		for i := 1; i <= 1000; i += 2 {
			ch <- i
			res := <-ch
			fmt.Println(res)
		}
		oddEnder <- true
	}(numCh)

	go func(ch chan int) {
		for i := 2; i <= 1000; i += 2 {
			res := <-ch
			fmt.Println(res)
			ch <- i
		}
		evenEnder <- true
	}(numCh)

	oddEnded := false
	evenEnded := false

	for {
		if evenEnded && oddEnded {
			close(numCh)
			break
		}
		select {
		case <-oddEnder:
			oddEnded = true
			close(oddEnder)
		case <-evenEnder:
			evenEnded = true
			close(evenEnder)
		}
	}

	return
}

/*
func deleteBookById(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	//
}
*/

func main() {
	simpleNumberPrinter()
}

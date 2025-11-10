package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func writer(id int, c <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Writer %d start.\n", id)
	filePath := "ints.txt"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}

	writer := bufio.NewWriter(file)

	for x := range c {
		idAsString := strconv.Itoa(id)
		xAsString := strconv.Itoa(x)
		fmt.Printf("Writer %d attempt to write: %s.\n", id, xAsString)
		_, err := writer.WriteString("Writer " + idAsString + ", " + xAsString + "\n")
		time.Sleep(time.Second)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}

	err2 := writer.Flush()
	if err2 != nil {
		fmt.Printf("Error during file flush: %v\n", err)
	}

	err3 := file.Close()
	if err3 != nil {
		fmt.Printf("Error during file close: %v\n", err)
	}
	fmt.Printf("Writer %d stop.\n", id)
}

func producer(id int, xs []int, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Producer %d start.\n", id)
	for _, x := range xs {
		fmt.Printf("Producer %d push: %d.\n", id, x)
		c <- x
		time.Sleep(time.Second)
	}
	fmt.Printf("Producer %d stop.\n", id)
}

func main() {
	fmt.Println("Hello, concur ints and file io.")

	// 1. Create channel for (native) ints
	c := make(chan int, 100)
	var wg0 sync.WaitGroup // 1x writer
	var wg sync.WaitGroup  // 2x producers

	// 2. Create 1 file writer (will act as sync mutex)
	wg0.Add(1) // increment goroutines count for wait group
	go writer(1, c, &wg0)

	// 3. Create 2 producers to push (native) ints to channel
	xs := []int{7, 2, 8, -9, 4, 0}
	fmt.Printf("input as xs(ints): %v\n", xs)
	numProducers := 2
	for i := 0; i < numProducers; i++ {
		wg.Add(1)
		if i == 0 {
			go producer(i, xs[:len(xs)/2], c, &wg)
		} else {
			go producer(i, xs[len(xs)/2:], c, &wg)
		}
	}

	// 4. wait for all producers to finish
	wg.Wait()
	close(c)

	// 5. wait for writer to finish
	wg0.Wait()

	fmt.Println("Good bye.")
}

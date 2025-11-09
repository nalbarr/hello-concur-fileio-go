package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func FileCloseCheck(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Received error for file close.:", err)
	}
}

func WriterFlushCheck(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Received error for file writer flush..:", err)
	}
}

func writer(id int, c <-chan int, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Goroutine %d start.\n", id)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	// NA.
	// - Remove linter warning
	// defer file.Close()
	// defer FileCloseCheck(file.Close())

	writer := bufio.NewWriter(file)
	// NA.
	// - Remove linter warning
	// defer writer.Flush()
	// defer WriterFlushCheck(writer.Flush())

	for x := range c {
		idAsString := strconv.Itoa(id)
		xAsString := strconv.Itoa(x)
		fmt.Printf("Goroutine %d attempt to write: %s.\n", id, xAsString)
		_, err := writer.WriteString("goroutine " + idAsString + ", " + xAsString + "\n")
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
	fmt.Printf("Goroutine %d stop.\n", id)
}

func main() {
	fmt.Println("Hello, concur ints and file io.")

	filePath := "ints.txt"
	xs := []int{7, 2, 8, -9, 4, 0}

	// 1. Create channel for (native) ints
	c := make(chan int)
	var wg sync.WaitGroup

	// 2. Create 2 goroutines as file writers
	for i := 1; i <= 3; i++ {
		wg.Add(1) // increment goroutines count for wait group
		go writer(i, c, filePath, &wg)
	}

	// 3. push (native) ints to channel, then close channel
	for _, x := range xs {
		c <- x
	}
	close(c)

	// 5. wait for goroutines to finish
	wg.Wait()

	fmt.Println("Good bye.")
}

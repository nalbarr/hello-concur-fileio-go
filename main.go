package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"strconv"
)

func writer(c <-chan int, filePath string, wg *sync.WaitGroup) {
	defer wg.Done() 

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for x := range c {
		xAsString := strconv.Itoa(x)
		_, err := writer.WriteString(xAsString + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
	fmt.Printf("Ints written to %s.\n", filePath)
}

func main() {
	fmt.Println("Hello, concur ints and file io.")

	c := make(chan int)
	var wg sync.WaitGroup

	filePath := "ints.txt"

	wg.Add(1)
	go writer(c, filePath, &wg)

	xs := []int{7, 2, 8, -9, 4, 0}

	// 1. push (native) ints to channel
	for _, x := range xs {
		c <- x 
	}

	// 2. close channel, signal no more data
	close(c)

	// 3. wait for goroutine to finish
	wg.Wait()

	fmt.Println("Good bye.")
}

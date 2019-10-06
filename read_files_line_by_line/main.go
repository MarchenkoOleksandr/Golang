package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type set struct {
	country string
	capital string
}

func main() {
	startTime := time.Now()
	countriesCh, capitalsCh, saveCh := make(chan string), make(chan string), make(chan set)

	go readFileLineByLine("Countries.txt", countriesCh)
	go readFileLineByLine("Capitals.txt", capitalsCh)
	go saveResultToFile("Set.txt", saveCh)

	for {
		country, capital := <-countriesCh, <-capitalsCh
		if country == "" && capital == "" {
			close(saveCh)
			fmt.Printf("\nA new file \"Set.txt\" is created.")
			break
		}
		saveCh <- set{country: country, capital: capital}
	}

	fmt.Println("\nThe program duration is", time.Since(startTime).Nanoseconds()/1000000, "milliseconds")
}

func readFileLineByLine(path string, ch chan string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		ch <- string(line)
	}
	close(ch)
}

func saveResultToFile(path string, ch chan set) {
	file, err := os.Create(path)

	if err != nil {
		fmt.Println("Unable to create file:", path, err)
		os.Exit(1)
	}
	defer file.Close()

	var index = 0
	for set := range ch {
		index++
		fmt.Printf("%d. %s - %s\n", index, set.country, set.capital)
		_, err := file.WriteString(set.country + " - " + set.capital + "\r\n")

		if err != nil {
			fmt.Println("Unable to write into file:", path, err)
			os.Exit(1)
		}
	}
}
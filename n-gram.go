package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Open the stops.txt file
	file, err := os.Open("stops.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the first line (header) and discard it
	scanner.Scan()

	// Loop through the remaining lines and split the stop names into n-grams
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		stopName := fields[2]

		// Split the stop name into n-grams
		n := 3 // n-gram size
		grams := make([]string, 0)
		for i := 0; i < len(stopName)-n+1; i++ {
			gram := stopName[i : i+n]
			grams = append(grams, gram)
		}

		fmt.Printf("%s: %v\n", stopName, grams)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

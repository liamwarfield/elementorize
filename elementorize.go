package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func loadElements() []string {
	// Open file and create scanner on top of it
	file, err := os.Open("elementlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	elements := make([]string, 0)
	// Default scanner is bufio.ScanLines. Lets use ScanWords.
	// Could also use a custom function of SplitFunc type

	// Scan for next token.
	for {
		success := scanner.Scan()
		if success == false {
			// False on error or EOF. Check error
			err = scanner.Err()
			if err == nil {
				log.Println("Scan completed and reached EOF")
				break
			} else {
				log.Fatal(err)
			}
		}
		fmt.Println("found:", scanner.Text())
		elements = append(elements, scanner.Text())
	}
	fmt.Println(elements)
	// Get data from scan with Bytes() or Text()

	// Call scanner.Scan() again to find next token
	return elements
}
func loadWords(c chan string) {
	// Open file and create scanner on top of it
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	// Default scanner is bufio.ScanLines. Lets use ScanWords.
	// Could also use a custom function of SplitFunc type

	// Scan for next token.
	for {
		success := scanner.Scan()
		if success == false {
			// False on error or EOF. Check error
			err = scanner.Err()
			if err == nil {
				log.Println("Scan completed and reached EOF")
				break
			} else {
				log.Fatal(err)
			}
		}
		fmt.Println("found:", scanner.Text())
		c <- scanner.Text()
	}
	// Get data from scan with Bytes() or Text()

	// Call scanner.Scan() again to find next token
	return
}
func elementorize(elements *[]string, words chan string, output chan string) {
	for {
		select {
		case word := <-words:
			word = strings.ToLower(word)
			fmt.Println("Processing: ", word)
			sbstr := ""
			output := ""
			for _, letter := range word {
				sbstr += string(letter)
				for _, element := range *elements {
					//fmt.Println(sbstr, strings.ToLower(element))
					if sbstr == strings.ToLower(element) {
						sbstr = ""
						output += "[" + element + "]"
						//fmt.Println(output)
					} else if len(sbstr) > 2 {
						output = ""
						break
					}
				}
				//fmt.Println(sbstr)
			}
			if output != "" {
				fmt.Println("processed: ", output)
			}
		default:
			fmt.Println("somethingbad")
		}
	}
}
func main() {
	wordPipeline := make(chan string, 100000000)
	outputPipeline := make(chan string, 10000)
	elements := make([]string, 0)

	elements = loadElements()
	go loadWords(wordPipeline)
	for i := 0; i < 2; i++ {
		go elementorize(&elements, wordPipeline, outputPipeline)
	}
	for {
	}
}

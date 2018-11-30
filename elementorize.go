package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type SafeClose struct {
	c   bool
	mux sync.Mutex
}

func loadElements() []string {
	// Open file and create scanner on top of it
	file, err := os.Open("elementlist.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	elements := make([]string, 0)
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
	return elements
}

func loadWords(c chan string) {
	// Open file and create scanner on top of it
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for {
		success := scanner.Scan()
		if success == false {
			// False on error or EOF. Check error
			err = scanner.Err()
			if err == nil {
				log.Println("Scan completed and reached EOF")
				close(c)
				break
			} else {
				log.Fatal(err)
			}
		}
		c <- scanner.Text()
	}
	return
}

func elementorize(elements *[]string, closed *SafeClose, words chan string, output chan string) {
	for {
		select {
		case word, ok := <-words:
			if !ok {
				if !closed.c {
					closed.mux.Lock()
					close(output)
					closed.c = true
					closed.mux.Unlock()
				}
				return
			}
			word = strings.ToLower(word)
			sbstr := ""
			outputstr := ""
			for _, letter := range word {
				sbstr += string(letter)
				for _, element := range *elements {
					if sbstr == strings.ToLower(element) {
						sbstr = ""
						outputstr += "[" + element + "]"
					} else if len(sbstr) > 2 {
						outputstr = ""
						break
					}
				}
			}
			if outputstr != "" && sbstr == "" {
				output <- word + " " + outputstr
			}
		default:
			fmt.Println("somethingbad")
		}
	}
}

func main() {
	wordPipeline := make(chan string, 10000)
	outputPipeline := make(chan string, 10000)
	elements := make([]string, 0)
	closed := SafeClose{c: false}

	elements = loadElements()
	go loadWords(wordPipeline)
	for i := 0; i < 4; i++ {
		go elementorize(&elements, &closed, wordPipeline, outputPipeline)
	}
	for i := range outputPipeline {
		fmt.Println("Elementalized", i)
	}
}

package main

import (
    "os"
    "log"
    "fmt"
    "bufio"
)

func loadElements() map[string]int {
    // Open file and create scanner on top of it
    file, err := os.Open("elementlist.txt")
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)
    elements := make(map[string]int)
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
        elements[scanner.Text()] += 1
    }
    fmt.Println(elements)
    // Get data from scan with Bytes() or Text()

    // Call scanner.Scan() again to find next token
    return elements
}
func loadWords (c chan string) {
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
func main () {
    wordPipeline := make(chan string, 10000)
    loadElements()
    go loadWords(wordPipeline)
    for{}
}

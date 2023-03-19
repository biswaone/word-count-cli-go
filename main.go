package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// Define flags
	wc := flag.Bool("w", false, "count words")
	lc := flag.Bool("l", false, "count lines")
	cc := flag.Bool("c", false, "count characters")

	// Parse flags
	flag.Parse()

	// Count requested statistics
	countStats(*wc, *lc, *cc, flag.Args())

}

func countStats(wc bool, lc bool, cc bool, args []string) {
	for _, arg := range args {
		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open file %q: %v\n", arg, err)
			continue // continue with the next file
		}
		defer file.Close()

		wordCount := 0
		lineCount := 0
		charCount := 0

		fmt.Println(*file)

		if cc {
			charater_scanner := bufio.NewScanner(file)
			//fmt.Println(*charater_scanner)
			charater_scanner.Split(bufio.ScanRunes)
			for charater_scanner.Scan() {
				charCount++
			}
		}
		//set file pointer to beginning of the file
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}
		// The line count here and using wc in linux may differ by 1
		// depending upon if a new empty line is present or not at the end of file
		// The wc counts the number of newline characters '\n' present
		// If the lastline doesn't have a newline character then it won't
		// count it as a line, whereas ScanLine function counts the last non-empty
		// line even if it has no newline
		if lc {
			line_scanner := bufio.NewScanner(file)
			line_scanner.Split(bufio.ScanLines)
			for line_scanner.Scan() {
				lineCount++
			}

		}
		//set file pointer to beginning of the file
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		if wc {
			word_scanner := bufio.NewScanner(file)
			//fmt.Println(&word_scanner)
			word_scanner.Split(bufio.ScanWords)
			for word_scanner.Scan() {
				wordCount++
			}

		}

		// print the results
		fmt.Printf("Word count: %d\n", wordCount)
		fmt.Printf("Line count: %d\n", lineCount)
		fmt.Printf("Character count: %d\n", charCount)
	}
}

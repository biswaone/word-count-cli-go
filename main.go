package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	// Define flags
	wc := flag.Bool("w", false, "count words")
	lc := flag.Bool("l", false, "count lines")
	cc := flag.Bool("c", false, "count characters")

	// Parse flags
	flag.Parse()

	// check if all the flags are false make it true
	if !*wc && !*lc && !*cc {
		*wc, *lc, *cc = true, true, true

	}

	// if no arguments are given take input from stdin
	if flag.NArg() < 1 {
		s := readStdin()
		lineCount, wordCount, charCount := countStatsfromStdin(*wc, *lc, *cc, s)
		printStats(lc, wc, cc, lineCount, wordCount, charCount, "")
		// if file given in argument
	} else {

		totalLineCount, totalWordCount, totalCharCount := 0, 0, 0

		for _, arg := range flag.Args() {
			file, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open file %q: %v\n", arg, err)
				continue // continue with the next file
			}
			defer file.Close()

			lineCount, wordCount, charCount := countStats(*wc, *lc, *cc, file)

			// print stats for a file
			printStats(lc, wc, cc, lineCount, wordCount, charCount, arg)

			if len(flag.Args()) > 1 {
				totalLineCount += lineCount
				totalWordCount += wordCount
				totalCharCount += charCount
			}

		}

		if len(flag.Args()) > 1 {
			// print the total count
			printStats(lc, wc, cc, totalLineCount, totalWordCount, totalCharCount, "total")
		}
	}

}

func countStats(wc bool, lc bool, cc bool, file *os.File) (int, int, int) {

	var err error
	lineCount := 0
	wordCount := 0
	charCount := 0

	if cc {
		charater_scanner := bufio.NewScanner(file)
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
		word_scanner.Split(bufio.ScanWords)
		for word_scanner.Scan() {
			wordCount++
		}

	}

	return lineCount, wordCount, charCount

}

func countStatsfromStdin(wc bool, lc bool, cc bool, s string) (int, int, int) {
	lineCount := 0
	wordCount := 0
	charCount := 0

	if cc {
		charCount = utf8.RuneCountInString(s)
	}

	if lc {
		lineCount = strings.Count(s, "\n")
	}

	if wc {
		wordCount = len(strings.Fields(s))
	}

	return lineCount, wordCount, charCount

}

func readStdin() string {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.Equal(line, []byte("\x04")) {
			break
		}
		buf.Write(line)
		buf.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return buf.String()
}

func printStats(lc, wc, cc *bool, lineCount, wordCount, charCount int, filename string) {
	output := ""
	if *lc {
		output += fmt.Sprintf(" %d ", lineCount)
	}

	if *wc {
		output += fmt.Sprintf(" %d ", wordCount)
	}
	if *cc {
		output += fmt.Sprintf(" %d ", charCount)
	}

	// print the results
	fmt.Println(output, filename)
}

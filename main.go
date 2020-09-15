package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	newLine = '\n'
	retrn   = '\r'
)

func main() {
	var (
		// relPath is the relative path (taken from the working directory) to the file to `tail`.
		relPath = flag.String("f", "", "file path")
		// c is the number of lines to print.
		c int
	)
	flag.IntVar(&c, "n", 10, "numer of lines")
	flag.Parse()

	// basePath is the absolute path of the working directory.
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}
	filePath := fmt.Sprintf("%s/%s", basePath, *relPath)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}
	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		log.Fatalf("Failed reading meta: %s", err)
	}

	var (
		fileSize = stats.Size()
		// lineCount is the number of lines processed.
		lineCount = 0
		// cursor is the representation of the read offset in the file.
		cursor int64 = 0
		// line is a buffer into which lines are written.
		line bytes.Buffer
		// linesBuf is a slice of lines which stores the text to be outputted.
		linesBuf = make([]string, c)
	)

	for lineCount < c {
		// Decreases the cursor (as it goes backwards) and sets the offset for the new read relative to the end.
		cursor--
		f.Seek(cursor, io.SeekEnd)

		// Reads one byte from the file.
		ch := make([]byte, 1)
		f.Read(ch)

		// If return '\r' or newline '\n' are found the loop has completed a line.
		//
		// The line bytes are reversed as it has been traversed backwards.
		// Then the resulting bytes are stored to linesBuf as a string and the buffer reseted.
		if cursor != -1 && (ch[0] == newLine || ch[0] == retrn) {
			lineCount++
			b := line.Bytes()
			for i := len(b)/2 - 1; i >= 0; i-- {
				opp := len(b) - 1 - i
				b[i], b[opp] = b[opp], b[i]
			}
			linesBuf[c-lineCount] = line.String()
			line.Reset()
			continue
		}

		// Write read byte to buffer
		line.WriteByte(ch[0])

		// Exit loop if the cursor has reached the beggining of the file
		if cursor == -fileSize {
			break
		}
	}

	for i, l := range linesBuf {
		if i == lineCount-1 && strings.HasSuffix(l, "\n") {
			// Trim newline from the last line
			fmt.Printf("%s", l)
		} else {
			fmt.Printf("%s\n", l)
		}
	}
}

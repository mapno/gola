package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	newLine = '\n'
	retrn    = '\r'
)

func main() {
	var (
		relPath = flag.String("file", "", "file path")
		c       int
	)
	flag.IntVar(&c, "c", 10, "numer of lines")
	flag.Parse()

	basePath, _ := os.Getwd()
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
		fileSize       = stats.Size()
		count          = 0
		cursor   int64 = 0
		line     bytes.Buffer
		lineBuf  = make([]string, c)
	)

	for count < c {
		cursor--
		f.Seek(cursor, io.SeekEnd)

		ch := make([]byte, 1)
		f.Read(ch)

		if cursor != -1 && (ch[0] == newLine || ch[0] == retrn) {
			count++
			b := line.Bytes()
			for i := len(b)/2 - 1; i >= 0; i-- {
				opp := len(b) - 1 - i
				b[i], b[opp] = b[opp], b[i]
			}
			lineBuf[c-count] = line.String()
			line.Reset()
			continue
		}

		line.WriteByte(ch[0])

		if cursor == -fileSize {
			break
		}
	}

	b := line.Bytes()
	for i := len(b)/2 - 1; i >= 0; i-- {
		opp := len(b) - 1 - i
		b[i], b[opp] = b[opp], b[i]
	}
	fmt.Printf("%s", line.String())

	for _, l := range lineBuf {
		fmt.Printf("%s\n", l)
	}
}

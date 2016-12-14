package main

import (
	"bufio"
	"crypto/sha1"
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("orduniq: ")
	bufoutSize := flag.Int("o", 1024, "output buffer size in `bytes`")
	bufinSize := flag.Int("i", 1024, "input buffer size in `bytes`")
	flag.Parse()

	var inputs []io.Reader

	argv := flag.Args()
	if len(argv) == 0 {
		inputs = []io.Reader{os.Stdin}
	} else {
		usedStdin := false
		for _, p := range argv {
			if p == "-" {
				if usedStdin {
					log.Panic("standard input specified more than once")
				}
				usedStdin = true

				inputs = append(inputs, os.Stdin)
				continue
			}

			fi, err := os.Open(p)
			if err != nil {
				panic(err)
			}
			defer fi.Close()

			inputs = append(inputs, fi)
		}
	}

	hashes := make(map[[sha1.Size]byte]struct{})
	input := bufio.NewReaderSize(io.MultiReader(inputs...), *bufinSize)

	bufout := bufio.NewWriterSize(os.Stdout, *bufoutSize)
	defer func() {
		ferr := bufout.Flush()
		if ferr != nil {
			log.Panic("unable to flush output buffer: ", ferr)
		}
	}()

	for {
		line, err := input.ReadBytes('\n')
		if err != nil {
			if err == io.EOF && len(line) == 0 {
				// Done
				return
			} else if err != io.EOF {
				log.Panic("error reading input line: ", err)
			}
		}

		data := line
		if err != nil {
			data = append(data, '\n')
		}
		sum := sha1.Sum(data)

		if _, ok := hashes[sum]; ok {
			continue
		}
		hashes[sum] = struct{}{}

		bufout.Write(data)
	}
}

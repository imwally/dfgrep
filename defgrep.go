package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sync/semaphore"
)

type Result struct {
	Path  string
	Line  string
	Index int
}

var (
	maxWorkers = runtime.GOMAXPROCS(0)
	sem        = semaphore.NewWeighted(int64(maxWorkers))
	out        = make(chan Result)
	ctx        = context.Background()
)

func SearchFile(search, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
		}

		index := SubStr(line, search)
		if index >= 0 {
			out <- Result{
				path,
				line,
				index,
			}
		}
	}

	return nil
}

func SearchPath(search, path string) error {
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if d.IsDir() {
			return nil
		}

		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}

		go func(search, path string) {
			defer sem.Release(1)
			err := SearchFile(search, path)
			if err != nil {
				log.Println(err)
			}
		}(search, path)

		return e
	})

	return err
}

// SubStr searches haystack for needle and returns the
// index of the first instance of needle, otherwise it
// returns -1.
func SubStr(haystack, needle string) int {
	// Needle too big, you'd see it.
	if len(needle) > len(haystack) {
		return -1
	}

	// The haystack is the needle, bruh.
	if needle == haystack {
		return 0
	}

	for i := 0; i < len(haystack); i++ {
		// Found starting character, so look for rest of substring.
		if haystack[i] == needle[0] {
			start := i
			for j := 0; j < len(needle); j++ {
				if haystack[i] == needle[j] {
					i++
					// Found substring.
					if j == len(needle)-1 {
						return start
					}
					// Hit end of haystack
					if i > len(haystack)-1 {
						return -1
					}
				} else {
					break
				}
			}
		}
	}

	return -1
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: <search string> path\n")
		os.Exit(1)
	}

	go func() {
		for result := range out {
			end := result.Index + 100
			if end > len(result.Line) {
				end = len(result.Line)
			}
			fmt.Printf("%s: %s", result.Path, result.Line[result.Index:end])
		}
	}()

	if err := SearchPath(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
}

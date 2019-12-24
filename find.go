package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
)

func main() {
	chunk := 1024 * 1024 * 10

	path := "TestData.txt"
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := int(fi.Size())

	mem, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}

	needles := make(chan int)
	var wg sync.WaitGroup
	done := false
	for i := 0; i < size; i += chunk {
		to := i + chunk
		if to > size {
			to = size
		}
		wg.Add(1)
		go func(from, to int) {
			for ; from < to; from += 2 {
				if mem[from] == '1' ||
					mem[from+1] == '1' {
					if mem[from+1] == '1' {
						needles <- from + 1
					} else {
						needles <- from
					}
					done = true
					close(needles)
				}
			}
			wg.Done()
		}(i, to)
		if done {
			break
		}
	}

	for needle := range needles {
		fmt.Printf("needle at %d\n", needle)
	}
	wg.Wait()
	err = syscall.Munmap(mem)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}

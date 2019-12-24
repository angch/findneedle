package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"unsafe"
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

	pattern := uint(0)
	for i := 0; i < 8; i++ {
		pattern = pattern<<8 + uint('0')
	}

	for i := 0; i < size; i += chunk {
		to := i + chunk
		if to > size {
			to = size
		}
		wg.Add(1)
		go func(from2, to2 int) {
			u := uintptr(unsafe.Pointer(&mem[from2]))
			z := u
			v := uintptr(unsafe.Pointer(&mem[to2-1]))
			for ; u < v-8; u += 8 {
				if *(*uint)(unsafe.Pointer(u)) != pattern {
					if !done {
						for start := int(u-z) + from2; start < to2; start++ {
							if mem[start] == '1' {
								needles <- start
								break
							}
						}
						done = true
						close(needles)
					}
					break
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

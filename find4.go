package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"io"
	"unsafe"
)


func find(mem []byte, offset int, length int, pattern uint64, needles chan <-int, wg *sync.WaitGroup, done *bool, pool *sync.Pool) {
	u := uintptr(unsafe.Pointer(&mem[0]))
	z := u
	v := uintptr(unsafe.Pointer(&mem[length-1]))

	for ; u < v-8; u += 8 {
		if *(*uint64)(unsafe.Pointer(u)) != pattern {
			if !*done {
				for start := int(u-z); start < length; start++ {
					if mem[start] == '1' {
						needles <- start+offset
						break
					}
				}
				*done = true
				close(needles)
			}
			break
		}
	}
	wg.Done()
	pool.Put(mem)
}


func main() {
	// Streaming chunks, in about 1 MB each, so we want to fill and start working even when the read hasn't finished.
	// Not too big, not too small.
	chunk := 1024 * 1024 

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

	var pool = sync.Pool{
		New: func() interface{} { return make([]byte, chunk) },
	 }

	 done := false
	 pattern := uint64(0)
	 for i := 0; i < 8; i++ {
		 pattern = pattern<<8 + uint64('0')
	 }

	 var wg sync.WaitGroup
	 needles := make(chan int)
	 for i := 0 ; i< size;{
		buf := pool.Get().([]byte) // reuse them, save some memory and allocs
		count, err := io.ReadAtLeast(f, buf, chunk)
		if err != nil {
			break
		}

		wg.Add(1)
		go find(buf, i, count, pattern, needles, &wg, &done, &pool)
		if done {
			break
		}
		i+= len(buf)
	}

	for needle := range needles {
		fmt.Printf("needle at %d\n", needle)
	}
	wg.Wait()
	f.Close()
}

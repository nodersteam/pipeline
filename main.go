package main

import (
	"fmt"
	"time"
)

const (
	BufferSize    = 10
	FlushInterval = 2 * time.Second
)

func generateNumbers(limit int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= limit; i++ {
			out <- i
		}
	}()
	return out
}
func filterNegative(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			if num >= 0 {
				out <- num
			}
		}
	}()
	return out
}

func filterNotDivisibleByThree(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			if num == 0 || num%3 == 0 {
				out <- num
			}
		}
	}()
	return out
}

type RingBuffer struct {
	buffer []int
	size   int
	head   int
	tail   int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]int, size),
		size:   size,
		head:   0,
		tail:   0,
	}
}

func (r *RingBuffer) Push(item int) {
	r.buffer[r.head] = item
	r.head = (r.head + 1) % r.size
	if r.head == r.tail {
		r.tail = (r.tail + 1) % r.size
	}
}

func (r *RingBuffer) Flush() []int {
	var result []int
	for r.tail != r.head {
		result = append(result, r.buffer[r.tail])
		r.tail = (r.tail + 1) % r.size
	}
	return result
}

func pipeline(limit int) <-chan int {
	buffer := NewRingBuffer(BufferSize)
	lastFlushTime := time.Now()

	out := make(chan int)

	numbers := generateNumbers(limit)

	numbers = filterNegative(numbers)

	numbers = filterNotDivisibleByThree(numbers)

	go func() {
		defer close(out)
		for num := range numbers {
			buffer.Push(num)

			if time.Since(lastFlushTime) > FlushInterval {
				for _, item := range buffer.Flush() {
					out <- item
				}
				lastFlushTime = time.Now()
			}
		}

		for _, item := range buffer.Flush() {
			out <- item
		}
	}()

	return out
}

func yield(item int) {
	fmt.Println(item)
}

func main() {
	for item := range pipeline(20) {
		yield(item)
	}
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type ConcurrentQueue struct {
	queue []int32
	mu    sync.Mutex
}

func NewConcurrentQueue() *ConcurrentQueue {
	return &ConcurrentQueue{
		queue: make([]int32, 0),
	}
}

func (q *ConcurrentQueue) Enqueue(val int32) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, val)
}

func (q *ConcurrentQueue) Dequeue() int32 {
	if len(q.queue) == 0 {
		panic("cannot dequeue from an empty queue")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	removed := q.queue[0]
	q.queue = q.queue[1:]
	return removed
}

func (q *ConcurrentQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue) == 0
}

func (q *ConcurrentQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}

var wg sync.WaitGroup

func main() {
	q := NewConcurrentQueue()

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			q.Enqueue(rand.Int31())
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(q.Size())

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			q.Dequeue()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(q.Size())
}

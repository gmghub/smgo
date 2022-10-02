package collector

import (
	"sync"
)

type RingBuffer struct {
	data    []interface{}
	size    int
	cursor  int
	written int
	lock    sync.Mutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]interface{}, size),
		size: size,
	}
}

func (b *RingBuffer) Put(e interface{}) {
	b.lock.Lock()
	b.data[b.cursor] = e
	b.cursor++
	if b.cursor >= b.size {
		b.cursor -= b.size
	}
	b.written++
	b.lock.Unlock()
}

func (b *RingBuffer) Get() interface{} {
	if b.written == 0 {
		return nil
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	cur := b.cursor - 1
	if cur < 0 {
		cur += b.size
	}
	return b.data[cur]
}

func (b *RingBuffer) GetN(n int) []interface{} {
	if n < 2 {
		n = 2
	}
	if n > b.size {
		n = b.size
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	if n > b.written {
		n = b.written
	}
	elems := make([]interface{}, n)
	cur := b.cursor
	for i := 0; i < n; i++ {
		cur--
		if cur < 0 {
			cur += b.size
		}
		elems[n-i-1] = b.data[cur]
	}

	return elems
}

package backup

import (
	"bytes"
	"sync"
)

type LimitedBuffer struct {
	buffer  bytes.Buffer
	maxSize int
	lock    sync.Mutex
	cond    *sync.Cond
	closed  bool
}

func NewLimitedBuffer(maxSize int) *LimitedBuffer {
	lb := &LimitedBuffer{
		maxSize: maxSize,
	}
	lb.cond = sync.NewCond(&lb.lock)
	lb.closed = false
	return lb
}

func (lb *LimitedBuffer) Close() {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	lb.cond.Broadcast()
	lb.closed = true
}

func (lb *LimitedBuffer) IsClosed() bool {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	return lb.closed
}

func (lb *LimitedBuffer) Write(data []byte) (int, error) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	remainingBytes := len(data)
	for lb.buffer.Len()+remainingBytes > lb.maxSize {
		if lb.maxSize-lb.buffer.Len() == 0 {
			lb.cond.Wait()
		}
		first := len(data) - remainingBytes

		last := min(len(data), first+(lb.maxSize-lb.buffer.Len()))

		n, _ := lb.buffer.Write(data[first:last])
		remainingBytes -= n
		lb.cond.Broadcast()
	}

	lb.buffer.Write(data[len(data)-remainingBytes:])
	lb.cond.Broadcast()

	return 0, nil
}

func (lb *LimitedBuffer) Len() int {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	return lb.buffer.Len()
}

func (lb *LimitedBuffer) Read(p []byte) (int, error) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	n, err := lb.buffer.Read(p)

	lb.cond.Broadcast()

	if lb.buffer.Len() == 0 && !lb.closed {
		lb.cond.Wait()
	}

	return n, err
}

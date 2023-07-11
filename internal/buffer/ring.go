package buffer

import (
	"container/ring"
	"fmt"
)

// RingBuffer is a ring buffer implementation, which stores keys and values
// and allows to retrieve the latest values in the order of insertion
// The implementation is not thread safe. It is the responsibility of the
// caller to ensure thread safety.
// When the buffer is full, the oldest key is removed. The buffer is
// implemented using a ring buffer and a map. The ring buffer is used to
// store the keys in the order of insertion. The map is used to get the
// values in O(1) time.
// When duplicate key is added, the value is updated and the element is
// moved to the head of the ring buffer.
type RingBuffer[K comparable, V any] struct {
	// use ring buffer, so we get elements in order of insertion
	// we only store keys as reference into the data map
	buffer *ring.Ring
	// use map to get elements in O(1)
	data map[K]V
}

// NewRingBuffer returns a new ring buffer
func NewRingBuffer[K comparable, V any](size int) *RingBuffer[K, V] {
	return &RingBuffer[K, V]{
		buffer: ring.New(size),
		data:   make(map[K]V, size),
	}
}

// Get returns the value associated with the key and whether it was found
func (rc *RingBuffer[K, V]) Get(key K) (V, bool) {
	value, ok := rc.data[key]
	return value, ok
}

// GetLatest returns number of the latest values added to the buffer
func (rc *RingBuffer[K, V]) GetLatest(n int) []V {
	// if n is greater than the number of inserted elements, return all elements
	currentLength := len(rc.data)
	if n > currentLength {
		n = currentLength
	}

	// get the latest n elements
	values := make([]V, n)
	for i, k := 0, rc.buffer; i < n; i, k = i+1, k.Prev() {
		values[i] = rc.data[k.Value.(K)]
	}

	return values
}

// Add adds the key value pair to the buffer
func (rc *RingBuffer[K, V]) Add(key K, value V) {
	// if the key already exists, make sure it is on the head
	if _, ok := rc.data[key]; ok {
		// if the key is not on the head, move it to the head
		if rc.buffer.Value.(K) != key {
			// move backwards and swap value with the next element
			// this is suboptimal approach, we could investigate
			// using `Move` and `Link` methods of the Ring
			previous := rc.buffer.Value.(K)
			end := false
			for k := rc.buffer.Prev(); ; k = k.Prev() {
				// if we reached to searched element, signal to break
				if k.Value.(K) == key {
					end = true
				}
				// swap the values
				curr := k.Value.(K)
				k.Value = previous
				previous = curr
				// break if we reached the element we were looking for
				if end {
					break
				}
			}
			// set the head to the key
			rc.buffer.Value = key
		}
		// update the value
		rc.data[key] = value
		return
	}

	// move the head of the ring to the next element
	rc.buffer = rc.buffer.Next()

	// if the ring is full, remove the oldest key
	if rc.buffer.Value != nil {
		delete(rc.data, rc.buffer.Value.(K))
	}

	// update the value
	rc.buffer.Value = key
	rc.data[key] = value
}

// Len returns the number of elements in the buffer
func (rc *RingBuffer[K, V]) Len() int {
	return len(rc.data)
}

// ringLen returns the number of elements in the buffer
// this method is costly, since it iterates over the ring,
// so it should only be used for testing
func (rc *RingBuffer[K, V]) ringLen() int {
	return rc.buffer.Len()
}

// printRingBuffer prints the buffer in the ring buffer
func (rc *RingBuffer[K, V]) printRingBuffer() {
	length := len(rc.data)
	for i, p := 0, rc.buffer; i < length; i, p = i+1, p.Prev() {
		fmt.Printf("%v \n", p.Value)
	}
}

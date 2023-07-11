package buffer

import (
	"testing"
)

// Test empty buffer has length 0 and no elements
func TestRingBuffer_Empty(t *testing.T) {
	rc := NewRingBuffer[int, string](3)

	// check the buffer is empty
	if rc.Len() != 0 {
		t.Error("expected length to be 0")
	}

	_, ok := rc.Get(1)
	if ok {
		t.Error("expected key to not be found")
	}
}

// Test element to buffer can be added and retrieved
func TestRingBuffer_AddAndGet(t *testing.T) {
	rc := NewRingBuffer[int, string](5)

	elements := []string{"one", "two", "three"}

	for i, element := range elements {
		// add element
		rc.Add(i, element)
		// check if element is returned
		value, ok := rc.Get(i)
		if !ok {
			t.Error("expected key to be found")
		}
		if value != element {
			t.Errorf("expected value '%s' to be returned", element)
		}
	}

	// assert length is 3
	if rc.Len() != 3 {
		t.Error("expected length to be 3")
	}
}

// Test latest elements can be retrieved
func TestRingBuffer_GetLatest(t *testing.T) {
	rc := NewRingBuffer[int, string](10)

	// add 10 elements
	elements := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	for i, element := range elements {
		rc.Add(i, element)
	}

	// check if latest 5 elements are returned
	retrieved := rc.GetLatest(5)
	if len(retrieved) != 5 {
		t.Error("expected 5 elements to be returned")
	}

	// check if elements are returned in correct order
	for i, element := range retrieved {
		if element != elements[len(elements)-i-1] {
			t.Error("expected elements to be returned in correct order")
		}
	}

	// getting more elements than inserted should return all elements
	retrieved = rc.GetLatest(20)
	if len(retrieved) != 10 {
		t.Error("expected 10 elements to be returned")
	}

	// assert length is 10
	if rc.Len() != 10 {
		t.Error("expected length to be 10")
	}
}

// Test elements are overwritten when buffer is full
func TestRingBuffer_ValuesAreRewritten(t *testing.T) {
	rc := NewRingBuffer[int, string](5)

	// add 10 elements
	elements := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	for i, element := range elements {
		rc.Add(i, element)
	}

	// assert first 5 elements are not present in buffer
	for i := 0; i < 5; i++ {
		_, ok := rc.Get(i)
		if ok {
			t.Error("expected key to not be found")
		}
	}

	// assert last 5 elements are present in buffer
	for i := 5; i < 10; i++ {
		value, ok := rc.Get(i)
		if !ok {
			t.Error("expected key to be found")
		}
		if value != elements[i] {
			t.Errorf("expected value '%s' to be returned", elements[i])
		}
	}

	// assert length is 5
	if rc.Len() != 5 {
		t.Error("expected length to be 5")
	}
}

// Test elements are overwritten when buffer is full
func TestRingBuffer_AddingSameKey(t *testing.T) {
	rc := NewRingBuffer[int, string](10)

	// add 10 elements
	elements := []string{"one", "two", "three", "four", "five"}
	for i, element := range elements {
		rc.Add(i, element)
	}

	// add element with key 2
	rc.Add(2, "new")

	// assert element was added and the order is correct
	retrieved := rc.GetLatest(5)
	newElements := []string{"one", "two", "four", "five", "new"}
	for i, element := range retrieved {
		if element != newElements[len(retrieved)-i-1] {
			t.Error("expected elements to be returned in correct order")
		}
	}

	// assert length is 5
	if rc.Len() != 5 {
		t.Error("expected length to be 5")
	}
}

package mapset

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Iterator struct {
	C    <-chan interface{}
	stop chan struct{}
}

// Stop stops the Iterator, no further elements will be received on C, C will be closed.
func (i *Iterator) Stop() {
	// Allows for Stop() to be called multiple times
	// (close() panics when called on already closed channel)
	defer func() {
		recover()
	}()

	close(i.stop)

	// Exhaust any remaining elements.
	for range i.C {
	}
}

// newIterator returns a new Iterator instance together with its item and stop channels.
func newIterator() (*Iterator, chan<- interface{}, <-chan struct{}) {
	itemChan := make(chan interface{})
	stopChan := make(chan struct{})
	return &Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
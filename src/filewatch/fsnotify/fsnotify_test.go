package fsnotify

import (
	"os"
	"testing"
	"time"
)

func TestEventStringWithValue(t *testing.T) {
	for opMask, expectedString := range map[Op]string{
		Chmod | Create: `"/usr/someFile": CREATE|CHMOD`,
		Rename:         `"/usr/someFile": RENAME`,
		Remove:         `"/usr/someFile": REMOVE`,
		Write | Chmod:  `"/usr/someFile": WRITE|CHMOD`,
	}{
		event := Event{Name:"/usr/someFile",Op: opMask}
		if event.String() != expectedString {
			t.Fatalf("Expected %s, got: %v",expectedString,event)
		}
	}
}

func TestEventOpStringWithValue(t *testing.T) {
	expectedOpString := "WRITE|CHMOD"
	event := Event{Name: "someFile", Op: Write|Chmod}
	if event.Op.String() != expectedOpString {
		t.Fatalf("Expected %s, got: %v", expectedOpString, event.Op.String())
	}
}

func TestEventOpStringWithNoValue(t *testing.T) {
	expectedOpString := ""
	event := Event{Name: "testFile", Op: 0}
	if event.Op.String() != expectedOpString {
		t.Fatalf("Expected %s, got: %v", expectedOpString, event.Op.String())
	}
}

// TestWatcherClose tests that the goroutine started by creating the watcher can be
// signalled to return at any time, even if there is no goroutine listening on the events
// or errors channels.
func TestWatcherClose(t *testing.T) {
	t.Parallel()

	name := tempMkFile(t, "")
	w := newWatcher(t)
	err := w.Add(name)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
	// Allow the watcher to receive the event.
	time.Sleep(time.Millisecond * 100)

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}
}
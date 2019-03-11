// Future represents an object that can be used to perform asynchronous
// tasks.  The constructor of the future will complete it, and listeners
// will block on getresult until a result is received.  This is different
// from a channel in that the future is only completed once, and anyone
// listening on the future will get the result, regardless of the number
// of listeners.

package futures

import (
  "fmt"
  "sync"
  "time"
)

type Completer <-chan interface{}

type Future struct {
  // because item can technically be nil and still be valid
  triggered bool
  item interface{}
  err error
  lock sync.Mutex
  wg sync.WaitGroup
}

// GetResult will immediately fetch the result if it exists
// or wait on the result until it is ready.
func (f *Future) GetResult() (interface{}, error) {
  f.lock.Lock()
  if f.triggered {
    f.lock.Unlock()
    return f.item, f.err
  }
  f.lock.Unlock()

  f.wg.Wait()  // 只有执行了 setItem后，才解出等待继续执行. 也就是 f都阻塞在了这里
  return f.item, f.err
}

// HasResult will return true iff the result exists
func (f *Future) HasResult() bool {
  f.lock.Lock()
  hasResult := f.triggered
  f.lock.Unlock()
  return hasResult
}

func (f *Future) setItem(item interface{}, err error) {
  f.lock.Lock()
  f.triggered = true
  f.item = item
  f.err = err
  f.lock.Unlock()
  f.wg.Done()  // 解除等待
}

func listenResult(f *Future, ch Completer, timeout time.Duration, wg *sync.WaitGroup) {
  wg.Done()
  t := time.NewTimer(timeout)  // time.Duration 原来是时间间隔
  select {
  case item := <-ch:
    f.setItem(item, nil)
    t.Stop()  // we want to trigger GC of this timer as soon as it's no longer needed
  case <-t.C:
    f.setItem(nil, fmt.Errorf(`timeout after %f seconds`, timeout.Seconds())) // 超时，报错。超时机制，保证了channel不会阻塞，goroutine不会泄露
  }
}

func New(completer Completer, timeout time.Duration) *Future {
  f := &Future{}
  f.wg.Add(1)
  var wg sync.WaitGroup
  wg.Add(1)
  go listenResult(f, completer, timeout, &wg)
  wg.Wait()
  return f
}


























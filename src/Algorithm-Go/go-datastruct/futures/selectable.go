package futures


import (
  "errors"
  "sync"
  "sync/atomic"
)

var ErrFutureCanceled = errors.New("future canceled")

// Selectable is a future with channel exposed for external `select`.
// Many simultaneous listeners may wait for result either with `f.Value()`
// or by selecting/fetching from `f.WaitChan()`, which is closed when future
// fulfilled.
// Selectable contains sync.Mutex, so it is not movable/copyable.

type Selectable struct {
  m sync.Mutex
  val interface{}
  err error  // sturct 上定义一个 error
  wait chan struct{}
  filled uint32
}

// NewSelectable returns new selectable future.
// Note: this method is for backward compatibility.
// You may allocate it directly on stack or embedding into larger structure

func NewSelectable() * Selectable {
  return &Selectable{}
}

// 为wait 赋值 chan struct{}, 然后返回
func (f *Selectable) wchan() <-chan struct{} {
  f.m.Lock()
  if f.wait == nil {
    f.wait = make(chan struct{})
  }
  ch := f.wait
  f.m.Unlock()
  return ch
}

// WaitChan returns channel, which is closed when future is fullfilled.
func (f *Selectable) WaitChan() <-chan struct{} {
  if atomic.LoadUint32(&f.filled) == 0 {
    <-f.wchan()
  }
  return f.wchan()
}

// GetResult waits for future to be fullfilled and returns value or error,
// whatever is set first
func (f *Selectable) GetResult() (interface{}, error) {
  if atomic.LoadUint32(&f.filled) == 0 {
    <-f.wchan()
  }
  return f.val, f.err
}

// Fill sets value for future, if it were not already fullfilled
// Returns error, if it were already set to future.
func (f *Selectable) Fill(v interface{}, e error) error {
  f.m.Lock()
  if f.filled == 0 {
    f.val = v
    f.err = e
    atomic.StoreUint32(&f.filled, 1)
    w := f.wait
    f.wait = closed
    if w != nil {
      close(w)   // close 操作， <-f.wchan() 操作就不再阻塞了
    }
  }
  f.m.Unlock()
  return f.err
}

func (f *Selectable) SetValue(v interface{}) error {
  return f.Fill(v, nil)
}

func (f *Selectable) SetError(e error) {
  f.Fill(nil, e)
}

// Cancel is alias for SetError(ErrFutureCanceled)
func (f *Selectable) Cancel() {
  f.SetError(ErrFutureCanceled)
}

var closed = make(chan struct{})

func init() {
  close(closed)
}
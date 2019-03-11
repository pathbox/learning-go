package futures

import (
  "fmt"
  "sync"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"
)

func TestSelectableGetResult(t *testing.T) {
  f := NewSelectable()
  var result interface{}
  var err error
  var wg sync.WaitGroup
  wg.Add(1)
  go func() {
    result, err = f.GetResult()
    wg.Done()
  }()

  f.SetValue(`test`)
  wg.Wait()

  assert.Nil(t, err)
  assert.Equal(t, `test`, result)

  // ensure we don't get paused on the next iteration.
  result, err = f.GetResult()

  assert.Equal(t, `test`, result)
  assert.Nil(t, err)
}

func TestSelectableSetError(t *testing.T) {
  f := NewSelectable()
  select {
  case <-f.WaitChan():
  case <-time.After(0):
    f.SetError(fmt.Errorf("timeout"))
  }

  result, err := f.GetResult()

  assert.Nil(t, result)
  assert.NotNil(t, err)
}

func BenchmarkSelectable(b *testing.B) {
  timeout := time.After(30 * time.Minute)
  var wg sync.WaitGroup

  b.ResetTimer()

  for i := 0; i < b.N; i++ {
    wg.Add(1)
    f := NewSelectable()
    go func() {
      select {
      case <-f.WaitChan():
      case <-timeout:
        f.SetError(fmt.Errorf("timeout"))
      }
      wg.Done()
    }()

    f.SetValue(`test`)
    wg.Wait()
  }
}
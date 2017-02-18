func doStuff() error {
  if someCondition {
    return errors.New("no space left on the device")
  } else {
    return errors.New("permission denied")
  }
}

var ErrNoSpaceLeft = errors.New("no space left on the device")
var ErrPermissionDenied = errors.New("permission denied")

func doStuff() error {
  if Condition {
    return ErrNoSpaceLeft
  } else {
    return ErrPermissionDenied
  }
}

err := doStuff()

if err == ErrPermissionDenied {
  // handle this particular error
}

var ErrBadRequest = errors.New("HTTP 400: Bad Request")
var ErrUnauthorized = errors.New("HTTP 401: Unauthorized")

type HTTPError struct {
  Code int
  Info string
}

func (h HTTPError) Error() string {
  return fmt.Sprintf("HTTP %d: %s", h.Code, h.Info)
}

func request() error {
  return HTTPError{404, "Not Found"}
}

func main() {
  err := request()
  if err != nil {
    if err.(HTTPError).Code == 404 {

      } else {

      }
  }
}

package main

import (
    "fmt"

    "github.com/juju/errgo"
)

type HTTPError struct {
    Code        int
    Description string
}

func (h HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", h.Code, h.Description)
}

func request() error {
    return errgo.Mask(HTTPError{404, "Not Found"})
}

func main() {
    err := request()

    fmt.Println(err.(errgo.Locationer).Location())

    realErr := err.(errgo.Wrapper).Underlying()

    if realErr.(HTTPError).Code == 404 {
        // handle a "not found" error
    } else {
        // handle a different error
    }
}

func f1() HTTPError { ... }
func f2() OSError { ... }

func main() {
    // err automatically declared as HTTPError
    err := f1()

    // OSError is a completely different type
    // The compiler does not allow this
    err = f2()
}
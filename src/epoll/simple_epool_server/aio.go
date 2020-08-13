package aio

// #include <errno.h>
// #include <fcntl.h>
// #include <sys/epoll.h>
import "C"

import (
    "syscall"
    "time"
    "unsafe"
)

type Poller int

func newPoller() (Poller, error) {
	fd, err := C.epoll_create1(C.O_CLOEXEC)
	if err != nil {
        return 0, err
    }
    return Poller(fd), nil
}
func (p Poller) Add(fd int, flags Flags) error {
    var ev C.struct_epoll_event
    if flags&In != 0 {
        ev.events |= C.EPOLLIN
    }
    if flags&Out != 0 {
        ev.events |= C.EPOLLOUT
    }
    if flags&Err != 0 {
        ev.events |= C.EPOLLERR
    }
    var dataFd = (*C.int)(unsafe.Pointer(&ev.data))
    *dataFd = C.int(fd)
    ok, err := C.epoll_ctl(C.int(p), C.EPOLL_CTL_ADD, C.int(fd), &ev)
    if ok < 0 && err != nil {
        if err == syscall.EEXIST {
            // Try MOD
            ok, err = C.epoll_ctl(C.int(p), C.EPOLL_CTL_MOD, C.int(fd), &ev)
        }
    }
    if ok >= 0 {
        err = nil
    }
    return err
}

func (p Poller) Delete(fd int) error {
    var ev C.struct_epoll_event
    // event must be non-NULL in kernels < 2.6.9
    ok, err := C.epoll_ctl(C.int(p), C.EPOLL_CTL_DEL, C.int(fd), &ev)
    if ok < 0 {
        return err
    }
    return nil
}

func (p Poller) Wait(duration time.Duration) ([]Event, error) {
    const maxEvents = 64
    inEvents := make([]C.struct_epoll_event, maxEvents)
    var timeout C.int
    if duration < 0 {
        timeout = -1
    } else {
        timeout = C.int(duration / time.Millisecond)
    }
    n, err := C.epoll_wait(C.int(p), (*C.struct_epoll_event)(unsafe.Pointer(&inEvents[0])), maxEvents, timeout)
    if err != nil {
        if err == syscall.EINTR {
            err = nil
        }
        return nil, err
    }
    events := make([]event, int(n))
    for ii := 0; ii < int(n); ii++ {
        inEvent := inEvents[ii]
        var flags Flags
        if inEvent.events&C.EPOLLIN != 0 {
            flags |= In
        }
        if inEvent.events&C.EPOLLOUT != 0 {
            flags |= Out
        }
        if inEvent.events&C.EPOLLERR != 0 {
            flags |= Err
        }
        fd := (*C.int)(unsafe.Pointer(&inEvent.data))
        events[ii] = Event{
            Fd:    int(*fd),
            Flags: flags,
        }
    }
    return events, nil
}
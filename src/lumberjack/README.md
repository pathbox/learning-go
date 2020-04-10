https://github.com/natefinch/lumberjack

Lumberjack is a Go package for writing logs to rolling files

import "gopkg.in/natefinch/lumberjack.v2"

```go
log.SetOutput(&lumberjack.Logger{
    Filename:   "/var/log/myapp/foo.log",
    MaxSize:    500, // megabytes
    MaxBackups: 3,
    MaxAge:     28, //days
    Compress:   true, // disabled by default
})

l := &lumberjack.Logger{}
log.SetOutput(l)
c := make(chan os.Signal, 1)
signal.Notify(c, syscall.SIGHUP)

go func() {
    for {
        <-c
        l.Rotate()
    }
}()
```
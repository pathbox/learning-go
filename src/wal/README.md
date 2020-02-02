https://github.com/tidwall/wal

Write ahead log for Go

High durability
Fast writes
Low memory footprint
Monotonic indexes
Log truncation from front or back.

```go
// open a new log file
l, _ := Open("mylog", nil)
// write some entries
l.Write(1, []byte("first entry"))
l.Write(2, []byte("first entry"))
l.Write(3, []byte("first entry"))

// read an entry
data, _ := l.Read(1)
println(string(data))  // output: first entry

// close the log
l.Close()
```
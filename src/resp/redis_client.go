// Redis Client interface with the Redis Serialization Protocol(RESP)
// Refs: https://www.redisgreen.net/blog/reading-and-writing-redis-protocol/

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	SIMPLE_STRING = '+'
	BULK_STRING   = '$'
	INTEGER       = ':'
	ARRAY         = '*'
	ERROR         = '-'
)

var (
	arrayPrefixSlice      = []byte{ARRAY}
	bulkStringPrefixSlice = []byte{BULK_STRING}
	lineEndingSlice       = []byte{'\r', '\n'}
	ErrInvalidSyntax      = errors.New("resp: invalid syntax")
)

type RESPWriter struct {
	*bufio.Writer
}

type RESPReader struct {
	*bufio.Reader
}

func NewRESPWriter(writer io.Writer) *RESPWriter {
	return &RESPWriter{
		Writer: bufio.NewWriter(writer),
	}
}

func NewRESPReader(reader io.Reader) *RESPReader {
	return &RESPReader{
		Reader: bufio.NewReaderSize(reader, 32*1024), //32 kilobytes
	}
}

func (w *RESPWriter) WriteCommand(args ...string) (err error) {
	// Write the array prefix and the number of arguments in the array
	w.Write(arrayPrefixSlice)
	w.Write([]byte(strconv.Itoa(len(args))))
	w.Write(lineEndingSlice)

	//write a bulk string for each argument
	for _, arg := range args {
		w.Write(bulkStringPrefixSlice)
		w.Write([]byte(strconv.Itoa(len(args))))
		w.Write(lineEndingSlice)
		w.WriteString(arg)
		w.Write(lineEndingSlice)
	}
	return w.Flush()
}

func (r *RESPReader) ReadObject() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case SIMPLE_STRING, INTEGER, ERROR:
		return line, nil
	case BULK_STRING:
		return r.readBulkString(line)
	case ARRAY:
		return r.readArray(line)
	default:
		return nil, ErrInvalidSyntax
	}
}

func (r *RESPReader) readLine() (line []byte, err error) {
	line, err = r.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	if len(line) > 1 && line[len(line)-2] == '\r' {
		return line, nil
	} else {
		// Line was too short or newline wasn't preceded by carriage return
		return nil, ErrInvalidSyntax
	}
}

func (r *RESPReader) readBulkString(line []byte) ([]byte, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}
	if count == -1 {
		return line, nil
	}

	buf := make([]byte, len(line)+count+2)
	copy(buf, line)
	_, err = r.Read(buf[len(line):])
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *RESPReader) getCount(line []byte) (int, error) {
	end := bytes.IndexByte(line, '\r')
	return strconv.Atoi(string(line[1:end]))
}

func (r *RESPReader) readArray(line []byte) ([]byte, error) {
	//Get number of array elements
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}

	// Read `count` number of RESP objects in the array.
	for i := 0; i < count; i++ {
		buf, err := r.ReadObject()
		if err != nil {
			return nil, err
		}
		line = append(line, buf...)
	}

	return nil, err
}

func main() {
	var buf bytes.Buffer
	writer := NewRESPWriter(&buf)
	writer.WriteCommand("GET", "foo")
	fmt.Println(string(buf.Bytes()))
}

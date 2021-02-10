package ringbuffer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"
)

// ErrIsEmpty 缓冲区为空
var ErrIsEmpty = errors.New("ring buffer is empty")

// RingBuffer 自动扩容循环缓冲区
type RingBuffer struct {
	buf     []byte
	size    int
	vr      int
	r       int // next position to read
	w       int // next position to write
	isEmpty bool
}

// New 返回一个初始大小为 size 的 RingBuffer
func New(size int) *RingBuffer {
	return &RingBuffer{
		buf:     make([]byte, size),
		size:    size,
		isEmpty: true,
	}
}

// NewWithData 特殊场景使用，RingBuffer 会持有data，不会自己申请内存去拷贝
func NewWithData(data []byte) *RingBuffer {
	return &RingBuffer{
		buf:  data,
		size: len(data),
	}
}

// VirtualFlush 刷新虚读指针
// VirtualXXX 系列配合使用
func (r *RingBuffer) VirtualFlush() {
	r.r = r.vr
	if r.r == r.w {
		r.isEmpty = true
	}
}

// VirtualRevert 还原虚读指针
// VirtualXXX 系列配合使用
func (r *RingBuffer) VirtualRevert() {
	r.vr = r.r
}

// VirtualRead 虚读，不移动 read 指针，需要配合 VirtualFlush 和 VirtualRevert 使用
// VirtualXXX 系列配合使用
func (r *RingBuffer) VirtualRead(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.isEmpty {
		return 0, ErrIsEmpty
	}
	n = len(p)
	if r.w > r.vr {
		if n > r.w-r.vr {
			n = r.w - r.vr
		}
		copy(p, r.buf[r.vr:r.vr+n])
		// move vr
		r.vr = (r.vr + n) % r.size
		if r.vr == r.w {
			r.isEmpty = true
		}
		return
	}
	if n > r.size-r.vr+r.w {
		n = r.size - r.vr + r.w
	}
	if r.vr+n <= r.size {
		copy(p, r.buf[r.vr:r.vr+n])
	} else {
		// head
		copy(p, r.buf[r.vr:r.size])
		// tail
		copy(p[r.size-r.vr:], r.buf[0:n-r.size+r.vr])
	}

	// move vr
	r.vr = (r.vr + n) % r.size
	return
}

// VirtualLength 虚拟长度，虚读后剩余可读数据长度
// VirtualXXX 系列配合使用
func (r *RingBuffer) VirtualLength() int {
	if r.w == r.vr {
		if r.isEmpty {
			return 0
		}
		return r.size
	}

	if r.w > r.vr {
		return r.w - r.vr
	}

	return r.size - r.vr + r.w
}

func (r *RingBuffer) RetrieveAll() {
	r.r = 0
	r.w = 0
	r.vr = 0
	r.isEmpty = true
}

func (r *RingBuffer) Retrieve(len int) {
	if r.isEmpty || len <= 0 {
		return
	}

	if len < r.Length() {
		r.r = (r.r + len) % r.size
		r.vr = r.r

		if r.w == r.r {
			r.isEmpty = true
		}
	} else {
		r.RetrieveAll()
	}
}

func (r *RingBuffer) Peek(len int) (first []byte, end []byte) {
	if r.isEmpty || len <= 0 {
		return
	}

	if r.w > r.r {
		if len > r.w-r.r {
			len = r.w - r.r
		}

		first = r.buf[r.r : r.r+len]
		return
	}

	if len > r.size-r.r+r.w {
		len = r.size - r.r + r.w
	}
	if r.r+len <= r.size {
		first = r.buf[r.r : r.r+len]
	} else {
		// head
		first = r.buf[r.r:r.size]
		// tail
		end = r.buf[0 : len-r.size+r.r]
	}
	return
}

func (r *RingBuffer) PeekAll() (first []byte, end []byte) {
	if r.isEmpty {
		return
	}

	if r.w > r.r {
		first = r.buf[r.r:r.w]
		return
	}

	first = r.buf[r.r:r.size]
	end = r.buf[0:r.w]
	return
}

func (r *RingBuffer) PeekUint8() uint8 {
	if r.Length() < 1 {
		return 0
	}

	f, e := r.Peek(1)
	if len(e) > 0 {
		return e[0]
	} else {
		return f[0]
	}
}

func (r *RingBuffer) PeekUint16() uint16 {
	if r.Length() < 2 {
		return 0
	}

	f, e := r.Peek(2)
	if len(e) > 0 {
		return binary.BigEndian.Uint16(copyByte(f, e))
	} else {
		return binary.BigEndian.Uint16(f)
	}
}

func (r *RingBuffer) PeekUint32() uint32 {
	if r.Length() < 4 {
		return 0
	}

	f, e := r.Peek(4)
	if len(e) > 0 {
		return binary.BigEndian.Uint32(copyByte(f, e))
	} else {
		return binary.BigEndian.Uint32(f)
	}
}

func (r *RingBuffer) PeekUint64() uint64 {
	if r.Length() < 8 {
		return 0
	}

	f, e := r.Peek(8)
	if len(e) > 0 {
		return binary.BigEndian.Uint64(copyByte(f, e))
	} else {
		return binary.BigEndian.Uint64(f)
	}
}

func (r *RingBuffer) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.isEmpty {
		return 0, ErrIsEmpty
	}
	n = len(p)
	if r.w > r.r {
		if n > r.w-r.r {
			n = r.w - r.r
		}
		copy(p, r.buf[r.r:r.r+n])
		// move readPtr
		r.r = (r.r + n) % r.size
		if r.r == r.w {
			r.isEmpty = true
		}
		r.vr = r.r
		return
	}
	if n > r.size-r.r+r.w {
		n = r.size - r.r + r.w
	}
	if r.r+n <= r.size {
		copy(p, r.buf[r.r:r.r+n])
	} else {
		// head
		copy(p, r.buf[r.r:r.size])
		// tail
		copy(p[r.size-r.r:], r.buf[0:n-r.size+r.r])
	}

	// move readPtr
	r.r = (r.r + n) % r.size
	if r.r == r.w {
		r.isEmpty = true
	}
	r.vr = r.r
	return
}

func (r *RingBuffer) ReadByte() (b byte, err error) {
	if r.isEmpty {
		return 0, ErrIsEmpty
	}
	b = r.buf[r.r]
	r.r++
	if r.r == r.size {
		r.r = 0
	}

	if r.w == r.r {
		r.isEmpty = true
	}
	r.vr = r.r
	return
}

func (r *RingBuffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	n = len(p)
	free := r.free()
	if free < n {
		r.makeSpace(n - free)
	}
	if r.w >= r.r {
		if r.size-r.w >= n {
			copy(r.buf[r.w:], p)
			r.w += n
		} else {
			copy(r.buf[r.w:], p[:r.size-r.w])
			copy(r.buf[0:], p[r.size-r.w:])
			r.w += n - r.size
		}
	} else {
		copy(r.buf[r.w:], p)
		r.w += n
	}

	if r.w == r.size {
		r.w = 0
	}

	r.isEmpty = false

	return
}

func (r *RingBuffer) WriteByte(c byte) error {
	if r.free() < 1 {
		r.makeSpace(1)
	}

	r.buf[r.w] = c
	r.w++

	if r.w == r.size {
		r.w = 0
	}

	r.isEmpty = false

	return nil
}

func (r *RingBuffer) Length() int {
	if r.w == r.r {
		if r.isEmpty {
			return 0
		}
		return r.size
	}

	if r.w > r.r {
		return r.w - r.r
	}

	return r.size - r.r + r.w
}

func (r *RingBuffer) Capacity() int {
	return r.size
}

func (r *RingBuffer) WriteString(s string) (n int, err error) {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return r.Write(*(*[]byte)(unsafe.Pointer(&h)))
}

// Bytes 返回所有可读数据，此操作不会移动读指针，仅仅是拷贝全部数据
func (r *RingBuffer) Bytes() (buf []byte) {
	if r.w == r.r {
		if !r.isEmpty {
			buf := make([]byte, r.size)
			copy(buf, r.buf)
			return buf
		}
		return
	}

	if r.w > r.r {
		buf = make([]byte, r.w-r.r)
		copy(buf, r.buf[r.r:r.w])
		return
	}

	buf = make([]byte, r.size-r.r+r.w)
	copy(buf, r.buf[r.r:r.size])
	copy(buf[r.size-r.r:], r.buf[0:r.w])
	return
}

func (r *RingBuffer) IsFull() bool {
	return !r.isEmpty && r.w == r.r
}

func (r *RingBuffer) IsEmpty() bool {
	return r.isEmpty
}

func (r *RingBuffer) Reset() {
	r.r = 0
	r.w = 0
	r.isEmpty = true
}

func (r *RingBuffer) String() string {
	return fmt.Sprintf("Ring Buffer: \n\tCap: %d\n\tReadable Bytes: %d\n\tWriteable Bytes: %d\n\tBuffer: %s\n", r.size, r.Length(), r.free(), r.buf)
}

func (r *RingBuffer) makeSpace(len int) {
	newSize := r.size + len
	newBuf := make([]byte, newSize)
	oldLen := r.Length()
	_, _ = r.Read(newBuf)

	r.w = oldLen
	r.r = 0
	r.size = newSize
	r.buf = newBuf
}

func (r *RingBuffer) free() int {
	if r.w == r.r {
		if r.isEmpty {
			return r.size
		}
		return 0
	}

	if r.w < r.r {
		return r.r - r.w
	}

	return r.size - r.w + r.r
}

func copyByte(f, e []byte) []byte {
	buf := make([]byte, len(f)+len(e))
	_ = copy(buf, f)
	_ = copy(buf[len(f):], e)
	return buf
}
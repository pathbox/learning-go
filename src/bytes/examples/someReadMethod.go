import "io"

// 每次读取一 line数据
func (b *Buffer) readSlice(delim byte) (line []byte, err error) {
	i := IndexByte(b.buf[b.off:], delim) // 找到从off开始到末尾出现的第一个delim的index值
	end := b.off + i + 1                 // + 1 是需要的操作因为[:] 实际是 [:) end 是 一行line数据的end位置
	if i < 0 {                           // 说明没有找到delim， 有可能这次已经读完了
		end = len(b.buf)
		err = io.EOF
	}

	line := b.buf[b.off:end]
	b.off = end
	b.lastRead = opRead
	return line, err

}

func (b *Buffer) ReadBytes(delim byte) (line []byte, err error) {
	slice, err := b.readSlice(delim)
	line = append(line, slice...) // 其实就是返回了 slice
	return line, err
}

func (b *Buffer) ReadString(delim byte) (line string, err error) {
	slice, err := b.readSlice(delim)
	return string(slice), err
}
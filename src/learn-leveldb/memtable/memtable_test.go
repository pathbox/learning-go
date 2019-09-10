package memtable

import (
	"fmt"
	"testing"

	"github.com/merlin82/leveldb/internal"
)

func Test_MemTable(t *testing.T) {
	memTable := New()
	memTable.Add(1234567, internal.TypeValue, []byte("aadsa34a"), []byte("bb23b3423"))
	value, _ := memTable.Get([]byte("aadsa34a"))
	fmt.Println(string(value))
	fmt.Println(memTable.ApproximateMemoryUsage())
}

// 在内存中存储，底层用眺表实现存储数据

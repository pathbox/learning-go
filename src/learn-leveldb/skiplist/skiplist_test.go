package skiplist

import (
	"fmt"
	"math/rand"
	"testing"

	"./utils"
)

func TestInsert(t *testing.T) {
	skiplist := New(utils.IntComparator)
	for i := 0; i < 10; i++{
		skiplist.Insert(rand.Int() % 10)
	}
	it := skiplist.NewIterator()
	for it.SeekToFirst(); it.Valid(); it.Next() {
		fmt.Println(it.Key())
	}
	fmt.Println()
	for it.SeekToLast(); it.Valid(); it.Prev() {
		fmt.Println(it.Key())
	}
}
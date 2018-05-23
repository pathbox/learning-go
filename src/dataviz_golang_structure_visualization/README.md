https://medium.com/@Arafat./introducing-dataviz-a-data-structure-visualization-library-for-golang-f6e60663bc9d

```go
package main

import (
	bheap "github.com/Arafatk/dataviz/trees/binaryheap"
)

func main() {
	heap := bheap.NewWithIntComparator()
	heap.Push(3)
	heap.Push(19)
	heap.Push(17)
	heap.Push(2)
	heap.Push(7)
	heap.Push(1)
	heap.Push(26)
	heap.Push(35)
	heap.Visualizer("heap.png")
}

package main

import (
	rbt "github.com/Arafatk/dataviz/trees/redblacktree"
)

func main() {
	tree := rbt.NewWithIntComparator()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite
	tree.Visualizer("out.png")
}
view raw
```
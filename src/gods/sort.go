package examples

import "github.com/emirpasic/gods/utils"

// SortExample to demonstrate basic usage of basic sort
func SortExample() {
	strings := []interface{}{}                  // []
	strings = append(strings, "d")              // ["d"]
	strings = append(strings, "a")              // ["d","a"]
	strings = append(strings, "b")              // ["d","a",b"
	strings = append(strings, "c")              // ["d","a",b","c"]
	utils.Sort(strings, utils.StringComparator) // ["a","b","c","d"]
}

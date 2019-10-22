// sort.Search  in golang
package sortexample
import (
	"fmt"
	"sort"
)

func GuessingGame() {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)
 		fmt.Scanf("%s", &s)
 		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}

func Search(n int, f func(int) bool) int {
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // 得到中间值 /2
		// i ≤ h < j
		if !f(h) { // 所以传入的f 要进行大小比较判断
			i = h + 1 // preserves f(i-1) == false 往后半区查找
		} else {
			j = h // preserves f(j) == true 往前半区查找
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

//found return binary search result with sort order
func (db *Db) found(b []byte, asc bool) int {
	db.sort() // 对keys [][]byte 进行排序，之后可以使用sort.Search 二分法查找值

	return sort.Search(len(db.keys), func(i int) bool { //二分法从keys中查找是否有b这个key，所以keys是有序的
		return bytes.Compare(db.keys[i], b) >= 0
	})

}
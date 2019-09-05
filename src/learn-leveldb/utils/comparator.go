package utils

type Comparator func(a, b interface{}) int

func IntComparator(a, b interface{}) int {
	aInt := a.(int)
	bInt := b.(int)
	return aInt - bInt
}

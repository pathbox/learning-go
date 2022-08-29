package hiset

// Set go set
type Set[T comparable] interface {
	// Add 添加
	Add(items ...T) bool

	// Remove 删除内容
	Remove(items ...T) bool

	// Contains 是否包含对应数据
	Contains(items ...T) bool

	// IsEmpty set是否为空
	IsEmpty() bool

	// Size 返回set长度
	Size() int

	// ToArray 返回slice数据
	ToArray() []T

	// Equals 两个set是否相等
	Equals(Set[T]) bool

	// Union 两个set取并集
	Union(Set[T]) Set[T]

	// Intersection 两个set取交集
	Intersection(Set[T]) Set[T]

	// Clear 清空set
	Clear()

	// MarshalJSON json序列化
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON 反序列化
	UnmarshalJSON(b []byte) error
}

func NewSet[T comparable](items ...T) Set[T] {
	s := newSet[T]()
	_ = s.Add(items...)
	return s
}

// NewThreadSafeSet 初始化线程安全Set
func NewThreadSafeSet[T comparable](items ...T) Set[T] {
	s := newThreadSafeSet[T]()
	_ = s.Add(items...)
	return s
}

func NewLinkedSet[T comparable](items ...T) Set[T] {
	s := newLinkedSet[T]()
	_ = s.Add(items...)
	return s
}

func union[T comparable](old1 Set[T], old2 Set[T], new Set[T]) {
	if old1 != nil {
		new.Add(old1.ToArray()...)
	}
	if old2 != nil {
		new.Add(old2.ToArray()...)
	}

}

func intersection[T comparable](old1 Set[T], old2 Set[T], new Set[T]) {
	if old1 == nil || old2 == nil {
		return
	}

	old1Array := old1.ToArray()
	for index, item := range old1Array {
		if old2.Contains(item) {
			// 这里为了防止for range item同一个地址的问题
			new.Add(old1Array[index])
		}
	}
}

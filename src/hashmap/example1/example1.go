package hashmap

import "fmt"

const maxCapacity int = 100     // bucket max cap
const loadFactor float32 = 0.75 // 负载因子

var nowCapacity int = 8 // 当前容量

type Entry struct {
	key   int
	value string
	next  *Entry
}

type HashMap struct {
	size    int
	buckets []Entry
}

func CreateHashMap() *HashMap {
	hm := &HashMap{}
	hm.size = 0
	hm.buckets = make([]Entry, nowCapacity, maxCapacity)
	return hm
}

func getHash(k int) int {
	return k % nowCapacity
}

func (hm *HashMap) GetSize() int {
	return hm.size
}

// 由key在hashmap中找到其Entry指针
func (hm *HashMap) GetEntry(k int) (*Entry, bool) {
	p := hm.buckets[getHash(k)].next
	for p != nil {
		if p.key == k {
			return p, true
		} else {
			p = p.next
		}
	}

	return nil, false
}

func (hm *HashMap) insert(e *Entry) {
	hash := getHash(e.key)
	// 从表头中的指针开始遍历
	var p *Entry = &hm.buckets[hash]
	for p.next != nil {
		// 如果找到相同的key，则覆盖其中的值，完成insert，return
		if p.next.key == e.key {
			p.next.value = e.value
			return
		} else {
			p = p.next
		}
	}

	p.next = e
	hm.size++
}

func (hm *HashMap) Put(k int, v string) {
	e := &Entry{k, v, nil}
	hm.insert(e)
	// 达到负载因子且还能扩容时，扩容并迁移数据
	if float32(hm.size)/float32(nowCapacity) >= loadFactor && nowCapacity < maxCapacity {
		if 2*nowCapacity > maxCapacity {
			nowCapacity = maxCapacity
		} else {
			nowCapacity = 2 * nowCapacity
		}

		newHm := CreateHashMap()
		var index int
		for index = 0; index < len(hm.buckets); index++ {
			p := hm.buckets[index].next
			for p != nil {
				pNext := p.next
				p.next = nil
				newHm.insert(p)
				p = pNext
			}
		}

		var b1 int
		var b2 int = len(newHm.buckets)

		for b1 = len(hm.buckets); b1 < b2; b++ {
			hm.buckets = append(hm.buckets, Entry{})
		}
		//移回数据
		for b1 = 0; b1 < b2; b1++ {
			hm.buckets[b1].next = newHm.buckets[b1].next
			newHm.buckets[b1].next = nil
		}
	}
}

// 删除指定Entry
func (hm *HashMap) DeleteEntry(e *Entry) {
	p := &hm.buckets[getHash(e.key)]
	for p != nil {
		if p.next == e {
			fmt.Println("正在删除")
			p.next = p.next.next
			fmt.Println("删除成功")

			hm.size--
			return
		} else {
			p = p.next
		}
	}
	fmt.Println("删除失败")
}

//删除指定key的Entry
func (hm *HashMap) Delete(k int) bool {
	e, ok := hm.GetEntry(k)
	if ok {
		fmt.Println("有该Entry,其hash为：", getHash(e.key), "值为：", e)
		hm.DeleteEntry(e)
		return true
	} else {
		fmt.Println("没有该Entry，删除失败")
		return false
	}
}

//遍历hashMap
func (hm *HashMap) Traverse() {
	var index int
	for index = 0; index < nowCapacity; index++ {
		p := hm.buckets[index].next
		if p == nil {
			fmt.Println(index, ":")
		} else {
			fmt.Print(index, ":")
		}
		for p != nil {
			fmt.Print("----->", p)
			p = p.next
			if p == nil {
				fmt.Println()
			}
		}
	}
	fmt.Println()
}

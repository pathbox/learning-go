https://github.com/iancoleman/orderedmap

最简单的有序map实现思考

type OrderedMap struct {
	keys   []string               // key数组
	values map[string]interface{} // 内部map
}

插入key时，将key加到keys的末尾

删除key时，也在keys数组中将其删除

func (o *OrderedMap) Delete(key string) {
	_, ok := o.values[key]
	if !ok {
		return
	}
	// remove from keys
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...) // 去除当前i的key
			break
		}
	}
	// remove from values
	delete(o.values, key)
}

只有当要顺序遍历OrderedMap时，先对keys数组进行排序，然后再按keys排序的结果进行遍历得到key，再去values map中得到对应值，这样每次写入的时候是不用排序的

另一种方案，keys数组按照二叉搜索树的方式存储，这样写入会是log(n),遍历的时候，对keys进行中序遍历log(n)，能得到排序的key
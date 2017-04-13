// 第一次初始session对象，这个session对象是A

// demo方法执行完毕会调用 session.Close()将A的mongodb连接释放掉。

// func demo() {
//     session, collection, err := GetCollection(DBNAME, COLLECTIONNAME)
//     defer session.Close()
//     session, collection, err = GetCollection(DBNAME, COLLECTIONNAME)
// }
// 这样便会出现连接“泄露”。

// 因为初始化了另一个sesseion对象B

// demo执行到第三句是，session的引用指向了对象B，而对象A的引用则丢失了，A成了野对象。

// demo执行完毕后，session.Close()释放了B的mongodb连接，但是A对象引用已经丢失,A的mongodb连接,永远不会得到释放。



// 垃圾回收可以把野对象回收了（内存回收），但是野对象所对应的socket并没有被释放（资源回收）
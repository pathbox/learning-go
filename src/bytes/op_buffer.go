func Action(w http.ResponseWriter, r *http.Request) {

    var result string
    // handle func
        // do something
    w.Write([]byte(result))  // 当result数据量比较大的时候,会导致内存的增长,比如 处理的是图片
    r.Body.Close()
}



buffer := bytes.NewBuffer(make([]byte, 0, resp.ContentLength)
buffer.ReadFrom(res.Body)
body := buffer.Bytes()

// buffer采用最小单位读，若不够，则继续申请2倍大的空间

buffer := bytes.NewBuffer(make([]byte, 0, 65536))
io.Copy(buffer, r.Body)
temp := buffer.Bytes()
length := len(temp)
var body []byte
//are we wasting more than 10% space?
if cap(temp) > (length + length / 10) {
  body = make([]byte, length)
  copy(body, temp)
} else {
  body = temp
}

A Memory Leak

/* Look, what's a memory leak within the context of a runtime that provides garbage collection?
 Typically it's either a rooted object, or a reference from a rooted object,
 which you haven't considered. This is obviously different as it's really extra memory
 you might not be aware of. Rooting the object might very well be intentional,
 but you don't realize just how much memory it is you've rooted. Sure,
 my ignorance is at least 75% to blame. Yet I can't help but shake the feeling that
 there's something too subtle about all of this. Any code can return something that
	looks and quacks like an array of 2 integers yet takes gigs of memory. Furthermore,
	bytes.MinRead as a global variable is just bad design. I can't imagine how many
	people think they've allocated X when they've really allocated X*2+512.

*/

// 这样就不会

// http://www.cnblogs.com/zhangqingping/p/4390913.html
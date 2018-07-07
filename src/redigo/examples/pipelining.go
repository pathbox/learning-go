// Connections support pipelining using the Send, Flush and Receive methods.
// Send(commandName string, args ...interface{}) error
// Flush() error
// Receive() (reply interface{}, err error)

c.Send("SET","foo","bar")
c.Send("GET", "foo")
c.Flush() // 将上面的两次Send 的 操作，作为一次请求发送给redis server
c.Receive() // reply from SET
v, err = c.Receive() // reply from GET

// 使用redis事务
c.Send("MULTI")
c.Send("INCR", "foo")
c.Send("INCR", "bar")
r, err := c.Do("EXEC") // 进行了 Flush 和 Receive， Do用的就是pipelining
fmt.Println(r) // prints [1, 1]

//The Do method combines the functionality of the Send, Flush and Receive methods. The Do method starts by writing the command and flushing the output buffer
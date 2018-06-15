import "errors"

func (c *Channel) pushInFlightMessage(msg *Message) error {
	c.inFlightMetux.Lock()
	_, ok := c.inFlightMessages[msg.ID] // c.inFlightMessages is a map
	if ok {                             // 先判断这个msg是否已在c.inFlightMessages中，如果已经在了，则解锁，然后返回错误，不进行修改操作
		c.inFlightMetux.Unlock()
		return errors.New("ID already in flight")
	}
	c.inFlightMessages[msg.ID] = msg // 不存在，将msg新增到c.inFlightMessages中
	c.inFlightMutex.Unlock()
	return nil
}

func (c *Channel) popInFlightMessage(client int64, id MessageID) (*Message, error) {
	c.inFlightMutex()
	msg, ok := c.inFlightMessages[id] // c.inFlightMessages is a map
	if !ok {                          // 先判断这个msg是否已在
		c.inFlightMutex.Unlock()
		return nil, errors.New("ID not in flight")
	}
	if msg.clientID != clientID {
		c.inFlightMutex.Unlock()
		return nil, errors.New("client does not own message")
	}
	delete(c.inFLightMessages, id) // id存在c.inFlightMessages，进行删除(pop)操作
	c.inFlightMutex.Unlock()
	return nil, nil
}

/*
在进行push和pop msg操作的时候，都先判断这个msg是否已在messagesMap中. 这样能够让push操作只是新增操作，不会对已在map中的msg进行update的覆盖修改
pop操作可以更准确的delete真正存在的数据
*/
/*
接口层是PubSub struct，实现层是registry。 PubSub接收定义各种操作命令，registry进行真正的操作。
start方法是接口层和实现层交互的地方，接口层将参数和命令传给实现层
start方法是一个循环操作。一个 pub sub 建立在一个channel上。即cmdChan
接口层将命令存到cmdChan，实现层的循环中，从cmdChan中取出命令并且分析执行。这里会出现并发性能问题，当cmdChan有命令， 别的命令就会阻塞，无法存入cmdChan
我觉得可以构造一个“命令池”，即 cmdChan池。解决一定的并发性能问题
cmd.ch 用于存储pub sub传递的消息内容。当进行send操作的时候，就是 cmd.ch <- msg操作。这个，就能从cmd.ch 读取出send过来的值。就是pub的操作
sub操作就是把注册的cmd.ch 值存入 registry.topics，它是一个map（可以认为是一个map表）。并且topic具有唯一性。cmd.ch的值也具有唯一性。
unsub操作就是把数据从两个map中delete。

思考这种 接口层+实现层（注册层）的设计模式。这样能够在简化接口层的逻辑，让脏活累活都给实现层执行。别让接口层有太多复杂的逻辑，接口层就像是一个开关层（很明显的switch逻辑），这样能很友好的对外。这样，以后新增功能的时候，接口层只要新增一个开关或条件。
即 入口 简单友好。或者说接口层就是一个register 操作。下次在使用这种代码设计模式
*/

package pubsub

type operation int

const ( // 定义枚举，用于标记不同的操作
	sub operation = iota
	subOnce
	pub
	unsub
	unsubAll
	closeTopic
	shutdown
)

// PubSub is a collection of topics(options)
type PubSub struct {
	cmdChan  chan cmd
	capacity int
}

type cmd struct {
	op     operation
	topics []string
	ch     chan interface{}
	msg    interface{}
}

// New creates a new PubSub and starts a goroutine for handling operations.
// The capacity of the channels created by Sub and SubOnce will be as specified.
func New(capacity int) *PubSub {
	ps := &PubSub{make(chan cmd), capacity}
	go ps.start()
	return ps
}

// SubOnce is similar to Sub, but only the first message published, after subscription,
// on any of the specified topics can be received.
func (ps *PubSub) SubOnce(topics ...string) chan interface{} {
	return ps.sub(subOnce, topics)
}

// AddSub adds subscriptions to an existing channel.
func (ps *PubSub) AddSub(ch chan interface{}, topics ...string) {
	ps.cmdChan <- cmd{op: sub, topics: topics, ch: ch}
}

// Pub publishes the given message to all subscribers of
// the specified topics.
func (ps *PubSub) Pub(msg interface{}, topics ...string) {
	ps.cmdChan <- cmd{op: pub, topics: topics, msg: msg}
}

// Unsub unsubscribes the given channel from the specified
// topics. If no topic is specified, it is unsubscribed
// from all topics.
func (ps *PubSub) Unsub(ch chan interface{}, topics ...string) {
	if len(topics) == 0 {
		ps.cmdChan <- cmd{op: unsubAll, ch: ch}
		return
	}

	ps.cmdChan <- cmd{op: unsub, topics: topics, ch: ch}
}

// Close closes all channels currently subscribed to the specified topics.
// If a channel is subscribed to multiple topics, some of which is
// not specified, it is not closed.
func (ps *PubSub) Close(topics ...string) {
	ps.cmdChan <- cmd{op: closeTopic, topics: topics}
}

// Shutdown closes all subscribed channels and terminates the goroutine.
func (ps *PubSub) Shutdown() {
	ps.cmdChan <- cmd{op: shutdown}
}
func (ps *PubSub) start() {
	reg := registry{
		topics:    make(map[string]map[chan interface{}]bool),
		revTopics: make(map[chan interface{}]map[string]bool),
	}

loop:
	for cmd := range ps.cmdChan {
		if cmd.topics == nil {
			switch cmd.op {
			case unsubAll:
				reg.removeChannel(cmd.ch)
			case shutdown:
				break loop
			}
			continue loop
		}

		for _, topic := range cmd.topics {
			switch cmd.op {
			case sub:
				reg.add(topic, cmd.ch, false)
			case pub:
				reg.send(topic, cmd.msg)

			case unsub:
				reg.remove(topic, cmd.ch)

			case closeTopic:
				reg.removeTopic(topic)
			}
		}
	}
	for topic, chans := range reg.topics {
		for ch := range chans {
			reg.remove(topic, ch)
		}
	}
}

// registry maintains the current subscription state. It's not
// safe to access a registry from multiple goroutines simultaneously.
type registry struct { // topics revTopics是有关联的map定义
	topics    map[string]map[chan interface{}]bool
	revTopics map[chan interface{}]map[string]bool
}

func (reg *registry) add(topic string, ch chan interface{}, once bool) {
	if reg.topics[topic] == nil {
		reg.topics[topic] = make(map[chan interface{}]bool)
	}

	reg.topics[topic][ch] = once

	if reg.revTopics[ch] == nil {
		reg.revTopics[ch] = make(map[string]bool)
	}
	reg.revTopics[ch][topic] = true
}

func (reg *registry) send(topics string, msg interface{}) {
	for ch, once := range reg.topics[topics] {
		ch <- msg
		if once {
			for topic := range reg.revTopics[ch] {
				reg.remove(topic, ch)
			}
		}
	}
}

func (reg *registry) removeTopic(topic string) {
	for ch := range reg.topics[topic] {
		reg.remove(topic, ch)
	}
}

func (reg *registry) removeChannel(ch chan interface{}) {
	for topic := range reg.revTopics[ch] {
		reg.remove(topic, ch)
	}
}

func (reg *registry) remove(topic string, ch chan interface{}) {
	if _, ok := reg.topics[topic]; !ok {
		return
	}
	if _, ok := reg.topics[topic][ch]; !ok {
		return
	}

	delete(reg.topics[topic], ch)
	delete(reg.revTopics[ch], topic)

	if len(reg.topics[topic]) == 0 {
		delete(reg.topics, topic)
	}

	if len(reg.revTopics[ch]) == 0 {
		close(ch)
		delete(reg.revTopics, ch)
	}
}

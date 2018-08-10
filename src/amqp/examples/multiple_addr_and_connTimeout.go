func subscriber() {
	rabbitmq, _ := GetConfigByKey("rabbitmq")
	mqStr := rabbitmq.(string)
	mqURLs := strings.Split(mqStr, ",")
	l := len(mqURLs)

	var conn *amqp.Connection
	var err error

	if l == 0 {
		ulog.Errorf("config rabbitmq blank")
		return
	}

	for i := 0; i < l; i++ {
		conn, err = amqp.DialConfig("amqp:///", amqp.Config{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout("tcp", mqURLs[i], 3*time.Second) // 每次3s超时
			},
		})

		if err != nil {
			ulog.Errorf("amqp.Dial Error: %s", err)
			if i == l {
				return // 所有addr地址都尝试过失败则返回
			}
			continue
		}

		goto BIND // 没有报错，则跳到bind部分
	}

BIND:
	ulog.Info("MQ Conn OK...")
	bindAndConsume(conn)
	defer conn.Close()
}
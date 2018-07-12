// This example acts as a bridge, shoveling all messages sent from the source
// exchange "log" to destination exchange "log".

// Confirming publishes can help from overproduction and ensure every message
// is delivered.

// Setup the source of the store and forward
source, err := amqp.Dial("amqp://source/")
if err != nil {
    log.Fatalf("connection.open source: %s", err)
}
defer source.Close()

chs, err := source.Channel()
if err != nil {
    log.Fatalf("channel.open source: %s", err)
}

if err := chs.ExchangeDeclare("log", "topic", true, false, false, false, nil); err != nil {
    log.Fatalf("exchange.declare destination: %s", err)
}

if _, err := chs.QueueDeclare("remote-tee", true, true, false, false, nil); err != nil {
    log.Fatalf("queue.declare source: %s", err)
}

if err := chs.QueueBind("remote-tee", "#", "logs", false, nil); err != nil {
    log.Fatalf("queue.bind source: %s", err)
}

shovel, err := chs.Consume("remote-tee", "shovel", false, false, false, false, nil)
if err != nil {
    log.Fatalf("basic.consume source: %s", err)
}

// Setup the destination of the store and forward
destination, err := amqp.Dial("amqp://destination/")
if err != nil {
    log.Fatalf("connection.open destination: %s", err)
}
defer destination.Close()

chd, err := destination.Channel()
if err != nil {
    log.Fatalf("channel.open destination: %s", err)
}

if err := chd.ExchangeDeclare("log", "topic", true, false, false, false, nil); err != nil {
    log.Fatalf("exchange.declare destination: %s", err)
}

// Buffer of 1 for our single outstanding publishing
confirms := chd.NotifyPublish(make(chan amqp.Confirmation, 1))

if err := chd.Confirm(false); err != nil {
    log.Fatalf("confirm.select destination: %s", err)
}

// Now pump the messages, one by one, a smarter implementation
// would batch the deliveries and use multiple ack/nacks
for {
    msg, ok := <-shovel
    if !ok {
        log.Fatalf("source channel closed, see the reconnect example for handling this")
    }

    err = chd.Publish("logs", msg.RoutingKey, false, false, amqp.Publishing{
        // Copy all the properties
        ContentType:     msg.ContentType,
        ContentEncoding: msg.ContentEncoding,
        DeliveryMode:    msg.DeliveryMode,
        Priority:        msg.Priority,
        CorrelationId:   msg.CorrelationId,
        ReplyTo:         msg.ReplyTo,
        Expiration:      msg.Expiration,
        MessageId:       msg.MessageId,
        Timestamp:       msg.Timestamp,
        Type:            msg.Type,
        UserId:          msg.UserId,
        AppId:           msg.AppId,

        // Custom headers
        Headers: msg.Headers,

        // And the body
        Body: msg.Body,
    })

    if err != nil {
        msg.Nack(false, false)
        log.Fatalf("basic.publish destination: %+v", msg)
    }

    // only ack the source delivery when the destination acks the publishing
    if confirmed := <-confirms; confirmed.Ack {
        msg.Ack(false) // 是 consumer 端进行Ack，告知publisher端
    } else {
        msg.Nack(false, false)
    }
}
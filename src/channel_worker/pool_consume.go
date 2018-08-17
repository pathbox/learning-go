package channel_worker

// consume takes delivered messages from the channel and manages a worker pool
// to process tasks concurrently
func (b *Broker) consume(deliveries <-chan amqp.Delivery, concurrency int, taskProcessor iface.TaskProcessor, amqpCloseChan <-chan *amqp.Error) error {
	pool := make(chan struct{}, concurrency) // make a chan pool with concurrency number

	// initialize worker pool with maxWorkers workers
	// if concurrency is zero, no limit
	// new a goroutine to send struct{}{} to the pool
	go func() {
		for i := 0; i < concurrency; i++ {
			pool <- struct{}{}
		}
	}()

	errorsChan := make(chan error)

	for {
		select {
		case amqpErr := <-amqpCloseChan:
			return amqpErr
		case err := <-errorsChan:
			return err
		case d := <-deliveries:
			if concurrency > 0 {
				// get worker from pool (blocks until one is avaliable)
			<-pool // if pool has struct{}{}, it can go, or it blocks here until pool chan has struct{}{}
 			}

			b.processingWG.Add(1)

			// new goroutine to consume the task , so multiple tasks can be processed concurrently limit concurrent number
			if err := b.consumeOne(d, taskProcessor); err != nil {
					errorsChan <- err
				}

				b.processingWG.Done()

				if concurrency > 0 {
					// once consume is finished, then give worker back to pool
					pool <- struct{}{}
				}
			}()

		case <-b.GetStop():
			return nil
		}
	}
}

// It implements a pool modle that limit consume process concurrently with  pool chanenl, the concurrent count of goroutinues is len(pool) pool size
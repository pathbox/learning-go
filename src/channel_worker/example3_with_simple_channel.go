package channel_worker

import "net/http"

// 使用channel 构建简单的队列模式

var Queue chan Payload

func init() {
	Queue = make(chan Payload, MAX_QUEUE)
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	//...
	for _, payload := range content.Payloads {
		Queue <- payload
	}
	// ...
}

func StartProcessor() {
	for {
		select {
		case job := <-Queue:
			job.payload.UploadToS3() // <- Still not good
		}
	}
}

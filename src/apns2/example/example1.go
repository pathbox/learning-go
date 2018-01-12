package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

func main() {
	cert, err := certificate.FromP12File("/home/user/ios_dev.p12", "123456")
	if err != nil {
		log.Panic("Cert Error: ", err)
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = "84de33c97e2c22dd3790ea693bca49bef6aa0bac0614ea554db1032eae7ca48f"
	// notification.Topic = "Morning"
	// notification.Payload = []byte(`{"aps":{"alert":"Hello Morning!"}}`)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := payload.NewPayload().Alert("Hello Morning!").Badge(1).Custom("key", "val")
	notification.Payload = payload

	client := apns2.NewClient(cert).Development()
	// res, err := client.Push(notification)
	res, err := client.PushWithContext(ctx, notification)

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}

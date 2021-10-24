package main

import (
	"encoding/json"
	"os"

	"github.com/yomorun/yomo"
	"github.com/yomorun/yomo/pkg/logger"
)

type noiseData struct {
	Noise float32 `json:"noise"`
	Time int64 `json:"time"`
	From string `json:"from"`
}

func main() {
	sfn := yomo.NewStreamFunction("Noise", yomo.WithZipperAddr("localhost:9000"))

	// set only monitoring data which tag=0x33
	sfn.SetObserveDataID(0x33)

	// set handler
	sfn.SetHandler(handler)

	err := sfn.Connect()
	if err != nil {
		logger.Errorf("[flow] connect err=%v", err)
		os.Exit(1)
	}

	select {}
}

func handler(data []byte) (byte, []byte) {
	var model noiseData
	if err != nil {
		logger.Errorf("[flow] json.Marshal err=%v", err)
		os.Exit(-2)
	} else {
		logger.Printf(">> [flow] got tag=0x33, data=%# x", model)
	}
	return 0x0, nil
}
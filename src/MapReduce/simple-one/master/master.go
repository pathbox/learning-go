package master

import (
	"log"
)

var (
	MapChanIn      chan MapInput // channel produced by master while consumed by mapper
	MapChanOut     chan string   // channel produced by mapper while consumed by master
	ReduceChanIn   chan string   // channel produced by master while consumed by reducer
	ReduceChanOut  chan string   // channel produced by reducer while consumed by master
	CombineChanIn  chan string   // channel produced by master while consumed by combiner
	CombineChanOut chan []Item   // channel produced by combiner while consumed by master
)

func Handle(inputArr []string, fileDir string) []Item {
	log.Println("handle called")
	const (
		mapperNumber  int = 5
		reducerNumber int = 2
	)
	MapChanIn = make(chan MapInput)
	MapChanOut = make(chan string)
	ReduceChanIn = make(chan string)
	ReduceChanOut = make(chan string)
	CombineChanIn = make(chan string)
	CombineChanOut = make(chan []Item)

	reduceJobNum := len(inputArr)
	combineJobNum := reducerNumber

	// start combiner
	go combiner()

	// start reducer
	for i := 1; i <= reducerNumber; i++ {
		go reducer(i, fileDir)
	}

	// start mapper
	for i := 1; i <= mapperNumber; i++ {
		go mapper(i, fileDir)
	}

	go func() {
		for i, v := range inputArr {
			MapChanIn <- MapInput{
				Filename: v,
				Nr:       i + 1,
			} // pass job to mapper
		}
		close(MapChanIn) // close map input channel when no more job
	}()

	var res []Item
outter:
	for {
		select {
		case v := <-MapChanOut:
			go func() {
				ReduceChanIn <- v
				reduceJobNum--
				if reduceJobNum <= 0 {
					close(ReduceChanIn)
				}
			}()
		case v := <-ReduceChanOut:
			go func() {
				CombineChanIn <- v
				combineJobNum--
				if combineJobNum <= 0 {
					close(CombineChanIn)
				}
			}()
		case v := <-CombineChanOut:
			res = v
			break outter
		}
	}
	close(MapChanOut)
	close(ReduceChanOut)
	close(CombineChanOut)

	return res

}

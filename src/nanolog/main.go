package main

import (
	// "bufio"
	"fmt"
	"github.com/ScottMansfield/nanolog"
	"github.com/ScottMansfield/nanolog/test/pkg"
	// "log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	logWorking nanolog.Handle
)

func init() {
	logWorking = nanolog.AddLogger("Worker %i, working on task %i, attempt %i.")

	// Set up nanolog writer
	nanologout, err := os.Create("foo.log")
	if err != nil {
		panic(err)
	}
	nanolog.SetWriter(nanologout)
}

func main() {
	wg := &sync.WaitGroup{}

	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		fmt.Println("======i: ", i)
		go func(id int, wg *sync.WaitGroup) {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for i := 0; i < 10; i++ {
				var v float64
				for j := 0; i < 10; j++ {
					nanolog.Log(logWorking, id, i, j)
					v = r.Float64()
				}

				s := fmt.Sprintf("%f", v)
				nanolog.Log(pkg.LogFinished, i, v, s)
			}

			wg.Done()
		}(i, wg)
	}

	wg.Wait()
	fmt.Println("Nanolog took", time.Since(start))
	nanolog.Flush()

	// start = time.Now()
	// regoutraw, err := os.Create("foo.log")
	// if err != nil {
	// 	panic(err)
	// }
	// regout := bufio.NewWriter(regoutraw)
	// reg := log.New(regout, "", 0)

	// for i := uint8(0); i < 100; i++ {
	// 	wg.Add(1)

	// 	go func(id uint8, wg *sync.WaitGroup) {
	// 		r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 		for i := 0; i < 10; i++ {
	// 			var v float64
	// 			for j := 0; j < 10; j++ {
	// 				reg.Printf("Worker %d, working on task %d, attempt %d.", id, i, j)

	// 				v = r.Float64()
	// 			}

	// 			s := fmt.Sprintf("%f", v)

	// 			reg.Printf("Finished task %d. Result was: %f, string version %s", i, v, s)
	// 		}

	// 		wg.Done()
	// 	}(i, wg)
	// }

	// wg.Wait()

	// fmt.Println("Stdlib log took", time.Since(start))

	// regout.Flush()

}

package fanout

import (
	"errors"
	"sync"
)

// feedInputs starts a goroutine to loop through inputs and send the
// input on the interface{} channel. If done is closed, feedInputs abandons its work.
func feedInputs(done <-chan int, inputs []interface{}) (<-chan interface{}, <-chan error) {
	inputsChan := make(chan interface{})

	errChan := make(chan error, 1)

	go func() {
		defer close(inputsChan)

		errChan <- func() error {
			for _, input := range inputs {
				select {
				case inputsChan <- input:
				case <-done:
					return errors.New("loop canceled")
				}
			}
			return nil // 这个nil 是返回给 errChan的
		}()
	}()
	return inputsChan, errChan // 从inputs 数组中读取值，存储到inputsChan中，等inputsChan的值被取走后，会继续往inputsChan中存入值，否则会阻塞

	// inputsChan 是chan类型,既可以作为<-chan 返回，也可以作为 chan<-返回
}

// Worker is the interface to be implemented when using this helper package
// If the Worker func needs to have multiple params, You can wrap them into one struct,
// Also for multiple result, You can wrap them into one result struct,
// In Worker, If it return any error, All the other workers will stop immediately.
// If you want to ignore Error in some of the workers, Then return nil error in your Worker func.
// Worker 是 func 类型,可以理解为 每个worker就是一个要执行的代码块，执行完后将返回值返回.这个由外部，将代码块作为参数传入到包中执行
// 这个包只提供了 创造goroutine 进行并发执行的框架，输入+输出
type Worker func(input interface{}) (interface{}, error)

// ParallelRun starts `workerNum` of goroutines immediately to consume the value of inputs, and provide input to `Worker` func.
// and run the `Worker`, If any worker finish, it will put the result value into a channel, then append to the results value.
// The func will block the execution and wait for all goroutines to finish, then return results all together.
func ParallelRun(workerNum int, w Worker, inputs []interface{}) ([]interface{}, error) {
	// closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.

	done := make(chan int)
	defer close(done)

	inputsc, errc := feedInputs(done, inputs)

	// start a fixed number of goroutines to do the worker
	c := make(chan resultWithError)
	var wg sync.WaitGroup
	wg.Add(workerNum)

	for i := 0; i < workerNum; i++ {
		go func() {
			work(done, inputsc, c, w)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	results := []interface{}{}
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		results = append(results, r.result)
	}

	if err := <-errc; err != nil {
		return nil, err
	}
	return results, nil
}

type resultWithError struct {
	result interface{}
	err    error
}

func work(done <-chan int, inputs <-chan interface{}, c chan<- resultWithError, w Worker) {
	for input := range inputs {
		re := resultWithError{}
		re.result, re.err = w(input)
		select {
		case c <- re:
		case <-done:
			return
		}
	}
}

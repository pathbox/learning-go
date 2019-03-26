package master

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vinllen/go-logger/logger"
)

type Item struct {
	key string
	val int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].val > pq[j].val
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func combiner() {
	mp := make(map[string]int) // store the frequence of words

	// read file and do combine
	for {
		val, ok := <-CombineChanIn
		if !ok {
			break
		}
		logger.Debug("combiner called")
		file, err := os.Open(val)
		if err != nil {
			errMsg := fmt.Sprintf("Read file(%s) error in combiner", val)
			logger.Error(errMsg)
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			str := scanner.Text()
			arr := strings.Split(str, " ")
			if len(arr) != 2 {
				errMsg := fmt.Sprintf("Read file(%s) error that len of line != 2(%s) in combiner", val, str)
				logger.Warn(errMsg)
				continue
			}
			v, err := strconv.Atoi(arr[1])
			if err != nil {
				errMsg := fmt.Sprintf("Read file(%s) error that line(%s) parse error in combiner", val, str)
				logger.Warn(errMsg)
				continue
			}
			mp[arr[0]] += v
		}
		file.Close()
	}

	// heap sort
	// pq := make(PriorityQueue, len(mp))
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	for k, v := range mp {
		node := &Item{
			key: k,
			val: v,
		}
		// logger.Debug(k, v)
		heap.Push(&pq, node)
	}

	res := []Item{}
	for i := 0; i < 10 && pq.Len() > 0; i++ {
		node := heap.Pop(&pq).(*Item)
		res = append(res, *node)
	}

	CombineChanOut <- res
}

package main

import (
	"errors"
	"log"
	"sync"
)

var errorNotExist = errors.New("not exist")

func main() {
	var wg sync.WaitGroup
	wg.Add(10)

	//模拟10个并发
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("key")
			if err != nil {
				log.Print(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}

//获取数据
func getData(key string) (string, error) {
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		//模拟从db中获取数据
		data, err = getDataFromDB(key)
		if err != nil {
			log.Println(err)
			return "", err
		}

		//TOOD: set cache
	} else if err != nil {
		return "", err
	}
	return data, nil
}

//模拟从cache中获取值，cache中无该值
func getDataFromCache(key string) (string, error) {
	return "", errorNotExist
}

//模拟从数据库中获取值
func getDataFromDB(key string) (string, error) {
	log.Printf("get %s from database", key)
	return "data", nil
}

/* 缓存击穿：缓存在某个时间点过期的时候，恰好在这个时间点对这个Key有大量的并发请求过来，这些请求发现缓存过期一般都会从后端DB加载数据并回设到缓存，这个时候大并发的请求可能会瞬间把后端DB压垮
其中通过 getData(key)方法获取数据，逻辑是：

先尝试从cache中获取

如果cache中不存在就从db中获取

我们模拟了10个并发请求，来同时调用 getData 函数，执行结果如下：

2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
2020/03/08 17:13:11 get key from database
2020/03/08 17:13:11 data
*/

package main

import (
	"fmt"
	"sync"
)

/*
*
验证并发锁的使用
*/
//var mu = sync.RWMutex
var numberMap = make(map[string]int)
var wg = sync.WaitGroup{}

func main() {
	numberMap["xiaoming"] = 100
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go getAndPrint("xiaoming")
	}

	wg.Wait()
}

func getAndPrint(name string) {
	if v, ok := numberMap[name]; ok {
		fmt.Printf("%v 的 年龄为 %v\n", name, v)
	} else {
		fmt.Printf("%v 未查询到年龄", name)
	}
	wg.Done()

}

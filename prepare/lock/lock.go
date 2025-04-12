package main

import (
	"fmt"
	"sync"
)

/*
*
验证并发锁的使用
*/
var mu = sync.RWMutex;
var numberMap = make(map[string]int)
func main() {
	numberMap["age"]100;
	for i := 0; i < 10; i++ {
		go getAndPrint("age")
	}

}

func getAndPrint(age string) {
	age := numberMap[age]
	fmt.Println(age)
}

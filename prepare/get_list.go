package main

import (
	"container/list"
	"fmt"
)

/*
*
验证List双向队列数据结构
*/
func main() {
	ll := list.New()
	ll2 := list.List{}
	fmt.Println(ll.Len())
	fmt.Println(ll2.Len())
}

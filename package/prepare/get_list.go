package main

import (
	"container/list"
	"fmt"
)

func main() {
	ll := list.New()
	ll2 := list.List{}
	fmt.Println(ll.Len())
	fmt.Println(ll2.Len())
}

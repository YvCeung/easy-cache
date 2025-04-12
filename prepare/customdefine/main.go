package main

import "fmt"

// MyInt 自定义类型相当于是一种新的类型
type MyInt int

// MyAliasInt 别名仍然还是原来的类型 只是起了个别名而已
type MyAliasInt = int

func main() {
	var i MyInt
	var j MyAliasInt
	fmt.Printf("i=%d\ni type is %T\n", i, i)
	fmt.Printf("j=%d\nj type is %T", j, j)
}

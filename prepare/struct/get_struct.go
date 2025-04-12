package main

import "fmt"

/*
*
验证结构体的初始化方式
*/
func main() {
	p := person{}
	p.name = "YvCeung"
	fmt.Println(p.name)
}

type person struct {
	name string
}

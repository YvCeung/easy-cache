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

	var p1 person
	p1.name = "LiHua"
	fmt.Printf("p1.name: %v", p1.name)
}

type person struct {
	name string
}

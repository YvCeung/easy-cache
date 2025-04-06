package main

import "fmt"

func main() {
	p := person{}
	p.name = "YvCeung"
	fmt.Println(p.name)
}

type person struct {
	name string
}

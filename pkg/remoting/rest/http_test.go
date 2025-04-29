package rest

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringSub(t *testing.T) {
	str := "helloworld"

	//参数代表下标 且左闭右开，如果两个数字相当的话怎返回空
	sub1 := str[1:len(str)]
	sub2 := str[1:]
	fmt.Println(sub1)

	sub3 := str[1:1]
	sub4 := str[2:2]
	if sub1 == sub2 {
		t.Log("sub1 and sub2")
	}

	if "" == sub3 {
		fmt.Println("sub3  is blank")

	}
	if sub3 == sub4 {
		t.Log("sub3 and sub4")
	}
}

func TestSplit(t *testing.T) {
	s := "/hello/world/nihao/"
	split := strings.Split(s, "/")
	for i, s2 := range split {
		fmt.Printf("index is %d,value is %s\n", i, s2)
	}

	s2 := "helloworld"
	i := strings.Split(s2, "/")
	if len(i) != 1 {
		t.Fatalf("split fail")
	}
	if i[0] != s2 {
		t.Fatalf("result not match")
	}
}

func TestSplitN(t *testing.T) {
	s := "hello/world/nihao"
	n := strings.SplitN(s, "/", 3)
	for i, s2 := range n {
		fmt.Printf("index is %d,value is %s\n", i, s2)
	}
}

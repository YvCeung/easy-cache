package consistenthash

import (
	"strconv"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	//自定义一个hash算法
	hash := Hash(func(data []byte) uint32 {
		v, _ := strconv.Atoi(string(data))
		return uint32(v)
	})
	//初始化
	consitentHash := New(3, hash)

	//虚拟节点 2，12，22，4，14，24，6，16，26，sort==> 2,4,6,12,16,22,26
	consitentHash.Add("2", "4", "6")

	//给出测试数据
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		//超出环上最大的值，验证取模是否有效
		"27": "2",
	}

	for k, v := range testCases {
		if consitentHash.GetNode(k) != v {
			t.Errorf("When find true node for key : %s, expect : %s, actual : %s", k, v, consitentHash.GetNode(k))
		}
	}

	consitentHash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if consitentHash.GetNode(k) != v {
			t.Errorf("When find true node for key : %s, expect : %s, actual : %s", k, v, consitentHash.GetNode(k))
		}
	}

}

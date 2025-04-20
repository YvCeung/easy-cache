package easycache

/*
用来封装一个只读的缓存值
这个 ByteView 是一个只读视图，它把缓存数据封装起来，防止外部修改，
同时提供基本的操作接口（长度、转字符串、复制）
*/
type ByteView struct {
	b []byte
}

// 实现了Value接口
func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// 返回一个缓存值的副本 避免原来的缓存值被外处修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	// dest src
	copy(c, b)
	return c
}

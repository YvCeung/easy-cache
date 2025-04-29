// Package rest 以http的方式实现网络通信
package rest

import (
	"fmt"
	"github.com/YvCeung/easy-cache/pkg/cache"
	"log"
	"net/http"
	"strings"
)

// 这个相当于Java中的习惯，会给每一个请求路径统一都加个前缀
const defaultBase = "/easycache/"

type HttpPool struct {
	//标志当前节点的地址 协议+ip+端口
	selfNode string
	//请求uri,里面包含一些业务参数
	baseBizPath string
}

func NewHttpPool(selfNode string) *HttpPool {
	return &HttpPool{
		selfNode:    selfNode,
		baseBizPath: defaultBase,
	}
}

// 日志打印
func (httpPool *HttpPool) Log(format string, v ...any) {
	log.Printf("[Server %s] %s", httpPool.selfNode, fmt.Sprintf(format, v...))
}

// 实现http包下的ServeHTTP方法以提供网络服务
func (httpPool *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//校验当前请求是否合法，如果不是跟缓存相关的，直接拒绝掉
	if !strings.HasPrefix(r.URL.Path, httpPool.baseBizPath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	httpPool.Log("Received request, info is %s  %s", r.Method, r.URL.Path)

	//r.URL.Path[len(httpPool.baseBizPath):]   ==》  /_mycache/<group>/<key>  ==》<group>/<key>
	parts := strings.SplitN(r.URL.Path[len(httpPool.baseBizPath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := cache.GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())

}

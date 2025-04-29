// Package rest 以http的方式实现网络通信
package rest

import (
	"fmt"
	"github.com/YvCeung/easy-cache/pkg/cache"
	"github.com/YvCeung/easy-cache/pkg/consistenthash"
	"github.com/YvCeung/easy-cache/pkg/multinode"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// 这个相当于Java中的习惯，会给每一个请求路径统一都加个前缀
const (
	defaultBase = "/easycache/"
	//默认50个虚拟节点
	defaultReplicas = 50
)

type HttpPool struct {
	//标志当前节点的地址 协议+ip+端口
	selfNode string
	//请求uri,里面包含一些业务参数
	baseBizPath string

	//下面的这几个属性都是为了节点选择的功能而做

	//锁
	mu sync.Mutex

	//一致性hash算法中的容器
	peers *consistenthash.ConsistentHash
	//真实节点的uri跟对应的*httpGetter之间的映射关系
	httpGetters map[string]*httpGetter
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

// 相当于是客户端，提供一个从服务端获取数据的能力
type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	requestPath := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	resp, err := http.Get(requestPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned: %v", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read responsebody error: %v", err)
	}
	return data, nil
}

// 可以理解为分布式节点dispatcher
func (httpPool *HttpPool) PickPeer(key string) (multinode.PeerGetter, bool) {
	httpPool.mu.Lock()
	defer httpPool.mu.Unlock()

	if node := httpPool.peers.GetNode(key); node != "" && node != httpPool.selfNode {
		httpPool.Log("Pick peer(node) %s", node)
		return httpPool.httpGetters[node], true
	}
	return nil, false
}

//还需要提供一个初始化HttpPool里面一致性哈希算法以及node跟Getter之间映射关系等数据的方法
func (httpPool *HttpPool) Set(nodes ...string) {
	httpPool.mu.Lock()
	defer httpPool.mu.Unlock()
	//初始化一致性哈希的容器
	httpPool.peers = consistenthash.New(defaultReplicas, nil)
	httpPool.peers.Add(nodes...)

	//构建getter的映射关系
	httpPool.httpGetters = make(map[string]*httpGetter, len(nodes))
	for _, node := range nodes {
		httpPool.httpGetters[node] = &httpGetter{
			//http:xxxx//easycache/
			baseURL: node + httpPool.baseBizPath,
		}
	}
}

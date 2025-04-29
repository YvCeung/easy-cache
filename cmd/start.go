package main

import (
	"flag"
	"fmt"
	"github.com/YvCeung/easy-cache/pkg/cache"
	"github.com/YvCeung/easy-cache/pkg/remoting/rest"
	"log"
	"net/http"
)

var mockDB = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	var port int
	var api bool
	//从command获取参数并绑定到当前变量上
	flag.IntVar(&port, "port", 8001, "EasyCache server port")
	flag.BoolVar(&api, "api", false, "Whether to enable the ApiServer identifier")

	//解析参数
	flag.Parse()

	apiAddr := "http://localhost:9999"
	serverAddrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var serverAddrs []string
	for _, serverAddr := range serverAddrMap {
		serverAddrs = append(serverAddrs, serverAddr)
	}

	//创建group，group相当于是分布式缓存各种组建的一个编排汇总
	group := createGroup()

	if api {
		go startAPIServer(apiAddr, group)
	}

	startCacheServer(serverAddrMap[port], serverAddrs, group)

}

//启动cache服务
func startCacheServer(serverPort string, serverAddrs []string, group *cache.Group) {
	httpPool := rest.NewHttpPool(serverPort)
	httpPool.Set(serverAddrs...)

	group.RegisterPeerPicker(httpPool)
	log.Println("EasyCache Server is running at", serverPort)
	//httpPool实现了server包下的ServeHTTP方法，所以在收到请求后会自动执行到实现后的方法
	log.Fatal(http.ListenAndServe(serverPort[7:], httpPool))

}

// 对外暴露API服务
func startAPIServer(apiAddr string, group *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//http://localhost:9999/api?key=Tom
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("API(frontend) server is running at", apiAddr)
	//http.ListenAndServe(apiAddr[7:], nil)会一直阻塞运行，直到出错才返回
	//Fatal记录日志然后退出运行
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func createGroup() *cache.Group {
	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		//这里其实就是在cache里面找不到值的时候，会回调的一段逻辑
		func(key string) ([]byte, error) {
			log.Println("Trigger the callback: From DB search,the key is", key)
			//此处mockDb的数据
			if v, ok := mockDB[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("key: %s not exist", key)
		}))
}

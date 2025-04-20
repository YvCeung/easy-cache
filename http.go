package easycache

import (
	"fmt"
	"log"
	"net/http"
)

const defaultBasePath = "/_easycache/"

type HTTPPool struct {
	self     string
	basePath string
}

// 实现了http包下的Handler接口下的ServeHTTP方法
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s]", p.self, fmt.Sprintf(format, v...))
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

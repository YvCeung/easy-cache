package main

import (
	"log"
	"net/http"
)

type Server string

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("hello world!"))
}

// 运行后 执行 curl http://localhost:8080
func main() {
	var s Server
	http.ListenAndServe("localhost:8080", &s)
}

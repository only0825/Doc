package main

import (
	"fmt"
	"net/http"
)

type server interface {
	route() func(pattern string, HandlerFunc http.HandlerFunc)
	start() func(address string) error
}

// 当一个结构体具备接口的所有的方法的时候，它就实现了这个接口
type webserver struct {
	name string
}

func (w *webserver) route(pattern string, HandlerFunc http.HandlerFunc) {
	http.HandleFunc(pattern, HandlerFunc)
}

func (w *webserver) start(address string) error {
	return http.ListenAndServe(address, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ciao %s", r.URL.Path[1:])
}

func main() {

	obj := &webserver{
		name: "openresty",
	}
	obj.route("/", home)
	obj.start("localhost:8099")

}

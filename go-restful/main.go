package main

import (
	"io"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
)

type User struct {
	Id   string
	Name string
}

// test data
var u = User{"1", "xcbeyond"}

func main() {
	// 创建一个Web Server
	wsContainer := restful.NewContainer()

	// 创建一个WebService
	apiV1Ws := new(restful.WebService)
	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(apiV1Ws)

	// 创建请求路由，即：请求URL
	apiV1Ws.Route(apiV1Ws.GET("/hello").To(hello))
	apiV1Ws.Route(apiV1Ws.GET("/user/{id}").To(queryUser))

	restful.Add(apiV1Ws)

	// 监听端口
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func hello(req *restful.Request, resp *restful.Response) {
	name := req.QueryParameter("name")
	io.WriteString(resp, name+",hello world!")
}

func queryUser(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	if id != u.Id {
		resp.WriteErrorString(http.StatusNotFound, "User could not found.")
	} else {
		resp.WriteEntity(u)
	}
}

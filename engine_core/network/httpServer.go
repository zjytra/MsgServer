/*
创建时间: 2022/3/12 17:26
作者: zjy
功能介绍:
http服务器
*/

package network

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"net/http"
	"sync"
)

type HttpServer struct {
	ip string
	server *http.Server
	wg sync.WaitGroup
}

//创建http服务
func NewHttpServer(port int32) *HttpServer  {
 	htpServer := new(HttpServer)
	htpServer.ip = fmt.Sprintf("0.0.0.0:%v", port)
	htpServer.server = &http.Server{Addr: htpServer.ip, Handler: nil}
	return htpServer
}



func (this *HttpServer) OnInit() bool {


	return true
}

func (this *HttpServer) OnStart()  {

	go func() {
		this.wg.Add(1)
		defer this.wg.Done()
		if err := this.server.ListenAndServe(); err != nil {
			fmt.Printf("start pprof failed on %s  erro = %v\n", this.ip,err)
		}

	}()

}


func (this *HttpServer) OnRelease()  {
	this.server.Close()
	this.wg.Wait()
	xlog.Debug("HttpServer Close")
}

func (this *HttpServer)WriteJson(writer http.ResponseWriter,jsonData interface{})  {
	writer.Header().Set("Content-Type","application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	data,erro := json.Marshal(jsonData)
	if erro != nil {
		xlog.Debug("HttpServer WriteJson %s",erro.Error())
	}
	writer.Write(data)
}


func (this *HttpServer)WriteProData(writer http.ResponseWriter,m proto.Message)  {
	writer.Header().Set("Content-Type","application/octet-stream")
	writer.WriteHeader(http.StatusOK)
	data,erro := json.Marshal(m)
	if erro != nil {
		xlog.Debug("HttpServer WriteProData %s",erro.Error())
	}
	writer.Write(data)
}

//注册处理函数
func (this *HttpServer) RegisterHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	 http.HandleFunc(pattern,handler)
}


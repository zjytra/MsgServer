/*
创建时间: 2021/3/6 23:45
作者: zjy
功能介绍:

*/

package network

import (
	"bytes"
	"encoding/json"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type HttpCb func(*HttpEvent,error)



var (
	evenPool sync.Pool
	httpQe  dispatch.WaitQueue
)


func InitHttpClient()  {
	evenPool.New = func() interface{} {
		return new(HttpEvent)
	}
	httpQe = dispatch.NewAsyncQueue(500,1)
	httpQe.Start()
}

func ReleaseHttp()  {
	if httpQe == nil {
		return
	}
	httpQe.Release()
}

type HttpEvent struct {
	Url   string //
	Mark   uint64 //
	ConnId uint32
	PostBytes []byte
	PostData url.Values
	Data   interface{}
	Body   string //返回的结果
	HttpType int32 //请求类型
	cb     HttpCb //回调
}

func (this *HttpEvent)AddData(key string,val string)  {
	this.PostData.Add(key,val)
}


func (this *HttpEvent) Execute(){
	//执行对应的函数
	var rsp *http.Response
	var inerr error
	if this.HttpType == 1 { // get
		rsp, inerr = http.Get(this.Url)
	}else if  this.HttpType == 2 { //post
		rsp, inerr = http.PostForm(this.Url,this.PostData)
	}else if this.HttpType == 3 { // put
		arr,errJson := json.Marshal(this.PostData)
		if errJson != nil {
			xlog.Error("HttpASyncPut %s ,封装json失败")
			return
		}
		body := bytes.NewReader(arr)
		request, err := http.NewRequest("PUT", this.Url, body)
		if err != nil {
			xlog.Warning("DoBytesPost err=%s url=%s", err, this.Url)
		}
		rsp, inerr = http.DefaultClient.Do(request)
		if inerr != nil {
			xlog.Warning("http.Do failed,[err=%s][url=%s]", inerr, this.Url)
		}
	}

	setBody(inerr, this, rsp)
	return
}

func setBody(inerr error, this *HttpEvent, rsp *http.Response)  {
	if inerr != nil {
		xlog.Warning("get %v 错误%v", this.Url, inerr)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		xlog.Warning("get %v 数据错误%v", this.Url, err)
	}
	if dispatch.MainQueue == nil || this.cb == nil {
		xlog.Debug("%s Get请求不用回调", this.Url)
		return
	}
	this.Body = string(body)
	cb := new(HttpEventCb)
	if inerr != nil {
		cb.Erro = inerr
	} else {
		cb.Erro = err
	}
	//向主队列投递
	erro := dispatch.MainQueue.AddEvent(cb)
	if erro != nil {
		xlog.Error("异步调用http%v", erro)
	}
	return
}


func (this *HttpEvent)EvenName() string {
	return "httpEvent"
}

type HttpEventCb struct {
	httpEvent *HttpEvent
	Erro error
}


func (this *HttpEventCb) Execute(){
	//执行对应的函数
	this.httpEvent.cb(this.httpEvent,this.Erro)
	evenPool.Put(this.httpEvent)
	return
}


func (this *HttpEventCb)EvenName() string {
	return "HttpEventCb"
}

func NewHttpEvent(url string,_cb HttpCb) *HttpEvent  {
	data := evenPool.Get()
	event,ok := data.(*HttpEvent)
	if !ok || event == nil {
		xlog.Error("NewHttpEvent 失败")
		return event
	}
	event.Url = url
	event.cb =_cb
	event.PostData = make(map[string][]string)
	return event
}

//异步调用
func HttpASyncGet(event *HttpEvent)  {
	if httpQe == nil {
		xlog.Warning("没有创建队列")
		return
	}
	event.HttpType = 1
	err :=  httpQe.AddEvent(event)
	if err != nil {
		xlog.Error("异步调用http%v",err)
	}
}



//异步调用Post
func HttpASyncPostForm(event *HttpEvent)  {
	if httpQe == nil {
		xlog.Warning("没有创建队列")
		return
	}
	event.HttpType = 2
	err :=  httpQe.AddEvent(event)
	if err != nil {
		xlog.Error("异步调用http%v",err)
	}
	if err != nil {
		xlog.Error("异步调用http%v",err)
	}

}

//异步调用PostJson
func HttpASyncPutJson(event *HttpEvent)  {
	if httpQe == nil {
		xlog.Warning("没有创建队列")
		return
	}
	event.HttpType = 3
	err :=  httpQe.AddEvent(event)
	if err != nil {
		xlog.Error("异步调用http%v",err)
	}
	if err != nil {
		xlog.Error("异步调用http%v",err)
	}

}

//body提交二进制数据
//func HttpASyncPostBytes(event *HttpEvent){
//	err :=  pool.Submit(func() {
//		body := bytes.NewReader(event.PostBytes)
//		request, err := http.NewRequest("POST", event.Url, body)
//		if err != nil {
//			xlog.Warning("DoBytesPost err=%s url=%s", err, event.Url)
//		}
//		request.Header.Set("Content-Type", "application/octet-stream")
//		resp, err2 := http.DefaultClient.Do(request)
//		if err2 != nil {
//			xlog.Warning("http.Do failed,[err=%s][url=%s]", err2, event.Url)
//			return
//		}
//		defer resp.Body.Close()
//		rb, err3 := ioutil.ReadAll(resp.Body)
//		if err3 != nil {
//			xlog.Warning("http.Do failed,[err=%s][url=%s]", err3, event.Url)
//		}
//		if dispatch.MainQueue == nil || event.cb == nil {
//			xlog.Debug("%s PostForm请求不用回调",event.Url)
//			return
//		}
//		event.Body = string(rb)
//		//向主队列投递
//		erro := dispatch.MainQueue.AddEvent(event)
//		if erro != nil {
//			xlog.Error("异步调用http%v",erro)
//		}
//	})
//
//	if err != nil {
//		xlog.Error("异步调用http%v",err)
//	}
//}
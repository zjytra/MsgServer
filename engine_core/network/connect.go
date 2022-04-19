// 1.封装tcp连接,该结构体负责写的工作，
// 2.读的协程 交给另一对象处理,方便其他对象知道
// 连接对象,网络连接是在Accept 或 作为服务器连接对象时单独开的协程
// 网络关闭
//    服务器主动关闭由逻辑线程通知写协程关闭tcp,读协程就会收到EOF,再关闭的时候向写协程通知,发现已经关闭就不再处理
//    客户端使连接被动关闭由读协程通知写协程关闭
//    连接与关闭在两个协程中处理
//    读又是单独的一个协程
package network

import (
	"errors"
	"fmt"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"net"
	"strings"
	"sync"
)

type Connect struct {
	conn         net.Conn
	connID       uint32      // 服务器创建连接时生成id
	sync.RWMutex             // 主要作用,防止向关闭后的通道中写入数据
	closeFlag    bool        // 检测关闭标志
	writeChan    chan []byte // 写的通道，我服务器写的消息先写入通道再用连接传出去
	// recMutex    sync.RWMutex // 接收锁
	lastRecTime int64 // 最後一次收包时间
	appID       int32 // 保存appID
	_msgParser  *MsgParser // 数据包解析对象
	onCloseFun  func()     //关闭连接通知方法
	isConnect   bool       //是否连接上
}


//公共初始化方法
func (this *Connect) initConnData(conn net.Conn, connId uint32,appID int32, msgParser *MsgParser,onCloseFun  func() ) error {
	if msgParser == nil {
		return errors.New("数据解析对象为null")
	}
	this.conn = conn
	this.connID = connId
	this.appID = appID
	cfg := this.GetNetCfg()
	if cfg == nil {
		return errors.New("网络配置数据为nil")
	}
	this.writeChan = make(chan []byte, cfg.Write_cap_num)
	this._msgParser = msgParser //使用外面传入的对象，所有连接共用一个对象，需要分开是才分开
	this.lastRecTime = timeutil.GetCurrentTimeS()
	this.onCloseFun = onCloseFun
	go this.sendData() // 写协程

	return nil
}

//获取连接id
func (this *Connect) GetConnID() uint32 {
	return this.connID
}

func (this *Connect) Destroy() {
	this.Lock()
	this.doDestroy()
	this.Unlock()
}

func (this *Connect) doDestroy() {
	erro := this.conn.(*net.TCPConn).SetLinger(0)
	if erro != nil {
		xlog.Error("doDestroy 错误 %v", erro)
	}
	erro = this.conn.Close()
	if erro != nil {
		xlog.Error("doDestroy 关闭连接错误 %v", erro)
	}

	if !this.closeFlag {
		this.closeFlag = true
		close(this.writeChan)
	}
}

func (this *Connect) Close() {
	this.Lock()
	defer this.Unlock()
	// 已经关闭
	if this.closeFlag {
		xlog.Debug("Connect 当前连接已经关闭")
		return
	}
	this.closeFlag = true
	this.doWrite(nil)
}

// b 字节数在其他协程不能修改,因为在另一个线程写
func (this *Connect) Write(b []byte) {
	this.Lock()
	// 已经关闭
	if this.closeFlag || b == nil {
		this.Unlock()
		xlog.Debug("当前连接状态 %v,", this.closeFlag)
		return
	}
	this.doWrite(b)
	this.Unlock()
}

func (this *Connect) doWrite(b []byte) {
	// 写的队列被撑满时
	if len(this.writeChan) == cap(this.writeChan) {
		xlog.Debug("close tcpConn: channel full")
		this.doDestroy() // 这里要主动断开避免阻塞 当前调用协程
		return
	}
	this.writeChan <- b
}

// 取通道的数据给连接
func (this *Connect) sendData() {
	// 这里接收写的通道，没有数据会一直阻塞，直到通道关闭 或收到nil数据
	for b := range this.writeChan {
		if b == nil {
			break
		}
		xlog.Debug("发送的数据 %v", b)
		_, err := this.conn.Write(b)
		if err != nil {
			break
		}
	}

	// 主动关闭,被动关闭关闭最后都要走这里
	erro := this.conn.Close()
	xlog.Debug("sendData this.conn 连接关闭")
	if erro != nil {
		xlog.Error("关闭连接错误 %v", erro)
	}
	// 已经关闭
	this.Lock()
	this.closeFlag = true
	this.Unlock()
	this.onCloseFun() //主要通知其他模块子类实现确保网络关闭只调一次
}

func (this *Connect) Read(b []byte) (int, error) {
	return this.conn.Read(b)
}

func (this *Connect) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *Connect) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

// 一起写多个数据包
// 每个包的数据 由GetOneMsgByteArr构建
func (this *Connect) WriteMsg(args ...[]byte)  {
	buf, erro := this._msgParser.MorePackageToOne(args...)
	if erro != nil {
		xlog.Debug("%v",erro)
		return
	}
	this.Write(buf)
}

// 多个包构成成一个包
func (this *Connect) ConnMorePackageToOne(args ...[]byte) ([]byte, error) {
	buf, erro := this._msgParser.MorePackageToOne(args...)
	if erro != nil {
		return nil, erro
	}
	return buf, nil
}

// 是否存活 没有存活会被提下线
func (this *Connect) IsAlive() bool {
	currentTime := timeutil.GetCurrentTimeS()
	temTime := this.lastRecTime + 5
	//5秒内没有明确身份就搞掉
	if currentTime > temTime && !this.isConnect {
		return false
	}
	cfg := this.GetNetCfg()
	if cfg  == nil {
		return  true
	}
	temTime = int64(cfg.Checklink_s) + this.lastRecTime
	return currentTime < temTime
}

//判断连接是否关闭
func (this *Connect) IsClose() bool {
	this.RLock()
	isClose := this.conn == nil || this.closeFlag
	this.RUnlock()
	return isClose
}


func (this *Connect) SetContAck() {
	this.isConnect = true
}

//获取网络配置,动态获取,方便更新
func (this *Connect) GetNetCfg() *csvdata.NetWorkCfg {
	cfg := csvdata.GetNetWorkCfgPtr(this.appID)
	if cfg == nil {
		xlog.Error("appid %v 未找到网络配置 ",this.appID)
		return nil
	}
	return cfg
}

func (this *Connect) GetCfgKind() int32 {
	cfg := this.GetNetCfg()
	if cfg == nil {
		return 0
	}
	return cfg.App_kind
}


func (this *Connect) RemoteAddrStr() string {
	if this.conn == nil {
		return ""
	}
	return this.conn.RemoteAddr().String()
}

func (this *Connect) RemoteAddrIp() string {
	if this.conn == nil {
		return ""
	}
	addr := this.conn.RemoteAddr().String()
	addrArr := strings.Split(addr,":")
	if addrArr == nil || len(addrArr) < 1 {
		return  ""
	}
	return addrArr[0]
}


func (this *Connect) RemoteAddrPort() string {
	if this.conn == nil {
		return ""
	}
	addr := this.conn.RemoteAddr().String()
	addrArr := strings.Split(addr,":")
	if addrArr == nil || len(addrArr) < 2 {
		return  ""
	}
	return addrArr[1]
}



func (this *Connect) CheckConn(maincmd uint32) error {
	if this.IsClose()  {
		return errors.New(fmt.Sprintf("消息id %d 连接为nil ", maincmd))
	}
	return nil
}

package network

import (
	"encoding/binary"
	"errors"
	"github.com/zjytra/MsgServer/devlop/xutil"
	"github.com/zjytra/MsgServer/devlop/xutil/mathutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"io"
	"math"
)

// 包格式
// --------------
//  | int32 包总长 | byte 加密key| uint16主命令 | uint16字命令 | int32datalen | Data |
// --------------
const(
	MIN_MSG_LEN = 9
)

type MsgParser struct {
	msgLenSize int32 // 包长度字节大小
	minMsgLen  int32 // 最小长度 包头长度
	maxMsgLen  int32 // 最大长度
	isEncrypt  bool   // 是否加密
	//加密key位置 在长度后
    KEY_POS int   //key的位置
}


var (
	ByteOrder binary.ByteOrder // 设置网址字节序
)
func init() {
	// 设置大端序列或小端序列
	if xutil.IsLittleEndian() {
		ByteOrder = binary.LittleEndian
	} else {
		ByteOrder = binary.BigEndian
	}
}

func NewMsgParser(maxMsgLen int32, isencrypt bool) *MsgParser {
	p := new(MsgParser)
	var head TcpMsgHead
	p.msgLenSize = 4
	//加密key的位置
	p.KEY_POS = int(p.msgLenSize + 1)
	p.minMsgLen = int32(binary.Size(head)) //包头需要长度  也是一个最小包的长度
	if p.minMsgLen > MIN_MSG_LEN { // 不能超过设置的
		p.minMsgLen = MIN_MSG_LEN
	}
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}
	if p.maxMsgLen > math.MaxInt32 {
		p.maxMsgLen = math.MaxInt32
	}
	p.isEncrypt = isencrypt
	return p
}

// goroutine safe
func (p *MsgParser) DoRead(conn Conner) ([]byte, error) {
	// 根据长度字节大小解析第一个长度
	msgLenBuf := make([]byte, p.msgLenSize) // TODO 可以优化接收消息的buf
	n, err :=  conn.Read(msgLenBuf)  //Peek只有这样用 bufio.NewReader(f).Peek(64)
	if err != nil {
		return nil, err
	}
	if p.msgLenSize != int32(n) {
		return nil, errors.New("读取包长度出错")
	}
	// parse len
	msgLen := int32(ByteOrder.Uint32(msgLenBuf))
	// check len
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}
	// Data 读取
	// 已经提取了消息长度的字节 剩余 加密key| 主命令 | 字命令 | datalen | Data | 这里扣除长度的字节
	backlen := msgLen - p.msgLenSize
	msgdata := make([]byte, backlen) // TODO 可以优化接收消息的buf
	// conn.DoRead(msgdata)

	//把剩余的读完
	readlen, errfull := io.ReadFull(conn,msgdata)
	xlog.Debug( " 收到原始消息 %v",msgdata)
	if errfull != nil {
		return nil, err
	}
	readlen = readlen
	unreadlen := mathutil.MaxInt32(0, backlen-int32(readlen))
	if unreadlen == 0 {
		return msgdata, nil
	}
	// // 出现断包再读一次  不处理断包避免不明连接占用连接
	// unreadData := msgdata[readlen:backlen:backlen]  //获取未填满的容器
	// rereadlen, err := tcpConn.DoRead(unreadData) // 未读数据
	// if err != nil {
	// 	return nil, err
	// }
	// //如果还不能填满长度
	// if uint32(rereadlen) != unreadlen {
	// 	return nil,errors.New("断包处理失败")
	// }
	// return msgdata, nil
	return nil, errors.New("读取数据出现断包")
}

// 加密key| 主命令 | 字命令 | datalen | Data |
func (p *MsgParser) UnpackOne(readb []byte) (maincmd uint32, msg []byte, err error) {
	msglen := len(readb) // 包头的长度已经被解掉
	if readb == nil || msglen == 0 {
		return
	}
	reader := NewNetBuffer(readb,0) //已经把消息长度的4个字节丢弃
	key := reader.ReadUint8()         // 解析加密key
	if p.isEncrypt {
		reader.setEncrypt(reader._readPos, key)         // 将数据 解密  用key解密
	}
	maincmd = reader.ReadUint32()
	if msglen <= int(p.minMsgLen - p.msgLenSize) {      //刚好解完
		return
	}
	msg = reader.ReadBytes()
	return
}


// 打单包
func (p *MsgParser) PackOne(maincmd uint32, msg []byte) ([]byte, error) {
	var msgLen,datalen int32
	msgLen = p.minMsgLen //最小长度
	if msg != nil && len(msg) > 0 { //空消息判断
		datalen = int32(len(msg))
		msgLen += 4 //还需要加上数据长度的大小
		msgLen += datalen
	}
	//包头长度加上 所有数据长度的字节数 还需要加上
	// check len
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}
	// 先查看是否 加密好计算长度
	var bkey byte
	if p.isEncrypt { // 只有加密了才生成
		bkey = byte(xutil.RandInterval(1, 254))
	}
	writeBuf := make([]byte, msgLen)
	writer := NewNetBuffer(writeBuf,int(p.msgLenSize)) //已经把加密部分丢弃
	writer.WriteUint32(uint32(msgLen)) //长度
	// 默认把加密key写上
	writer.WriteUInt8(bkey)
	writer.WriteUint32(maincmd) // 主命令
	if datalen > 0 { //有数据才写
		writer.WriteBytes(msg)     // 最后写proto数据
	}
	if p.isEncrypt {
		writer.setEncrypt(p.KEY_POS,bkey)		// 将key后面的数据加密
	}
	return writeBuf, nil
}

// 将多个包合并成一个
func (p *MsgParser) MorePackageToOne(args ...[]byte) ( []byte,error ){
	if args == nil {
		return  nil,errors.New("MorePackageToOne args is nil")
	}
	var msgLen int32
	// 计算消息长度
	for i := 0; i < len(args); i++ {
		if args[i] == nil {
			continue
		}
		msgLen += int32(len(args[i]))
	}
	// check len
	if msgLen > p.maxMsgLen {
		return nil,errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil,errors.New("message too short")
	}
	// 构建所有数据包
	buf := make([]byte,msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		if args[i] == nil {
			continue
		}
		copy(buf[l:], args[i])
		l += len(args[i])
	}
	return buf,nil
}


// // param |加密key| 主命令 | 字命令 | datalen | Data |
// func (p *MsgParser) UnpackOne(readb []byte) (maincmd uint16, subcmd uint16, msg []byte, err error) {
// 	msglen := len(readb)
// 	if readb == nil || msglen == 0 {
// 		return
// 	}
// 	parseBuf := make([]byte, msglen-1) // 存儲key 后面的数据
// 	copy(parseBuf, readb[1:])          // 将后面的命令放到一个桶中
// 	key := readb[0]                    // 加密key
// 	if key > 0 { //由加密key才解密
// 		setEncrypt(parseBuf, key)          // 将数据 解密 第一个是key 用key解密
// 	}
//
// 	// 保存数据到缓冲区
// 	reader := p.tcpConn.reader
// 	reader.Reset()
// 	reader.Write(parseBuf)
// 	err = binary.DoRead(reader, ByteOrder, &maincmd) // 解析主命令
// 	if err != nil {
// 		return
// 	}
// 	err = binary.DoRead(reader, ByteOrder, &subcmd) // 子命令
// 	if err != nil {
// 		return
// 	}
// 	var datalen uint32
// 	err = binary.DoRead(reader, ByteOrder, &datalen)
// 	if err != nil {
// 		return
// 	}
// 	msg = make([]byte, datalen)
// 	err = binary.DoRead(reader, ByteOrder, &msg)
// 	return
// }
//
// // 打单包
// func (p *MsgParser) PackOne(maincmd uint32, msg []byte) ([]byte, error) {
// 	datalen := uint32(len(msg))
// 	// 加密key
// 	bkey := byte(xutil.RandInterval(1, 254))
// 	// head := TcpMsgHead{
// 	// 	MKey: bkey,
// 	// 	MainCmd: maincmd,
// 	// 	SubCmd:  subcmd,
// 	// 	Datalen: datalen,
// 	// }
// 	msgLen := p.minMsgLen + datalen
// 	// check len
// 	if msgLen > p.maxMsgLen {
// 		return nil, errors.New("message too long")
// 	} else if msgLen < p.minMsgLen {
// 		return nil, errors.New("message too short")
// 	}
// 	writer := p.tcpConn.writer
// 	writer.Reset()
// 	// 为了加密一个一个写
// 	erro := binary.Write(writer, ByteOrder, maincmd) // 写头
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	erro = binary.Write(writer, ByteOrder, subcmd) // 写头
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	erro = binary.Write(writer, ByteOrder, datalen) // 写头
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	erro = binary.Write(writer, ByteOrder, msg) // 写数据
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	databuf := make([]byte, writer.Len())
// 	copy(databuf, writer.Bytes())
// 	setEncrypt(databuf, bkey)                    // 将数据 加密
// 	writer.Reset()                               // 重置buf变量
// 	p.WriteLen(writer, msgLen)                   // 写长度  长度与key 不加密
// 	erro = binary.Write(writer, ByteOrder, bkey) // 写加密数据
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	erro = binary.Write(writer, ByteOrder, databuf) // 最后写数据
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	return writer.Bytes(), erro
// }
// func (p *MsgParser) WriteLen(writer *bytes.Buffer, msgLen uint32) {
// 	switch p.msgLenSize {
// 	case 2: // uint16
// 		binary.Write(writer, ByteOrder, uint16(msgLen))
// 	case 4: // uint32
// 		binary.Write(writer, ByteOrder, uint32(msgLen))
// 	}
// }

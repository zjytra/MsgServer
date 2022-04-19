//网络数据传输的buf封装
package network

import (
	"encoding/binary"
	"github.com/zjytra/MsgServer/devlop/xutil/mathutil"
	"math"
)

type NetBuffer struct {
	ByteOrder binary.ByteOrder // 设置网址字节序
	data      []byte
	_readPos  int //读的位置
	_writePos int //写的位置
	msgLensize int
}

func NewNetBuffer(buffer []byte,lenSize int) *NetBuffer {
	netbuf := new(NetBuffer)
	netbuf.ByteOrder = ByteOrder
	netbuf.data = buffer
	netbuf.msgLensize = lenSize
	return netbuf
}

//获取buffer原始数据
func (this *NetBuffer) SetData(bufer []byte) {
	this._readPos = 0
	this._writePos = 0
	this.data = bufer
}

//获取buffer原始数据
func (this *NetBuffer) GetData() []byte {
	return this.data
}

//剩余写的空间
func (this *NetBuffer) Space() int {
	return len(this.data) - this._writePos
}

//获取拷贝数据
func (this *NetBuffer) GetCpData() []byte {
	buf := make([]byte, len(this.data))
	copy(buf, this.data)
	return buf
}

func (this *NetBuffer) ReadInt8() int8 {
	data := this.data[this._readPos]
	this._readPos++
	return int8(data)
}

// uint8 = byte
func (this *NetBuffer) ReadUint8() uint8 {
	data := this.data[this._readPos]
	this._readPos++
	return data
}

func (this *NetBuffer) ReadInt16() int16 {
	data := this.ByteOrder.Uint16(this.data[this._readPos : this._readPos+2])
	this._readPos += 2
	return int16(data)
}

func (this *NetBuffer) ReadUint16() uint16 {
	data := this.ByteOrder.Uint16(this.data[this._readPos : this._readPos+2])
	this._readPos += 2
	return data
}
func (this *NetBuffer) ReadInt32() int32 {
	data := this.ByteOrder.Uint32(this.data[this._readPos : this._readPos+4])
	this._readPos += 4
	return int32(data)
}

func (this *NetBuffer) ReadUint32() uint32 {
	data := this.ByteOrder.Uint32(this.data[this._readPos : this._readPos+4])
	this._readPos += 4
	return data
}

func (this *NetBuffer) ReadInt64() int64 {
	data := this.ByteOrder.Uint64(this.data[this._readPos : this._readPos+8])
	this._readPos += 8
	return int64(data)
}
func (this *NetBuffer) ReadUint64() uint64 {
	data := this.ByteOrder.Uint64(this.data[this._readPos : this._readPos+8])
	this._readPos += 8
	return data
}

func (this *NetBuffer) ReadFloat() float32 {
	data := this.ByteOrder.Uint32(this.data[this._readPos : this._readPos+4])
	floatData := math.Float32frombits(data)
	this._readPos += 4
	return floatData
}

func (this *NetBuffer) ReadDouble() float64 {
	data := this.ByteOrder.Uint64(this.data[this._readPos : this._readPos+8])
	floatData := math.Float64frombits(data)
	this._readPos += 8
	return floatData
}

func (this *NetBuffer) ReadBytes() []byte {
	readlen := this.ReadInt32() //先读长度
	if readlen == 0 {
		return nil
	}
	buf := make([]byte, readlen)
	this.ReadToBuf(buf)
	return buf
}

func (this *NetBuffer) ReadToBuf(buf []byte) {
	buflen := len(buf)
	for i := 0; i < buflen; i++ {
		buf[i] = this.data[this._readPos]
		this._readPos++
	}
}

func (this *NetBuffer) ReadString() string {
	readlen := this.ReadInt32() //先读长度
	if readlen == 0 {
		return ""
	}
	buf := make([]byte, readlen)
	this.ReadToBuf(buf)
	return string(buf)
}

func (this *NetBuffer) WriteInt8(data int8) {
	this.data[this._writePos] = uint8(data)
	this._writePos++
}

// uint8 = byte
func (this *NetBuffer) WriteUInt8(data uint8) {
	this.data[this._writePos] = data
	this._writePos++
}

func (this *NetBuffer) WriteInt16(data int16) {
	this.ByteOrder.PutUint16(this.data[this._writePos:this._writePos+2], uint16(data))
	this._writePos += 2
}

func (this *NetBuffer) WriteUint16(data uint16) {
	this.ByteOrder.PutUint16(this.data[this._writePos:this._writePos+2], data)
	this._writePos += 2
}

func (this *NetBuffer) WriteInt32(data int32) {
	this.ByteOrder.PutUint32(this.data[this._writePos:this._writePos+4], uint32(data))
	this._writePos += 4
}

func (this *NetBuffer) WriteUint32(data uint32) {
	this.ByteOrder.PutUint32(this.data[this._writePos:this._writePos+4], data)
	this._writePos += 4
}

func (this *NetBuffer) WriteInt64(data int64) {
	this.ByteOrder.PutUint64(this.data[this._writePos:this._writePos+8], uint64(data))
	this._writePos += 8
}
func (this *NetBuffer) WriteUint64(data uint64) {
	this.ByteOrder.PutUint64(this.data[this._writePos:this._writePos+8], data)
	this._writePos += 8
}

func (this *NetBuffer) WriteFloat(floatData float32) {
	data := math.Float32bits(floatData)
	this.ByteOrder.PutUint32(this.data[this._writePos:this._writePos+4], data)
	this._writePos += 4
}

func (this *NetBuffer) WriteDouble(floatData float64) {
	data := math.Float64bits(floatData)
	this.ByteOrder.PutUint64(this.data[this._writePos:this._writePos+8], data)
	this._writePos += 8
}

func (this *NetBuffer) WriteBytes(data []byte) {
	if data == nil {
		return
	}
	datalen := len(data)
	if datalen == 0 {
		return
	}
	this.WriteInt32(int32(datalen)) //先写长度
	this.WriteDataTobuf(data)
}

func (this *NetBuffer) WriteDataTobuf(data []byte) {
	buflen := len(data)
	for i := 0; i < buflen; i++ {
		this.data[this._writePos] = data[i]
		this._writePos++
	}
}

func (this *NetBuffer) WriteString(data string) {
	datalen := len(data)
	if datalen == 0 {
		return
	}
	bufdata := []byte(data)
	this.WriteInt32(int32(len(bufdata))) //先写长度
	this.WriteDataTobuf(bufdata)
}


// 加密函数
func (this *NetBuffer)setEncrypt(keyIndex int, bkey byte) {
	if	bkey == 0 {
		return
	}
	datlen := len(this.data)
	encryptlen := datlen - 1 - this.msgLensize //把不加密的减了
	encryptlen = mathutil.MinInt(encryptlen,50) //加密优化，避免加密太长
	for i := 0; i < encryptlen; i++ {
		if (keyIndex >= datlen) {
			break
		}
		this.data[keyIndex] = this.data[keyIndex] ^ bkey
		keyIndex++
	}
}
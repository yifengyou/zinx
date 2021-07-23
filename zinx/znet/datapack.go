package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

//封包、拆包的具体模块
type DataPack struct {
}

// 拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节） + ID unit32（4字节）
	return 8
}

// 封包方法
//|datalen|msgID/data/
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放byte自己的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 将datalen写进databuf中
	// 二进制写法
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 将MsgId写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data数据写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法
// 先读头，再读内容。根据头部数据长度读取内容
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息得到datalen和MsgID
	msg := &Message{}

	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断datalen是否超出最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		fmt.Println("msg.DataLen", msg.DataLen)
		fmt.Println("maxPackageSize", utils.GlobalObject.MaxPackageSize)
		return nil, errors.New("too large message data receive!")
	}

	return msg, nil
}

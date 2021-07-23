package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

/*
	连接模块
*/
type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 当前连接ID
	ConnID uint32

	// 当前连接状态
	isClosed bool

	// 当前连接锁绑定的处理业务方法
	//handleAPI ziface.HandlerFunc

	// 告知当前连接已经退出/停止的channel
	ExitChan chan bool

	//该连接处理的方法Router
	Router ziface.IRouter
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit,remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("Recve buf err", err)
		//	continue
		//}
		//// 调用当前连接所绑定的HandlerAPI
		//if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("Co", c.ConnID, "handle is err", err)
		//	break
		//}

		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的MsgHead，8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read from client error!", err)
			break
		}
		// 拆包，得到msgID 和 msgDatalen 放在一个msg对象中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据datalen再次读取Data，放在msg.Data字段中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg err ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		// 从路由中，找到注册绑定的Conn对应的router调用
	}
}

// 启动连接，让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID=", c.ConnID)
	// 启动从当前连接的读数据的业务
	// TODO 启动从当前连接写数据的业务
	go c.StartReader()
	// 阻塞
}

// 停止当前连接，结束当前连接工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop() ... ConnID=", c.ConnID)

	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	c.Conn.Close()

	// 回收资源
}

// 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//func (c *Connection) Send(data []byte) error {
//	panic("implement me")
//}

// 提供一个SendMsg方法，将要发送给客户端的数据先封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send message")
	}
	// 将data进行封包,MsgDtaLen/MsgID/MsgData
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("package message error", err, " msgID=", msgId)
		return errors.New("package message err")
	}
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Wirte msg id=", msgId, "error", err)
		return errors.New("conn Write err")
	}
	return nil
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:   conn,
		ConnID: connID,
		//handleAPI: callback_api,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

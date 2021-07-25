package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
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
	// handleAPI ziface.HandlerFunc

	// 告知当前连接已经退出/停止的channel
	// 由Reader告知Writer退出，读都不行肯定要关写
	ExitChan chan bool

	// 该连接处理的方法Router
	// Router ziface.IRouter

	// 消息的管理MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	// 添加无缓冲管道用于读写goroutine通信
	msgChan chan []byte

	// zinx_v0.9添加连接管理模块
	// 当前conn隶属于那个server
	TcpServer ziface.IServer

	// 连接属性集合
	property map[string]interface{}

	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

// 初始化连接模块的方法
// func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		//handleAPI: callback_api,
		//Router:   router,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	// zinx_v0.9添加连接管理模块
	// 将conn加入到ConnMaanger中
	// 通过c找到server，绑定关系
	c.TcpServer.GetConnMgr().Add(c)

	return c
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
		//go func(request ziface.IRequest) {
		//c.Router.PreHandle(request)
		//c.Router.Handle(request)
		//c.Router.PostHandle(request)
		//}(&req)

		// 从路由中，找到注册绑定的Conn对应的router调用
		// 根据绑定好的MsgID找到对应的处理API业务
		//go c.MsgHandler.DoMsgHandler(&req)

		// 判断是否开启工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启工作池机制
			// 放到工作池中处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)

		}
	}
}

/*
	写消息GoRoutine，专门发送给客户端消息的模块
	客户端将要写的消息发送给写者就行
*/
func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")

	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")
	// 不断阻塞等待channel消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有消息写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err!", err)
				return
			}
		case <-c.ExitChan:
			// reader告诉退出writer
			return
		}
	}
}

// 启动连接，让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID=", c.ConnID)
	// 启动从当前连接的读数据的业务
	go c.StartReader()
	// 启动从当前连接写数据的业务
	go c.StartWriter()

	//zinx_v0.9添加连接管理模块
	// 按照开发者传递进来的创建连接之后需要调用的处理业务，执行对应的hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止当前连接，结束当前连接工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop() ... ConnID=", c.ConnID)

	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//zinx_v0.9添加连接管理模块
	// 按照开发者传递进来的销毁连接之后需要调用的处理业务，执行对应的hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 告知writer关闭
	c.ExitChan <- true

	//zinx_v0.9添加连接管理模块
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

	// 读者肩负写任务
	//if _, err := c.Conn.Write(binaryMsg); err != nil {
	//	fmt.Println("Wirte msg id=", msgId, "error", err)
	//	return errors.New("conn Write err")
	//}

	// 将数据发给写客户端
	c.msgChan <- binaryMsg

	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 添加属性
	c.property[key] = value
}

// 获取连接里属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found!")
	}
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 删除属性
	delete(c.property, key)
}

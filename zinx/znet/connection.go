package znet

import (
	"fmt"
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
	handleAPI ziface.HandlerFunc

	// 告知当前连接已经退出/停止的channel
	ExitChan chan bool
}


// 连接的读业务方法
func (c *Connection) StartReader(){
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit,remote addr is",c.RemoteAddr().String() )
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt,err := c.Conn.Read(buf)
		if err!=nil {
			fmt.Println("Recve buf err", err)
			continue
		}
		// 调用当前连接所绑定的HandlerAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("Co", c.ConnID, "handle is err", err)
			break
		}
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

func (c *Connection) Send(data []byte) error {
	panic("implement me")
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandlerFunc) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		handleAPI: callback_api,
		isClosed: false,
		ExitChan: make(chan bool,1),
	}
	return c
}

package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

//iServer的接口实现，定义一个Server的服务模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

// 定义当前客户端连接所绑定的handle api（目前这个handle是写死的，以后优化应该有用户去自定义这个handle）
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显业务
	fmt.Println("[Conn Handle] CallbackToClient ...")
	if _,err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP:%s, Port %d is starting", s.IP, s.Port)
	go func() {
		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error! ", err)
			return
		}
		// 2. 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server success", s.Name, "Listenning...")
		var cid uint32
		cid = 0
		// 3. 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {

			// 如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, cid, CallBackToClient  )
			cid ++

			// 启动当前连接业务处理
			go dealConn.Start()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器资源，状态或者已经开辟的连接信息，进行停止或者回收

}

//运行服务
func (s *Server) Serve() {
	// 启动server服务功能
	s.Start()
	// 做一些启动服务器之后额外的业务

	// 调用Start之后进入阻塞状态，否则程序会直接退出
	select {}
}

/*
	初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8888,
	}
	return s
}

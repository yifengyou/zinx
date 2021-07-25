package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test PreRouter
//func (this *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
//	if err != nil {
//		fmt.Println("call back before ping err", err)
//	}
//}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	//_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	//if err != nil {
	//	fmt.Println("call back ping err", err)
	//}

	// 先读取客户端数据，再回写ping...ping
	fmt.Println("recv from client : msgID=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// Test PostHandle
//func (this *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router Handle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
//	if err != nil {
//		fmt.Println("call back after ping err", err)
//	}
//}

// hello test 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	//_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	//if err != nil {
	//	fmt.Println("call back ping err", err)
	//}

	// 先读取客户端数据，再回写ping...ping
	fmt.Println("recv from client : msgID=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("hello welcome to zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建连接之执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is called ...")
	if err := conn.SendMsg(202, []byte("DoCennection BEGIN")); err != nil {
		fmt.Println(err)
	}

}

// 关闭连接之执行钩子函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is called ...")
	fmt.Println("conn ID=", conn.GetConnID(), "is lost...")
}

func main() {
	// 创建一个server句柄
	s := znet.NewServer("[zinx V0.10]")

	//zinx_v0.9添加连接管理模块
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 添加一个Router，自定义的Router
	//s.AddRouter(&PingRouter{})
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//启动server
	s.Serve()
}

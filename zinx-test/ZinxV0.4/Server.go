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
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping err", err)
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping err", err)
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping err", err)
	}
}

func main() {
	// 创建一个server句柄
	s := znet.NewServer("[zinx V0.1]")

	// 添加一个Router，自定义的Router
	s.AddRouter(&PingRouter{})

	//启动server
	s.Serve()
}

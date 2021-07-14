package main

import "zinx/znet"

var (
	fuck = make([]int, 5)
)


func main() {

	// 创建一个server句柄
	s := znet.NewServer("[zinx V0.1]")
	//启动server
	s.Serve()
}

package ziface

//定义一个服务器接口
type IServer interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	//启动服务
	Serve()
}
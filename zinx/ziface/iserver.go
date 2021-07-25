package ziface

//定义一个服务器接口
type IServer interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	//启动服务
	Serve()

	//路由功能：给当前服务注册一个路由方法，供客户端的连接处理使用
	//AddRouter(router IRouter)
	AddRouter(msgID uint32, router IRouter)
	//zinx_v0.9添加连接管理模块
	GetConnMgr() IConnManager
	// 注册OnStart钩子函数方法
	SetOnConnStart(func (connection IConnection))
	// 注册OnStop钩子函数方法
	SetOnConnStop(func (connection IConnection))
	// 调用OnStart钩子函数方法
	CallOnConnStart(connection IConnection)
	// 调用OnStop钩子函数方法
	CallOnConnStop(connection IConnection)

}

package ziface


/*
	消息管理抽象层
 */
type IMsgHandle interface {
	// 调度/执行路由器Router消息处理方法
	DoMsgHandler(request IRequest)

	// 为消息添加具体的处理逻辑 添加路由器
	AddRouter(msgID uint32, router IRouter)
}


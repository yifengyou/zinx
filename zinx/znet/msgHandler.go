package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

/*
	消息处理模块的实现
*/
type MsgHandle struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 调度/执行路由器Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1. 从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), "is NOT FOUND! need register")
	}
	// 2. 根据map中找到route，调用三个方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑 添加路由器
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1. 判断，当前msg绑定的API处理方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID=" + strconv.Itoa(int(msgID)))
	}

	// 2. 添加msg与API的对应关系
	// map写不加锁？？？
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID=", msgID, "success!")
}

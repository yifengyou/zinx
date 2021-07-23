package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
	消息处理模块的实现
*/
type MsgHandle struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter

	// 消息队列，负责取任务的消息队列
	TaskQueue []chan ziface.IRequest

	// 业务工作池的worker数量
	WorkerPoolSize uint32
}

// 初始化/创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 将消息交给TaskQueue，由worker执行
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1. 消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	// 基本平均分配法则
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		"request MsgID=", request.GetMsgID(),
		"to WorkerID", workerID)
	// 2. 将消息发送给对应worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
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

// 启动一个worker工作池（开启工作池动作只能发生一次）
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 当前的worker对应的channel消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker((uint32)(i), mh.TaskQueue[i])
	}
}

// 启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID uint32, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=", workerID, "is started...")

	// 不间断阻塞等待对应消息到来
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，指定当前Request绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}

}

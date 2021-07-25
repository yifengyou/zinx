package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/*
	当前连接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理连接的集合
	connLock    sync.RWMutex                  // 保护连接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源，添加则加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到ConnMager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID=", conn.GetConnID(), "connection add to ConnManager successfully: conn num =", connMgr.Len())
}

// 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源，添加则加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID=", conn.GetConnID(), "connection remove from ConnManager successfully: conn num =", connMgr.Len())

}

// 根据commID获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源，添加则加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not Found!")
	}

}

// 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源，添加则加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除connection并停止conn的工作
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()

		// 删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All connection successful, conn num = ", connMgr.Len())
}
package znet

import "zinx/ziface"

// 实现router时，先嵌入这个BaseRouter基类，然后根据需求对这个基类方法进行重写
type BaseRouter struct{}

// 这里之所以BaseRouter方法都为空
// 目的是因为有的Router不希望有PreHandle、PostHandle这两个业务
// 所有Router全部继承BaseRouter的好处就是，不需要实现PreHandler方法
// 在处理conn业务之前的钩子方法
func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

// 在处理conn业务的主方法hook
func (b *BaseRouter) PHandle(reqest ziface.IRequest) {}

// 在处理conn业务之后的钩子方法hook
func (b *BaseRouter) PPostHandle(request ziface.IRequest) {}

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinx/ziface"
)

/*
	定义存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server相关配置
	*/
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            // 当前服务器端口号
	Name      string         // 当前服务器名称

	/*
		Zinx 相关配置
	*/
	Version        string // 当前Zinx的版本号
	MaxConn        int    // 当前服务器允许的最大连接数
	MaxPackageSize uint32 // 当前Zinx框架数据包的最大值
}

/*
	定义一个全局的对外Globalobj对象
 */

var GlobalObject *GlobalObj


/*
	从zinx.json去加载用户自定义参数
 */
func (g *GlobalObj) Reload(){
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
	fmt.Println("load conf/zinx.json success!")
}
/*
	导入包时候自动执行init方法
	初始化当前的GlobalOjbect
 */
func init() {
	GlobalObject = &GlobalObj{
		Name : "ZinxServerApp",
		Version : "v0.6",
		TcpPort: 8999,
		Host: "0.0.0.0",
		MaxConn: 1000,
		MaxPackageSize: 4096,
	}
	// 应该尝试从conf/zinx.json配置文件中加载用户自定义参数
	GlobalObject.Reload()
}
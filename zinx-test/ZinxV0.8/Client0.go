package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

/*
	模拟客户端
*/

func main() {
	fmt.Println("Client0 start")
	time.Sleep(time.Second)
	// 1. 创建tcp连接，得到conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("dail error", err)
		return
	}
	// 2. 链接调用Write写数据
	for {

		//_, err := conn.Write([]byte("Hello Zinx V0.1..."))
		//if err != nil {
		//	fmt.Println("write failed!", err)
		//	return
		//}
		//
		//buf := make([]byte, 512)
		//cnt, err := conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read buf err")
		//	return
		//}
		//fmt.Printf("server call back : %s , cnt = %d\n", string(buf[:cnt]), cnt)

		// 发送封包的Message消息 MsgID:0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxV0.7 client0 Test Message")))
		if err != nil {
			fmt.Println("Pack error", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		//服务器应该回复一个message数据，MsgID: 1 ping...ping...ping

		//1. 先读取流中的head部分，得到ID和datalen

		//再根据datalen进行第二次读取，将data读出来
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}
		//2.将二进制head拆包到Message消息中
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msghead error", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// msg中有数据，第二次读取开始.
			// 将msgHead转为Message类型
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.DataLen)
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Println("--->Recv Server Msg : ID=", msg.Id, ", len=", msg.DataLen, ",data=", string(msg.Data))
		}
		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}

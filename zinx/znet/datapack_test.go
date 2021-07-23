package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是负责测试datapack拆包、封包的单元测试
// 功能测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	// 1. 创建socket TCP套接字
	listenner, err := net.Listen("tcp", "0.0.0.0:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return

	}
	// 创建一个go承载负责从客户端处理业务
	go func() {
		// 2. 从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accetp err", err)
			}
			go func(conn net.Conn) {
				// 处理客户端请求
				// ----》 拆包过程
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg是有数据的，需要进行第二次读取
						// 第二次读，从conn读，根据head中的datalen再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen值再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err", err)
							return
						}
						//完整的消息已经读取完毕，打印出来
						fmt.Println("-->Recv MsgId:", msg.Id, "datalen=", msg.DataLen, "data=", string(msg.Data))
					}
				}

			}(conn)

		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "0.0.0.0:7777")
	if err != nil {
		fmt.Println("client dail err", err)
		return
	}
	// 发包过程，创建一个封包对象
	dp := NewDataPack()
	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'z', 'i', 'n', 'x', '!'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg 1 err", err)
		return
	}
	// 封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', '0', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg 2 err", err)
		return
	}
	// 将两个包粘到一起
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	//客户端阻塞
	select {}
}

package znet

type Message struct {
	Id      uint32 // 消息ID
	DataLen uint32 //消息长度
	Data    []byte //消息内容
}

// 创建Message消息包方法
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息ID
func (m *Message) SetMsgID(Id uint32) {
	m.Id = Id
}

// 设置消息内容
func (m *Message) SetData(Data []byte) {
	m.Data = Data
}

// 设置消息长度
func (m *Message) SetDataLen(DataLen uint32) {
	m.DataLen = DataLen
}

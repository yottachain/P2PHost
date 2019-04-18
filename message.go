package host

import inet "github.com/libp2p/go-libp2p-net"

// Msg 消息
type Msg struct {
	Content []byte
	MsgType string
	conn    inet.Conn
}

// Conn 返回连接信息
// 包括 RemotePeer, RemotePubkey 等
func (msg Msg) Conn() inet.Conn {
	return msg.conn
}

// MsgData 消息数据
type MsgData []byte

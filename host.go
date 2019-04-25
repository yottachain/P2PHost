package host

import (
	"context"
	"encoding/gob"
	"fmt"
	"io/ioutil"

	"github.com/yottachain/P2PHost/util"

	ci "github.com/libp2p/go-libp2p-crypto"

	"github.com/libp2p/go-libp2p-peer"

	"github.com/libp2p/go-libp2p-peerstore"

	"github.com/libp2p/go-libp2p-host"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

// PrivKey 私钥
type PrivKey ci.PrivKey

// PubKey 公钥
type PubKey ci.PubKey

// Host 接口
// [update] 所有id string 替换成peer.ID
type Host interface {
	ID() peer.ID
	Addrs() []string
	Peerstore() peerstore.Peerstore
	// Start()

	// 连接节点
	// id 节点ID
	// addrs 节点地址
	Connect(id peer.ID, addrs []string) error
	// 断开来接
	DisConnect(id peer.ID) error
	// 发送消息
	// id 发送目标节点id
	// msgType 消息类型
	SendMsg(id peer.ID, msgType string, msg MsgData) ([]byte, error)
	NewStream(id peer.ID, msgType string) (inet.Stream, error)
	// 注册回调函数
	// msgType 消息类型
	// MsgHandlerFunc 消息处理函数
	RegisterHandler(msgType string, MassageHandler MsgHandlerFunc)
	// 注销回调函数
	unregisterHandler(msgType string)
	// 关闭节点所有连接
	Close()
}

type hst struct {
	lhost host.Host
}

// MsgHandler 消息处理器
type MsgHandler interface {
	Process(msgType string, msg []byte)
}

// MsgHandlerFunc 消息处理函数
type MsgHandlerFunc func(msg Msg) []byte

// Start 启动节点 【暂时不需要，未实现】
func (h hst) Start() {
}

// ID
// [update] 为了和libp2p工具集统一 返回值改成了peer.id类型。可以 将string 强制转换成 peer.id peer.ID("BP_1")
func (h hst) ID() peer.ID {
	return h.lhost.ID()
}

func (h hst) Addrs() []string {
	maddrs := h.lhost.Addrs()
	addrs := make([]string, len(maddrs))
	for k, ma := range maddrs {
		addrs[k] = ma.String()
	}
	return addrs
}

func (h hst) Peerstore() peerstore.Peerstore {
	return h.Peerstore()
}

// Connect 连接节点
func (h hst) Connect(id peer.ID, addrs []string) error {
	maddrs, err := util.StringListToMaddrs(addrs)
	if err != nil {
		return err
	}
	info := peerstore.PeerInfo{
		ID:    id,
		Addrs: maddrs,
	}
	h.lhost.Connect(context.Background(), info)
	return nil
}

// DisConnect 断开连接
func (h hst) DisConnect(id peer.ID) error {
	h.lhost.Peerstore().ClearAddrs(id)
	return nil
}

// SendMsg 发送消息
//
// id 节点id，msgType 消息类型， msg 消息数据 字节集
// 远程节点返回内容将通过返回值返回
func (h hst) SendMsg(peerID peer.ID, msgType string, msg MsgData) ([]byte, error) {
	pid := protocol.ID(fmt.Sprintf("%s", msgType))
	stm, err := h.lhost.NewStream(context.Background(), peerID, pid)
	if err != nil {
		fmt.Println(h.lhost.Peerstore().Addrs(peerID))
		return nil, err
	}
	ed := gob.NewEncoder(stm)
	ed.Encode(msg)
	return ioutil.ReadAll(stm)
}
func (h hst) NewStream(peerID peer.ID, msgType string) (inet.Stream, error) {
	return h.lhost.NewStream(context.Background(), peerID, protocol.ID(msgType))
}

// RegisterHandler 注册消息回调函数
func (h hst) RegisterHandler(msgType string, MassageHandler MsgHandlerFunc) {
	pid := protocol.ID(msgType)
	h.lhost.SetStreamHandler(pid, func(stm inet.Stream) {
		defer stm.Close()
		var msgData []byte
		dd := gob.NewDecoder(stm)
		dd.Decode(&msgData)
		stm.Write(MassageHandler(Msg{
			msgData,
			msgType,
			stm.Conn(),
		}))
	})
}

// unregisterHandler 移除消息处理器
func (h hst) unregisterHandler(msgType string) {
	h.lhost.RemoveStreamHandler(protocol.ID(msgType))
}

// Close 关闭
func (h hst) Close() {
	h.lhost.Close()
}

// NewHost 创建节点
//
// addrs 监听地址
// privKey 私钥
func NewHost(addrs []string, privKey PrivKey) (Host, error) {

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(addrs...),
		libp2p.EnableRelay(circuit.OptHop, circuit.OptDiscovery),
	}
	if privKey != nil {
		opts = append(opts, libp2p.Identity(privKey))
	}
	h, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	return &hst{
		h,
	}, nil
}

// ListenAddrStrings 返回监听地址
func ListenAddrStrings(addrs ...string) []string {
	return addrs
}

// WarpHost 包装host
func WarpHost(host host.Host) Host {
	return &hst{
		host,
	}
}

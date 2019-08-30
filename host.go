package host

import (
	"context"
	"crypto/tls"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ripemd160"

	connmgr "github.com/libp2p/go-libp2p-connmgr"
	crypto "github.com/libp2p/go-libp2p-crypto"
	multiaddr "github.com/multiformats/go-multiaddr"
	"golang.org/x/net/http2"

	peer "github.com/libp2p/go-libp2p-peer"
	base58 "github.com/mr-tron/base58"

	"github.com/libp2p/go-libp2p-peerstore"

	"github.com/libp2p/go-libp2p-host"

	rl "github.com/juju/ratelimit"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
	ma "github.com/multiformats/go-multiaddr"
)

// Host 接口
type Host interface {
	ID() string
	Addrs() []string
	// Start()
	Connect(id string, addrs []string) error
	DisConnect(id string) error
	SendMsg(id string, msgType string, msg []byte) ([]byte, error)
	RegisterHandler(msgType string, MassageHandler MsgHandlerFunc)
	UnregisterHandler(msgType string)
	Close()
}

type hst struct {
	lhost       host.Host
	client      http.Client
	ratelimiter *rl.Bucket
}

// // MsgHandler 消息处理器
// type MsgHandler interface {
// 	Process(msgType string, msg []byte, publicKey string)
// }

// MsgHandlerFunc 消息处理函数
type MsgHandlerFunc func(msgType string, msg []byte, publicKey string) ([]byte, error)

var ratelimit int64
var callbackPort int
var connectTimeout int
var readTimeout int
var writeTimeout int
var connUpperLimit int
var connLowerLimit int
var enablePprof bool
var pprofPort int

func init() {
	ratelimitstr := os.Getenv("P2PHOST_RATELIMIT")
	rls, err := strconv.Atoi(ratelimitstr)
	if err != nil {
		ratelimit = 8000
	} else {
		ratelimit = int64(rls)
	}
	callbackPortstr := os.Getenv("P2PHOST_CALLBACKPORT")
	cbp, err := strconv.Atoi(callbackPortstr)
	if err != nil {
		callbackPort = 18999
	} else {
		callbackPort = cbp
	}

	connectTimeoutstr := os.Getenv("P2PHOST_CONNECTTIMEOUT")
	cto, err := strconv.Atoi(connectTimeoutstr)
	if err != nil {
		connectTimeout = 30
	} else {
		connectTimeout = cto
	}
	readTimeoutstr := os.Getenv("P2PHOST_READTIMEOUT")
	rto, err := strconv.Atoi(readTimeoutstr)
	if err != nil {
		readTimeout = 0
	} else {
		readTimeout = rto
	}
	writeTimeoutstr := os.Getenv("P2PHOST_WRITETIMEOUT")
	wto, err := strconv.Atoi(writeTimeoutstr)
	if err != nil {
		writeTimeout = 0
	} else {
		writeTimeout = wto
	}
	connUpperLimitStr := os.Getenv("P2PHOST_CONNUPPERLIMIT")
	cul, err := strconv.Atoi(connUpperLimitStr)
	if err != nil {
		connUpperLimit = 2000
	} else {
		connUpperLimit = cul
	}
	connLowerLimitStr := os.Getenv("P2PHOST_CONNLOWERLIMIT")
	cll, err := strconv.Atoi(connLowerLimitStr)
	if err != nil {
		connLowerLimit = 1000
	} else {
		connLowerLimit = cll
	}
	enablePprofStr := os.Getenv("P2PHOST_ENABLEPPROF")
	ep, err := strconv.ParseBool(enablePprofStr)
	if err != nil {
		enablePprof = false
	} else {
		enablePprof = ep
	}
	pprofPortStr := os.Getenv("P2PHOST_PPROFPORT")
	pp, err := strconv.Atoi(pprofPortStr)
	if err != nil {
		pprofPort = 6161
	} else {
		pprofPort = pp
	}
	if enablePprof {
		go func() {
			http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", pprofPort), nil)
		}()
	}
	// setupSigusr1Trap()
}

// func setupSigusr1Trap() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, syscall.SIGUSR1)
// 	go func() {
// 		for range c {
// 			DumpStacks()
// 		}
// 	}()
// }
// func DumpStacks() {
// 	buf := make([]byte, 1<<20)
// 	buf = buf[:runtime.Stack(buf, true)]
// 	fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===", buf)
// }

// Start 启动节点 【暂时不需要，未实现】
func (h hst) Start() {
}
func (h hst) ID() string {
	id := h.lhost.ID().Pretty()
	return string(id)
}

func (h hst) Addrs() []string {
	maddrs := h.lhost.Addrs()
	addrs := make([]string, len(maddrs))
	for k, ma := range maddrs {
		addrs[k] = ma.String()
	}
	return addrs
}

// Connect 连接节点
func (h hst) Connect(id string, addrs []string) error {
	maddrs, err := StringListToMaddrs(addrs)
	if err != nil {
		return err
	}
	pid, err := peer.IDB58Decode(id)
	if err != nil {
		return err
	}
	info := peerstore.PeerInfo{
		pid,
		maddrs,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*time.Duration(connectTimeout))
	defer cancle()
	h.lhost.Connect(ctx, info)
	return nil
}

// DisConnect 断开连接
func (h hst) DisConnect(id string) error {
	pid, err := peer.IDB58Decode(id)
	if err != nil {
		return err
	}
	h.lhost.Peerstore().ClearAddrs(pid)
	return nil
}

// Msg 消息
type Msg []byte

// SendMsg 发送消息
//
// id 节点id，msgType 消息类型， msg 消息数据 字节集
// 远程节点返回内容将通过返回值返回
func (h hst) SendMsg(id string, msgType string, msg []byte) ([]byte, error) {
	pid := protocol.ID(msgType)
	peerID, err := peer.IDB58Decode(id)
	if err != nil {
		return nil, err
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*time.Duration(connectTimeout))
	defer cancle()
	stm, err := h.lhost.NewStream(ctx, peerID, pid)
	if err != nil {
		//fmt.Println(h.lhost.Peerstore().Addrs(peerID))
		return nil, err
	}
	defer stm.Close()
	if writeTimeout > 0 {
		stm.SetWriteDeadline(time.Now().Add(time.Duration(writeTimeout) * time.Second))
	}
	if readTimeout > 0 {
		stm.SetReadDeadline(time.Now().Add(time.Duration(readTimeout) * time.Second))
	}
	ed := gob.NewEncoder(stm)
	err = ed.Encode(msg)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(stm)
}

// RegisterHandler 注册消息回调函数
func (h hst) RegisterHandler(msgType string, MessageHandler MsgHandlerFunc) {
	if MessageHandler == nil {
		MessageHandler = func(msgType string, msg []byte, publicKey string) ([]byte, error) {
			//log.Printf("##### %s Receive Message: Type: 0x%s, Public Key: %s", time.Now().Format("2006-01-02 15:04:05"), hex.EncodeToString(msg[0:2]), publicKey)
			resp, err := h.client.PostForm(fmt.Sprintf("http://127.0.0.1:%d", callbackPort),
				url.Values{"type": {msgType}, "data": {hex.EncodeToString(msg)}, "pubkey": {publicKey}})
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return hex.DecodeString(string(body))
		}
	}
	pid := protocol.ID(msgType)
	h.lhost.SetStreamHandler(pid, func(stm inet.Stream) {
		//s := time.Now()
		defer stm.Close()
		h.ratelimiter.Wait(1)
		// if c := h.ratelimiter.TakeAvailable(1); c == 0 {
		// 	stm.Reset()
		// 	log.Println("Qos backoff request")
		// 	return
		// }
		var msg Msg
		dd := gob.NewDecoder(stm)
		dd.Decode(&msg)
		pkarr, err := stm.Conn().RemotePublicKey().Raw()
		if err != nil {
			stm.Reset()
			log.Println(err)
			return
		}
		hasher := ripemd160.New()
		hasher.Write(pkarr)
		sum := hasher.Sum(nil)
		pkarr = append(pkarr, sum[0:4]...)
		resp, err := MessageHandler(msgType, msg, base58.Encode(pkarr))
		if err != nil {
			stm.Reset()
			log.Println(err)
			return
		}
		_, err = stm.Write(resp)
		if err != nil {
			stm.Reset()
			log.Println(err)
			return
		}
		//log.Println(fmt.Sprintf("###### process message from %s, cost %f", base58.Encode(pkarr), time.Now().Sub(s).Seconds()))
	})
}

// unregisterHandler 移除消息处理器
func (h hst) UnregisterHandler(msgType string) {
	h.lhost.RemoveStreamHandler(protocol.ID(msgType))
}

// Close 关闭
func (h hst) Close() {
	h.lhost.Close()
}

// NewHost 创建节点
func NewHost(privateKey string, listenAddrs ...string) (Host, error) {
	maddrs, err := StringListToMaddrs(listenAddrs)

	addrs := libp2p.ListenAddrs(maddrs...)
	connMgr := connmgr.NewConnManager(connUpperLimit, connLowerLimit, 0)
	opts := []libp2p.Option{
		addrs,
		libp2p.NATPortMap(),
		libp2p.EnableRelay(circuit.OptHop, circuit.OptDiscovery),
		libp2p.ConnectionManager(connMgr),
		libp2p.AddrsFactory(func(addrs []ma.Multiaddr) []ma.Multiaddr { return addrs }),
	}
	if privateKey != "" {
		privbytes, err := base58.Decode(privateKey)
		if err != nil {
			return nil, errors.New("bad format of private key,Base58 format needed")
		}
		priv, err := crypto.UnmarshalSecp256k1PrivateKey(privbytes[1:33])
		if err != nil {
			return nil, errors.New("bad format of private key")
		}
		opts = append(opts, libp2p.Identity(priv))
	} else {
		peer, err := RandomPeer()
		if err != nil {
			return nil, err
		}
		opts = append(opts, libp2p.Identity(peer.Priv))
	}
	h, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	ratelimiter := rl.NewBucketWithRate(float64(ratelimit), ratelimit)
	log.Printf("##### initializing rate limit to %d per second.", ratelimit)
	return &hst{
		h,
		http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
				StrictMaxConcurrentStreams: false,
			},
		},
		ratelimiter,
	}, nil
}

func StringListToMaddrs(addrs []string) ([]multiaddr.Multiaddr, error) {
	maddrs := make([]multiaddr.Multiaddr, len(addrs))
	for k, addr := range addrs {
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return maddrs, err
		}
		maddrs[k] = maddr
	}
	return maddrs, nil
}

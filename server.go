package server

import "C"
import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	host "github.com/yottachain/YTHost"
	hst "github.com/yottachain/YTHost/hostInterface"
	"os"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mr-tron/base58"
	ma "github.com/multiformats/go-multiaddr"
	pb "github.com/yottachain/P2PHost/pb"
	"github.com/yottachain/YTHost/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

// Server implemented server API for P2PHostServer service.
type Server struct {
	Host hst.Host
	Hc   Hclient
}

const GETTOKEN = 50311
var wt int
var ct int

func init() {
	writetimeout := os.Getenv("P2PHOST_WRITETIMEOUT")
	wt = 5000
	if writetimeout == "" {
		wt = 5000
	}else {
		wto, err := strconv.Atoi(writetimeout)
		if err != nil {
			wt = 5000
		}else {
			wt = wto
		}
	}

	conntimeout := os.Getenv("P2PHOST_CONNECTTIMEOUT")
	ct = 15000
	if conntimeout == "" {
		ct = 15000
	}else {
		cto, err := strconv.Atoi(conntimeout)
		if err != nil {
			ct = 15000
		}else {
			ct = cto
		}
	}
}

// ID implemented ID function of P2PHostServer
func (server *Server) ID(ctx context.Context, req *pb.Empty) (*pb.StringMsg, error) {
	return &pb.StringMsg{Value: server.Host.Config().ID.String()}, nil
}

// Addrs implemented Addrs function of P2PHostServer
func (server *Server) Addrs(ctx context.Context, req *pb.Empty) (*pb.StringListMsg, error) {
	maddrs := server.Host.Addrs()
	addrs := make([]string, len(maddrs))
	for k, madd := range maddrs {
		addr := madd.String()
		addrs[k] = addr
	}
	return &pb.StringListMsg{Values: addrs}, nil
}

// Connect implemented Connect function of P2PHostServer
func (server *Server) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.Empty, error) {
	maddrs, _ := stringListToMaddrs(req.GetAddrs())
	ID, err := peer.Decode(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(ct))
	defer cancel()

	_, err = server.Host.ClientStore().Get(ctx, ID, maddrs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Empty{}, nil
}

// DisConnect implemented DisConnect function of P2PHostServer
func (server *Server) DisConnect(ctx context.Context, req *pb.StringMsg) (*pb.Empty, error) {
	ID, err := peer.Decode(req.GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	
	err = server.Host.ClientStore().Close(ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Empty{}, nil
}

// SendMsg implemented SendMsg function of P2PHostServer
func (server *Server) SendMsg(ctx context.Context, req *pb.SendMsgReq) (*pb.SendMsgResp, error) {
	msid := req.GetMsgid()[:2:2]
	bytebuff := bytes.NewBuffer(msid)
	var tmp uint16
	err := binary.Read(bytebuff, binary.BigEndian, &tmp)

	msgId := int32(tmp)
	
	ID, err := peer.Decode(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(wt))
	if msgId == GETTOKEN {
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*1000)
	}
	defer cancel()

	//tstart := time.Now()
	bytes, err := server.Host.SendMsg(ctx, ID, msgId, req.GetMsg())
	//interval := time.Now().Sub(tstart).Milliseconds()
	//if msgId == GETTOKEN {
	//	lg.Info.Printf("grpc server send [get token] time:%d\n", interval)
	//}else {
	//	lg.Info.Printf("grpc server send [not get token] time:%d\n", interval)
	//}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.SendMsgResp{Value: bytes}, nil
}

// RegisterHandler implemented RegisterHandler function of P2PHostServer
func (server *Server) RegisterHandler(ctx context.Context, req *pb.StringMsg) (*pb.Empty, error) {
	server.Host.RegisterHandler(0x0, server.Hc.MessageHandler)
	return &pb.Empty{}, nil
}

// UnregisterHandler implemented UnregisterHandler function of P2PHostServer
func (server *Server) UnregisterHandler(ctx context.Context, req *pb.StringMsg) (*pb.Empty, error) {
	server.Host.RemoveGlobalHandler()
	return &pb.Empty{}, nil
}

// Close implemented Close function of P2PHostServer
func (server *Server) Close(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	//server.Host.Close()
	return &pb.Empty{}, nil
}

func stringListToMaddrs(addrs []string) ([]ma.Multiaddr, error) {
	maddrs := make([]ma.Multiaddr, len(addrs))
	for k, addr := range addrs {
		maddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return maddrs, err
		}
		maddrs[k] = maddr
	}
	return maddrs, nil
}

func NewServer(port string, priKey string) (*Server, error){
	srv := Server{}
	pt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	ma, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", pt))
	privbytes, err := base58.Decode(priKey)
	if err != nil {
		return nil, err
	}
	pk, err := crypto.UnmarshalSecp256k1PrivateKey(privbytes[1:33])
	if err != nil {
		return nil, err
	}

	srv.Host, err = host.NewHost(option.ListenAddr(ma), option.Identity(pk))
	if err != nil {
		return nil, err
	}
	go srv.Host.Accept()

	srv.Hc, err = NewHclient()

	return &srv, nil
}

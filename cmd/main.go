package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	host "github.com/yottachain/P2PHost"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"

	pb "github.com/yottachain/P2PHost/pb"
)

const P2PHOST_ETCD_PREFIX = "/p2phost/"
const P2PHOST_PORT = P2PHOST_ETCD_PREFIX + "port"
const P2PHOST_PRIVKEY = P2PHOST_ETCD_PREFIX + "privkey"

const USER_MSG = "/user/0.0.2"
const BPNODE_MSG = "/bpnode/0.0.2"
const NODE_MSG = "/node/0.0.2"

func main() {
	etcdHostname := os.Getenv("ETCDHOSTNAME")
	if etcdHostname == "" {
		etcdHostname = "etcd-svc"
	}
	etcdPortStr := os.Getenv("ETCDPORT")
	etcdPort, err := strconv.Atoi(etcdPortStr)
	if err != nil {
		etcdPort = 2379
	}
	log.Printf("ETCD URL: %s:%d\n", etcdHostname, etcdPort)

	callbackHostname := os.Getenv("P2PHOST_CALLBACKHOSTNAME")
	if callbackHostname == "" {
		callbackHostname = "ytsn-server"
	}
	log.Printf("Callback hostname: %s\n", callbackHostname)
	callbackPortStr := os.Getenv("P2PHOST_CALLBACKPORT")
	callbackPort, err := strconv.Atoi(callbackPortStr)
	if err != nil {
		callbackPort = 18999
	}
	log.Printf("Callback port: %d\n", callbackPort)
	clnt, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%d", etcdHostname, etcdPort)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln("connect etcd failed, err: ", err)
	}
	log.Println("connect etcd success")
	defer clnt.Close()

	for {
		time.Sleep(time.Second * 1)
		resp, err := clnt.Get(context.Background(), P2PHOST_PRIVKEY)
		if err != nil {
			log.Printf("get %s failed, err: %s\n", P2PHOST_PRIVKEY, err)
			continue
		}
		if len(resp.Kvs) == 0 {
			log.Printf("get %s failed, no content\n", P2PHOST_PRIVKEY)
			continue
		}
		p2pPrivkey := resp.Kvs[0].Value
		log.Printf("Read P2P private key from ETCD: %s\n", p2pPrivkey)

		resp, err = clnt.Get(context.Background(), P2PHOST_PORT)
		if err != nil {
			log.Printf("get %s failed, err: %s\n", P2PHOST_PORT, err)
			continue
		}
		if len(resp.Kvs) == 0 {
			log.Printf("get %s failed, no content\n", P2PHOST_PORT)
			continue
		}
		p2pPortStr := resp.Kvs[0].Value
		p2pPort, err := strconv.Atoi(string(p2pPortStr))
		if err != nil {
			log.Printf("parse %s failed, err: %s\n", P2PHOST_PORT, err)
			continue
		}
		log.Printf("Read P2P port from ETCD: %d\n", p2pPort)

		h, err := host.NewHost(string(p2pPrivkey), fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", p2pPort))
		if err != nil {
			log.Fatalf("create p2phost instance failed, err: %s\n", err)
		}
		log.Printf("create p2phost success, listening address is %s\n", fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", p2pPort))
		h.ConfigCallback(callbackHostname, int32(callbackPort))
		h.RegisterHandler(USER_MSG, nil)
		h.RegisterHandler(BPNODE_MSG, nil)
		h.RegisterHandler(NODE_MSG, nil)
		log.Printf("configure callback handler successful.")

		server := &host.Server{Host: h}

		p2pGRPCPortStr := os.Getenv("P2PHOST_GRPCPORT")
		p2pGRPCPort, err := strconv.Atoi(p2pGRPCPortStr)
		if err != nil {
			p2pGRPCPort = 11002
		}
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", p2pGRPCPort))
		if err != nil {
			log.Fatalf("failed to listen GRPC port %d: %v", p2pGRPCPort, err)
		}
		log.Printf("GRPC address: 0.0.0.0:%d\n", p2pGRPCPort)
		grpcServer := grpc.NewServer()
		pb.RegisterP2PHostServer(grpcServer, server)
		grpcServer.Serve(lis)
		log.Printf("GRPC server started.")
		break
	}
}

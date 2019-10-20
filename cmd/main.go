package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	host "github.com/yottachain/P2PHost"
	ytcrypto "github.com/yottachain/YTCrypto"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"

	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

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

func main9() {
	h, err := host.NewHost("5HtM6e3mQNLEu2TkQ1ZrbMNpRQiHGsKxEsLdxd9VsdCmp1um8QH", "/ip4/0.0.0.0/udp/9002/quic")
	if err != nil {
		panic(fmt.Sprintf("new host error: %s", err))
	} else {
		fmt.Printf("new host error: %s", h.ID())
	}
	err = h.Connect("16Uiu2HAmRpgTQLrmiDU2i5sWEWewXFa611pbkZ9ReHtQhDYw9MUt", []string{"/ip4/152.136.17.121/tcp/9001"})
	if err != nil {
		panic(fmt.Sprintf("new host error: %s", err))
	}
	addrs := h.Addrs()
	for addr := range addrs {
		fmt.Println(addr)
	}
}

func main1() {
	h, err := host.NewHost("5HtM6e3mQNLEu2TkQ1ZrbMNpRQiHGsKxEsLdxd9VsdCmp1um8QH", "/ip4/0.0.0.0/tcp/7777")
	if err != nil {
		panic(fmt.Sprintf("new host error: %s", err))
	} else {
		fmt.Printf("HOST ID: %s\n", h.ID())
	}
	addrs := h.Addrs()
	for _, addr := range addrs {
		fmt.Println(addr)
	}
	h.RegisterHandler("/user/0.0.1", func(proto string, data []byte, pubkey string) ([]byte, error) {
		str := string(data[:])
		fmt.Println("Receive: " + str)
		return []byte("Echo: " + str), nil
	})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("Got signal:", s)
}

func main6() {
	ch := make(chan int)
	var start time.Time
	for i := range [1]int{} {
		ii := i
		pk, _ := ytcrypto.CreateKey()
		go func() {
			// h, err := host.NewHost("5KQKydL7TuRwjzaFSK4ezH9RUXWuYHW1yYDp5CmQfsfTuu9MBLZ", "/ip4/0.0.0.0/tcp/9999")
			h, err := host.NewHost(pk, fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 40001+ii))
			if err != nil {
				panic(fmt.Sprintf("new host error: %s\n", err))
			} else {
				fmt.Printf("new host: %s\n", h.ID())
			}
			// addrs := h.Addrs()
			// for _, addr := range addrs {
			// 	fmt.Println(addr)
			// }
			err = h.Connect("16Uiu2HAm44FX3YuzGXJgHMqnyMM5zCzeT6PUoBNZkz66LutfRREM", []string{fmt.Sprintf("/ip4/127.0.0.1/tcp/8888")})
			if err != nil {
				panic(fmt.Sprintf("========================== connect error: %s\n", err))
			}
			start = time.Now()
			for i := range [10000]int{} {
				go (func(index int) {
					_, err := h.SendMsg("16Uiu2HAm44FX3YuzGXJgHMqnyMM5zCzeT6PUoBNZkz66LutfRREM", "/user/0.0.1", []byte(fmt.Sprintf("%s%d\n", "testaaa", index)))
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					//fmt.Printf("%s\n", string(ret))
					ch <- index
				})(i)
			}

			// err = h.Connect("16Uiu2HAm44FX3YuzGXJgHMqnyMM5zCzeT6PUoBNZkz66LutfRREM", []string{fmt.Sprintf("/ip4/192.168.3.107/tcp/8888")})
			// if err != nil {
			// 	panic(fmt.Sprintf("========================== connect error: %s\n", err))
			// }
			// go func() {
			// 	for i := range [10000]int{} {
			// 		ret, err := h.SendMsg("16Uiu2HAm44FX3YuzGXJgHMqnyMM5zCzeT6PUoBNZkz66LutfRREM", "/user/0.0.1", []byte(fmt.Sprintf("%s%d", "testa", i)))
			// 		if err != nil {
			// 			fmt.Println(err.Error())
			// 			return
			// 		}
			// 		// fmt.Printf("%s\n", string(ret))
			// 		ch <- index
			// 	}
			// }()
			// h.DisConnect("16Uiu2HAm44FX3YuzGXJgHMqnyMM5zCzeT6PUoBNZkz66LutfRREM")
			// h.Close()
		}()
	}
	total := 0
	for {
		<-ch
		total++
		if total == 10000 {
			break
		}
	}
	fmt.Printf("finish: %f", time.Now().Sub(start).Seconds())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("Got signal:", s)
}

func main14() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello h2c")
	})

	h2s := &http2.Server{}

	h1s := &http.Server{
		Addr:    ":8972",
		Handler: h2c.NewHandler(mux, h2s),
	}
	log.Fatal(h1s.ListenAndServe())
}

func main11() {
	client := http.Client{
		// Skip TLS dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	resp, err := client.Get("http://localhost:8972")
	if err != nil {
		log.Fatal(fmt.Errorf("error making request: %v", err))
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Proto)
}

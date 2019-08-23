package main

import (
	"fmt"
	"os"
	"os/signal"
	"yottachain/p2phost"

	ytcrypto "github.com/yottachain/YTCrypto"
)

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

func main() {
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

func main1() {
	for i := range [1000]int{} {
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
			err = h.Connect("16Uiu2HAm91ZquksU6MdT74uTEVLb4vwixjLwEBWPXJsfeUMheBY8", []string{fmt.Sprintf("/ip4/10.2.0.153/tcp/7777")})
			if err != nil {
				panic(fmt.Sprintf("========================== connect error: %s\n", err))
			}

			// for i := range [10000]int{} {
			// 	go (func(index int) {
			// 		ret, err := h.SendMsg("16Uiu2HAm44FX3YuzGXJgHMqnyMM5zCzeT6PUoBNZkz66LutfRREM", "/test/v0.01", []byte(fmt.Sprintf("%s%d\n", "testaaa", index)))
			// 		if err != nil {
			// 			fmt.Println(err.Error())
			// 			return
			// 		}
			// 		fmt.Printf("%s\n", string(ret))
			// 	})(i)
			// }

			//go func() {
			for i := range [100000]int{} {
				ret, err := h.SendMsg("16Uiu2HAm91ZquksU6MdT74uTEVLb4vwixjLwEBWPXJsfeUMheBY8", "/user/0.0.1", []byte(fmt.Sprintf("%s%d\n", "testa", i)))
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Printf("%s\n", string(ret))
			}
			//}()
			h.DisConnect("16Uiu2HAm44FX3YuzGXJgHMqnyMM5zCzeT6PUoBNZkz66LutfRREM")
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("Got signal:", s)
}

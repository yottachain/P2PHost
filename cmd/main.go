package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	rl "github.com/juju/ratelimit"
	"github.com/yottachain/P2PHost"
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

func main() {
	// arr, _ := base58.Decode("69YnNn4L6FxonZxQGYVyKvmqVPBmpPjuxxs5NHfZXoeSnFdipR")
	// hasher := ripemd160.New()
	// hasher.Write(arr[0 : len(arr)-4])
	// sum := hasher.Sum(nil)
	// fmt.Println(base58.Encode(append(arr[0:len(arr)-4], sum[0:4]...)))
	bucket := rl.NewBucketWithRate(10, 10)
	//time.Sleep(5000)
	for i := 0; i < 100; i++ {
		bucket.Wait(1)
		fmt.Println(i)
	}
}

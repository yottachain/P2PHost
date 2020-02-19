package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/yottachain/P2PHost/pb"
	"github.com/yottachain/YTHost/service"
	"log"
	"sync"
	"testing"
)

func TestConnSend(t *testing.T){
	srv, err := NewServer("9876", "5JhXaYtCgA7eW9HAq5LAPqJ2FQ7xt68qnc9VRCumpv24D6pX1sL")
	if err != nil {
		log.Println(err.Error())
	}

	srv.Host.RegisterGlobalMsgHandler(func(requestData []byte, head service.Head) (bytes []byte, e error) {
		fmt.Println(fmt.Sprintf("msg is %s", string(requestData)))
		return []byte("111111111111"), nil
	})

	//srv.Host.RegisterGlobalMsgHandler(srv.Hc.MessageHandler)

	srv1, err := NewServer("0", "5JhXaYtCgA7eW9HAq5LAPqJ3337xt68qnc9VRCumpv24D6pX1sL")
	if err != nil {
		log.Println(err.Error())
	}
	sID := srv.Host.Config().ID.String()
	fmt.Println(sID)

	maddrs := srv.Host.Addrs()
	addrs := make([]string, len(maddrs))
	for k, m := range maddrs{
		addrs[k] = m.String()
	}

	connReq := pb.ConnectReq {
		Id: *proto.String(srv.Host.Config().ID.String()),
		Addrs: addrs,
	}
	_, err = srv1.Connect(context.Background(), &connReq)
	if err != nil {
		t.Fatal(err)
	}
	/*clt, _ := srv1.Host.Connect(context.Background(), srv.Host.Config().ID, srv.Host.Addrs())
	if res, err := clt.SendMsg(context.Background(), 0x0, []byte("22222222223333333333333333")); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(string(res))
	}*/

	var buffer = make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, uint16(0))

	sendMsg := pb.SendMsgReq{
		Id: *proto.String(srv.Host.Config().ID.String()),
		//Id: srv.Host.Config().ID.String(),
		Msgid: buffer,
		Msg: []byte("dasfasdkfdas"),
	}

	wg := sync.WaitGroup{}
	wg.Add(100000)
	for i := 0; i < 100000; i++ {
		go func() {
			defer wg.Done()
			resp, err := srv1.SendMsg(context.Background(), &sendMsg)
			if err == nil {
				fmt.Println(string(resp.Value))
			}
		}()
	}

	wg.Wait()

}

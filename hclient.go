package server

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	rl "github.com/juju/ratelimit"
	"github.com/mr-tron/base58"
	lg "github.com/yottachain/P2PHost/log"
	"github.com/yottachain/YTHost/service"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Hclient interface {
	ConfigCallback(host string, port int)
	Callbackinit() (error)
	GetHost() (string)
	GetProt() (int)
	GetClient() (http.Client)
	MessageHandler(requestData []byte, head service.Head) ([]byte, error)
}

type hc struct {
	client       http.Client
	callbackHost string
	callbackPort int
	ratelimiter  *rl.Bucket
}

var gcallbackHost string
var gcallbackPort int
var ratelimit int64

func init() {
	ratelimitstr := os.Getenv("P2PHOST_RATELIMIT")
	rls, err := strconv.Atoi(ratelimitstr)
	if err != nil {
		ratelimit = 8000
	} else {
		ratelimit = int64(rls)
	}

	gcallbackHost = os.Getenv("P2PHOST_CALLBACKHOST")
	if gcallbackHost == "" {
		gcallbackHost = "127.0.0.1"
	}
	callbackPortstr := os.Getenv("P2PHOST_CALLBACKPORT")
	cbp, err := strconv.Atoi(callbackPortstr)
	if err != nil {
		gcallbackPort = 18999
	} else {
		gcallbackPort = cbp
	}
}

func (h *hc) ConfigCallback(host string, port int) {
	h.callbackHost = host
	h.callbackPort = port
}

func (h *hc)Callbackinit()(error){
	h.ConfigCallback(gcallbackHost, gcallbackPort)
	return nil
}

func (h *hc) GetHost()(string){
	return h.callbackHost
}

func (h *hc) GetProt()(int){
	return h.callbackPort
}

func (h *hc) GetClient()(http.Client){
	return h.client
}


func (h *hc) MessageHandler(requestData []byte, head service.Head) ([]byte, error) {
	h.ratelimiter.Wait(1)
	lg.Info.Printf("cru available %d\n", h.ratelimiter.Available())

	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int16(head.MsgId))
	binary.Write(buf, binary.BigEndian, requestData)

	callbackURL := fmt.Sprintf("http://%s:%d", h.callbackHost, h.callbackPort)
	//pkarr, err := head.RemotePubKey.Raw()
	pkarr := head.RemotePubKey

	//--------
	hasher := ripemd160.New()
	hasher.Write(pkarr)
	sum := hasher.Sum(nil)
	pkarr = append(pkarr, sum[0:4]...)

	publicKey := base58.Encode(pkarr)

	resp, err := h.client.PostForm(callbackURL, url.Values{"data": {hex.EncodeToString(buf.Bytes())}, "pubkey": {publicKey}})
	if err != nil {
		log.Printf("Receive message error: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(body))
}

func NewHclient()(*hc, error){
	hcli := new(hc)
	hcli.client = http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
			StrictMaxConcurrentStreams: false,
		},
		Timeout: 60*time.Second,
	}

	rratelimiter := rl.NewBucketWithRate(float64(ratelimit), ratelimit)
	hcli.ratelimiter = rratelimiter

	hcli.Callbackinit()

	return hcli, nil
}


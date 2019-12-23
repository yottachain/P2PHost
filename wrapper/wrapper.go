package main

/*
#cgo CFLAGS: -std=c99

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

typedef struct idret {
	char *id;
	char *error;
} idret;

typedef struct addrsret {
	char **addrs;
	int size;
	char *error;
} addrsret;

typedef struct sendmsgret {
	char *msg;
	long long size;
	char *error;
} sendmsgret;

typedef sendmsgret* (*msghandler)(char*, char*, long long, char*);

extern char* StartWrp(int port, char *privkey);
extern idret* IDWrp();
extern addrsret* AddrsWrp();
extern char* CloseWrp();
extern char* ConnectWrp(char *nodeID, char **addrs, int size);
extern sendmsgret* SendMsgWrp(char *nodeID, char *msgid, char *msg, long long size);
extern char* RegisterHandlerWrp(char *msgType, void *handler);
extern char* UnregisterHandlerWrp(char *msgType);
extern void FreeString(void *ptr);
extern void FreeIDRet(idret *ptr);
extern void FreeAddrsRet(addrsret *ptr);
extern void FreeSendMsgRet(sendmsgret *ptr);
extern sendmsgret* CreateSendMsgRet(char *msg, long long size, char *err);

static char** makeCharArray(int size) {
	char **ret = (char**)malloc(sizeof(char*) * size);
	memset(ret, 0 , sizeof(char*) * size);
	return ret;
}

static void setArrayString(char **a, char *s, int n) {
    a[n] = s;
}

static void freeCharArray(char **a, int size) {
    int i;
    for (i = 0; i < size; i++) {
		free(a[i]);
		a[i] = NULL;
	}
    free(a);
}

static sendmsgret* executeHandler(msghandler handler, char *msgType, char *data, long long size, char *pubkey) {
	return (*handler)(msgType, data, size, pubkey);
}

//以下为测试用代码
static sendmsgret* msgprocessor(char* msgType, char* data, long long size, char* pubkey) {
	printf("Type: %s\n", msgType);
	printf("Pubkey: %s\n", pubkey);
	for (int i=0; i<size; i++) {
		    printf("%c", data[i]);
	}
	puts("\n");
	char *retdata = (char*)malloc(sizeof(char) * (size + 5));
	strncpy(retdata, data, size);
	strncpy(retdata+size, ": ack", 5);
	sendmsgret *ret = CreateSendMsgRet(retdata, size + 5, NULL);
	return ret;
}

static void sstart() {
 	int port = 7999;
 	char *privkey = "5JSgrkY3jawhV1yTj3HiGJ643TeDFJdbEV3JA2akMJA3LLZFxn1";
	char *err = StartWrp(port, privkey);
	if (err != NULL) {
		printf("error: %s\n", err);
		free(err);
		err = NULL;
		return;
	}

	idret *retp = IDWrp();
	if (retp->error != NULL) {
		printf("Error: %s\n", retp->error);
		free(retp->error);
		free(retp);
		retp = NULL;
		return;
	}
	printf("ID: %s\n", retp->id);
	FreeIDRet(retp);
	retp = NULL;

	addrsret *retp2 = AddrsWrp();
	if (retp2->error != NULL) {
		printf("Error: %s\n", retp2->error);
		FreeAddrsRet(retp2);
		return;
	}
	for (int i=0; i<retp2->size; i++) {
		printf("Addr%d: %s\n", i, (retp2->addrs)[i]);
	}
	FreeAddrsRet(retp2);

	char* error = RegisterHandlerWrp("test", &msgprocessor);
	if (err != NULL) {
		printf("error: %s\n", err);
		free(err);
		err = NULL;
		return;
	}

	//CloseWrp();
	while(1)
		sleep(10000);
	printf("server end!!!!!");
}

static void cstart() {
	int port = 9998;
 	char *privkey = "5JhXaYtCgA7eW9HAq5LAPqJ2FQ7xt68qnc9VRCumpv24D6pX1sL";
	char *err = StartWrp(port, privkey);
	if (err != NULL) {
		printf("error: %s\n", err);
		free(err);
		err = NULL;
		return;
	}

	idret *retp = IDWrp();
	if (retp->error != NULL) {
		printf("Error: %s\n", retp->error);
		free(retp->error);
		free(retp);
		retp = NULL;
		return;
	}
	printf("ID: %s\n", retp->id);
	FreeIDRet(retp);
	retp = NULL;

	addrsret *retp2 = AddrsWrp();
	if (retp2->error != NULL) {
		printf("Error: %s\n", retp2->error);
		FreeAddrsRet(retp2);
		return;
	}
	for (int i=0; i<retp2->size; i++) {
		printf("Addr%d: %s\n", i, (retp2->addrs)[i]);
	}
	FreeAddrsRet(retp2);

	char *addrs[1] = {"/ip4/127.0.0.1/tcp/7999"};
	err = ConnectWrp("16Uiu2HAmAvd2jETZcJL3pwqaRBU9UP6bhZLXRPqFNyiSVZRBftxJ", addrs, 1);
	if (err != NULL) {
		printf("error: %s\n", err);
		free(err);
		err = NULL;
		return;
	}

	char data[12] = {'s','e','n','d',' ','m','e','s','s','a','g','e'};
	char msid[2] ;
	msid[0] = 1;
	msid[1] = 1;
	sendmsgret* retp3 = SendMsgWrp("16Uiu2HAmAvd2jETZcJL3pwqaRBU9UP6bhZLXRPqFNyiSVZRBftxJ", msid, data, 12);
	if (retp3->error != NULL) {
		printf("error: %s\n", retp3->error);
		FreeSendMsgRet(retp3);
		return;
	}
	puts("Received: ");
	for (int i=0; i<retp3->size; i++) {
		printf("%c", (retp3->msg)[i]);
	}
	puts("\n");
	CloseWrp();
}
*/
import "C"
import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/graydream/YTHost/option"
	"github.com/libp2p/go-libp2p-core/peer"
	base58 "github.com/mr-tron/base58"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"time"
	"unsafe"

	hst "github.com/graydream/YTHost"
	host "github.com/graydream/YTHost/hostInterface"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	ma "github.com/multiformats/go-multiaddr"
	p2ph "github.com/yottachain/P2PHost"
)

var p2phst host.Host
//var p2pcli *client.YTHostClient
var mu sync.Mutex
var p2phcli p2ph.Hclient

//export StartWrp
func StartWrp(port C.int, privkey *C.char) *C.char {
	mu.Lock()
	defer mu.Unlock()
	if p2phst != nil {
		return C.CString("p2phost has started")
	}

	pt := int(port)
	ma, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", pt))

	pks := C.GoString(privkey)
	privbytes, err := base58.Decode(pks)
	if err != nil {
		return C.CString("bad format of private key,Base58 format needed")
	}
	pk, err := crypto.UnmarshalSecp256k1PrivateKey(privbytes[1:33])
	if err != nil {
		return C.CString("bad format of private key")
	}

	p2phst, err = hst.NewHost(option.ListenAddr(ma), option.Identity(pk))
	if err != nil {
		return C.CString(err.Error())
	}

	p2phcli, err = p2ph.NewHclient()
	if err != nil {
		return C.CString(err.Error())
	}

	go p2phst.Accept()

	return nil
}

//export IDWrp
func IDWrp() *C.idret {
	if p2phst == nil {
		return CreateIDRet(nil, C.CString("p2phost has not started"))
	}
	id := p2phst.Config().ID.String()

	var cID *C.char
	if id != "" {
		cID = C.CString(id)
	}
	retp := CreateIDRet(cID, nil)
	return retp
}

//export AddrsWrp
func AddrsWrp() *C.addrsret {
	if p2phst == nil {
		return CreateAddrsRet(nil, 0, C.CString("p2phost has not started"))
	}

	maddrs := p2phst.Addrs()
	addrs := make([]string, len(maddrs))
	for k, m := range maddrs {
		addrs[k] = m.String()
	}

	var caddrs **C.char
	if addrs != nil && len(addrs) > 0 {
		caddrs = C.makeCharArray(C.int(len(addrs)))
		for i, s := range addrs {
			C.setArrayString(caddrs, C.CString(s), C.int(i))
		}
	}
	retp := CreateAddrsRet(caddrs, C.int(len(addrs)), nil)
	return retp
}

//export CloseWrp
func CloseWrp() *C.char {
	if p2phst == nil {
		return C.CString("p2phost has not started")
	}
	//p2phst.Close()
	return nil
}

//export ConnectWrp
func ConnectWrp(nodeID *C.char, addrs **C.char, size C.int) *C.char {
	if p2phst == nil {
		return C.CString("p2phost has not started")
	}
	length := int(size)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(addrs))[:length:length]
	gaddrs := make([]string, length)
	for i, s := range tmpslice {
		gaddrs[i] = C.GoString(s)
	}
	nodeIdStr := C.GoString(nodeID)

	maddrs, err := stringListToMaddrs(gaddrs)

	conntimeout := os.Getenv(" P2PHOST_CONNECTTIMEOUT")
	ct := 60
	if conntimeout == "" {
		ct = 60
	}else {
		ct, err = strconv.Atoi(conntimeout)
		if err != nil {
			ct = 60
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ct))
	defer cancel()
	//p2pcli, err = p2phst.Connect(ctx, peer.ID(nodeIdStr), maddrs)
	_, err = p2phst.ClientStore().Get(ctx, peer.ID(nodeIdStr), maddrs)

	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export DisconnectWrp
func DisconnectWrp(nodeID *C.char) *C.char {
	if p2phst == nil {
		return C.CString("p2phost has not started")
	}

	//err := p2phst.ClientStore().Close(peer.ID(C.GoString(nodeID)))
	//if err != nil {
	//	return C.CString(err.Error())
	//}

	return nil
}

//export SendMsgWrp
func SendMsgWrp(nodeID *C.char, msgid *C.char, msg *C.char, size C.longlong) *C.sendmsgret {
	if p2phst == nil {
		return CreateSendMsgRet(nil, 0, C.CString("p2phost has not started"))
	}

	if unsafe.Pointer(nodeID) == nil {
		return CreateSendMsgRet(nil, 0, C.CString("nodeid is nil when send msg"))
	}
	nodeIDStr := C.GoString(nodeID)


	msid := (*[2]byte)(unsafe.Pointer(msgid))[:2:2]
	bytebuff := bytes.NewBuffer(msid)
	var tmp uint16
	err := binary.Read(bytebuff, binary.BigEndian, &tmp)

	msgId := int32(tmp)

	msgSlice := (*[1 << 30]byte)(unsafe.Pointer(msg))[:int64(size):int64(size)]
	conntimeout := os.Getenv(" P2PHOST_WRITETIMEOUT")
	ct := 60
	if conntimeout == "" {
		ct = 60
	}else {
		ct, err = strconv.Atoi(conntimeout)
		if err != nil {
			ct = 60
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ct))
	defer cancel()
	ret, err := p2phst.SendMsg(ctx, peer.ID(nodeIDStr), msgId, msgSlice)
	//ret, err := p2phst.SendMsg(context.Background(), peer.ID(nodeIDStr), 0x11, msgSlice)
	if err != nil {
		return CreateSendMsgRet(nil, C.longlong(0), C.CString(err.Error()))
	}
	return CreateSendMsgRet(C.CString(string(ret)), C.longlong(len(ret)), nil)
}

//export RegisterHandlerWrp
func RegisterHandlerWrp(msgType *C.char, f unsafe.Pointer) *C.char {
	if p2phst == nil {
		return C.CString("p2phost has not started")
	}

	if p2phcli == nil {
		return C.CString("p2phcli has not created")
	}

	/*
	MessageHandler := func(requestData []byte, head service.Head) ([]byte, error){
		fmt.Println(string(requestData))
		return []byte("ok!!!"), nil
	}
	*/

	//p2phst.RegisterHandler(0x14, MessageHandler)
	p2phst.RegisterGlobalMsgHandler(p2phcli.MessageHandler)
	return nil
}

//export UnregisterHandlerWrp
func UnregisterHandlerWrp(msgType *C.char) *C.char {
	if p2phst == nil {
		return C.CString("p2phost has not started")
	}
	//gmsgType := C.GoString(msgType)
	//p2phst.RemoveHandler(gmsgType)
	p2phst.RemoveGlobalHandler()
	return nil
}

//export FreeString
func FreeString(ptr unsafe.Pointer) {
	C.free(ptr)
}

//export FreeIDRet
func FreeIDRet(ptr *C.idret) {
	if ptr != nil {
		if (*ptr).id != nil {
			C.free(unsafe.Pointer((*ptr).id))
			(*ptr).id = nil
		}
		if (*ptr).error != nil {
			C.free(unsafe.Pointer((*ptr).error))
			(*ptr).error = nil
		}
		C.free(unsafe.Pointer(ptr))
	}
	//C.freeIdRet((*C.idret)(ptr))
}

//export FreeAddrsRet
func FreeAddrsRet(ptr *C.addrsret) {
	if ptr != nil {
		if (*ptr).addrs != nil {
			C.freeCharArray((*ptr).addrs, (*ptr).size)
			(*ptr).addrs = nil
		}
		if (*ptr).error != nil {
			C.free(unsafe.Pointer((*ptr).error))
			(*ptr).error = nil
		}
		C.free(unsafe.Pointer(ptr))
	}
	//C.freeAddrsRet((*C.addrsret)(ptr))
}

//export FreeSendMsgRet
func FreeSendMsgRet(ptr *C.sendmsgret) {
	if ptr != nil {
		if (*ptr).msg != nil {
			C.free(unsafe.Pointer((*ptr).msg))
			(*ptr).msg = nil
		}
		if (*ptr).error != nil {
			C.free(unsafe.Pointer((*ptr).error))
			(*ptr).error = nil
		}
		C.free(unsafe.Pointer(ptr))
	}

	//C.freeSendMsgRet((*C.sendmsgret)(ptr))
}

//export CreateIDRet
func CreateIDRet(id *C.char, err *C.char) *C.idret {
	ptr := (*C.idret)(C.malloc(C.size_t(unsafe.Sizeof(C.idret{}))))
	C.memset(unsafe.Pointer(ptr), 0, C.size_t(unsafe.Sizeof(C.idret{})))
	if id != nil {
		(*ptr).id = id
	}
	if err != nil {
		(*ptr).error = err
	}
	return ptr
}

//export CreateAddrsRet
func CreateAddrsRet(addrs **C.char, size C.int, err *C.char) *C.addrsret {
	ptr := (*C.addrsret)(C.malloc(C.size_t(unsafe.Sizeof(C.addrsret{}))))
	C.memset(unsafe.Pointer(ptr), 0, C.size_t(unsafe.Sizeof(C.addrsret{})))
	if addrs != nil {
		(*ptr).addrs = addrs
		(*ptr).size = size
	}
	if err != nil {
		(*ptr).error = err
	}
	return ptr
}

//export CreateSendMsgRet
func CreateSendMsgRet(msg *C.char, size C.longlong, err *C.char) *C.sendmsgret {
	ptr := (*C.sendmsgret)(C.malloc(C.size_t(unsafe.Sizeof(C.sendmsgret{}))))
	C.memset(unsafe.Pointer(ptr), 0, C.size_t(unsafe.Sizeof(C.sendmsgret{})))
	if msg != nil {
		(*ptr).msg = msg
		(*ptr).size = size
	}
	if err != nil {
		(*ptr).error = err
	}
	return ptr
}

//export CreateSendMsgRet2
func CreateSendMsgRet2(msg *C.char, size C.longlong, err *C.char) *C.sendmsgret {
	ptr := (*C.sendmsgret)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_sendmsgret{}))))
	C.memset(unsafe.Pointer(ptr), 0, C.size_t(unsafe.Sizeof(C.struct_sendmsgret{})))
	if msg != nil {
		msgCp := C.malloc((C.size_t(size)))
		C.memcpy(msgCp, unsafe.Pointer(msg), C.size_t(size))
		(*ptr).msg = (*C.char)(msgCp)
		(*ptr).size = size
	}
	if err != nil {
		(*ptr).error = err
	}
	return ptr
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

func main() {
	//分别在不同进程启动cstart和sstart方法来模拟服务端和客户端
	C.sstart()
	//C.cstart()
}

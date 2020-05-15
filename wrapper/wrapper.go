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
 	char *privkey = "16Uiu2HAmPR1qWUmFLatKf8QmHtJ3fkQpjP4tSa99wYbWvcvkzwYw";
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
 	char *privkey = "16Uiu2HAmPR1qWUmFLatKf8QmHtJ3fkQpjP4tSa99wYbWvcvkzwYw";
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
	err = ConnectWrp("16Uiu2HAmPR1qWUmFLatKf8QmHtJ3fkQpjP4tSa99wYbWvcvkzwYw", addrs, 1);
	if (err != NULL) {
		printf("error: %s\n", err);
		free(err);
		err = NULL;
		return;
	}

	char data[12] = {'s','e','n','d',' ','m','e','s','s','a','g','e'};
	char msid[2] ;
	msid[0] = 0;
	msid[1] = 0;
	sendmsgret* retp3 = SendMsgWrp("16Uiu2HAmPR1qWUmFLatKf8QmHtJ3fkQpjP4tSa99wYbWvcvkzwYw", msid, data, 12);
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
	"github.com/libp2p/go-libp2p-core/peer"
	base58 "github.com/mr-tron/base58"
	//"github.com/prometheus/common/log"
	lg "github.com/yottachain/P2PHost/log"
	"github.com/yottachain/P2PHost/pb"
	"github.com/yottachain/YTHost/option"
	"net"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"time"
	"unsafe"

	crypto "github.com/libp2p/go-libp2p-core/crypto"
	ma "github.com/multiformats/go-multiaddr"
	p2ph "github.com/yottachain/P2PHost"
	hst "github.com/yottachain/YTHost"
	"github.com/yottachain/YTHost/clientPool"
	host "github.com/yottachain/YTHost/hostInterface"
	"google.golang.org/grpc"
)

var p2phst host.Host
var mu sync.Mutex
var p2phcli p2ph.Hclient
var CliPool *clientPool.ClientPool

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

	p2phst, err = hst.NewHost(option.ListenAddr(ma), option.OpenPProf("0.0.0.0:10000"), option.Identity(pk))
	if err != nil {
		return C.CString(err.Error())
	}

	p2phcli, err = p2ph.NewHclient()
	if err != nil {
		return C.CString(err.Error())
	}

	go p2phst.Accept()

	p2phst.RegisterGlobalMsgHandler(p2phcli.MessageHandler)
	lg.Info.Printf("configure callback handler successful.")

	nodelist := hst.GetACNodeList()
	CliPool = clientPool.NewPool(p2phst, nodelist)

	server := &p2ph.Server{Host: p2phst, Hc: p2phcli, CliPool:CliPool}

	p2pGRPCPortStr := os.Getenv("P2PHOST_GRPCPORT")
	p2pGRPCPort, err := strconv.Atoi(p2pGRPCPortStr)
	if err != nil {
		p2pGRPCPort = 11002
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", p2pGRPCPort))
	if err != nil {
		lg.Info.Fatalf("failed to listen GRPC port %d: %v", p2pGRPCPort, err)
		return C.CString(err.Error())
	}
	lg.Info.Printf("GRPC address: 0.0.0.0:%d\n", p2pGRPCPort)
	grpcServer := grpc.NewServer()
	pb.RegisterP2PHostServer(grpcServer, server)

	go func(ser *grpc.Server) {
		err = grpcServer.Serve(lis)
		if err == nil {
			lg.Info.Printf("GRPC server started.")
		}else {
			lg.Info.Printf("GRPC server start fail.")
		}
	}(grpcServer)

	return nil
}

/*func StartWrpBak(port C.int, privkey *C.char) *C.char {
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

	p2phst, err = hst.NewHost(option.ListenAddr(ma), option.OpenPProf("0.0.0.0:10000") ,option.Identity(pk))
	if err != nil {
		return C.CString(err.Error())
	}

	p2phcli, err = p2ph.NewHclient()
	if err != nil {
		return C.CString(err.Error())
	}

	go p2phst.Accept()

	return nil
}*/

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

	maddrs, err := stringListToMaddrs(gaddrs)

	conntimeout := os.Getenv("P2PHOST_CONNECTTIMEOUT")
	ct := 60
	if conntimeout == "" {
		ct = 60
	}else {
		ct, err = strconv.Atoi(conntimeout)
		if err != nil {
			ct = 60
		}
	}

	nodeIdStr := C.GoString(nodeID)
	ID, err := peer.Decode(nodeIdStr)
	if err != nil {
		return C.CString(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ct))
	defer cancel()
	_, err = p2phst.ClientStore().Get(ctx, ID, maddrs)
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

	nodeIdStr := C.GoString(nodeID)
	ID, err := peer.Decode(nodeIdStr)
	if err != nil {
		return C.CString(err.Error())
	}
	err = p2phst.ClientStore().Close(ID)
	if err != nil {
		return C.CString(err.Error())
	}

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

	c_msg := (*[1 << 30]byte)(unsafe.Pointer(msg))[:int64(size):int64(size)]
	s := int64(size)
	//多了一次copy
	msgSlice := make([]byte, s)
	copy(msgSlice, c_msg)

	conntimeout := os.Getenv("P2PHOST_WRITETIMEOUT")
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

	ID, err := peer.Decode(nodeIDStr)
	if err != nil {
		return CreateSendMsgRet(nil, C.longlong(0), C.CString(err.Error()))
	}

	ret, err := p2phst.SendMsg(ctx, ID, msgId, msgSlice)

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


	/*MessageHandler := func(requestData []byte, head service.Head) ([]byte, error){
		fmt.Println(string(requestData))
		return []byte("ok!!!"), nil
	}

	p2phst.RegisterGlobalMsgHandler(MessageHandler)*/
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

////export GetOptNodes
//func GetOptNodes(ids **C.char, size C.int) **C.char{
//	if p2phst == nil {
//		return nil
//	}
//
//	length := int(size)
//	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(ids))[:length:length]
//	iids := make([]string, length)
//	for i, s := range tmpslice {
//		iids[i] = C.GoString(s)
//	}
//
//	oids := p2phst.Optmizer().Get2(iids...)
//	ptr := (**C.char)(C.malloc(C.size_t(unsafe.Sizeof(*C.char))*len(oids)))
//	C.memset(unsafe.Pointer(ptr), 0, C.size_t(unsafe.Sizeof(*C.char))*len(oids))
//
//	for i, s := range oids {
//		ptr[i] = C.CString(s)
//	}
//
//	return ptr
//}

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
	//C.sstart()
	C.cstart()
}

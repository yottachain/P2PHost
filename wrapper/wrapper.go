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
extern sendmsgret* SendMsgWrp(char *nodeID, char *msgType, char *msg, long long size);
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
 	int port = 8888;
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

	sleep(10000);
}

static void cstart() {
	int port = 9999;
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

	char *addrs[1] = {"/ip4/127.0.0.1/tcp/8888"};
	err = ConnectWrp("16Uiu2HAmAvd2jETZcJL3pwqaRBU9UP6bhZLXRPqFNyiSVZRBftxJ", addrs, 1);
	if (err != NULL) {
		printf("error: %s\n", err);
		free(err);
		err = NULL;
		return;
	}

	char data[12] = {'s','e','n','d',' ','m','e','s','s','a','g','e'};
	sendmsgret* retp3 = SendMsgWrp("16Uiu2HAmAvd2jETZcJL3pwqaRBU9UP6bhZLXRPqFNyiSVZRBftxJ", "test", data, 12);
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
	"fmt"
	_ "net/http/pprof"
	"sync"
	"unsafe"

	host "github.com/yottachain/P2PHost"
)

var p2phst host.Host
var mu sync.Mutex

//export StartWrp
func StartWrp(port C.int, privkey *C.char) *C.char {
	mu.Lock()
	defer mu.Unlock()
	if p2phst != nil {
		return C.CString("p2phost has started")
	}

	pk := C.GoString(privkey)
	pt := int(port)
	var err error
	p2phst, err = host.NewHost(pk, fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", pt))
	//_ = p2phst
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export IDWrp
func IDWrp() *C.idret {
	if p2phst == nil {
		return CreateIDRet(nil, C.CString("p2phost has not started"))
	}
	id := p2phst.ID()
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
	addrs := p2phst.Addrs()
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
	p2phst.Close()
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
	err := p2phst.Connect(nodeIdStr, gaddrs)
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
	err := p2phst.DisConnect(C.GoString(nodeID))
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export SendMsgWrp
func SendMsgWrp(nodeID *C.char, msgType *C.char, msg *C.char, size C.longlong) *C.sendmsgret {
	if p2phst == nil {
		return CreateSendMsgRet(nil, 0, C.CString("p2phost has not started"))
	}
	nodeIDStr := C.GoString(nodeID)
	msgTypeStr := C.GoString(msgType)
	msgSlice := (*[1 << 30]byte)(unsafe.Pointer(msg))[:int64(size):int64(size)]
	ret, err := p2phst.SendMsg(nodeIDStr, msgTypeStr, msgSlice)
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
	// handler := func(msgType string, data []byte, pubkey string) ([]byte, error) {
	// 	cmsgType := C.CString(msgType)
	// 	cpubkey := C.CString(pubkey)
	// 	cdata := C.CBytes(data)
	// 	csize := C.longlong(len(data))
	// 	defer C.free(unsafe.Pointer(cmsgType))
	// 	defer C.free(unsafe.Pointer(cpubkey))
	// 	defer C.free(unsafe.Pointer(cdata))
	// 	ret := C.executeHandler((*[0]byte)(f), cmsgType, (*C.char)(cdata), csize, cpubkey)
	// 	defer FreeSendMsgRet(ret)
	// 	if ret.error != nil {
	// 		return nil, errors.New(C.GoString(ret.error))
	// 	}
	// 	retdata := C.GoBytes(unsafe.Pointer(ret.msg), C.int(ret.size))
	// 	return retdata, nil
	// }
	p2phst.RegisterHandler(C.GoString(msgType), nil)
	return nil
}

//export UnregisterHandlerWrp
func UnregisterHandlerWrp(msgType *C.char) *C.char {
	if p2phst == nil {
		return C.CString("p2phost has not started")
	}
	gmsgType := C.GoString(msgType)
	p2phst.UnregisterHandler(gmsgType)
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

func main() {
	//分别在不同进程启动cstart和sstart方法来模拟服务端和客户端
	C.cstart()
	//C.cstart()
}

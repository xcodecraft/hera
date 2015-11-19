package qbus

/*
#include <stdlib.h>
#include <assert.h>

#ifndef __CUTILS_H__
#define __CUTILS_H__

struct Slice {
    int length;
    char* data;
};

struct StringVector {
    int length;
    struct Slice data[0];
};

#endif


#ifndef __CCONSUMER4STORM__
#define __CCONSUMER4STORM__

#include <sys/types.h>
#include <stdint.h>

struct Consumer4StormHandle;

struct Consumer4StormHandle* newConsumer4Storm(const char* zk_cluster, const char *topic, const char* group, const char* conf_path);
void deleteConsumer4Storm(struct Consumer4StormHandle* handle);
int getErrno(struct Consumer4StormHandle* handle);
struct StringVector* nextMessages(struct Consumer4StormHandle* handle);
uint8_t commitOffset(struct Consumer4StormHandle* handle);

#endif
struct Slice* getSlice(struct StringVector* sv, int offset) {
    assert(sv != NULL && offset < sv->length);

    return &sv->data[offset];
}
*/
import "C"

import(
    "unsafe"
)

type Consumer4Storm struct {
    handle *C.struct_Consumer4StormHandle
}

func (this *Consumer4Storm) Init(zk_cluster_ string, topic_ string, group_ string, conf_path_ string) {
    zk_cluster := C.CString(zk_cluster_)
    topic := C.CString(topic_)
    group := C.CString(group_)
    conf_path := C.CString(conf_path_)
    defer C.free(unsafe.Pointer(zk_cluster))
    defer C.free(unsafe.Pointer(topic))
    defer C.free(unsafe.Pointer(group))
    defer C.free(unsafe.Pointer(conf_path))

    this.handle = C.newConsumer4Storm(zk_cluster, topic, group, conf_path)
}

func (this Consumer4Storm) GetErrno() int {
    return int(C.getErrno(this.handle))
}

func (this Consumer4Storm) CommitOffset() bool {
    ret := C.commitOffset(this.handle)
    return ret == C.uint8_t(1)
}


func (this Consumer4Storm) NextMessages() [][]byte {
    var ret [][]byte
    var sv *C.struct_StringVector

    sv = C.nextMessages(this.handle)
    if sv == nil {
        return nil
    }
    for i := 0; i < int(sv.length); i++ {
        sl := C.getSlice(sv, C.int(i))
        if sl == nil {
            continue
        }
        buf := C.GoBytes(unsafe.Pointer(sl.data), sl.length)
        ret = append(ret, buf)
    }
    C.free(unsafe.Pointer(sv))

    return ret
}


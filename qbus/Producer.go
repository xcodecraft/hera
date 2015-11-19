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

#ifndef __CPRODUCER_H__
#define __CPRODUCER_H__

#include <sys/types.h>
#include <stdint.h>

struct ProducerHandle;

struct ProducerHandle* getInstance(const char* zk_cluster, const char* conf_path, uint8_t safe_send);
uint8_t setSendAck(struct ProducerHandle* handle, uint8_t safe_send);
uint8_t setSendTimeout(struct ProducerHandle* handle, int seconds);
uint8_t setRecvTimeout(struct ProducerHandle* handle, int seconds);
uint8_t setConnMax(struct ProducerHandle* handle, int max);
uint8_t asyncSend(struct ProducerHandle* handle, const struct StringVector* cmessages, const char* ctopic, char** cerrstr, const struct Slice* csequence);
uint8_t syncSend(struct ProducerHandle* handle, const struct StringVector* cmessages, const char* ctopic, char** cerrstr, int flag, const struct Slice* csemantic_key, uint8_t zflag);

#endif

struct StringVector *allocStringVector(int len) {
    struct StringVector* ret = (struct StringVector *)malloc(sizeof(struct StringVector)+len*sizeof(struct Slice));
    if (ret != NULL)
        ret->length = len;
    return ret;
}

void insertSlice(struct StringVector* sv, int offset, char *addr, int len) {
    assert(sv != NULL && offset < sv->length);
    sv->data[offset].length = len;
    sv->data[offset].data = addr;
}

void cleanStringVector(struct StringVector *sv) {
    if (sv == NULL) return;
    int i;
    for (i = 0; i < sv->length; i++)
        sv->data[i].data = NULL;
    free(sv);
}

void fillSlice(struct Slice* slice, int len, char *data) {
    slice->length = len;
    slice->data = data;
}
*/
import "C"

import (
    "errors"
    "unsafe"
)

const (
    RANDOM_SEND = 1
    AFFINITY_SEND = 2
    SEMANTIC_SEND = 3
    BATCH_SEND = 8
)

type Producer struct {
    handle *C.struct_ProducerHandle
}

func convertBoolToUint8(flag bool) C.uint8_t {
    if (flag) {
        return C.uint8_t(1)
    } else {
        return C.uint8_t(0)
    }
}

func (this *Producer) Init(zk_cluster_ string, conf_path_ string, safe_send_ bool) bool {
    zk_cluster := C.CString(zk_cluster_)
    conf_path := C.CString(conf_path_)
    safe_send := convertBoolToUint8(safe_send_)
    defer C.free(unsafe.Pointer(zk_cluster))
    defer C.free(unsafe.Pointer(conf_path))
    this.handle = C.getInstance(zk_cluster, conf_path, safe_send)

    return this.handle != nil
}

func (this Producer) SetSendAck(safe_send_ bool) bool {
    var safe_send, ret C.uint8_t
    safe_send = convertBoolToUint8(safe_send_)
    ret = C.setSendAck(this.handle, safe_send)

    return ret == C.uint8_t(1)
}

func (this Producer) SetSendTimeout(seconds int) bool {
    var ret C.uint8_t
    ret = C.setSendTimeout(this.handle, C.int(seconds))

    return ret == C.uint8_t(1)
}

func (this Producer) SetRecvTimeout(seconds int) bool {
    var ret C.uint8_t
    ret = C.setRecvTimeout(this.handle, C.int(seconds))

    return ret == C.uint8_t(1)
}

func (this Producer) SetConnMax(max int) bool {
    var ret C.uint8_t
    ret = C.setConnMax(this.handle, C.int(max))

    return ret == C.uint8_t(1)
}

func (this Producer) SyncSend(buffers_ [][]byte, topic_ string, flag_ int, semantic_key_ []byte, zflag_ bool) (bool, error){
    topic := C.CString(topic_)
    var errstr *C.char
    var zflag C.uint8_t

    zflag = convertBoolToUint8(zflag_)

    defer C.free(unsafe.Pointer(topic))

    var sv *C.struct_StringVector
    sv = C.allocStringVector(C.int(len(buffers_)))
    if sv == nil {
        return false, errors.New("SyncSend failed when allocating memory for StringVector")
    }

    for i := 0; i < len(buffers_); i++ {
        C.insertSlice(sv, C.int(i), (*C.char)(unsafe.Pointer(&buffers_[i][0])), C.int(len(buffers_[i])))
    }

//    var semantic_key C.struct_Slice
//    C.fillSlice((*C.struct_Slice)(unsafe.Pointer(&semantic_key)),C.int(len(semantic_key_)), (*C.char)(unsafe.Pointer(&semantic_key_[0])))
    var ret C.uint8_t
    if semantic_key_ == nil || len(semantic_key_) == 0{
        ret = C.syncSend(this.handle, sv, topic, (**C.char)(unsafe.Pointer(&errstr)), C.int(flag_), nil, zflag)
    } else {
        var semantic_key C.struct_Slice
        semantic_key.length = C.int(len(semantic_key_))
        semantic_key.data = (*C.char)(unsafe.Pointer(&semantic_key_[0]))
        ret = C.syncSend(this.handle, sv, topic, (**C.char)(unsafe.Pointer(&errstr)), C.int(flag_), &semantic_key, zflag)
    }
    C.cleanStringVector(sv)

    if ret == C.uint8_t(0) {
        defer C.free(unsafe.Pointer(errstr))
        return false, errors.New(C.GoString(errstr))
    } else {
        return true, nil
    }
}

func (this Producer) RandomSend(buffers_ [][]byte, topic_ string) (bool, error) {
    return this.SyncSend(buffers_, topic_, RANDOM_SEND, []byte(""), false)
}

func (this Producer) SequencedAsyncSend(buffers_ [][]byte, topic_ string, sequence_ []byte) (bool, error) {
    topic := C.CString(topic_)
    defer C.free(unsafe.Pointer(topic))

    var errstr *C.char
    var sv *C.struct_StringVector

    sv = C.allocStringVector(C.int(len(buffers_)))
    if sv == nil {
        return false, errors.New("AsyncSend failed when allocating memory for StringVector")
    }

    for i := 0; i < len(buffers_); i++ {
        C.insertSlice(sv, C.int(i), (*C.char)(unsafe.Pointer(&buffers_[i][0])), C.int(len(buffers_[i])))
    }
    var ret C.uint8_t
    if sequence_ == nil || len(sequence_) == 0 {
        ret = C.asyncSend(this.handle, sv, topic, (**C.char)(unsafe.Pointer(&errstr)), nil)
    } else {
        var sequence C.struct_Slice
        sequence.length = C.int(len(sequence_))
        sequence.data = (*C.char)(unsafe.Pointer(&sequence_[0]))
        ret = C.asyncSend(this.handle, sv, topic, (**C.char)(unsafe.Pointer(&errstr)), &sequence)
    }

    if ret == C.uint8_t(0) {
        defer C.free(unsafe.Pointer(errstr))
        return false, errors.New(C.GoString(errstr))
    } else {
        return true, nil
    }

}


func (this Producer) AsyncSend(buffers_ [][]byte, topic_ string) (bool, error) {
    return this.SequencedAsyncSend(buffers_, topic_, []byte(""))
}

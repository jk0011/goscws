package goscws

/*
#cgo LDFLAGS: -lscws
#include "goscws.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

const (
	SCWS_XDICT_XDB     = C.SCWS_XDICT_XDB     // 表示直接读取 xdb 文件
	SCWS_XDICT_MEM     = C.SCWS_XDICT_MEM     // 表示将 xdb 文件全部加载到内存中，以 XTree 结构存放，可用异或结合另外2个使用。
	SCWS_XDICT_TXT     = C.SCWS_XDICT_TXT     // 表示要读取的词典文件是文本格式，可以和后2项结合用
	SCWS_MULTI_SHORT   = C.SCWS_MULTI_SHORT   // 短词
	SCWS_MULTI_DUALITY = C.SCWS_MULTI_DUALITY // 二元（将相邻的2个单字组合成一个词）
	SCWS_MULTI_ZMAIN   = C.SCWS_MULTI_ZMAIN   // 重要单字
	SCWS_MULTI_ZALL    = C.SCWS_MULTI_ZALL    // 全部单字
)

type Result struct {
	Word []byte
	Idf  float32
	Attr []byte
}

type TopWord struct {
	Word   []byte
	Weight float32
	Times  int
	Attr   []byte
}

type GoScws struct {
	scws C.PSTScws
	text []byte
}

func NewScws() *GoScws {
	gs := &GoScws{}
	gs.scws = C.NewScws()
	return gs
}

func (this *GoScws) ForkScws(gs *GoScws) *GoScws {
	ret := C.ForkScws(this.scws)
	if ret == nil {
		return nil
	}
	return &GoScws{scws: ret, text: gs.text}
}

func (this *GoScws) DeleteScws() {
	C.DeleteScws(this.scws)
}

func (this *GoScws) SetCharset(cs string) {
	C.SetCharset(this.scws, C.CString(cs))
}

func (this *GoScws) SetDict(fPath string, mode int) error {
	var ret C.int
	ret = C.SetDict(this.scws, C.CString(fPath), C.int(mode))
	if ret != 0 {
		return errors.New("set dict error")
	}
	return nil
}

func (this *GoScws) SetIgnore(yes int) error {
	var ret C.int
	ret = C.SetIgnore(this.scws, C.int(yes))
	if ret != 0 {
		return errors.New("set ignore error")
	}
	return nil
}

func (this *GoScws) SetDuality(yes int) error {
	var ret C.int
	ret = C.SetDuality(this.scws, C.int(yes))
	if ret != 0 {
		return errors.New("set duality error")
	}
	return nil
}

func (this *GoScws) SetRule(fPath string) error {
	var ret C.int
	ret = C.SetRule(this.scws, C.CString(fPath))
	if ret != 0 {
		return errors.New("set rule error")
	}
	return nil
}

func (this *GoScws) SetMulti(mode int) error {
	var ret C.int
	ret = C.SetMulti(this.scws, C.int(mode))
	if ret != 0 {
		return errors.New("set multi error")
	}
	return nil
}

func (this *GoScws) SendText(text []byte, len int) error {
	var ret C.int
	ret = C.SendText(this.scws, C.CString(string(text)), C.int(len))
	this.text = text
	if ret != 0 {
		return errors.New("send text error")
	}
	return nil
}

func (this *GoScws) HasWord(attr string) int {
	var ret C.int
	ret = C.HasWord(this.scws, C.CString(attr))
	return int(ret)
}

func (this *GoScws) GetResult() *Result {
	res := &Result{}
	r := C.GetResult(this.scws)
	if r == nil {
		return nil
	}
	res.Idf = float32(r.idf)
	len := r.off + C.int(r.len)
	res.Word = this.text[r.off:len]
	c_char1 := *(*C.char)(unsafe.Pointer(&r.attr[0]))
	c_char2 := *(*C.char)(unsafe.Pointer(&r.attr[1]))
	res.Attr = append(res.Attr, byte(c_char1))
	if C.int(c_char2) != 0 {
		res.Attr = append(res.Attr, byte(c_char2))
	}
	return res
}

func (this *GoScws) GetTops(others ...interface{}) *TopWord {
	res := &TopWord{}
	var limit int = 10
	var attr []byte = nil
	for _, value := range others {
		switch value := value.(type) {
		case int:
				limit = value
		case string:
				attr = []byte(value)
		}
	}
	r := C.GetTops(this.scws, C.int(limit), C.CString(string(attr)))
	if r == nil {
		return nil
	}
	word := C.GoString(r.word)
	res.Word = []byte(word)
	res.Weight = float32(r.weight)
	c_char1 := *(*C.char)(unsafe.Pointer(&r.attr[0]))
	c_char2 := *(*C.char)(unsafe.Pointer(&r.attr[1]))
	res.Attr = append(res.Attr, byte(c_char1))
	if C.int(c_char2) != 0 {
		res.Attr = append(res.Attr, byte(c_char2))
	}
	res.Times = int(r.times)
	return res
}

func (this *GoScws) GetWords(attr []byte) *TopWord {
	res := &TopWord{}
	r := C.GetWords(this.scws, C.CString(string(attr)))
	if r == nil {
		return nil
	}
	word := C.GoString(r.word)
	res.Word = []byte(word)
	res.Weight = float32(r.weight)
	c_char1 := *(*C.char)(unsafe.Pointer(&r.attr[0]))
	c_char2 := *(*C.char)(unsafe.Pointer(&r.attr[1]))
	res.Attr = append(res.Attr, byte(c_char1))
	if C.int(c_char2) != 0 {
		res.Attr = append(res.Attr, byte(c_char2))
	}
	res.Times = int(r.times)
	return res
}

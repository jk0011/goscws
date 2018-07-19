package goscws

/*
#cgo CFLAGS: -I /usr/local/scws/include 
#cgo LDFLAGS: -L /usr/local/scws/lib -lscws 
#include "goscws.h"
*/
import "C"
// 拷贝scws源码到这里
// 拷贝编译好的libscws.so到这里
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

// 分配或初始化与 GoScws 系列操作的 `scws_st` 对象。
func NewScws() *GoScws {
	gs := &GoScws{}
	gs.scws = C.NewScws()
	return gs
}

// 在已有 GoScws 对象上产生一个分支，可以独立用于某个协程分词,
//以便共享内存词典、规则集，但它继承并共享父对象词典、
//规则集资源。同样需要调用 `DeleteScws()` 来释放对象。
//在该分支对象上重设词典、规则集不会影响父对象及其它分支。
func (this *GoScws) ForkScws(gs *GoScws) *GoScws {
	ret := C.ForkScws(this.scws)
	if ret == nil {
		return nil
	}
	return &GoScws{scws: ret, text: gs.text}
}

// 释放 GoScws 系列操作的 `scws_st` 对象及对象内容，
//同时也会释放已经加载的词典和规则。
func (this *GoScws) DeleteScws() {
	C.DeleteScws(this.scws)
}

// 设定当前 scws 所使用的字符集, 缺省gbk
func (this *GoScws) SetCharset(cs string) {
	C.SetCharset(this.scws, C.CString(cs))
}


// 添加词典文件到当前 scws 对象
//SCWS_XDICT_TXT 读取的词典为文本格式
//SCWS_XDICT_XDB 直接读取xdb文件
//SCWS_XDICT_MEM 加载xdb文件直接到内存中
//先释放已加载的所有词典，在加载指定的词典。
func (this *GoScws) SetDict(fPath string, mode int) error {
	var ret C.int
	ret = C.SetDict(this.scws, C.CString(fPath), C.int(mode))
	if ret != 0 {
		return errors.New("set dict error")
	}
	return nil
}

// 不会释放已加载的词典，若已经加载过该词典，则新加入的词典具有更高的优先权。
func (this *GoScws) AddDict(fPath string, mode int) error {
        var ret C.int
	ret = C.AddDict(this.scws, C.CString(fPath), C.int(mode))
	if ret != 0 {
		return errors.New("add dict error")
	}
	return nil
}

// 设定规则集文件
func (this *GoScws) SetRule(fPath string) error {
	var ret C.int
	ret = C.SetRule(this.scws, C.CString(fPath))
	if ret != 0 {
		return errors.New("set rule error")
	}
	return nil
}

// 设定分词结果是否忽略标点符号 1 忽略 0 不忽略
func (this *GoScws) SetIgnore(yes int) error {
	var ret C.int
	ret = C.SetIgnore(this.scws, C.int(yes))
	if ret != 0 {
		return errors.New("set ignore error")
	}
	return nil
}

// 是否将闲散文字自动以二字分词法聚合
//如果为 1 表示执行二分聚合，0 表示不处理，缺省为 0
func (this *GoScws) SetDuality(yes int) error {
	var ret C.int
	ret = C.SetDuality(this.scws, C.int(yes))
	if ret != 0 {
		return errors.New("set duality error")
	}
	return nil
}

// 是否针对长词符合切分
//SCWS_MULTI_SHORT   短词
//SCWS_MULTI_DUALITY 二元（将相邻的2个单字组合成一个词）
//SCWS_MULTI_ZMAIN   重要单字
//SCWS_MULTI_ZALL    全部单字
//例如 mode := SCWS_MULTI_SHORT | SCWS_MULTI_DUALITY | SCWS_MULTI_ZMAIN
func (this *GoScws) SetMulti(mode int) error {
	var ret C.int
	ret = C.SetMulti(this.scws, C.int(mode))
	if ret != 0 {
		return errors.New("set multi error")
	}
	return nil
}

// 返回需要的mode 
//例如 [SCWS_MULTI_SHORT, SCWS_MULTI_DUALITY, SCWS_MULTI_ZMAIN]
//返回 ret = SCWS_MULTI_SHORT | SCWS_MULTI_DUALITY | SCWS_MULTI_ZMAIN
func (this *GoScws)GetMulti(param []int32) C.int {
	var ret C.int
	for _, v := range param {
		ret |= C.int(v)
	}
	return  ret
}


// 这个函数应在 `GetResult()` 和 `GetTops()` `HasWord()` `GetWords` 之前调用,防止出错
//连续调用后会覆盖之前的设定；故不应在多次的 GetResult
//循环中再调用 SendText() 以免出错。
func (this *GoScws) SendText(text []byte, len int) error {
	var ret C.int
	ret = C.SendText(this.scws, C.CString(string(text)), C.int(len))
	this.text = text
	if ret != 0 {
		return errors.New("send text error")
	}
	return nil
}

// 返回值 如果有返回 1 没有则返回 0
//参数 attr 用来描述要排除或参与的统计词汇词性，多个词性之间用逗号隔开
//当以~开头时表示统计结果中不包含这些词性，否则表示必须包含，传入 nil 表示统计全部词性
func (this *GoScws) HasWord(attr string) bool {
	var ret C.int
	ret = C.HasWord(this.scws, C.CString(attr))
	if int(ret) == 1 {
	        return true       
	} else {
	        return false 
	}
}

// 返回分词结果集 Result，无返回nil
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

// 返回指定的关键词表统计集，系统会自动根据词语出现的次数及其 idf 值计算排名
//参数 limit 指定取回数据的最大条数，若传入值为0或负数，则自动重设为10
//参数 attr 用来描述要排除或参与的统计词汇词性，多个词性之间用逗号隔开
func (this *GoScws) GetTops(others ...interface{}) *TopWord {
	res := &TopWord{}
	var limit int = 10
	var attr []byte = nil
	for _, value := range others {
		switch value := value.(type) {
		case int:
		                // limit
				limit = value
		case string:
		                // attr
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

// 参数 attr 是一系列词性组成的字符串，各词性之间以半角的逗号隔开
//这表示返回的词性必须在列表中，如果以~开头，则表示取反，词性必须不在列表中，若为空则返回全部词
//返回值 成功返回符合要求词汇组成的TopWord，返回 nil。返回的词汇包含的键值参见 `GetResult()`
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

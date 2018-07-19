#ifndef _GO_SCWS_H_
#define _GO_SCWS_H_

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <scws/scws.h>
//
typedef struct STScws *PSTScws;

typedef void (*fnDelete)(PSTScws );

typedef void (*fnSetCharset)(PSTScws , const char *);

typedef int (*fnSetDict)(PSTScws , const char *, int);

typedef int (*fnSetRule)(PSTScws , const char *);

typedef int (*fnAddDict)(PSTScws , const char *, int);

typedef int (*fnSetIgnore)(PSTScws , int);

typedef int (*fnSetMulti)(PSTScws , int);

typedef int (*fnSetDuality)(PSTScws , int);

typedef int (*fnSetDebug)(PSTScws , int);

typedef int (*fnSendText)(PSTScws , const char *, int);

typedef scws_res_t (*fnGetResult)(PSTScws );

typedef int (*fnFreeResult)(PSTScws , scws_res_t);

typedef scws_top_t (*fnGetTops)(PSTScws , int, char *);

typedef int (*fnFreeTops)(PSTScws , scws_top_t);

typedef scws_top_t (*fnGetWords)(PSTScws , char *);

typedef int (*fnHasWord)(PSTScws , char *);

typedef PSTScws (*fnForkScws)(PSTScws );

struct STScws {
    fnDelete fn_delete;
    fnSetCharset fn_setCharset;
    fnSetDict fn_setDict;
    fnSetRule fn_setRule;
    fnAddDict fn_addDict;
    fnSetIgnore fn_setIgnore;
    fnSetMulti fn_setMulti;
    fnSetDuality fn_setDuality;
    fnSetDebug fn_setDebug;
    fnSendText fn_sendText;
    fnGetResult fn_getResult;
    fnFreeResult fn_freeResult;
    fnGetTops fn_getTops;
    fnFreeTops fn_freeTops;
    fnGetWords fn_getWords;
    fnHasWord fn_hasWord;
    fnForkScws fn_forkScws;
    scws_t fd_scws;
    scws_res_t fd_pTmpRes;
    scws_res_t fd_pRes;
    scws_res_t fd_pCurRes;
    scws_top_t fd_pTmpTop;
    scws_top_t fd_pTops;
    scws_top_t fd_pCurTops;
    scws_top_t fd_pTmpWords;
    scws_top_t fd_pWords;
    scws_top_t fd_pCurWords;
};


void DeleteScws(PSTScws pstScws);

void SetCharset(PSTScws pstScws, const char *cs);

int SetDict(PSTScws pstScws, const char *fpath, int mode);

int SetRule(PSTScws pstScws, const char *fpath);

int AddDict(PSTScws pstScws, const char *fpath, int mode);

int SetIgnore(PSTScws pstScws, int yes);

int SetMulti(PSTScws pstScws, int mode);

int SetDuality(PSTScws pstScws, int yes);

int SetDebug(PSTScws pstScws, int yes);

int SendText(PSTScws pstScws, const char *text, int len);

scws_res_t GetResult(PSTScws pstScws);

int FreeResult(PSTScws pstScws, scws_res_t res);

scws_top_t GetTops(PSTScws pstScws, int limit, char *attr);

int FreeTops(PSTScws pstScws, scws_top_t top);

scws_top_t GetWords(PSTScws pstScws, char *attr);

int HasWord(PSTScws pstScws, char *attr);

PSTScws  ForkScws(PSTScws pstScws);

PSTScws InitScws();

//分配或初始化与 scws 系列操作的 `scws_st` 对象。该函数将自动分配、初始化、并返回新对象的指针。
// 只能通过调用 `scws_free()` 来释放该对象
PSTScws NewScws() {
    PSTScws pS = InitScws();
    pS->fd_scws = scws_new();
    if (!pS->fd_scws) {
        printf("ERROR: cann't init the scws!\n");
        exit(-1);
    }
    return pS;
}

//在已有 scws 对象上产生一个分支，可以独立用于某个协程分词,以便共享内存词典、规则集，但它继承并共享父对象词典、
//   规则集资源。同样需要调用 `scws_free()` 来释放对象。在该分支对象上重设词典、规则集不会影响父对象及其它分支。
PSTScws  ForkScws(PSTScws pstScws) {
    PSTScws ps = InitScws();
    scws_t sw = scws_fork(pstScws->fd_scws);
    if (sw == NULL) {
        return NULL;
    }
    ps->fd_scws = sw;
    return ps;
}

PSTScws InitScws() {
    PSTScws pS = (PSTScws) malloc(sizeof(struct STScws));
    memset(pS, 0, sizeof(struct STScws));
    pS->fn_delete = DeleteScws;
    pS->fn_setCharset = SetCharset;
    pS->fn_setDict = SetDict;
    pS->fn_setRule = SetRule;
    pS->fn_addDict = AddDict;
    pS->fn_setIgnore = SetIgnore;
    pS->fn_setMulti = SetMulti;
    pS->fn_setDuality = SetDuality;
    pS->fn_setDebug = SetDebug;
    pS->fn_sendText = SendText;
    pS->fn_getResult = GetResult;
    pS->fn_freeResult = FreeResult;
    pS->fn_getTops = GetTops;
    pS->fn_freeTops = FreeTops;
    pS->fn_getWords = GetWords;
    pS->fn_hasWord = HasWord;
    pS->fn_forkScws = ForkScws;
    return pS;
}

//释放 scws 操作句柄及对象内容，同时也会释放已经加载的词典和规则
void DeleteScws(PSTScws pstScws) {
    scws_free(pstScws->fd_scws);
    pstScws->fd_scws = NULL;
    free(pstScws);
}

//设定当前 scws 所使用的字符集, 缺省gbk
void SetCharset(PSTScws pstScws, const char *cs) {
    scws_set_charset(pstScws->fd_scws, cs);
}

//添加词典文件到当前 scws 对象
//SCWS_XDICT_TXT 读取的词典为文本格式
//SCWS_XDICT_XDB 直接读取xdb文件
//SCWS_XDICT_MEM 加载xdb文件直接到内存中
//先释放已加载的所有词典，在加载指定的词典
int SetDict(PSTScws pstScws, const char *fpath, int mode) {
    int nRet = -1; //-1失败 0 成功
    pstScws->fd_scws->d = NULL;
    if (nRet = scws_set_dict(pstScws->fd_scws, fpath, mode) != 0) {
        printf("scws: 加载词典失败。");
    }
    return nRet;
}

//不会释放已加载的词典，若已经加载过该词典，则新加入的词典具有更高的优先权。
int AddDict(PSTScws pstScws, const char *fpath, int mode) {
    int nRet = -1; //-1失败 0 成功
    pstScws->fd_scws->d = NULL;
    if (nRet = scws_add_dict(pstScws->fd_scws, fpath, mode) != 0) {
        printf("scws: 添加词典失败。");
    }
    return nRet;
}

//设定规则集文件 -1 失败 0 成功
int SetRule(PSTScws pstScws, const char *fpath) {
    scws_set_rule(pstScws->fd_scws, fpath);
    if (pstScws->fd_scws->r == NULL) {
        return -1;
    }
    return 0;
}

//设定分词结果是否忽略标点符号 1 忽略 0 不忽略
int SetIgnore(PSTScws pstScws, int yes) {
    scws_set_ignore(pstScws->fd_scws, yes);
    return 0;
}

//是否针对长词符合切分
//SCWS_MULTI_SHORT   短词
//SCWS_MULTI_DUALITY 二元（将相邻的2个单字组合成一个词）
//SCWS_MULTI_ZMAIN   重要单字
//SCWS_MULTI_ZALL    全部单字
int SetMulti(PSTScws pstScws, int mode) {
    scws_set_multi(pstScws->fd_scws, mode);
    return 0;
}

//是否将闲散文字自动以二字分词法聚合
// yes** 如果为 1 表示执行二分聚合，0 表示不处理，缺省为 0
int SetDuality(PSTScws pstScws, int yes) {
    scws_set_duality(pstScws->fd_scws, yes);
    return 0;
}

//打印使用的是 `fprintf(stderr, ...)` 故不要随便用，并且只有编译时加入 --enable-debug 选项才有效
//设定分词时对于疑难多路径综合分词时，是否打印出各条路径的情况
int SetDebug(PSTScws pstScws, int yes) {
    scws_set_debug(pstScws->fd_scws, yes);
}


//这个函数应在 `scws_get_result()` 和 `scws_get_tops()` 之前调用,防止出错
//cws 结构内部维护着该字符串的指针和相应的偏移及长度，连续调用后会覆盖之前的设定；故不应在多次的 scws_get_result
//循环中再调用 scws_send_text() 以免出错
int SendText(PSTScws pstScws, const char *text, int len) {
    scws_send_text(pstScws->fd_scws, text, len);
    return 0;
}

//返回分词结果集，无返回NULL
scws_res_t GetResult(PSTScws pstScws) {
    if (pstScws->fd_pTmpRes == NULL) {
        pstScws->fd_pTmpRes = (scws_res_t)malloc(sizeof(struct scws_result));
    } else {
        free(pstScws->fd_pTmpRes);
        pstScws->fd_pTmpRes = (scws_res_t)malloc(sizeof(struct scws_result));
    }
    if (pstScws->fd_pCurRes == NULL) {
        pstScws->fd_pRes = pstScws->fd_pCurRes = scws_get_result(pstScws->fd_scws);
        if (pstScws->fd_pCurRes == NULL) {
            if (pstScws->fd_pTmpRes != NULL) {
                free(pstScws->fd_pTmpRes);
                pstScws->fd_pTmpRes = NULL;
            }
            return NULL;
        }
        memcpy(pstScws->fd_pTmpRes, pstScws->fd_pCurRes, sizeof(struct scws_result));
        if (pstScws->fd_pCurRes->next == NULL) {
            pstScws->fd_pCurRes = NULL;
            FreeResult(pstScws, pstScws->fd_pRes);
        } else {
            pstScws->fd_pCurRes = pstScws->fd_pCurRes->next;
        }
    } else {
        memcpy(pstScws->fd_pTmpRes, pstScws->fd_pCurRes, sizeof(struct scws_result));
        pstScws->fd_pCurRes = pstScws->fd_pCurRes->next;
        if (pstScws->fd_pCurRes == NULL) {
            FreeResult(pstScws, pstScws->fd_pRes);
        }
    }
    return pstScws->fd_pTmpRes;
}

//释放分词结果集，注意必须传入结果集头指针。
int FreeResult(PSTScws pstScws, scws_res_t res) {
    scws_free_result(res);
    return 0;
}

//参数 attr** 是一系列词性组成的字符串，各词性之间以半角的逗号隔开
//这表示返回的词性必须在列表中，如果以~开头，则表示取反，词性必须不在列表中，若为空则返回全部词
//返回值** 成功返回符合要求词汇组成的数组，返回 false。返回的词汇包含的键值参见 `scws_get_result`
scws_top_t GetWords(PSTScws pstScws, char *attr) {
    if (pstScws->fd_pTmpWords == NULL) {
        pstScws->fd_pTmpWords = (scws_top_t)malloc(sizeof(struct scws_topword));
    } else {
        free(pstScws->fd_pTmpWords);
        pstScws->fd_pTmpWords = (scws_top_t)malloc(sizeof(struct scws_topword));
    }

    if (pstScws->fd_pCurWords == NULL && pstScws->fd_pWords == NULL) {
        pstScws->fd_pWords =  pstScws->fd_pCurWords = scws_get_words(pstScws->fd_scws, attr);
    }

    if (pstScws->fd_pCurWords == NULL) {
        if (pstScws->fd_pTmpWords != NULL) {
            free(pstScws->fd_pTmpWords);
            pstScws->fd_pTmpWords = NULL;
        }

        FreeTops(pstScws, pstScws->fd_pWords);
        pstScws->fd_pWords = NULL;
        return NULL;
    }
    memcpy((char *)pstScws->fd_pTmpWords, (char *)pstScws->fd_pCurWords, sizeof(struct scws_topword));
    pstScws->fd_pCurWords = pstScws->fd_pCurWords->next;
    return pstScws->fd_pTmpWords;
}

//返回指定的关键词表统计集，系统会自动根据词语出现的次数及其 idf 值计算排名
//参数 limit 指定取回数据的最大条数，若传入值为0或负数，则自动重设为10
//参数 attr 用来描述要排除或参与的统计词汇词性，多个词性之间用逗号隔开
scws_top_t GetTops(PSTScws pstScws, int limit, char *attr) {
    if (pstScws->fd_pTmpTop == NULL) {
        pstScws->fd_pTmpTop = (scws_top_t)malloc(sizeof(struct scws_topword));
    } else {
        free(pstScws->fd_pTmpTop);
        pstScws->fd_pTmpTop = (scws_top_t)malloc(sizeof(struct scws_topword));
    }

    if (pstScws->fd_pCurTops == NULL && pstScws->fd_pTops == NULL) {
        pstScws->fd_pTops = pstScws->fd_pCurTops = scws_get_tops(pstScws->fd_scws, limit, attr);
    }
    if (pstScws->fd_pCurTops == NULL) {
        if (pstScws->fd_pTmpTop != NULL) {
            free(pstScws->fd_pTmpTop);
            pstScws->fd_pTmpTop = NULL;
        }
        FreeTops(pstScws, pstScws->fd_pTops);
        pstScws->fd_pTops = NULL;
        return NULL;
    }
    memcpy((char *)pstScws->fd_pTmpTop, (char *)pstScws->fd_pCurTops, sizeof(struct scws_topword));
    pstScws->fd_pCurTops = pstScws->fd_pCurTops->next;
    return pstScws->fd_pTmpTop;
}

//根据词表集的链表头释放词表集
int FreeTops(PSTScws pstScws, scws_top_t top) {
    scws_free_tops(top);
    return 0;
}

//返回值 如果有返回 1 没有则返回 0
//参数 xattr** 用来描述要排除或参与的统计词汇词性，多个词性之间用逗号隔开
//当以~开头时表示统计结果中不包含这些词性，否则表示必须包含，传入 NULL 表示统计全部词性
int HasWord(PSTScws pstScws, char *attr) {
    return scws_has_word(pstScws->fd_scws, attr);
}

#endif
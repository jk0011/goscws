## goscws 

scws是一个很好的中文分词库，由于项目中用到scws分词但go语言没有很好用的绑定，于是就自己
写了goscws，并且他支持协程分词。


## 例子：

根据 SendText 设定的文本内容，返回一系列切好的词汇

```
	gs := goscws.NewScws()
	gs.SetCharset("utf8")
	err := gs.SetDict("/usr/local/scws/etc/dict.utf8.xdb", goscws.SCWS_XDICT_MEM)
	if err != nil {
		log.Println(err)
	}
	err = gs.SetRule("/usr/local/scws/etc/rules.utf8.ini")
	if err != nil {
		log.Println(err)
	}
	err = gs.SetIgnore(0)
	if err != nil {
		log.Println(err)
	}
	text := []byte("Hello, 我名字叫李那曲是一个中国人, 我有时买Q币来玩, 我还听说过C#语言。。")
	gs.SendText(text, len(text))
	for  res := gs.GetResult(); res != nil;  res = gs.GetResult() {
			log.Println(string(res.Word), res.Idf, string(res.Attr))
	}
	gs.DeleteScws()
```

​	



根据 SendText 设定的文本内容，返回系统计算出来的最关键词汇列表

​	

```
	gs := goscws.NewScws()
	gs.SetCharset("utf8")
	err := gs.SetDict("/usr/local/scws/etc/dict.utf8.xdb", goscws.SCWS_XDICT_MEM)
	if err != nil {
		log.Println(err)
	}
	err = gs.SetRule("/usr/local/scws/etc/rules.utf8.ini")
	if err != nil {
		log.Println(err)
	}
	err = gs.SetIgnore(0)
	if err != nil {
		log.Println(err)
	}
	text := []byte("Hello, 我名字叫李那曲是一个中国人, 我有时买Q币来玩, 我还听说过C#语言。。")
	gs.SendText(text, len(text))
	for res := gs.GetTops(); res != nil; res = gs.GetTops() {
		log.Println(string(res.Word), res.Weight, res.Times,
			string(res.Attr))
	}
	gs.DeleteScws()
```



根据SendText 设定的文本内容，返回系统中词性符合要求的关键词汇

```
gs := goscws.NewScws()
	gs.SetCharset("utf8")
	err := gs.SetDict("/usr/local/scws/etc/dict.utf8.xdb", goscws.SCWS_XDICT_MEM)
	if err != nil {
		log.Println(err)
	}
	err = gs.SetRule("/usr/local/scws/etc/rules.utf8.ini")
	if err != nil {
		log.Println(err)
	}
	err = gs.SetIgnore(0)
	if err != nil {
		log.Println(err)
	}
	text := []byte("Hello, 我名字叫李那曲是一个中国人, 我有时买Q币来玩, 我还听说过C#语言。。")
	gs.SendText(text, len(text))
	for res := gs.GetWords([]byte("n")); res != nil; res = gs.GetWords(
		[]byte("n")) {
		log.Println(string(res.Word), res.Weight, res.Times,
			string(res.Attr))
	}
	gs.DeleteScws()
```


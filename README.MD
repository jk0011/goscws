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
	text := []byte(`陈凯歌并不是《无极》的唯一著作权人，一部电影的整体版权归电影制片厂所有。

一部电影的作者包括导演、摄影、编剧等创作人员，这些创作人员对他们的创作是有版权的。不经过制片人授权，其他人不能对电影做拷贝、发行、反映，不能通>过网络来传播，既不能把电影改编成小说、连环画等其他艺术形式发表，也不能把一部几个小时才能放完的电影改编成半个小时就能放完的短片。

著作权和版权在我国是同一个概念，是法律赋予作品创作者的专有权利。所谓专有权利就是没有经过权利人许可又不是法律规定的例外，要使用这个作品，就必须>经过作者授权，没有授权就是侵权。

一九八零年春天`)
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
	text := []byte(`陈凯歌并不是《无极》的唯一著作权人，一部电影的整体版权归电影制片厂所有。

一部电影的作者包括导演、摄影、编剧等创作人员，这些创作人员对他们的创作是有版权的。不经过制片人授权，其他人不能对电影做拷贝、发行、反映，不能通>过网络来传播，既不能把电影改编成小说、连环画等其他艺术形式发表，也不能把一部几个小时才能放完的电影改编成半个小时就能放完的短片。

著作权和版权在我国是同一个概念，是法律赋予作品创作者的专有权利。所谓专有权利就是没有经过权利人许可又不是法律规定的例外，要使用这个作品，就必须>经过作者授权，没有授权就是侵权。

一九八零年春天`)
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
	text := []byte(`陈凯歌并不是《无极》的唯一著作权人，一部电影的整体版权归电影制片厂所有。

一部电影的作者包括导演、摄影、编剧等创作人员，这些创作人员对他们的创作是有版权的。不经过制片人授权，其他人不能对电影做拷贝、发行、反映，不能通>过网络来传播，既不能把电影改编成小说、连环画等其他艺术形式发表，也不能把一部几个小时才能放完的电影改编成半个小时就能放完的短片。

著作权和版权在我国是同一个概念，是法律赋予作品创作者的专有权利。所谓专有权利就是没有经过权利人许可又不是法律规定的例外，要使用这个作品，就必须>经过作者授权，没有授权就是侵权。

一九八零年春天`)
	gs.SendText(text, len(text))
	for res := gs.GetWords([]byte("n")); res != nil; res = gs.GetWords(
		[]byte("n")) {
		log.Println(string(res.Word), res.Weight, res.Times,
			string(res.Attr))
	}
	gs.DeleteScws()
```


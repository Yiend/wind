package models

import "time"

//小说基本信息
type BookInfo struct {
	BookName, //书名
	Author,   //作者
	Sort,     //分类
	Content string       //内容介绍
	UpdateTime time.Time //更新时间
	Tags       []string  //标签
}

type BookChapter struct {
	ChapterName,
	ChapterContent string
	ChapterUpdateTime time.Time
}

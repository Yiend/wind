//整个爬取数据的引擎驱动
//引擎中需要不同的种子,进行爬取操作
//但必须要有一个或者多个种子.如:https://www.ozixue.com
//engine收到种子时是不会马上执行解析的,它需要一个任务队列来维护
package engine

import "github.com/PuerkitoBio/goquery"

//Url and parse,url is a seeds
//ParseFunc is url parse function
type Request struct {
	Url       string
	ParseFunc func(*goquery.Document) ParseResult
}

type ParseResult struct {
	Requests []Request     //需要访问的URL
	Items    []interface{} //访问url得到的数据
}
//TODO 用于test
func NilParser(*goquery.Document) ParseResult {
	return ParseResult{}
}

package parser

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"github.com/wind/engine"
	"github.com/wind/models"
	"strings"
	"time"
)

var (
	urlSort = "http://all.17k.com"
	counter = 2
)

//parse free book list and page down
func ParseFreeBookAll(doc *goquery.Document) engine.ParseResult {
	result := engine.ParseResult{}
	doc.Find(".inner").Each(func(i int, selection *goquery.Selection) {
		bookUrl, _ := selection.Find("li strong a").Attr("href")
		bookName := selection.Find("li strong a").Text()
		result.Items = append(result.Items, bookName)
		result.Requests = append(result.Requests, engine.Request{
			Url:       bookUrl,
			ParseFunc: ParseBookInfo,
		})
	})
	//page down
	result.Requests = append(result.Requests, engine.Request{
		Url:       urlSort + "/lib/book/2_0_0_0_0_0_1_0_" + strconv.Itoa(counter) + ".html",
		ParseFunc: ParseFreeBookAll,
	})
	//TODO This is not the best.
	counter++
	return result
}

//
func ParseBookInfo(doc *goquery.Document) engine.ParseResult {
	result := engine.ParseResult{}
	info := models.BookInfo{}

	bookInfo := doc.Find(".BookInfo")
	bookUrl, _ := bookInfo.Find(".Props .Bar a").Attr("href")
	info.BookName = bookInfo.Find(".Info h1 a").Text()
	info.Author = doc.Find(".bRight .AuthorInfo .author .name").Text()

	info.UpdateTime, _ = strToTime(bookInfo.Find("dt.tit em").Text())

	info.Content = bookInfo.Find("dd div.cont a").Eq(0).Text()
	info.Sort = bookInfo.Find("dd div.cont a").Eq(1).Text()
	size := bookInfo.Find("dd div.cont tr.label td a span").Size()
	for i := 0; i < size; i++ {
		text := bookInfo.Find("dd div.cont tr.label td a span").Eq(i).Text()
		if text != " " {
			info.Tags = append(info.Tags[:], replace(replace(text, " "), "\n"))
		}
	}

	result.Items = append(result.Items, info)
	result.Requests = append(result.Requests, engine.Request{
		Url:       bookUrl,
		ParseFunc: engine.NilParser,
	})
	return result
}

func strToTime(strTime string) (time.Time, error) {
	s := replace(replace(replace(strTime, " "), "\n"), "更新:")
	bytes := []byte(s)
	strTime = string(bytes[0:10]) + " " + string(bytes[10:])
	return time.Parse("2006-01-02 15:04:05", strTime+":00")
}

func replace(s, s2 string) string {
	return strings.Replace(s, s2, "", -1)
}

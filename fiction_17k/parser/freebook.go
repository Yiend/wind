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
		ParseFunc: ParseChapterList,
	})
	return result
}

func ParseChapterList(doc *goquery.Document) engine.ParseResult {
	result := engine.ParseResult{}
	size := doc.Find(".Volume dd a").Size()
	urls := make([]string, 0)
	for i := 0; i < size; i++ {
		val, _ := doc.Find(".Volume dd a").Eq(i).Attr("href")
		urls = append(urls[:], "http://www.17k.com"+val)
	}
	for _, url := range urls {
		result.Requests = append(result.Requests, engine.Request{
			Url:       url,
			ParseFunc: ParseChapterContent,
		})
	}
	return result
}

func ParseChapterContent(doc *goquery.Document) engine.ParseResult {
	result := engine.ParseResult{}
	chapter := models.BookChapter{}
	date := doc.Find(".readAreaBox .chapter_update_time").Text()

	chapter.ChapterUpdateTime, _ = time.Parse("2006-01-02 15:04:05", date)
	ret, _ := doc.Find(".readAreaBox div.p").Html()
	chapter.ChapterName = replace(doc.Find(".readAreaBox h1").Text()," ")

	s := replace(ret, " ")
	chapter.ChapterContent = replace(s, `本书首发来自17K小说网，第一时间看正版内容！<br/><br/>
<divclass="author-say"></div>
<!--二维码广告Start-->
<divclass="qrcode">
<imgsrc="http://img.17k.com/images/ad/qrcode.jpg"alt="wap_17K"width="96"height="118"/>
<ul>
<li>下载17K客户端，《`+ doc.Find(".area div.infoPath span").Prev().Text()+ `》最新章节无广告纯净阅读。</li>
<li>17K客户端专享，签到即送VIP，免费读全站。</li>
</ul>
</div>
<!--二维码广告End-->
<divclass="chapter_text_ad"id="BAIDU_933954"></div>`)
	result.Items = append(result.Items, chapter)
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

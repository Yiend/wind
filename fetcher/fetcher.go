//TODO 基础模型,待改进
//把url中的文本拉取下来进行转码
//转码后的文本是UTF-8
package fetcher

import (
	"net/http"
	"fmt"
	"golang.org/x/text/transform"
	"bufio"
	"golang.org/x/text/encoding"
	"github.com/gpmgo/gopm/modules/log"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/net/html/charset"
	"github.com/PuerkitoBio/goquery"
)

//需要爬取内容的URL
//返回具体的内容
func Fetch(url string) (*goquery.Document, error) {
	//打开一个网页URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//判断头部信息
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}
	//拿到Document中的元素
	return goquery.NewDocumentFromReader(getEncodingRespBody(resp))
}

//获取确认编码后的response body
func getEncodingRespBody(resp *http.Response) *transform.Reader {
	return transform.NewReader(resp.Body,
		determineEncoding(bufio.NewReader(resp.Body)).NewDecoder())
}

//处理编码
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Error("Fetcher determineEncoding method Error: ", err)
		return unicode.UTF8
	}
	//确认编码
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

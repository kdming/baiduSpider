package spider

import (
	"encoding/csv"
	"fmt"
	"net/url"
	url2 "net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 爬取结果
type spiderResult struct {
	title    string // 标题
	abstract string // 简介
	keyword  string // 关键字
	url      string // 路径
	domain   string // 域名
	fullText string // 完整内容
	date     string // 爬取日期
}

// 搜索url
const serverUrl = "http://www.baidu.com/s?ie=UTF-8"

// channel
var c = make(chan int)

// 根据关键字搜索，获取搜索结果
func getSearchResult(keyword string, pageIndex int, searchResult *[]spiderResult) {
	// 搜索关键字获取html内容
	searchUrl := serverUrl + "&rn=10&pn=" + strconv.Itoa((10 * (pageIndex - 1))) + "&wd=" + url.QueryEscape(keyword)
	htmlBody := GetHtmlBody(searchUrl)
	fmt.Println(searchUrl, pageIndex)
	if htmlBody != nil {
		defer htmlBody.Close()
		doc, err := goquery.NewDocumentFromReader(htmlBody)
		if err != nil {
			fmt.Println(keyword, "获取搜索结果失败")
		} else {
			// 解析html内容获取指定字段
			doc.Find(".result.c-container").Each(func(index int, s *goquery.Selection) {
				res := spiderResult{}
				res.keyword = keyword
				res.url = s.Find("h3.t a").AttrOr("href", "")
				u, _ := url2.Parse(res.url)
				res.domain = u.Host
				res.title = s.Find("h3.t a").Text()
				res.abstract = s.Find("div.c-abstract").Text()
				(*searchResult) = append((*searchResult), res)
			})
		}
		c <- pageIndex
	}
}

// 获取搜索结果详情
func getResDetail(searchResult *spiderResult, index int) {
	if &searchResult != nil && searchResult.url != "" {
		htmlBody := GetHtmlBody(searchResult.url)
		if htmlBody != nil {
			defer htmlBody.Close()
			doc, err := goquery.NewDocumentFromReader(htmlBody)
			if err != nil {
				fmt.Println(searchResult.keyword, "获取搜索结果失败")
				searchResult.fullText = ""
			} else {
				searchResult.fullText = strings.ReplaceAll(doc.Find("body").Text(), "\t", "")
			}
		}
	}
	fmt.Println(searchResult.keyword, index, "页爬取完毕")
	c <- index
}

// 保存搜索结果为csv
func saveResult(searchResult []spiderResult, fileName string) error {
	// 创建文件
	pwd, _ := os.Getwd()
	fileUrl := pwd + "/static/result/" + fileName
	f, err := os.Create(fileUrl)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return err
	}
	defer f.Close()
	// 写入UTF-8 BOM，防止中文乱码
	f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	// 写入csv列头
	w.Write([]string{"标题", "简介", "关键字", "url", "域名", "文本内容", "日期"})
	w.Flush()
	// 写入数据详情
	for i := 0; i < len(searchResult); i += 1 {
		res := searchResult[i]
		t := reflect.TypeOf(res)
		v := reflect.ValueOf(res)
		detail := []string{}
		for k := 0; k < t.NumField(); k += 1 {
			var val string
			if v.Field(k).String() == "" {
				val = ","
				continue
			}
			val = strings.ReplaceAll(v.Field(k).String(), "\n", "") + ","
			detail = append(detail, val)
		}
		// 写入日期
		date := time.Now().Format("2006-01-02 15:04:05")
		detail = append(detail, date)
		w.Write(detail)
		w.Flush()
	}
	return nil
}

// 任务调度
func Schedule() {
	jobs := GetJobs()
	for keyword, num := range jobs {
		searchResult := &[]spiderResult{}
		// 多线程搜索，默认搜索12页
		limit := 12
		index := 1
		for {
			for i := index; i < limit; i += 1 {
				go getSearchResult(keyword, i, searchResult)
			}
			for i := index; i < limit; i += 1 {
				<-c
			}
			if len(*searchResult) < num {
				index = limit
				limit += 1
				continue
			} else {
				break
			}
		}
		fmt.Println(len(*searchResult))
		// 多线程获取页面具体内容
		limit = 20
		index = 0
		for {
			for i := index; i < limit; i += 1 {
				result := &(*searchResult)[i]
				if result == nil {
					continue
				}
				go getResDetail(result, i)
			}
			for i := index; i < limit; i += 1 {
				<-c
			}
			if limit < len(*searchResult) {
				index = limit
				if (limit + 20) > len(*searchResult) {
					limit = len(*searchResult)
				} else {
					limit += 20
				}
			} else {
				break
			}
		}
		saveResult(*searchResult, keyword+".csv")
		fmt.Println(keyword, "爬取完毕")
	}
}

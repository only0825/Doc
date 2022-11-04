package main

import (
	"crawler_news/common"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

// http://www.bd54.org.cn
func CrawlerBd54(domain string) {
	// Request the HTML page.
	res, err := http.Get(domain)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	news, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	news.Find(".block-headlines div div a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		if i < 10 {
			url, _ := s.Attr("href")
			fmt.Println(url)
			res, err := http.Get(domain + url)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}
			// Load the HTML document
			detail, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			var article Article
			detail.Find("article").Each(func(j int, sc *goquery.Selection) {
				// 文章标题
				title := sc.Find("h1").Text()
				article.Title = title
				article.Md5 = common.MD5(title)

				// 查找并保存图片到本地
				sc.Find(".news-content p img").Each(func(k int, image *goquery.Selection) {
					imgUrl, _ := image.Attr("src")
					imgPath, _ := common.DownImage(domain+imgUrl, savePath)
					article.ImgPath = imgPath
				})

				// 文章内容
				content, _ := sc.Find(".news-content").Html()
				article.Content = strings.TrimSpace(content)

				// 日期
				sc.Find(".news-info .news-time").Each(func(q int, time *goquery.Selection) {
					article.Date = strings.Trim(time.Text(), " 发布时间：")
				})
				fmt.Println(article.Date)

				// 作者
				sc.Find(".news-info .news-hits").Each(func(q int, time *goquery.Selection) {
					if q == 1 {
						article.Author = time.Text()
					}
				})
			})
			//
			//sql, _, _ := dialect.From(table).Select(goqu.COUNT("*").As("id")).Where(goqu.Ex{
			//	"md5": article.Md5,
			//}).ToSQL()
			//var id int
			//DB.Get(&id, sql)
			//// 如果文章已存在则不写入
			//if id == 0 {
			//	article.CreateTime = time.Now().Unix()
			//	query, _, _ := dialect.Insert(table).Rows(&article).ToSQL()
			//	// 存库
			//	_, err = DB.Exec(query)
			//	if err != nil {
			//		fmt.Println("insert failed query: ", query)
			//		fmt.Println("insert failed error: ", err.Error())
			//	} else {
			//		fmt.Println("insert success !")
			//		common.SaveSql(query)
			//	}
			//} else {
			//	fmt.Printf("《%s》 文章已存在！", article.Title)
			//}
		}
	})
}

package main

import (
	"crawler_news/common"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/doug-martin/goqu/v9"
	"log"
	"net/http"
	"time"
)

// https://www.wutiyu.com
func CrawlerWutiyu(domain string) {
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
	news.Find(".newlist li h3 a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		url, _ := s.Attr("href")

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
		if i < total {

			detail.Find(".detail_main").Each(func(j int, sc *goquery.Selection) {
				title := sc.Find("h3").Text()
				article.Title = title
				article.Md5 = common.MD5(title)
				// 查找并保存图片到本地
				sc.Find("img").Each(func(k int, sc2 *goquery.Selection) {
					imgUrl, _ := sc2.Attr("src")
					imgPath, _ := common.DownImage(imgUrl, savePath)
					article.ImgPath = imgPath
				})

				// 文章内容
				content, _ := sc.Find(".detail_article").Html()
				article.Content = content

				// 作者和日期
				sc.Find(".meta span").Each(func(q int, scMeta *goquery.Selection) {
					if q == 0 {
						author := scMeta.Text()
						//fmt.Println(author)
						article.Author = author
					}
					if q == 1 {
						date := scMeta.Text()
						//fmt.Println(date)
						article.Date = date
					}
				})
			})

			sql, _, _ := dialect.From(table).Select(goqu.COUNT("*").As("id")).Where(goqu.Ex{
				"md5": article.Md5,
			}).ToSQL()
			var id int
			DB.Get(&id, sql)
			// 如果文章已存在则不写入
			if id == 0 {
				article.CreateTime = time.Now().Unix()
				query, _, _ := dialect.Insert(table).Rows(&article).ToSQL()
				// 存库
				_, err = DB.Exec(query)
				if err != nil {
					fmt.Println("insert failed query: ", query)
					fmt.Println("insert failed error: ", err.Error())
				} else {
					fmt.Println("insert success !")
					common.SaveSql(query)
				}
			} else {
				fmt.Printf("《%s》 文章已存在！", article.Title)
			}

		}
	})
}

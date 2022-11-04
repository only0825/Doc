package main

import (
	"crawler_news/common"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/doug-martin/goqu/v9"
	"log"
	"math/rand"
	"strings"
	"time"
)

// https://2022.7m.com.cn/gb/news
func Crawler7m(htmlContent string) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".n_list .new_box").Each(func(i int, s *goquery.Selection) {
		newsId, _ := s.Find("a").Attr("onclick")
		newsId = strings.Trim(newsId, "newsList.toDetail(")
		newsId = strings.Trim(newsId, ")")

		// 将拿到的ID拼上文章详情页地址
		newsContent, _ := common.GetHttpHtmlContent("https://2022.7m.com.cn/gb/news_detail/"+newsId,
			// 第二个参数 selector 即为我们爬取的数据对应的html选择器
			"#js_txt > p",
			//从body里面获取数据
			"document.querySelector(\"body\")")
		detail, err := goquery.NewDocumentFromReader(strings.NewReader(newsContent))
		if err != nil {
			log.Fatal(err)
		}

		var article Article
		detail.Find(".bl_content").Each(func(j int, sc *goquery.Selection) {
			// 文章标题
			title := sc.Find("#js_title").Text()
			article.Title = title
			article.Md5 = common.MD5(title)

			// 查找并保存图片到本地
			imgUrl, _ := sc.Find("#js_txt p img").Attr("src")
			imgPath, _ := common.DownImage("https:"+imgUrl, savePath)
			article.ImgPath = imgPath

			// 文章内容
			content, _ := sc.Find("#js_txt").Html()
			article.Content = strings.TrimSpace(content)

			// 日期
			article.Date = time.Now().Format("2006-01-02 15:04:05")

			// 作者
			authors := []string{"预言家", "锦书美人", "长辞", "毛豆花花", "小米豆", "蓝蓝小瓶"}
			rand.Seed(time.Now().Unix()) //Seed生成的随机数
			index := rand.Intn(6-0) + 0
			article.Author = authors[index]
		})

		var id int
		sql, _, _ := dialect.From(table).Select(goqu.COUNT("*").As("id")).Where(goqu.Ex{
			"md5": article.Md5,
		}).ToSQL()
		DB.Get(&id, sql)
		// 如果文章已存在则不写入
		if id == 0 {
			article.CreateTime = time.Now().Unix()
			query, _, _ := dialect.Insert(table).Rows(&article).ToSQL()
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
	})
}

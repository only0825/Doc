package main

import (
	"crawler_news/common"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	"log"
)

var (
	total    = 30
	DB       *sqlx.DB
	dialect  = goqu.Dialect("mysql")
	table    = "hn_crawler_news"
	savePath = "/news/"
)

func main() {

	db, err := sqlx.Connect("mysql", "u_hainiu:hainiu@2022@(43.135.76.182:3306)/db_hainiu")
	if err != nil {
		log.Fatalln(err)
	}

	DB = db

	//CrawlerWutiyu("https://www.wutiyu.com")
	//CrawlerBd54("http://www.bd54.org.cn")

	htmlContent, _ := common.GetHttpHtmlContent("https://2022.7m.com.cn/gb/news",
		// 第二个参数 selector 即为我们爬取的数据对应的html选择器
		"#js_newsList > div > div > div:nth-child(2) > a",
		//从body里面获取数据
		"document.querySelector(\"body\")")
	Crawler7m(htmlContent)
}

type Article struct {
	Title      string `json:"title" db:"title"`
	Md5        string `json:"md5" db:"md5"`
	ImgPath    string `json:"img_path" db:"img_path"`
	Content    string `json:"content" db:"content"`
	Author     string `json:"author" db:"author"`
	Date       string `json:"date" db:"date"`
	CreateTime int64  `json:"create_time" db:"create_time"`
}

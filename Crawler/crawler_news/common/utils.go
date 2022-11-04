package common

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"time"
)

func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func SaveSql(query string) {
	var str = query + ";\n"
	var date = time.Now().Format("2006-01-02")
	var filename = "./news-" + date + ".sql"
	var f *os.File
	var err1 error
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
	} else {
		f, err1 = os.Create(filename) //创建文件
	}
	defer f.Close()
	//n, err1 := f.Write([]byte(str)) //写入文件(字节数组)
	//
	//fmt.Printf("写入 %d 个字节n", n)
	n, err1 := f.WriteString(str) //写入文件(字符串)
	if err1 != nil {
		panic(err1)
	}
	fmt.Printf("写入 %d 个字节\n", n)
	f.Sync()
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

/**
关于 chromedp 涉及的接口如下介绍几个
名字			说明
Navigate	进入某个页面
Run			运行各类操作
Screenshot	截屏
Click		模拟鼠标点击
WaitVisible	等候某元素出现
ActionFunc	执行自定义函数
SendKeys	模拟键盘输入
*/
// chromedp工具：动态获取网站上的HTML（模拟浏览器访问）
func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//创建一个上下文，超时时间为40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		return "", err
	}
	//log.Println(htmlContent)

	return htmlContent, nil
}

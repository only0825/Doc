package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	_ "go.uber.org/automaxprocs"
	"log"
	"match/middleware"
	"match/model"
	"match/router"
	"os"
	"strings"
)

var (
	gitReversion   = ""
	buildTime      = ""
	buildGoVersion = ""
)

func main() {

	//argc := len(os.Args)
	//if argc != 3 {
	//	fmt.Printf("%s <etcds> <cfgPath>\r\n", os.Args[0])
	//	return
	//}
	//
	//cfg := conf{}
	//endpoints := strings.Split(os.Args[1], ",")
	//
	//apollo.New(endpoints)
	//apollo.Parse(os.Args[2], &cfg)
	//apollo.Close()

	mt := new(model.MetaTable)
	// 连接数据库
	db, err := sqlx.Connect("mysql", "root:password@(localhost:3306)/hainiu")
	if err != nil {
		log.Fatalln(err)
	}

	// 将DB的连接实例赋值给MatchDB
	mt.MatchDB = db

	bin := strings.Split(os.Args[0], "/")
	mt.Program = bin[len(bin)-1]

	model.Constructor(mt)
	//session.New(mt.MerchantRedis, mt.Prefix)

	defer func() {
		model.Close()
		mt = nil
	}()

	b := router.BuildInfo{
		GitReversion:   gitReversion,
		BuildTime:      buildTime,
		BuildGoVersion: buildGoVersion,
	}
	app := router.SetupRouter(b)
	srv := &fasthttp.Server{
		Handler:            middleware.Use(app.Handler),
		ReadTimeout:        router.ApiTimeout,
		WriteTimeout:       router.ApiTimeout,
		Name:               "matchApi",
		MaxRequestBodySize: 51 * 1024 * 1024,
	}
	fmt.Printf("gitReversion = %s\r\nbuildGoVersion = %s\r\nbuildTime = %s\r\n", gitReversion, buildGoVersion, buildTime)
	fmt.Println("matchApi running", 6031)
	if err := srv.ListenAndServe(":6031"); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	_ "go.uber.org/automaxprocs"
	"goblog/middleware"
	"goblog/model"
	"goblog/router"
	"log"
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

	db, err := sqlx.Connect("mysql", "root:123456@(localhost:3306)/blog_db")
	if err != nil {
		log.Fatalln(err)
	}
	mt.MerchantDB = db

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
		Name:               "logApi",
		MaxRequestBodySize: 51 * 1024 * 1024,
	}
	fmt.Printf("gitReversion = %s\r\nbuildGoVersion = %s\r\nbuildTime = %s\r\n", gitReversion, buildGoVersion, buildTime)
	fmt.Println("logApi running", 6031)
	if err := srv.ListenAndServe(":6031"); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

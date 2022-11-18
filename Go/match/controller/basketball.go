package controller

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"time"
)

type Basketball struct{}

// Schedule 篮球赛事
func (that *Basketball) Schedule(ctx *fasthttp.RequestCtx) {

	leagueId := string(ctx.PostArgs().Peek("leagueId"))
	matchId := string(ctx.PostArgs().Peek("matchId"))
	date := string(ctx.PostArgs().Peek("date"))
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	url := "http://interface.bet007.com"
	key := "E6DDD3614E584010"
	api := "/basketball/schedule.aspx"

	apiUrl := url + api + "?date=" + date + "&key=" + key
	//if leagueId != '' {
	//	$api_url = $url.$api.'?leagueId='.$leagueId;
	//}
	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	//fmt.Println(res)
	fmt.Println(leagueId)
	fmt.Println(matchId)
	//
	//if username == "" || password == "" {
	//	helper.Print(ctx, false, helper.ParamErr)
	//}
	//
	//isSet := model.UserExist(username)
	//if isSet != 0 {
	//	helper.Print(ctx, false, helper.UserExist)
	//}
	//
	//uid, err := model.UserReg(username, password)
	//if err != nil {
	//	helper.Print(ctx, false, err.Error())
	//}
	//
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	fmt.Println(res.StatusCode)
	if res.StatusCode == 200 {
		fmt.Println("ok")
	}
	//helper.Print(ctx, true, res)
}

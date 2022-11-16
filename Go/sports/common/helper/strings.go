package helper

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type Response struct {
	State bool        `json:"state"`
	Data  interface{} `json:"data"`
}

func Print(ctx *fasthttp.RequestCtx, state bool, data interface{}) {

	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")

	res := Response{
		State: state,
		Data:  data,
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		ctx.SetBody([]byte(err.Error()))
	}

	ctx.SetBody(bytes)
}

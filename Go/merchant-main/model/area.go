package model

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"time"
)

func areaIp(client *fasthttp.PipelineClient, uri string) (string, error) {

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetHost("ipapi.co")
	req.SetRequestURI(uri)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	err := client.DoTimeout(req, resp, 5*time.Second)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return "", fmt.Errorf("%d", resp.StatusCode())
	}

	//fmt.Println("Hello, 世界 = ", string(resp.Body()))
	return string(resp.Body()), nil
}

func Area(ips []string) (string, error) {

	uniq := map[string]bool{}
	client := fasthttp.PipelineClient{
		Addr:          "ipapi.co",
		MaxBatchDelay: 1 * time.Second,
		MaxConns:      20,
		IsTLS:         true,
	}

	arr := new(fastjson.Arena)
	aa := arr.NewObject()

	for _, val := range ips {
		if _, ok := uniq[val]; ok {
			continue
		}

		uniq[val] = true
		uri := fmt.Sprintf("https://ipapi.co/%s/json/", val)
		data, err := areaIp(&client, uri)
		if err == nil {
			v, err := fastjson.Parse(data)
			if err == nil {
				aa.Set(val, v)
			}
		}
	}

	buf := aa.String()
	arr = nil

	return buf, nil
}

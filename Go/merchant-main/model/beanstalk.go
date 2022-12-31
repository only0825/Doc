package model

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/valyala/fasthttp"
)

func BeanPut(name string, param map[string]interface{}) error {

	m := &fasthttp.Args{}
	for k, v := range param {
		if _, ok := v.(string); ok {
			m.Set(k, v.(string))
		}
	}

	topic := meta.Prefix + "_" + name + "s"
	err := meta.MerchantMQ.SendAsync(ctx,
		func(c context.Context, result *primitive.SendResult, e error) {
			if e != nil {
				fmt.Printf("receive message error: %s\n", e.Error())
			} else {
				fmt.Printf("send message success: result=%s\n", result.String())
			}
		}, primitive.NewMessage(topic, m.QueryString()))

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	}

	fmt.Println("beanPut topic = ", topic)
	fmt.Println("beanPut param = ", m.String())
	fmt.Println("beanPut err = ", err)

	return err
}

func TransAg(name string, param map[string]interface{}) error {

	m := &fasthttp.Args{}
	for k, v := range param {
		if _, ok := v.(string); ok {
			m.Set(k, v.(string))
		}
	}

	topic := name
	err := meta.MerchantMQ.SendAsync(ctx,
		func(c context.Context, result *primitive.SendResult, e error) {
			if e != nil {
				fmt.Printf("receive TransAg error: %s\n", e.Error())
			} else {
				fmt.Printf("send message success: result=%s\n", result.String())
			}
		}, primitive.NewMessage(topic, m.QueryString()))

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	}

	fmt.Println("TransAg topic = ", topic)
	fmt.Println("TransAg param = ", m.String())
	fmt.Println("TransAg err = ", err)

	return err
}

func BeanPutDelay(name string, param map[string]interface{}, delay int) error {

	//fmt.Printf("BeanPut topic: %s, param : %#v, delay : %d\n", name, param, delay)
	m := &fasthttp.Args{}
	for k, v := range param {
		if _, ok := v.(string); ok {
			m.Set(k, v.(string))
		}
	}

	topic := meta.Prefix + "_" + name
	tube := &beanstalk.Tube{Conn: meta.MerchantBean, Name: topic}
	_, err := tube.Put(m.QueryString(), 1, time.Duration(delay)*time.Second, 10*time.Minute)
	if err != nil {
		return err
	}

	return nil
}

/*
//BeanBetPut 注单
func BeanBetPut(name string, param map[string]interface{}, delay int) (string, error) {

	m := &fasthttp.Args{}
	for k, v := range param {
		if _, ok := v.(string); ok {
			m.Set(k, v.(string))
		}
	}

	c, err := meta.BeanBetPool.Get()
	if err != nil {
		return "sys", err
	}

	if conn, ok := c.(*beanstalk.Conn); ok {

		tube := &beanstalk.Tube{Conn: conn, Name: name}
		_, err = tube.Put(m.QueryString(), 1, time.Duration(delay)*time.Second, 10*time.Minute)
		if err != nil {
			_ = meta.BeanPool.Put(c)
			return "sys", err
		}
	}

	//将连接放回连接池中
	return "", meta.BeanPool.Put(c)
}
*/

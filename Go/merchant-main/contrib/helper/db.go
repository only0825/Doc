package helper

import "reflect"

/**
 * @Description: 根据model结构获取db查询字段
 * @Author: sar
 * @Date: 2021/04/02
 * @LastEditTime: 2021/04/02
 * @LastEditors: sar
 **/
func EnumFields(obj interface{}) []interface{} {

	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Struct {
		return nil
	}

	var fields []interface{}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if field := f.Tag.Get("db"); field != "" && field != "-" {
			fields = append(fields, field)
		}
	}

	return fields
}

func EnumRedisFields(obj interface{}) []string {

	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Struct {
		return nil
	}

	var fields []string
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if field := f.Tag.Get("db"); field != "" && field != "-" {
			fields = append(fields, field)
		}
	}

	return fields
}
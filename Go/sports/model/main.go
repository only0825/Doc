package model

import (
	g "github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type MetaTable struct {
	MerchantDB *sqlx.DB
	Program    string
}

var (
	loc     *time.Location
	meta    *MetaTable
	dialect = g.Dialect("mysql")
)

func Constructor(mt *MetaTable) {

	meta = mt
	loc, _ = time.LoadLocation("Asia/Shanghai")
}

//func pushLog(err error, code string) error {
//
//	_, file, line, _ := runtime.Caller(1)
//	paths := strings.Split(file, "/")
//	l := len(paths)
//	if l > 2 {
//		file = paths[l-2] + "/" + paths[l-1]
//	}
//	path := fmt.Sprintf("%s:%d", file, line)
//
//	ts := time.Now()
//	id := helper.GenId()
//
//	fields := g.Record{
//		"id":       id,
//		"content":  err.Error(),
//		"project":  meta.Program,
//		"flags":    code,
//		"filename": path,
//		"ts":       ts.In(loc).UnixMilli(),
//	}
//	query, _, _ := dialect.Insert("goerror").Rows(&fields).ToSQL()
//	//fmt.Println(query)
//	_, err1 := meta.MerchantTD.Exec(query)
//	if err1 != nil {
//		fmt.Println("insert SMS = ", err1.Error(), fields)
//	}
//
//	note := fmt.Sprintf("Hệ thống lỗi %s", id)
//	return errors.New(note)
//}

func Close() {
}

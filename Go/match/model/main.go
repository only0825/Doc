package model

import (
	g "github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type MetaTable struct {
	MatchDB *sqlx.DB
	Program string
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

func Close() {
}

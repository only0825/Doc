package model

import (
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

var DB *gorm.DB

// var Ctx = context.Background()
var Rdbc *redis.ClusterClient
var rdb *redis.Client

package model

import (
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Rdb *redis.ClusterClient

//var Rdb *redis.Client

package configs

import (
	"github.com/spf13/viper"
	"os"
)

type Mysql struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Db       string `yaml:"db"`
	Charset  string `yaml:"charset"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Timeout  int    `yaml:"timeout"`
	Select   int    `yaml:"select"`
}

type Server struct {
	Name           string
	Scheme         string
	Ip             string
	Port           int
	HeartBeatTime  int
	InChanSize     int
	OutChanSize    int
	WorkerPoolSize uint64
	MaxWorkTaskLen int
	MaxConn        int
}

type AppConfig struct {
	Mysql        Mysql
	Redis        Redis
	RedisCluster map[int]string
	Cache        string
	Server       Server
	WsKey        string
}

var Conf *AppConfig

func LoadConfig(fileName string) error {
	wd, err := os.Getwd() // 获取当前文件路径
	if err != nil {
		return err
	}

	c := &AppConfig{}
	v := viper.New()
	v.SetConfigName(fileName) //这里就是上面我们配置的文件名称，不需要带后缀名
	v.AddConfigPath(wd)       //文件所在的目录路径
	v.SetConfigType("yml")    //这里是文件格式类型

	err = v.ReadInConfig()
	if err != nil {
		return err
	}

	configs := v.AllSettings()
	for k, val := range configs {
		v.SetDefault(k, val)
	}

	err = v.Unmarshal(c) //反序列化至结构体
	if err != nil {
		return err
	}

	Conf = c
	return nil
}

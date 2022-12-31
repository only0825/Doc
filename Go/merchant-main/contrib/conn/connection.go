package conn

import (
	"context"
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/meilisearch/meilisearch-go"

	//mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/olivere/elastic/v7"
	"github.com/panjf2000/ants/v2"
	cpool "github.com/silenceper/pool"
	_ "github.com/taosdata/driver-go/v2/taosRestful"
	"log"
	//"strings"
	"time"
)

var ctx = context.Background()

func InitTD(dsn string, maxIdleConn, maxOpenConn int) *sqlx.DB {

	db, err := sqlx.Connect("taosRestful", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Second * 30)
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func InitDB(dsn string, maxIdleConn, maxOpenConn int) *sqlx.DB {

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Second * 30)
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func InitRedisSentinel(dsn []string, psd, name string, db int) *redis.Client {

	reddb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    name,
		SentinelAddrs: dsn,
		Password:      psd, // no password set
		DB:            db,  // use default DB
		DialTimeout:   10 * time.Second,
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		PoolSize:      500,
		PoolTimeout:   30 * time.Second,
		MaxRetries:    2,
		IdleTimeout:   5 * time.Minute,
	})
	pong, err := reddb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("InitRedisSentinel failed: %s", err.Error())
	}
	fmt.Println(pong, err)

	return reddb
}

func InitRedis(dsn string, psd string, db int) *redis.Client {

	reddb := redis.NewClient(&redis.Options{
		Addr:         dsn,
		Password:     psd, // no password set
		DB:           db,  // use default DB
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     500,
		PoolTimeout:  30 * time.Second,
		MaxRetries:   2,
		IdleTimeout:  5 * time.Minute,
	})
	pong, err := reddb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("InitRedis failed: %s", err.Error())
	}
	fmt.Println(pong, err)

	return reddb
}

func InitRedisSentinelRead(dsn []string, psd, name string, db int) *redis.ClusterClient {

	reddb := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:     name,
		SentinelAddrs:  dsn,
		Password:       psd, // no password set
		DB:             db,  // use default DB
		DialTimeout:    10 * time.Second,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		PoolSize:       500,
		PoolTimeout:    30 * time.Second,
		MaxRetries:     2,
		IdleTimeout:    5 * time.Minute,
		SlaveOnly:      true,
		RouteByLatency: true,
	})
	pong, err := reddb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("InitRedisSentinelRead failed: %s", err.Error())
	}
	fmt.Println(pong, err)

	return reddb
}

func InitRedisCluster(dsn []string, psd string) *redis.ClusterClient {

	reddb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        dsn,
		Password:     psd, // no password set
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     500,
		PoolTimeout:  30 * time.Second,
		MaxRetries:   2,
		IdleTimeout:  5 * time.Minute,
	})

	pong, err := reddb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("initRedisSlave failed: %s", err.Error())
	}
	fmt.Println(pong, err)

	return reddb
}

func InitMeiliSearch(url, key string) *meilisearch.Client {

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   url,
		APIKey: key,
	})
	_, err := client.GetKeys()
	if err != nil {
		log.Fatalf("InitMeiliSearch failed: %s", err.Error())
	}

	return client
}

func InitES(url []string, username, password string) *elastic.Client {

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(url...),
		elastic.SetBasicAuth(username, password))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func InitBeanstalk(beanstalkConn string, initialCap, maxIdle, maxCap int) cpool.Pool {

	factory := func() (interface{}, error) { return beanstalk.Dial("tcp", beanstalkConn) }
	closed := func(v interface{}) error { return v.(*beanstalk.Conn).Close() }
	poolConfig := &cpool.Config{
		InitialCap:  initialCap, // 资源池初始连接数
		MaxIdle:     maxIdle,    // 最大空闲连接数
		MaxCap:      maxCap,     // 最大并发连接数
		Factory:     factory,
		Close:       closed,
		IdleTimeout: 15 * time.Second,
	}

	beanPool, err := cpool.NewChannelPool(poolConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return beanPool
}

func InitFluentd(addr string, port int) *fluent.Fluent {

	c := fluent.Config{
		FluentPort:   port,
		FluentHost:   addr,
		Async:        true,
		MaxRetry:     3,
		AsyncConnect: false,
	}

	zlog, err := fluent.New(c)
	if err != nil {
		log.Fatalln(err)
	}

	return zlog
}

func InitRoutinePool() *ants.Pool {
	// TODO 现在先写默认值，后续优化为根据配置文件动态更新数量
	pool, err := ants.NewPool(5000)
	if err != nil {
		log.Fatalln(err)
	}
	return pool
}

// 创建nats.io链接
func InitNatsIO(urls []string, name, password string) *nats.Conn {

	opts := nats.Options{
		Servers:        urls,
		User:           name,
		Password:       password,
		AllowReconnect: true,
		MaxReconnect:   10,
		PingInterval:   5 * time.Second,
		ReconnectWait:  5 * time.Second,
		Timeout:        5 * time.Second,
	}

	nc, err := opts.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	/*
		nc, err := nats.Connect("nats://10.170.0.9:4242")
		if err != nil {
			log.Fatalln(err)
		}
	*/

	return nc
}

func InitIpDB(path string) *ip2location.DB {

	db, err := ip2location.OpenDB(path)
	if err != nil {
		log.Fatalf("initIPBin failed: %s", err.Error())
	}

	return db
}

/*
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

// 连接mqtt
func InitMqttService(addrs []string, clientID, username, password string) mqtt.Client {

	clientOptions := mqtt.NewClientOptions().
		SetClientID(clientID).
		SetUsername(username).
		SetPassword(password).
		SetCleanSession(false).
		SetAutoReconnect(true).
		SetKeepAlive(120 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetWriteTimeout(5 * time.Second).
		SetMaxReconnectInterval(10 * time.Second)

	for _, v := range addrs {
		clientOptions.AddBroker(v)
	}

	clientOptions.OnConnect = connectHandler
	clientOptions.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(clientOptions)
	if conn := client.Connect(); conn.WaitTimeout(time.Duration(10)*time.Second) && conn.Wait() && conn.Error() != nil {
		log.Fatalf("token: %s-%s", strings.Join(addrs, ","), conn.Error())
	}
	return client
}
*/

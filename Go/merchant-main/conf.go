package main

type conf struct {
	Lang           string `json:"lang"`
	Prefix         string `json:"prefix"`
	EsPrefix       string `json:"es_prefix"`
	PullPrefix     string `json:"pull_prefix"`
	IsDev          bool   `json:"is_dev"`
	IndexUrl       string `json:"index_url"`
	GcsDoamin      string `json:"gcs_doamin"`
	AutoCommission bool   `json:"auto_commission"`
	Sock5          string `json:"sock5"`
	RPC            string `json:"rpc"`
	Fcallback      string `json:"fcallback"`
	AutoPayLimit   string `json:"autoPayLimit"`
	Nats           struct {
		Servers  []string `json:"servers"`
		Username string   `json:"username"`
		Password string   `json:"password"`
	} `json:"nats"`
	Rocketmq   []string `json:"rocketmq"`
	Beanstalkd struct {
		Addr    string `json:"addr"`
		MaxIdle int    `json:"maxIdle"`
		MaxCap  int    `json:"maxCap"`
	} `json:"beanstalkd"`
	Db struct {
		Master struct {
			Addr        string `json:"addr"`
			MaxIdleConn int    `json:"max_idle_conn"`
			MaxOpenConn int    `json:"max_open_conn"`
		} `json:"master"`
		Report struct {
			Addr        string `json:"addr"`
			MaxIdleConn int    `json:"max_idle_conn"`
			MaxOpenConn int    `json:"max_open_conn"`
		} `json:"report"`
		Bet struct {
			Addr        string `json:"addr"`
			MaxIdleConn int    `json:"max_idle_conn"`
			MaxOpenConn int    `json:"max_open_conn"`
		} `json:"bet"`
		Tidb struct {
			Addr        string `json:"addr"`
			MaxIdleConn int    `json:"max_idle_conn"`
			MaxOpenConn int    `json:"max_open_conn"`
		} `json:"tidb"`
	} `json:"db"`
	Td struct {
		Log struct {
			Addr        string `json:"addr"`
			MaxIdleConn int    `json:"max_idle_conn"`
			MaxOpenConn int    `json:"max_open_conn"`
		} `json:"log"`
		Message struct {
			Addr        string `json:"addr"`
			MaxIdleConn int    `json:"max_idle_conn"`
			MaxOpenConn int    `json:"max_open_conn"`
		} `json:"message"`
	} `json:"td"`
	Redis struct {
		Addr     []string `json:"addr"`
		Password string   `json:"password"`
	} `json:"redis"`
	Pika struct {
		Addr     []string `json:"addr"`
		Password string   `json:"password"`
		Sentinel string   `json:"sentinel"`
		Db       int      `json:"db"`
	} `json:"pika"`
	Es struct {
		Host     []string `json:"host"`
		Username string   `json:"username"`
		Password string   `json:"password"`
	} `json:"es"`
	AccessEs struct {
		Host     []string `json:"host"`
		Username string   `json:"username"`
		Password string   `json:"password"`
	} `json:"access_es"`
	Port struct {
		Game     string `json:"game"`
		Member   string `json:"member"`
		Promo    string `json:"promo"`
		Merchant string `json:"merchant"`
		Finance  string `json:"finance"`
	} `json:"port"`
}

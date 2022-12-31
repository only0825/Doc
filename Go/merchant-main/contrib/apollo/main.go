package apollo

import (
	"fmt"
	"time"

	//"errors"
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/goccy/go-json"
	"github.com/pelletier/go-toml"
)

var (
	conn           *clientv3.Client
	dialTimeout    = 5 * time.Second
	requestTimeout = 8 * time.Second
)

func New(endpoints []string) {

	var err error
	conn, err = clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   endpoints,
	})

	if err != nil {
		panic(err)
	}
}

func Close() {
	conn.Close()
}

func isDigit(s string) bool {

	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func Parse(key string, v interface{}) error {

	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	kv := clientv3.NewKV(conn)
	gr, _ := kv.Get(ctx, key)
	if gr == nil || len(gr.Kvs) == 0 {
		return fmt.Errorf("No more '%s'", key)
	}

	return json.Unmarshal(gr.Kvs[0].Value, v)
}

func ParseToml(key string, filter bool) (map[string]map[string]interface{}, error) {

	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	kv := clientv3.NewKV(conn)
	gr, _ := kv.Get(ctx, key)
	if gr == nil || len(gr.Kvs) == 0 {
		return nil, fmt.Errorf("No more '%s'", key)
	}

	recs := map[string]map[string]interface{}{}
	config, err := toml.LoadBytes(gr.Kvs[0].Value)
	if err != nil {
		return recs, err
	}

	keys := config.Keys()
	for _, val := range keys {

		if filter && isDigit(val) {

			tree := config.Get(val).(*toml.Tree)
			recs[val] = tree.ToMap()
		} else {
			tree := config.Get(val).(*toml.Tree)
			recs[val] = tree.ToMap()
		}
	}

	return recs, nil
}


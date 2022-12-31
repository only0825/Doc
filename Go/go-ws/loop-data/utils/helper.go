package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
)

func Md5String(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		log.Fatal(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}

/*
Redis 操作：
   err := rdb.Set(ctx, "key", "value", 0).Err()
   if err != nil {
       panic(err)
   }

	// 过期时间
	err := rdbc.Set(ctx, "key111", "value111", time.Duration(20)*time.Second).Err()
	if err != nil {
		panic(err)
	}

   val, err := rdb.Get(ctx, "key").Result()
   if err != nil {
       panic(err)
   }
   fmt("key", val)

   val2, err := rdb.Get(ctx, "key2").Result()
   if err == redis.Nil {
       fmt("key2 does not exist")
   } else if err != nil {
       panic(err)
   } else {
       fmt("key2", val2)
   }
   // Output: key value
   // key2 does not exist


		// 将json消息转为数组
		//err = json.Unmarshal(msg, &m)
		//if err != nil {
		//	fmt("err :", err)
		//	break
		//}
*/

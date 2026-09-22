package widget

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
	"encoding/json"
	"fmt"
)



type RedisCach struct {
	Rdb *redis.Client
}



type RedisError struct {
	msg string
}

func (err RedisError) Error() string {
	return err.msg
}





func NewRedisCach(addr string , password string) RedisCach{
	rc := RedisCach{
		Rdb: redis.NewClient(&redis.Options{
			Addr: addr,
			Password: password,
			DB: 0,
			Protocol: 2,
		}),
	}
	return rc
}

func (rc RedisCach) Set (ctx context.Context , key string , value interface{} , exp_time time.Duration) error {
	to , cancel := context.WithTimeout(ctx , time.Second)

	fmt.Println(key)
	defer cancel()

	vj , err := json.Marshal(value)
	if err != nil{
		return RedisError{msg:"unable to marshal data"}
	}

	return rc.Rdb.Set(to , key , vj , exp_time).Err()


}

func (rc RedisCach) Get (ctx context.Context , key string , out any) (error) {
	timeout , cancel := context.WithTimeout(ctx , time.Second)
	fmt.Println(key)
	defer cancel()

	val , err := rc.Rdb.Get(timeout , key).Result()
	if err == redis.Nil {
		return RedisError{msg : "key dosnt exist or expired"}
	} else if err != nil{
		return RedisError{msg : "unexpected error"}
	}

	err = json.Unmarshal([]byte(val) , out)
	if err != nil {
		return RedisError{"unable to unmarshal"}
	}

	return nil
}
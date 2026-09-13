package widget

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
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

func (rc RedisCach) Set (ctx context.Context , key , value string , exp_time int) error {
	timeout , cancel := context.WithTimeout(ctx , time.Second)

	defer cancel()

	return rc.Rdb.Set(timeout , key , value , time.Duration(exp_time)).Err()

}

func (rc RedisCach) Get (ctx context.Context , key string ) (string , error) {
	timeout , cancel := context.WithTimeout(ctx , time.Second)

	defer cancel()

	val , err := rc.Rdb.Get(timeout , key).Result()
	if err != redis.Nil {
		return "" , RedisError{msg : "key dosnt exist or expired"}
	} else if err != nil{
		return "" , RedisError{msg : "unexpected error"}
	}
	return val , nil
}
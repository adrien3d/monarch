package store

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	Pool *redis.Pool
}

func GetRedis() Redis {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	return Redis{pool}
}

// Gets the object from the redis store, if not found, returns an err
func (r *Redis) GetValueForKey(key string, result interface{}) error {
	redisConn := r.Pool.Get()
	defer redisConn.Close()

	data, err := redisConn.Do("GET", key)
	if err != nil {
		return err
	}

	if data == nil {
		return errors.New("The given key was not found in the redis store")
	}

	// Unmarshall the json
	if err := json.Unmarshal(data.([]byte), &result); err != nil {
		return err
	}

	return nil
}

func (r *Redis) SetValueForKey(key string, value interface{}) error {
	redisConn := r.Pool.Get()
	defer redisConn.Close()

	// Marshals the value to JSON format
	marshaledData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Sets the value to the redis keystore
	if _, err := redisConn.Do("SET", key, marshaledData); err != nil {
		return err
	}

	// Sets the expiration
	/*if _, err := redisConn.Do("EXPIRE", key, r.Config.GetInt("redis_cache_expiration")*60); err != nil {
		return err
	}*/

	return nil
}

// InvalidateObject an object in the redis store.
func (r *Redis) InvalidateObject(id string) error {
	redisConn := r.Pool.Get()
	defer redisConn.Close()

	if _, err := redisConn.Do("DEL", id); err != nil {
		return err
	}

	return nil
}

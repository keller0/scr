package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func Ping() error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func Set(key string, value []byte) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func Exists(key string) (bool, error) {

	conn := Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func Delete(key string) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func GetKeys(pattern string) ([]string, error) {

	conn := Pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func Incr(counterKey string) (int, error) {

	conn := Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", counterKey))
}

func IsValidApiKey(key string) bool {
	conn := Pool.Get()
	defer conn.Close()

	r, _ := conn.Do("HGET", apihash, key)
	return r != nil
}

func IncrApiKey(key string) (int, error) {
	conn := Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("HINCRBY", apihash, key, 1))
}

func LimitApiKeyDura(key, dur string) (int, error) {
	conn := Pool.Get()
	defer conn.Close()
	tmpkey := key + "tmp" + dur
	// if tmpkey not exist redis will set it to 1
	ok, err := redis.Bool(conn.Do("EXISTS", tmpkey))
	if err != nil {
		return 0, fmt.Errorf("error checking if key %s exists: %v", tmpkey, err)
	}
	if ok != true {
		// set api tmp count
		_, err := conn.Do("SET", tmpkey, 1, "ex", dur)
		if err != nil {
			return 0, fmt.Errorf("error set %s exists: %v", tmpkey, err)
		}
		return 0, nil
	}

	return redis.Int(conn.Do("INCR", tmpkey))
}

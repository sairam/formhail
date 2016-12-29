package main

import (
	"fmt"
	"log"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	"github.com/garyburd/redigo/redis"
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 1000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			// redis.DialURL("redis://:secrets@example.com:1234/9?foo=bar&baz=qux")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

var redisPool = newRedisPool()

type redisDB struct {
	// conn redis.Conn
}

func encodeStruct(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

func decodeStruct(b []byte, i interface{}) error {
	return msgpack.Unmarshal(b, i)
}

func (db *redisDB) load(key string, id string, i interface{}) bool {
	c := redisPool.Get()
	defer c.Close()

	data, err := c.Do("GET", fmt.Sprintf("%s:%s", key, id))
	if data == nil || err != nil {
		log.Println("error getting key", key, id)
		return false
	}

	err = decodeStruct(data.([]byte), i)
	if err != nil {
		log.Println("error", err)
		return false
	}
	return true
}

func (db *redisDB) save(key string, id string, i interface{}) bool {
	c := redisPool.Get()
	defer c.Close()

	d, err := encodeStruct(i)
	if err != nil {
		return false
	}
	_, err = c.Do("SET", fmt.Sprintf("%s:%s", key, id), d)
	if err != nil {
		return false
	}

	return true
}

func (db *redisDB) autoincr(key string) int64 {
	c := redisPool.Get()
	defer c.Close()

	out, err := c.Do("INCR", fmt.Sprintf("autoincr:%s", key))

	if err != nil {
		return 0
	}

	return out.(int64)
}

func (db *redisDB) setKeyValue(key, value string) {
	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("SET", key, value)
	if err != nil {
		log.Printf("error setting value %s | %s", key, value)
	}
}

func (db *redisDB) getKeyValue(key string) string {
	c := redisPool.Get()
	defer c.Close()

	out, err := c.Do("GET", key)
	if out == nil || err != nil {
		return ""
	}
	return string(out.([]byte))
}

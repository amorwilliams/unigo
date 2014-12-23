package redis

import (
	"testing"
	"time"
)

import (
	"github.com/amorwilliams/unigo/cache"
	"github.com/garyburd/redigo/redis"
)

func TestRedisCache(t *testing.T) {
	bm, err := cache.NewCache("redis", `{"conn": "127.0.0.1:6379"}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("amor", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("amor") {
		t.Error("check err")
	}

	time.Sleep(11 * time.Second)

	if bm.IsExist("amor") {
		t.Error("check err")
	}
	if err = bm.Put("amor", 1, 10); err != nil {
		t.Error("set Error", err)
	}

	if v, _ := redis.Int(bm.Get("amor"), err); v != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("amor"); err != nil {
		t.Error("Incr Error", err)
	}

	if v, _ := redis.Int(bm.Get("amor"), err); v != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("amor"); err != nil {
		t.Error("Decr Error", err)
	}

	if v, _ := redis.Int(bm.Get("amor"), err); v != 1 {
		t.Error("get err")
	}
	bm.Delete("amor")
	if bm.IsExist("amor") {
		t.Error("delete err")
	}
	//test string
	if err = bm.Put("amor", "author", 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("amor") {
		t.Error("check err")
	}

	if v, _ := redis.String(bm.Get("amor"), err); v != "author" {
		t.Error("get err")
	}
	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}
}

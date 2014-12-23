package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("amor", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("amor") {
		t.Error("check err")
	}

	if v := bm.Get("amor"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("amor") {
		t.Error("check err")
	}

	if err = bm.Put("amor", 1, 10); err != nil {
		t.Error("set Error", err)
	}

	if err = bm.Incr("amor"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("amor"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("amor"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("amor"); v.(int) != 1 {
		t.Error("get err")
	}
	bm.Delete("amor")
	if bm.IsExist("amor") {
		t.Error("delete err")
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"/cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("amor", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("amor") {
		t.Error("check err")
	}

	if v := bm.Get("amor"); v.(int) != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("amor"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("amor"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("amor"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("amor"); v.(int) != 1 {
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

	if v := bm.Get("amor"); v.(string) != "author" {
		t.Error("get err")
	}
}

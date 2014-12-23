package cache

import (
	"fmt"
)

type Cache interface {
	// get cached value by key.
	Get(key string) interface{}
	// set cached value with key and expire time.
	Put(key string, val interface{}, timeout int64) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// clear all cache.
	ClearAll() error
	// start gc routine based on config string settings.
	StartAndGC(config string) error
}

var adapters = make(map[string]Cache)

func Register(name string, adapter Cache) {
	if adapter == nil {
		panic("cache: Register adapter is nill")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func NewCache(adapterName, config string) (adapter Cache, err error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknow adapter name %q (forgot to import?)", adapterName)
		return
	}
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}

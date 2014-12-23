package config

import (
	"fmt"
)

// ConfigContainer defines how to get and set value from configuration raw data.
type ConfigContainer interface {
	Set(key, val string) error   // support section::key type in given key when using ini type.
	String(key string) string    // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	Strings(key string) []string //get string slice
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultval string) string      // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultStrings(key string, defaultval []string) []string //get string slice
	DefaultInt(key string, defaultval int) int
	DefaultInt64(key string, defaultval int64) int64
	DefaultBool(key string, defaultval bool) bool
	DefaultFloat(key string, defaultval float64) float64
	DIY(key string) (interface{}, error)
	GetSection(section string) (map[string]string, error)
	SaveConfigFile(filename string) error
}

type Config interface {
	Parse(key string) (ConfigContainer, error)
	ParseData(data []byte) (ConfigContainer, error)
}

var adapters = make(map[string]Config)

// Register makes a config adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Config) {
	if adapter == nil {
		panic("config: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("config: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// adapterName is ini/json/xml/yaml.
// filename is the config file path.
func NewConfig(adapterName, fileaname string) (ConfigContainer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
	}
	return adapter.Parse(fileaname)
}

// adapterName is ini/json/xml/yaml.
// data is the config data.
func NewConfigData(adapterName string, data []byte) (ConfigContainer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
	}
	return adapter.ParseData(data)
}

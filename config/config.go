package config

import (
	"context"
	"fmt"
	"io"

	"github.com/goccy/go-json"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

type Config struct {
	v gjson.Result
}

func Load(ctx context.Context, r io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(r)
	var d map[string]any
	if err := decoder.Decode(&d); nil != err {
		return nil, fmt.Errorf("config: failed to decode config file yaml content: %v", err)
	}
	b, err := json.MarshalContext(ctx, d, json.UnorderedMap())
	if nil != err {
		return nil, fmt.Errorf("config: failed to prepare configuration file content for internal usage: %v", err)
	}
	return &Config{v: gjson.ParseBytes(b)}, nil
}

func (c *Config) Int(path string) int {
	res := c.v.Get(path)
	if !res.Exists() {
		panic(fmt.Sprintf("path %s does not exist in config", path))
	}
	if t := res.Type; t == gjson.Number {
		return int(res.Num)
	} else {
		panic(fmt.Sprintf("path %s value does not match expected type, expected %s, got %s", path, gjson.Number, t))
	}
}

func (c *Config) Str(path string) string {
	res := c.v.Get(path)
	if !res.Exists() {
		panic(fmt.Sprintf("path %s does not exist in config", path))
	}
	if t := res.Type; t == gjson.String {
		return res.Str
	} else {
		panic(fmt.Sprintf("path %s value does not match expected type, expected %s, got %s", path, gjson.String, t))
	}
}

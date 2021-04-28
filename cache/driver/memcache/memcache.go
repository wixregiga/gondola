// +build !appengine

package memcache

import (
	"net"
	"strings"
	"time"

	"gondola/cache/driver"
	"gondola/config"

	"github.com/rainycape/memcache"
)

type memcacheDriver struct {
	*memcache.Client
}

func (c *memcacheDriver) Set(key string, b []byte, timeout int) error {
	item := memcache.Item{Key: key, Value: b, Expiration: int32(timeout)}
	return c.error(c.Client.Set(&item))
}

func (c *memcacheDriver) Get(key string) ([]byte, error) {
	item, err := c.Client.Get(key)
	if err != nil {
		return nil, c.error(err)
	}
	if item != nil {
		return item.Value, nil
	}
	return nil, nil
}

func (c *memcacheDriver) GetMulti(keys []string) (map[string][]byte, error) {
	results, err := c.Client.GetMulti(keys)
	if err != nil {
		return nil, c.error(err)
	}
	value := make(map[string][]byte, len(results))
	for k, v := range results {
		value[k] = v.Value
	}
	return value, nil
}

func (c *memcacheDriver) Delete(key string) error {
	return c.error(c.Client.Delete(key))
}

func (c *memcacheDriver) Connection() interface{} {
	return c.Client
}

func (c *memcacheDriver) Flush() error {
	return c.Client.Flush(0)
}

func (c *memcacheDriver) error(err error) error {
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil
		}
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			// Don't log these errors since they're so frequent
			// in memcache that they end up generating a lot
			// of logs
			return nil
		}
	}
	return err
}

func memcacheOpener(url *config.URL) (driver.Driver, error) {
	hosts := strings.Split(url.Value, ",")
	conns := make([]string, len(hosts))
	for ii, v := range hosts {
		conns[ii] = driver.DefaultPort(v, 11211)
	}
	client, err := memcache.New(conns...)
	if err != nil {
		return nil, err
	}
	if tm, ok := url.Fragment.Int("timeout"); ok {
		client.SetTimeout(time.Millisecond * time.Duration(tm))
	}
	if maxIdle, ok := url.Fragment.Int("max_idle"); ok {
		client.SetMaxIdleConnsPerAddr(maxIdle)
	}
	return &memcacheDriver{Client: client}, nil
}

func init() {
	driver.Register("memcache", memcacheOpener)
}

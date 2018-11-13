package cache

import (
	"config"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/nuveo/log"

	"github.com/SimonWaldherr/golibs/cache"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
)

type iDistributedCache interface {
	Get(key string) (valid bool, value string)
	Set(key string, value interface{}, duration time.Duration)
	Remove(key string)
	HGet(key string, field string) (valid bool, value string)
	HSet(key string, field string, value interface{})
	HRemove(key string, field string)
	HClear(key string, prefix string)
	Clear()
}

//MemoryCache Object
type MemoryCache struct {
	mcache  *cache.Cache
	dCacher iDistributedCache
	lock    sync.RWMutex
	lock2   sync.RWMutex
	lock3   sync.RWMutex
}

//GetMemoryItem is used to get a cache value with a expire callback function, function can be null
func (cache *MemoryCache) GetMemoryItem(key string, callBack func() interface{}, duration time.Duration) interface{} {
	item := cache.mcache.Get(key)

	if item == nil && callBack != nil {
		cache.lock.Lock()
		defer cache.lock.Unlock()
		item = cache.mcache.Get(key)
		if item == nil && callBack != nil {
			item = callBack()
			log.Printf("cache: expired callback function called, key:%s", key)
			if duration == -1 {
				duration = cache.mcache.Expiration
			}
			cache.mcache.SetWithDuration(key, item, time.Now(), duration)
		}

	}
	return item
}

func (cache *MemoryCache) SetMemoryItem(key string, value interface{}, duration time.Duration) {
	cache.mcache.SetWithDuration(key, value, time.Now(), duration)
}

func (cache *MemoryCache) RemoveMemoryItem(key string) {
	cache.mcache.Delete(key)
}

//GetMemoryDistributedItem is used to get distributed cache item
func (cache *MemoryCache) GetMemoryDistributedItem(key string, callBack func() interface{}, duration time.Duration) (value interface{}) {
	if cache.dCacher == nil {
		panic("cache: dCacher is null! You should call NewCacheDistributed first.")
	}
	b, t1 := cache.dCacher.Get(key)
	sv := key + "$version"
	t2 := cache.mcache.Get(sv)

	if !b || t1 != t2 {
		cache.lock2.Lock()
		defer cache.lock2.Unlock()

		b, t1 = cache.dCacher.Get(key)
		t2 = cache.mcache.Get(sv)

		if !b || t1 != t2 {

			cache.mcache.Delete(key)

			if callBack != nil {
				if duration == -1 {
					duration = cache.mcache.Expiration
				} else if duration == 0 {
					duration = time.Minute * math.MaxInt16
				}
				if !b {
					nver := uuid.Must(uuid.NewV4()).String()
					if duration == -1 {
						duration = cache.mcache.Expiration
					}
					cache.dCacher.Set(key, nver, duration)
					cache.mcache.SetWithDuration(sv, nver, time.Now(), duration)
				} else {
					cache.mcache.SetWithDuration(sv, t1, time.Now(), duration)
				}
			}
		}
	}
	value = cache.GetMemoryItem(key, callBack, duration)
	return value
}

//RemoveMemoryDistributedItem is used to remove distributed cache item
func (cache *MemoryCache) RemoveMemoryDistributedItem(key string) {
	cache.mcache.Delete(key)
	cache.mcache.Delete(key + "$version")
	cache.dCacher.Remove(key)
}

func (cache *MemoryCache) GetMemoryDistributedHashItem(key string, field string, callBack func() interface{}, duration time.Duration) (value interface{}) {
	if cache.dCacher == nil {
		panic("cache: dCacher is null! You should call NewCacheDistributed first.")
	}
	b, t1 := cache.dCacher.HGet(key, field)
	slKey := key + "$" + field
	sv := slKey + "$version"
	t2 := cache.mcache.Get(sv)

	if !b || t1 != t2 {
		cache.lock3.Lock()
		defer cache.lock3.Unlock()

		b, t1 := cache.dCacher.HGet(key, field)
		t2 = cache.mcache.Get(sv)

		if !b || t1 != t2 {
			cache.mcache.Delete(slKey)
			if callBack != nil {
				if duration == -1 {
					duration = cache.mcache.Expiration
				}
				if !b {
					nver := uuid.Must(uuid.NewV4()).String()
					if duration == -1 {
						duration = cache.mcache.Expiration
					}
					cache.dCacher.HSet(key, field, nver)
					cache.mcache.SetWithDuration(sv, nver, time.Now(), duration)
				} else {
					cache.mcache.SetWithDuration(sv, t1, time.Now(), duration)
				}
			}
		}
	}
	value = cache.GetMemoryItem(slKey, callBack, duration)
	return value
}

func (cache *MemoryCache) SetMemoryDistributedHashItem(key string, field string, value interface{}) {
	slKey := key + "$" + field
	sv := slKey + "$version"

	cache.mcache.Set(slKey, value)

	b, v := cache.dCacher.HGet(key, field)
	if !b {
		cache.lock3.Lock()
		defer cache.lock3.Unlock()
		if b, v = cache.dCacher.HGet(key, field); !b {
			v = uuid.Must(uuid.NewV4()).String()
			cache.dCacher.HSet(key, field, v)
		}
	}
	cache.mcache.Set(sv, v)
}

func (cache *MemoryCache) RemoveMemoryDistributedHashItem(key string, field string) {
	slKey := key + "$" + field
	sv := slKey + "$version"
	cache.mcache.Delete(slKey)
	cache.mcache.Delete(sv)
	cache.dCacher.HRemove(key, field)
}

func (cache *MemoryCache) GetDistributedItem(key string) (bool, string) {
	return cache.dCacher.Get(key)
}

func (cache *MemoryCache) SetDistributedItem(key string, value string, duration time.Duration) {
	cache.dCacher.Set(key, value, duration)
}

func (cache *MemoryCache) RemoveDistributedItem(key string) {
	cache.dCacher.Remove(key)
}

func (cache *MemoryCache) GetDistributedHashItem(key string, field string) (bool, string) {
	return cache.dCacher.HGet(key, field)
}

func (cache *MemoryCache) SetDistributedHashItem(key string, field string, value string) {
	cache.dCacher.HSet(key, field, value)
}

func (cache *MemoryCache) RemoveDistributedHashItem(key string, field string) {
	cache.dCacher.HRemove(key, field)
}

func (cache *MemoryCache) ClearDistributedHashItem(key string, prefix string) {
	cache.dCacher.HClear(key, prefix)
}

func (cache *MemoryCache) ClearCache() {
	cache.mcache.Clear()
	if cache.dCacher != nil {
		cache.dCacher.Clear()
	}
}

//NewCache "defaultDuration"=-1 means nerver expired.
func NewCache(defaultDuration, cleanInterval time.Duration) *MemoryCache {
	if defaultDuration == -1 {
		defaultDuration = math.MaxInt16 * time.Minute
	}
	c := cache.New2(defaultDuration, cleanInterval, func(key string, value interface{}) {
		log.Printf("cache: key %s is expired.", key)
	})
	m := &MemoryCache{}
	m.mcache = c
	return m
}

//NewCacheDistributed with redis."defaultDuration"=-1 means nerver expired. "masterName" can be empty if not use sentinal mode.
func NewCacheDistributed(defaultDuration, cleanInterval time.Duration, conf config.RedisCfg) *MemoryCache {
	m := NewCache(defaultDuration, cleanInterval)
	var d iDistributedCache
	if conf.MasterName == "" {
		d = iDistributedCache(newCacher(strings.Split(conf.Hosts, ",")[0], conf.Password, conf.DBNumber, conf.IdleTimeout, conf.PoolSize, conf.MinIdleConns))
	} else {
		d = iDistributedCache(newSentinelCacher(strings.Split(conf.Hosts, ","), conf.Password, conf.MasterName, conf.DBNumber, conf.IdleTimeout, conf.PoolSize, conf.MinIdleConns))
	}
	m.dCacher = d
	return m
}

type redisCacher struct {
	client *redis.Client
}

var _ iDistributedCache = (*redisCacher)(nil)

func newSentinelCacher(addrs []string, pwd string, masterName string, dbNumber int, idleTimeout int, poolSize int, minIdleConns int) *redisCacher {
	return &redisCacher{
		client: redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: addrs,
			Password:      pwd,
			MasterName:    masterName,
			DB:            dbNumber,
			IdleTimeout:   time.Second * time.Duration(idleTimeout),
			PoolSize:      poolSize,
			MinIdleConns:  minIdleConns,
		})}
}
func newCacher(addrs string, pwd string, dbNumber int, idleTimeout int, poolSize int, minIdleConns int) *redisCacher {
	return &redisCacher{
		client: redis.NewClient(&redis.Options{
			Addr:         addrs,
			Password:     pwd,
			DB:           dbNumber,
			IdleTimeout:  time.Second * time.Duration(idleTimeout),
			PoolSize:     poolSize,
			MinIdleConns: minIdleConns,
		})}
}

func (r *redisCacher) Get(key string) (valid bool, value string) {
	value, err := r.client.Get(key).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	return err == nil, value
}

func (r *redisCacher) Set(key string, value interface{}, duration time.Duration) {
	err := r.client.Set(key, value, duration).Err()
	if err != nil {
		panic(err)
	}
}

func (r *redisCacher) Remove(key string) {
	err := r.client.Del(key).Err()
	if err != nil {
		panic(err)
	}
}

func (r *redisCacher) HGet(key string, field string) (valid bool, value string) {
	value, err := r.client.HGet(key, field).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	return err == nil, value
}

func (r *redisCacher) HSet(key string, field string, value interface{}) {
	err := r.client.HSet(key, field, value).Err()
	if err != nil {
		panic(err)
	}
}

func (r *redisCacher) HRemove(key string, field string) {
	err := r.client.HDel(key, field).Err()
	if err != nil {
		panic(err)
	}
}

func (r *redisCacher) HClear(key string, prefix string) {
	keys, err := r.client.HKeys(key).Result()
	if err != nil {
		panic(err)
	}
	for _, field := range keys {
		if strings.HasPrefix(field, prefix) {
			err := r.client.HDel(key, field).Err()
			if err != nil {
				panic(err)
			}
		}
	}
}

func (r *redisCacher) Clear() {
	err := r.client.FlushDBAsync().Err()
	if err != nil {
		panic(err)
	}
}

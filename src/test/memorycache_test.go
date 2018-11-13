package test

import (
	"config"
	"framework/cache"
	"framework/utils"
	"log"
	"testing"
	"time"
)

func Benchmark_memcache(t *testing.B) {
	conf := config.RedisCfg{
		Hosts:      "10.202.80.117:17900,10.202.80.118:17900",
		Password:   "",
		MasterName: "master7900",
	}
	m := cache.NewCacheDistributed(time.Minute*60, time.Minute*10, conf)
	done := make(chan int, t.N)
	defer close(done)
	for index := 0; index < t.N; index++ {
		go func(idx int) {
			val := m.GetMemoryDistributedItem("test1", func() interface{} {
				return time.Now().Nanosecond()
			}, time.Second)
			log.Printf("index:%d,val:%v", idx, val)
			done <- val.(int)
		}(index)
	}
	for index := 0; index < t.N; index++ {
		<-done
	}
}
func Benchmark_memcache2(t *testing.B) {
	conf := config.RedisCfg{
		Hosts:    "10.202.80.117:7900,10.202.80.118:7900",
		Password: "",
	}
	m := cache.NewCacheDistributed(time.Minute*60, time.Minute*10, conf)
	for index := 0; index < t.N; index++ {
		val := m.GetMemoryDistributedItem("test2", func() interface{} {
			return utils.JoinString("-", time.Now().Nanosecond(), index)
		}, time.Second)
		log.Printf("index:%d,val:%v", index, val)
	}
}

package cache

import (
	"config"
	"time"
)

type ICacheStore interface {
	Put(tableName, key string, value interface{}) error
	Get(tableName, key string) (interface{}, error)
	Del(tableName, key string) error

	PutD(tableName, key string, value string) error
	GetD(tableName, key string) (string, bool)
	DelD(tableName, key string) error
	Clear(tableName, prefix string) error
}

type DBCacheStore struct {
	memory          *MemoryCache
	defaultDuration time.Duration
}

var _ ICacheStore = (*DBCacheStore)(nil)

func NewDBCacheStore(defaultDuration time.Duration, conf config.RedisCfg) *DBCacheStore {
	m := NewCacheDistributed(defaultDuration, 0, conf)
	dbc := &DBCacheStore{m, defaultDuration}
	return dbc
}

func (m *DBCacheStore) Get(tableName, key string) (interface{}, error) {
	val := m.memory.GetMemoryDistributedHashItem(tableName, key, nil, 0)

	return val, nil
}

func (m *DBCacheStore) Put(tableName, key string, value interface{}) error {
	m.memory.SetMemoryDistributedHashItem(tableName, key, value)
	return nil
}

func (m *DBCacheStore) Del(tableName, key string) error {
	m.memory.RemoveMemoryDistributedHashItem(tableName, key)
	return nil
}

func (m *DBCacheStore) GetD(tableName, key string) (string, bool) {
	ok, val := m.memory.GetDistributedHashItem(tableName, key)
	return val, ok
}

func (m *DBCacheStore) PutD(tableName, key string, value string) error {
	m.memory.SetDistributedHashItem(tableName, key, value)
	return nil
}

func (m *DBCacheStore) DelD(tableName, key string) error {
	m.memory.RemoveDistributedHashItem(tableName, key)
	return nil
}

func (m *DBCacheStore) Clear(tableName, prefix string) error {
	m.memory.ClearDistributedHashItem(tableName, prefix)
	return nil
}

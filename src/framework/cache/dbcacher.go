package cache

import (
	"fmt"
	"hash/crc32"

	"github.com/go-xorm/core"
)

type DBCacher struct {
	store ICacheStore
}

var _ core.Cacher = (*DBCacher)(nil)

func NewDBCacher(store ICacheStore) *DBCacher {
	return &DBCacher{store}
}

func (c *DBCacher) GetIds(tableName, sql string) interface{} {
	key := getSqlKey(tableName, sql)
	val, ok := c.store.GetD(tableName, key)
	if !ok {
		return nil
	}
	return val
}

func (c *DBCacher) GetBean(tableName string, id string) interface{} {
	key := getBeanKey(tableName, id)
	val, _ := c.store.Get(tableName, key)
	return val
}

func (c *DBCacher) PutIds(tableName, sql string, ids interface{}) {
	key := getSqlKey(tableName, sql)
	c.store.PutD(tableName, key, ids.(string))
}

func (c *DBCacher) PutBean(tableName string, id string, obj interface{}) {
	key := getBeanKey(tableName, id)
	c.store.Put(tableName, key, obj)
}

func (c *DBCacher) DelIds(tableName, sql string) {
	key := getSqlKey(tableName, sql)
	c.store.DelD(tableName, key)
}

func (c *DBCacher) DelBean(tableName string, id string) {
	key := getBeanKey(tableName, id)
	c.store.Del(tableName, key)
}

func (c *DBCacher) ClearIds(tableName string) {
	key := getTableSqlKey(tableName)
	c.store.Clear(tableName, key)
}

func (c *DBCacher) ClearBeans(tableName string) {
	key := getTableBeanKey(tableName)
	c.store.Clear(tableName, key)
}

func getBeanKey(tableName string, id string) string {
	return fmt.Sprintf("dbcacher:bean:%s:%s", tableName, id)
}
func getTableBeanKey(tableName string) string {
	return fmt.Sprintf("dbcacher:bean:%s:", tableName)
}

func getSqlKey(tableName string, sql string) string {
	// hash sql to minimize key length
	crc := crc32.ChecksumIEEE([]byte(sql))
	return fmt.Sprintf("dbcacher:sql:%s:%d", tableName, crc)
}
func getTableSqlKey(tableName string) string {
	// hash sql to minimize key length
	return fmt.Sprintf("dbcacher:sql:%s:", tableName)
}

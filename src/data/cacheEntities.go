package data

import (
	"datamodels"
	"log"
)

func CacheEntities(db *Database, dbReadOnly *Database) {
	CacheEntity(new(datamodels.User), db, dbReadOnly)
	CacheEntity(new(datamodels.UserDetail), db, dbReadOnly)
	log.Println("CacheEntities: finished")
}

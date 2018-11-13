package data

import (
	"fmt"
	"framework/cache"
	"framework/globals"
	"math"
	"strings"

	_ "github.com/denisenkom/go-mssqldb" //mssql
	_ "github.com/go-sql-driver/mysql"   //mysql
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq" //pg
	"github.com/nuveo/log"
)

type Database struct {
	engin   *xorm.Engine
	session *xorm.Session
}

func NewDB(driverName string, dataSourceName string) (*Database, error) {
	var err error
	var engin *xorm.Engine
	if strings.Contains(dataSourceName, ",") {
		sources := strings.Split(dataSourceName, ",")
		var enginGroup *xorm.EngineGroup
		enginGroup, err = xorm.NewEngineGroup(driverName, sources)
		engin = enginGroup.Engine
	} else {
		engin, err = xorm.NewEngine(driverName, dataSourceName)
	}
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	if globals.Config.AppIsDebug {
		engin.ShowSQL()
	}
	db := &Database{engin, nil}
	return db, nil
}

func (db *Database) NewTransaction() *Database {
	dbnew := &Database{engin: db.engin}
	dbnew.session = db.engin.NewSession()
	dbnew.session.Begin()
	return dbnew
}
func (db *Database) Commit() {
	db.session.Commit()
	db.session.Close()
	db.session = nil
}
func (db *Database) RollBack() {
	db.session.Rollback()
	db.session.Close()
	db.session = nil
}

func (db *Database) List(slicePtr interface{}, boolColumn string, queryObject ...interface{}) error {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}

	for i := len(queryObject) - 1; i > 0; i++ {
		session = session.Or(queryObject[i])
	}
	if boolColumn != "" {
		if strings.Contains(boolColumn, ",") {
			session = session.UseBool(strings.Split(boolColumn, ",")...)
		} else {
			session = session.UseBool(boolColumn)
		}
	}
	err := session.Find(slicePtr, queryObject...)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func (db *Database) ListBy(slicePtr interface{}, query string, params ...interface{}) error {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	err := session.Where(query, params...).Find(slicePtr)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func (db *Database) ListByObjects(slicePtr interface{}, cols string, start int, limit int, orderby string, isdecending bool, joins *[]Join, queryObject ...interface{}) error {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	if limit == -1 {
		limit = math.MaxUint32
	}
	session = session.Cols(cols).Limit(limit, start)
	if len(orderby) > 0 {
		if isdecending {
			session = session.Desc(orderby)
		} else {
			session = session.Asc(orderby)
		}
	}
	if joins != nil {
		for _, j := range *joins {
			session = session.Join(j.Operator, j.TableOrName, j.JoinCondition)
		}
	}
	err := session.Find(slicePtr, queryObject...)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func (db *Database) ListByCondition(slicePtr interface{}, cols string, start int, limit int, orderby string, isdecending bool, joins *[]Join, query string, params ...interface{}) error {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	if limit == -1 {
		limit = math.MaxUint32
	}
	session = session.Cols(cols).Limit(limit, start).Where(query, params...)
	if len(orderby) > 0 {
		if isdecending {
			session = session.Desc(orderby)
		} else {
			session = session.Asc(orderby)
		}
	}
	if joins != nil {
		for _, j := range *joins {
			session = session.Join(j.Operator, j.TableOrName, j.JoinCondition)
		}
	}
	err := session.Find(slicePtr)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func (db *Database) Get(dest interface{}) (bool, error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	b, err := session.Get(dest)
	if err != nil {
		log.Errorln(err)
	}
	return b, err
}

func (db *Database) GetByID(ID interface{}, dest interface{}) (bool, error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	b, err := session.ID(ID).Get(dest)
	if err != nil {
		log.Errorln(err)
	}
	return b, err
}

func (db *Database) GetByOrder(dest interface{}, orderby string, isdecending bool) (bool, error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	if len(orderby) > 0 {
		if isdecending {
			session = session.Desc(orderby)
		} else {
			session = session.Asc(orderby)
		}
	}
	b, err := session.Get(dest)
	if err != nil {
		log.Errorln(err)
	}
	return b, err
}

func (db *Database) Create(data ...interface{}) (affect int64, err error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	affect, err = session.Insert(data...)
	if err != nil {
		log.Errorln(err)
	}
	return
}

func (db *Database) Update(data interface{}, isAllCols bool, conditionObject ...interface{}) (affect int64, err error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	if isAllCols {
		session = session.AllCols()
	}
	affect, err = session.Update(data, conditionObject...)
	if err != nil {
		log.Errorln(err)
	}
	return
}

func (db *Database) UpdateByID(ID interface{}, data interface{}, isAllCols bool) (affect int64, err error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	if isAllCols {
		session = session.AllCols()
	}
	affect, err = session.ID(ID).Update(data)
	if err != nil {
		log.Errorln(err)
	}
	return
}

func (db *Database) UpdateByCondition(data interface{}, isAllCols bool, query string, params ...interface{}) (affect int64, err error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	if isAllCols {
		session = session.AllCols()
	}
	affect, err = session.Where(query, params...).Update(data)
	if err != nil {
		log.Errorln(err)
	}
	return
}

func (db *Database) Delete(conditionObject interface{}) (affect int64, err error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	affect, err = session.Delete(conditionObject)
	if err != nil {
		log.Errorln(err)
	}
	return
}

func (db *Database) DeleteByCondition(conditionObject interface{}, query string, params ...interface{}) (affect int64, err error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	affect, err = session.Where(query, params...).Delete(conditionObject)
	if err != nil {
		log.Errorln(err)
	}
	return
}

func (db *Database) Exec(query string, params ...interface{}) (int64, error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	qs := []interface{}{query}
	if params != nil {
		qs = append(qs, params...)
	}
	result, err := session.Exec(qs...)
	if err != nil {
		log.Errorln(err)
	}
	n, _ := result.RowsAffected()
	return n, err
}

func (db *Database) Query(slicePtr interface{}, selectquery string, params ...interface{}) error {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	err := session.Sql(selectquery, params...).Find(slicePtr)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func (db *Database) Count(dest interface{}, query string, params ...interface{}) (n int64, err error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	if _, ok := dest.(string); ok {
		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s", dest)
		if len(query) > 0 {
			sql += fmt.Sprintf(" WHERE %s", query)
		}
		_, err = session.Sql(sql, params...).Get(&n)
	} else {
		n, err = session.Where(query, params...).Count(dest)
	}
	if err != nil {
		log.Errorln(err)
	}
	return
}

func (db *Database) Sum(dest interface{}, colName string, query string, params ...interface{}) (float64, error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	n, err := session.Where(query, params...).Sum(dest, colName)
	if err != nil {
		log.Errorln(err)
	}
	return n, err
}

func (db *Database) AddOrUpdate(ID interface{}, dest interface{}) (affect int64, err error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	//Hack xorm
	empty := ""
	b, err := session.Table(dest).ID(ID).Exist(&empty)
	if err != nil {
		log.Errorln(err)
	} else {
		if b {
			affect, err = session.ID(ID).Update(dest)
		} else {
			affect, err = session.ID(ID).Insert(dest)
		}
	}
	return
}

func (db *Database) Exists(dest interface{}, query string, params ...interface{}) (bool, error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	b, err := session.Where(query, params...).Exist(dest)
	if err != nil {
		log.Errorln(err)
	}
	return b, err
}

func (db *Database) IsTableExist(dest interface{}) (bool, error) {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	b, err := session.IsTableExist(dest)
	if err != nil {
		log.Errorln(err)
	}
	return b, err
}

func (db *Database) DropTable(dest interface{}) error {
	var session *xorm.Session
	if db.session == nil {
		session = db.engin.NewSession()
		defer session.Close()
	} else {
		session = db.session
	}
	err := session.DropTable(dest)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func (db *Database) Sync(entity interface{}) error {
	err := db.engin.Sync2(entity)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func CacheEntity(entity interface{}, dbs ...*Database) {
	conf := globals.Config.Redis
	dbcache := cache.NewDBCacheStore(-1, conf)
	cacher := cache.NewDBCacher(dbcache)
	for _, db := range dbs {
		db.engin.MapCacher(entity, cacher)
	}
}

// func (db *Database) CacheEntity(entity interface{}) {
// 	dbcache := xorm.NewMemoryStore()
// 	cacher := xorm.NewLRUCacher(dbcache, 100000)
// 	db.engin.MapCacher(entity, cacher)
// }
func (db *Database) ClearCacheEntity(entity interface{}) {
	db.engin.ClearCache(entity)
}

//Join expression struct
type Join struct {
	//Operator:INNER, LEFT OUTER, CROSS
	Operator string
	//TableOrName:tablename or bean
	TableOrName interface{}
	//JoinCondition:users.id=user_detail.userid
	JoinCondition string
}

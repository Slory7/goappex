package main

import (
	"appstart"
	"config"
	"data"
	"data/migration"
	"data/migration/migrations"
	"data/repositories"
	"flag"
	"framework/cache"
	"framework/globals"
	"framework/validates"
	"time"

	"github.com/kataras/iris"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/crgimenes/goconfig/json"
	"github.com/nuveo/log"
)

var (
	rollbackVersionID = flag.String("rollback", "", "Rollback migration version id")
)

func main() {

	flag.Parse()

	//Config
	globals.Config = config.GetConfig(globals.GetEnvironment())
	conf := globals.Config

	//Cache
	globals.Cache = cache.NewCacheDistributed(time.Minute*120, time.Minute*5, conf.Redis)

	//validator
	globals.Validator = validates.NewValidator()

	//db
	log.Println("db type:", conf.DBType)

	db, err := data.NewDB(conf.DBType, conf.DBConnectionString)
	if err != nil {
		log.Fatal("db: ", err)
	}
	sDBReadConnectionString := conf.DBReadOnlyConnectionString
	if len(sDBReadConnectionString) == 0 {
		sDBReadConnectionString = conf.DBConnectionString
	}
	dbReadOnly, err := data.NewDB(conf.DBType, sDBReadConnectionString)
	if err != nil {
		log.Fatal("db: ", err)
	}

	//db.Sync(new(datamodels.User))

	//db migration
	mig := migration.New(db, &migration.Options{TableName: "appversions", IDColumnName: "versionid"}, migrations.MigrationVersions)
	mig.SetInitSchema(migrations.InitMigration())

	//rollback command
	if len(*rollbackVersionID) > 0 {
		rmig := mig.GetMigration(*rollbackVersionID)
		if rmig == nil {
			log.Fatal("rollback migration not exists: ", *rollbackVersionID)
		}
		err = mig.RollbackMigration(rmig)
		if err != nil {
			log.Fatal("rollback migration error: ", err)
		}
	} else {
		err = mig.Migrate()
		if err != nil {
			log.Fatal("migration error: ", err)
		}
	}

	//data cache
	data.CacheEntities(db, dbReadOnly)

	repo := repositories.NewRepository(db)
	repoReadOnly := repositories.NewRepositoryReadOnly(dbReadOnly)

	//IoC
	appstart.RegisterIoC(repo, repoReadOnly)

	app := iris.New()

	//routes
	appstart.ConfigureRoutes(app)

	//curl localhost:8181/debug/pprof/trace?seconds=10 > trace.out
	//go tool trace goappex.exe trace.out
	//http://www.sharelinux.com/2017/03/22/Golang%E4%B9%8Bprofiler%E5%92%8Ctrace%E5%B7%A5%E5%85%B7/
	if globals.Config.AppIsDebug {
		go func() {
			http.ListenAndServe("localhost:8181", http.DefaultServeMux)
		}()
	}

	app.Run(
		// Start the web server at localhost:8080
		iris.Addr(conf.Addr),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)
}

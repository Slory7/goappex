package migration

import (
	"data"
	"errors"
	"fmt"
)

type MigrateFunc func(*data.Database) error

// RollbackFunc is the func signature for rollbacking.
type RollbackFunc func(*data.Database) error

// InitSchemaFunc is the func signature for initializing the schema.
type InitSchemaFunc func(*data.Database) (string, error)

// Options define options for all migrations.
type Options struct {
	// TableName is the migration table.
	TableName string
	// IDColumnName is the name of column where the migration id will be stored.
	IDColumnName string
}

// Migration represents a database migration (a modification to be made on the database).
type Migration struct {
	// ID is the migration identifier. Usually a timestamp like "201601021504".
	ID string
	//Description is the version memo
	Description string
	// Migrate is a function that will br executed while running this migration.
	Migrate MigrateFunc
	// Rollback will be executed on rollback. Can be nil.
	Rollback RollbackFunc
}

// Migrate represents a collection of all migrations of a database schema.
type Migrate struct {
	db         *data.Database
	options    *Options
	migrations []*Migration
	initSchema *Migration
}

var (
	// DefaultOptions can be used if you don't want to think about options.
	DefaultOptions = &Options{
		TableName:    "migrations",
		IDColumnName: "id",
	}

	// ErrRollbackImpossible is returned when trying to rollback a migration
	// that has no rollback function.
	ErrRollbackImpossible = errors.New("It's impossible to rollback this migration")

	// ErrNoMigrationDefined is returned when no migration is defined.
	ErrNoMigrationDefined = errors.New("No migration defined")

	// ErrMissingID is returned when the ID od migration is equal to ""
	ErrMissingID = errors.New("Missing ID in migration")

	// ErrNoRunnedMigration is returned when any runned migration was found while
	// running RollbackLast
	ErrNoRunnedMigration = errors.New("Could not find last runned migration")
)

// New returns a new Gormigrate.
func New(db *data.Database, options *Options, migrations []*Migration) *Migrate {
	return &Migrate{
		db:         db,
		options:    options,
		migrations: migrations,
	}
}

// SetInitSchema sets a function that is run if no migration is found.
// The idea is preventing to run all migrations when a new clean database
// is being migrating. In this function you should create all tables and
// foreign key necessary to your application.
func (m *Migrate) SetInitSchema(initSchema *Migration) {
	m.initSchema = initSchema
}

// Migrate executes all migrations that did not run yet.
func (m *Migrate) Migrate() error {
	if err := m.createMigrationTableIfNotExists(); err != nil {
		return err
	}

	if m.initSchema != nil && m.isFirstRun() {
		if err := m.runMigration(m.initSchema); err != nil {
			return err
		}
	}

	for _, migration := range m.migrations {
		if m.initSchema == nil || m.initSchema.ID < migration.ID {
			if err := m.runMigration(migration); err != nil {
				return err
			}
		}
	}
	return nil
}

// RollbackLast undo the last migration
func (m *Migrate) RollbackLast() error {
	if len(m.migrations) == 0 {
		return ErrNoMigrationDefined
	}

	lastRunnedMigration, err := m.getLastRunnedMigration()
	if err != nil {
		return err
	}

	if err := m.RollbackMigration(lastRunnedMigration); err != nil {
		return err
	}
	return nil
}

func (m *Migrate) getLastRunnedMigration() (*Migration, error) {
	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]
		run, err := m.migrationDidRun(migration)
		if err != nil {
			return nil, err
		} else if run {
			return migration, nil
		}
	}
	return nil, ErrNoRunnedMigration
}

// RollbackMigration undo a migration.
func (m *Migrate) RollbackMigration(mig *Migration) error {
	if mig.Rollback == nil {
		return ErrRollbackImpossible
	}

	if err := mig.Rollback(m.db); err != nil {
		return err
	}

	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", m.options.TableName, m.options.IDColumnName)
	if _, err := m.db.Exec(sql, mig.ID); err != nil {
		return err
	}
	return nil
}

func (m *Migrate) GetMigration(id string) (mig *Migration) {
	if m.initSchema != nil && m.initSchema.ID == id {
		return m.initSchema
	}
	for _, v := range m.migrations {
		if v.ID == id {
			mig = v
			break
		}
	}
	return mig
}

// func (m *Migrate) runInitSchema() error {
// 	baseVer, err := m.initSchema(m.db)
// 	if err != nil {
// 		return err
// 	}

// 	for _, migration := range m.migrations {
// 		if baseVer < migration.ID {
// 			if err := m.runMigration(migration); err != nil {
// 				return err
// 			}
// 		} else if err := m.insertMigration(migration.ID, migration.Description); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (m *Migrate) runMigration(migration *Migration) error {
	if len(migration.ID) == 0 {
		return ErrMissingID
	}

	run, err := m.migrationDidRun(migration)
	if err != nil {
		return err
	}

	if !run {
		if err := migration.Migrate(m.db); err != nil {
			return err
		}

		if err := m.insertMigration(migration.ID, migration.Description); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrate) createMigrationTableIfNotExists() error {
	exists, err := m.db.IsTableExist(m.options.TableName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	sql := fmt.Sprintf("CREATE TABLE %s (%s VARCHAR(255) PRIMARY KEY, description VARCHAR(1000))", m.options.TableName, m.options.IDColumnName)
	if _, err := m.db.Exec(sql); err != nil {
		return err
	}
	return nil
}

func (m *Migrate) migrationDidRun(mig *Migration) (bool, error) {
	count, err := m.db.Count(m.options.TableName, fmt.Sprintf("%s = ?", m.options.IDColumnName), mig.ID)
	return count > 0, err
}

func (m *Migrate) isFirstRun() bool {
	count, _ := m.db.Count(m.options.TableName, "")

	return count == 0
}

func (m *Migrate) insertMigration(id string, description string) error {
	sql := fmt.Sprintf("INSERT INTO %s (%s,description) VALUES (?,?)", m.options.TableName, m.options.IDColumnName)
	_, err := m.db.Exec(sql, id, description)
	return err
}

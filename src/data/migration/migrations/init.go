package migrations

import (
	"data"
	"data/migration"
	m "datamodels"

	"github.com/nuveo/log"
)

func InitMigration() *migration.Migration {
	ver := "v201810111656"
	var mig = &migration.Migration{
		ID:          ver,
		Description: `Init everything`,

		Migrate: func(db *data.Database) error {
			log.Printf("init migrating: %s\n", ver)

			if err := db.Sync(new(m.User)); err != nil {
				return err
			}

			// sql := `ALTER TABLE child_table
			// ADD CONSTRAINT constraint_fk
			// FOREIGN KEY (c1)
			// REFERENCES parent_table(p1)
			// ON DELETE CASCADE;`
			// if _, err := db.Exec(sql); err != nil {
			// 	return err
			// }

			return nil
		},

		Rollback: func(db *data.Database) error {
			log.Printf("init rollback: %s\n", ver)

			//sql := "DROP TABLE `user`"
			if err := db.DropTable(new(m.User)); err != nil {
				return err
			}

			return nil
		},
	}
	return mig
}

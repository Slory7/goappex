package migrations

import (
	"data"
	"data/migration"
	"datamodels"

	"github.com/nuveo/log"
)

func v201810311029() *migration.Migration {
	ver := "v201810311029"
	var mig = &migration.Migration{
		ID:          ver,
		Description: `Add UserDetail and UserLogins table and user(DeletedAt)`,

		Migrate: func(db *data.Database) error {
			log.Printf("migrating: %s\n", ver)
			err := db.Sync(new(datamodels.User))
			if err == nil {
				err = db.Sync(new(datamodels.UserDetail))
			}
			if err == nil {
				err = db.Sync(new(datamodels.UserLogins))
			}
			if err == nil {
				sql1 := `ALTER TABLE user_detail 
				ADD CONSTRAINT fk_userdetail_userid FOREIGN KEY (userid) REFERENCES users (id)
				ON DELETE CASCADE;`
				sql2 := `ALTER TABLE user_logins 
				ADD CONSTRAINT fk_userlogins_userid FOREIGN KEY (userid) REFERENCES users (id)
				ON DELETE CASCADE;`
				_, err = db.Exec(sql1)
				_, err = db.Exec(sql2)
			}
			return err
		},

		Rollback: func(db *data.Database) error {
			log.Printf("rollback: %s\n", ver)
			sql := "ALTER TABLE `users` DROP COLUMN `deleted_at` "
			_, err := db.Exec(sql)
			if err == nil {
				err = db.DropTable(new(datamodels.UserDetail))
			}
			if err == nil {
				err = db.DropTable(new(datamodels.UserLogins))
			}
			return err
		},
	}
	return mig
}

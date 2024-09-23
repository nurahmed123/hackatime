package migrations

import (
	"github.com/hackclub/hackatime/config"
	"gorm.io/gorm"
)

func init() {
	const name = "20210221-add_created_date_column"
	f := migrationFunc{
		name: name,
		f: func(db *gorm.DB, cfg *config.Config) error {
			if hasRun(name, db) {
				return nil
			}

			if err := db.Exec("UPDATE heartbeats SET created_at = time").Error; err != nil {
				return err
			}

			setHasRun(name, db)
			return nil
		},
	}

	registerPostMigration(f)
}

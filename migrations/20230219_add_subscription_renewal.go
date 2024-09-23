package migrations

import (
	"log/slog"

	"github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/models"
	"gorm.io/gorm"
)

func init() {
	const name = "20230219-add_subscription_renewal"
	f := migrationFunc{
		name: name,
		f: func(db *gorm.DB, cfg *config.Config) error {
			if hasRun(name, db) {
				return nil
			}

			migrator := db.Migrator()

			if migrator.HasColumn(&models.User{}, "subscription_renewal") {
				slog.Info("running migration", "name", name)

				if err := db.Exec("UPDATE users SET subscription_renewal = subscribed_until WHERE subscribed_until is not null").Error; err != nil {
					return err
				}
			}

			setHasRun(name, db)
			return nil
		},
	}

	registerPostMigration(f)
}

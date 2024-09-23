package migrations

import (
	"log/slog"

	"github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/models"
	"gorm.io/gorm"
)

func init() {
	const name = "20221016-drop_rank_column"
	f := migrationFunc{
		name: name,
		f: func(db *gorm.DB, cfg *config.Config) error {
			if hasRun(name, db) {
				return nil
			}

			migrator := db.Migrator()

			if migrator.HasTable(&models.LeaderboardItem{}) && migrator.HasColumn(&models.LeaderboardItem{}, "rank") {
				slog.Info("running migration", "name", name)

				if err := migrator.DropColumn(&models.LeaderboardItem{}, "rank"); err != nil {
					slog.Warn("failed to drop column", "column", "rank", "error", err)
				}
			}

			setHasRun(name, db)
			return nil
		},
	}

	registerPostMigration(f)
}

package migrate

import (
	"Edos_Docer/config"
	"Edos_Docer/iternal/domain"
	"github.com/gofrs/uuid"
	"gopkg.in/gormigrate.v1"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Migrate запустите миграцию для всех объектов и добавьте для них ограничения
// создаем таблицы и закидываем в бд тут
func Migrate() {
	db := config.DB
	regID, _ := uuid.NewV4()
	infoCSVfileID, _ := uuid.NewV4()
	videoInfoID, _ := uuid.NewV4()
	// создаем объект миграции данная строка всегда статична (всегда такая)
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			// id всех миграций кторые были проведены
			ID: regID.String(),
			// переписываем так при создании таблицы изменяется только структура которую мы передаем
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.User{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("users").Error
				if err != nil {
					return err
				}
				return nil
			},
		}, {
			// id всех миграций кторые были проведены
			ID: infoCSVfileID.String(),
			// переписываем так при создании таблицы изменяется только структура которую мы передаем
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.PostCSV{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("settings_csv").Error
				if err != nil {
					return err
				}
				return nil
			},
		}, {
			// id всех миграций кторые были проведены
			ID: videoInfoID.String(),
			// переписываем так при создании таблицы изменяется только структура которую мы передаем
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.VideoInfo{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("video_infos").Error
				if err != nil {
					return err
				}
				return nil
			},
		},
	})

	err := m.Migrate()
	if err != nil {
		log.WithField("component", "migration").Panic(err)
	}

	if err == nil {
		log.WithField("component", "migration").Info("Migration did run successfully")
	} else {
		log.WithField("component", "migration").Infof("Could not migrate: %v", err)
	}
}

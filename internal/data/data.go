package data

import (
	"gchat/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewUserRepo, NewChatRepo, NewMessageRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}

func NewDB(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic("failed to connect database 2")
	}

	sqlDb.SetConnMaxLifetime(time.Second)
	InitDB(db)
	return db
}

func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate( /*&User{}, &Chat{},*/ &Message{}, &UserChatRela{}, &MessageGrouping{}); err != nil {
		panic(err)
	}
}

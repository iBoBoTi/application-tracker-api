package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/iBoBoTi/ats/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	GormDB *gorm.DB
}

func GetDB(c config.Config) *Database {
	db := &Database{}
	db.Init(c)
	return db
}

func (d *Database) Init(c config.Config) {
	d.GormDB = getPostgresDB(c)

	if err := migrate(d.GormDB); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}
}

func getPostgresDB(c config.Config) *gorm.DB {
	log.Printf("Connecting to postgres: %+v", c)
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Africa/Lagos", c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
	log.Println("database dsn", postgresDSN)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level Info, Silent, Warn, Error
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	gormConfig := &gorm.Config{
		Logger: newLogger,
	}
	if c.Environment == "production" {
		gormConfig = &gorm.Config{}
	}
	postgresDB, err := gorm.Open(postgres.Open(postgresDSN), gormConfig)
	if err != nil {
		log.Fatal(err)
	}
	return postgresDB
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &JobPosting{}, &Applicant{})
	if err != nil {
		return fmt.Errorf("migrations error: %v", err)
	}

	return nil
}

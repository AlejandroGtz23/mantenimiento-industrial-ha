package database

import (
	"fiber3-app/backend/config"
	"log"
	"time"

	firebird "github.com/flylink888/gorm-firebird"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB gorm connector
var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(firebird.Open(config.Env.DNS), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Error),
		NamingStrategy: firebird.NamingStrategy{},
	})

	if err != nil {
		log.Fatal("Error de conexion: ", err.Error())
	}

	// Pool conn config
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error al obtener la conexión SQL: ", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("¡Conectado a Firebird 5 con GORM!")
}

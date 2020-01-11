package gorest

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

// DB wrapper object
type DB struct {
	NativeDB *gorm.DB
	log      *log.Entry
	// log      *logging.Logger
	// config   *ContentDBConfig
}

var g_DB *DB

func GetDB() *DB {
	return g_DB
}

func NewDb() *DB {
	g_DB = &DB{
		log: log.WithFields(log.Fields{
			"key": "db",
		}),
	}

	if !g_DB.init() {
		g_DB.log.Errorln("Unable to init database")
		return nil
	}

	return g_DB
}

func (db *DB) init() bool {
	var err error
	db.NativeDB, err = gorm.Open("mysql", Get("mysql"))

	if err != nil {
		db.log.Fatalln("Error to init database. Error: %s \n", err)
		return false
	}

	if !db.ValidationQuery() {
		db.NativeDB = nil
		return false
	}

	db.log.Infoln("Connected to DB.")
	return true
}

func (db *DB) ValidationQuery() bool {
	err := db.NativeDB.Select("1").Error
	if err != nil {
		db.log.Errorln("PostgreSQL is down!")
		return false
	}
	db.log.Infoln("PostgreSQL is running!")
	return true
}

func (db *DB) Close() {
	if err := db.NativeDB.Close(); err != nil {
		db.log.Fatalf("Error to close database connection, Error: %s \n", err)
		return
	}
	db.log.Infoln("DB closed")
}

func (db *DB) CreateFakeData() {
	/*
		for _, user := range entity.UserGetFakeData() {
			err := db.NativeDB.Insert(user)
			if err != nil {
				db.log.Errorln(err)
			}
		}
		db.log.Infoln("Created fake data for development")
	*/
}

func (db *DB) CreateSchema(entities ...interface{}) error {
	for _, entity := range entities {
		err := db.NativeDB.AutoMigrate(entity).Error
		if err != nil {
			db.log.Errorln("Error to create schema: %s", err)
			return err
		}
	}

	db.log.Infoln("Created schema")
	return nil
}

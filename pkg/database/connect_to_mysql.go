package database

import (
	"database/sql"
	"github.com/dienggo/diego/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

// connectToMySql : Simplify method and reusable method to connect with mysql driver
func connectToMySql(username string, password string, host string, port string, name string, useTimeStamp string) *gorm.DB {
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + name + "?parseTime=" + strings.ToLower(useTimeStamp)
	mainDsn = dsn
	sqlDB, err1 := sql.Open(MYSQL, dsn)

	if err1 != nil {
		log.Fatal("Fail open connection to mysql ", err1.Error())
	}

	gormDB, err2 := gorm.Open(mysql.New(mysql.Config{
		Conn:                     sqlDB,
		DisableDatetimePrecision: false,
	}), &gorm.Config{})

	if err2 != nil {
		log.Fatal("Fail connect to mysql", err2.Error())
	}

	return gormDB
}

// mysqlOpenConnection : To open connection with mysql driver
func mysqlOpenConnection() *dbc {
	if isConnected {
		return dbConnection
	}

	var db = config.Database()

	gormDB := connectToMySql(db.Username, db.Password, db.Host, db.Port, db.Name, db.UseTimestamp)

	var replicas []*gorm.DB

	// Do connection on replicas mysql database
	if len(db.Replicas) > 0 {
		for _, replica := range db.Replicas {
			if replica.Connection == MYSQL {
				replicaConnect := connectToMySql(replica.Username, replica.Password, replica.Host, replica.Port, replica.Name, replica.UseTimestamp)
				replicas = append(replicas, replicaConnect)
			}
		}
	}

	// Set initiator variable flag to connected for singleton reason
	isConnected = true
	dbConnection = &dbc{
		main:     gormDB,
		replicas: replicas,
	}

	return dbConnection
}

package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/eveisesi/skillz/internal/mysql"
	driver "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func buildMySQL() {

	m := cfg.MySQL

	config := driver.Config{
		User:                 m.User,
		Passwd:               m.Pass,
		Net:                  "tcp",
		Addr:                 m.Host,
		DBName:               m.DB,
		Loc:                  time.UTC,
		Timeout:              time.Second,
		ReadTimeout:          time.Second,
		WriteTimeout:         time.Second,
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Panicf("[MySQL Connect] Failed to connect to mysql server: %s", err)
	}

	db.SetConnMaxIdleTime(time.Second * 5)
	db.SetConnMaxLifetime(time.Second * 30)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	err = db.Ping()
	if err != nil {
		log.Panicf("[MySQL Connect] Failed to ping mysql server: %s", err)
	}

	dbConn = sqlx.NewDb(db, "mysql")
	mysqlClient = mysql.NewQueryLogger(dbConn, logger)

}

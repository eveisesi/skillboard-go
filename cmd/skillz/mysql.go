package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func buildMySQL() *sqlx.DB {

	m := cfg.MySQL

	config := mysql.Config{
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
	db.SetMaxIdleConns(15)

	err = db.Ping()
	if err != nil {
		log.Panicf("[MySQL Connect] Failed to ping mysql server: %s", err)
	}

	return sqlx.NewDb(db, "mysql")

}

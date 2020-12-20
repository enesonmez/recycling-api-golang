package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Veritabanı bilgilerini nesne değişkenleri ile eşitler.
func (db *Db) Assign() {
	db.dbHost = "localhost"
	db.dbPort = 5432
	db.dbUserName = "postgres"
	db.dbPass = "123456789"
	db.dbName = "upcycling"
}

// Veritabanı bağlantısı
func (dbms *Db) Connect() (*sql.DB, error) {
	dbms.Assign()
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbms.dbHost, dbms.dbPort, dbms.dbUserName, dbms.dbPass, dbms.dbName)
	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Veritabanı bağlantısını keser
func (dbms *Db) Close(db *sql.DB) {
	db.Close()
}

// Veritabanı tablolarını oluşturma
func (dbms *Db) CreateTables(db *sql.DB) error {
	users := "CREATE TABLE IF NOT EXISTS users (uID SERIAL PRIMARY KEY NOT NULL, firstName VARCHAR (50) NOT NULL, lastName VARCHAR (60) NOT NULL, phoneNumber VARCHAR (18) NOT NULL UNIQUE, email VARCHAR (50) NOT NULL UNIQUE, password VARCHAR (256) NOT NULL, gender VARCHAR (15) NOT NULL, birthDay TIMESTAMP DEFAULT NULL, recordTime TIMESTAMP NOT NULL, isVerifyEmail boolean NOT NULL, isBlock boolean NOT NULL);"
	_, err := db.Exec(users)
	if err != nil {
		return err
	}
	return nil
}

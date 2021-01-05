package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	//. "oyeco-api/models/config"

	_ "github.com/lib/pq"
)

// Veritabanı bilgilerini nesne değişkenleri ile eşitler.
func (db *Db) Assign() { // localhost postgres 123456789 upcycling => johnny.heliohost.org enesonme_local A25.pSt13Cd; enesonme_upcycling
	/*conf := new(Config) // Konfigürasyon dosyasındaki değişkenleri kullanmak için nesne kullanılmıştır.
	conf.ConfigRead()*/
	db.dbHost = os.Getenv("DBHost")
	db.dbPort, _ = strconv.Atoi(os.Getenv("DBPort"))
	db.dbUserName = os.Getenv("DBUsername")
	db.dbPass = os.Getenv("DBPassword")
	db.dbName = os.Getenv("DBName")
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

// Veritabanı tablolarını oluşturma
func (dbms *Db) CreateTables(db *sql.DB) error {
	users := "CREATE TABLE IF NOT EXISTS users (uID SERIAL PRIMARY KEY NOT NULL, firstName VARCHAR (50) NOT NULL, lastName VARCHAR (60) NOT NULL, phoneNumber VARCHAR (18) NOT NULL UNIQUE, email VARCHAR (50) NOT NULL UNIQUE, password VARCHAR (256) NOT NULL, gender VARCHAR (15) NOT NULL, birthDay TIMESTAMP DEFAULT NULL, recordTime TIMESTAMP NOT NULL, isVerifyEmail boolean NOT NULL, isBlock boolean NOT NULL);"
	address := "CREATE TABLE IF NOT EXISTS address (aID SERIAL PRIMARY KEY NOT NULL, fullAddress TEXT NOT NULL, district TEXT NOT NULL, city TEXT NOT NULL, postcode VARCHAR (50) NOT NULL, userID INT, FOREIGN KEY (userID) REFERENCES users(uID));"
	_, err := db.Exec(users)
	_, errAddress := db.Exec(address)
	if err != nil || errAddress != nil {
		fmt.Println(errAddress.Error())
		os.Exit(1)
		return err
	}
	return nil
}

package db

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"server/internal/constants"

	_ "github.com/go-sql-driver/mysql"
)

// Configuration struct for getting data from config.json file.
type Configuration struct {
	MYSQL_DATABASE string
	MYSQL_USER     string
	MYSQL_PASSWORD string
	MYSQL_SERVER   string
	MYSQL_PORT     string
	JWT_SECRET     string
}

type Database struct {
	db *sql.DB
}

var (
	ConfigFilePath = flag.String("configfilepath", "../config.json", "JSON Config file path")
)

func NewDatabase() (*Database, error) {

	// reading json file
	configfile, err := os.Open(*ConfigFilePath)
	if err != nil {
		return nil, err
	}
	var config Configuration

	//decoding json file and checking for error
	decoder := json.NewDecoder(configfile)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, err
	}
	// closing configuration file at end of program
	defer configfile.Close()

	constants.JWT_SECRET = config.JWT_SECRET

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true",
		config.MYSQL_USER, config.MYSQL_PASSWORD, config.MYSQL_SERVER, config.MYSQL_PORT, config.MYSQL_DATABASE)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var envFile string = "../.env"
var serverAddr, dbAddr string

func main() {

	env, err := godotenv.Read(envFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	serverAddr = env["SERVER_ADDR"]
	if serverAddr == "" {
		serverAddr = "localhost:8080"
		log.Println("Fallback server-addr")
	}
	dbAddr = env["MAIN_DB_ADDR"]
	if dbAddr == "" {
		dbAddr = "root:root@tcp(localhost:3306)/test_db?allowNativePasswords=false&checkConnLiveness=false&maxAllowedPacket=0"
		log.Println("Fallback db-addr")
	}

	dbconf := dbConfig{
		addr:        dbAddr,
		MaxOpenConn: 5,
		MaxIdleConn: 5,
		MaxIdleTime: "5m",
	}

	config := config{
		addr:   serverAddr,
		dbConf: dbconf,
	}

	db, err := sql.Open("mysql", dbAddr)
	if err != nil {
		log.Println(err.Error())
	}

	app := application{
		conf: config,
		db:   db,
	}

	log.Fatal(app.run(app.mount()))

}

package main

import (
	"database/sql"
	"log"
	"yomiya/backend/api/auth"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var envFile string = ".env"
var serverAddr, dbAddr string

func main() {

	env, err := godotenv.Read(envFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	serverAddr = env["SERVER_ADDR"]
	auth.Domain = serverAddr
	if serverAddr == "" {
		serverAddr = "localhost:8080"
		log.Println("Fallback server-addr")
	}
	dbAddr = env["MAIN_DB_ADDR"]
	if dbAddr == "" {
		dbAddr = "root:root@tcp(localhost:3306)/test_db?allowNativePasswords=false&checkConnLiveness=false&maxAllowedPacket=0"
		log.Println("Fallback db-addr")
	}
	auth.SigningKey = env["JWT_KEY"]
	if auth.SigningKey == "" {
		log.Fatal("Signing key can't be empty")
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

	log.Printf("Connected to DB at [%s]", dbAddr)

	app := application{
		conf: config,
		db:   db,
	}

	log.Printf("Server open at [%s]", serverAddr)

	log.Fatal(app.run(app.mount()))

}

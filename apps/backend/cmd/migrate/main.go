package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL обязательна")
	}

	dir := flag.String("dir", "migrations", "каталог миграций")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("укажите команду: up, down, status, create")
	}
	command := args[0]

	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("подключение к БД: %v", err)
	}
	defer sqlDB.Close()

	if err := goose.Run(command, sqlDB, *dir, args[1:]...); err != nil {
		log.Fatalf("goose %s: %v", command, err)
	}

	fmt.Printf("goose %s: выполнено\n", command)
}

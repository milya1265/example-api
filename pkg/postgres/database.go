package postgres

import (
	"database/sql"
	"example1/config"
	"example1/internal/logger"
	"fmt"
	_ "github.com/lib/pq"
)

func NewDatabaseClient(sc config.StorageConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		sc.Username, sc.Password, sc.Host, sc.Port, sc.Database, sc.SSLMode)
	fmt.Println(dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	return db, nil
}

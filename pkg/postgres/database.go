package postgres

import (
	"database/sql"
	"example1/config"
	"example1/pkg/logger"
	"fmt"
	_ "github.com/lib/pq"
)

func NewDatabaseClient(sc config.StorageConfig) (*sql.DB, error) {
	log := logger.Get()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		sc.Username, sc.Password, sc.Host, sc.Port, sc.Database, sc.SSLMode)

	log.Debug("dsn --> ", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

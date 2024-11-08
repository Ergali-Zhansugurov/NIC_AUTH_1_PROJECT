package postgres

import (
	"awesomeProject4/user-auth-service/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	DB *sqlx.DB
}

// ConnectPostgres создает подключение к PostgreSQL
func ConnectPostgres(cfg *config.Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf("postgres://%s", cfg.DBUri)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{DB: db}, nil
}

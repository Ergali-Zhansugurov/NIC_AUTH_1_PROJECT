package postgres

import (
	"awesomeProject4/user-auth-service/internal/config"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	DB *sqlx.DB
}

// ConnectPostgres создает подключение к PostgreSQL
func ConnectPostgres(cfg *config.Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf("postgres://postgres:%v@%v:%v/%v", cfg.DBPassword, cfg.DBHOST, cfg.DBPort, cfg.DBNAME)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{DB: db}, nil
}

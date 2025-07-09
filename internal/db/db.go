package db

import (
	"database/sql"
	"fmt"
	"socialNetworkOtus/internal/config"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
)

func NewDatabase(cfg *config.Config) (*goqu.Database, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	sqldb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return goqu.New("postgres", sqldb), nil
}

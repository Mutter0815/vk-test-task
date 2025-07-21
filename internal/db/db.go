package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Mutter0815/marketplace/internal/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func Connect(cfg *config.Config, log *zap.SugaredLogger) (*sql.DB, error) {
	log.Infow("connecting to database", "host", cfg.DBHost)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Errorw("sql.Open failed", "err", err)
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err == nil {
			log.Info("database connection established")
			return db, nil
		}
		time.Sleep(1 * time.Second)
	}
	log.Errorw("db.Ping timeout", "err", err)
	return nil, fmt.Errorf("db.Ping timeout: %w", err)
}

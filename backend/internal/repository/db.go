package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"time"

	"backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

type Tx = pgx.Tx

func NewDB(cfg *config.Config) *DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	log.Println("DB connected")

	return &DB{
		Pool: pool,
	}
}

func (db *DB) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return db.Pool.Begin(ctx)
}

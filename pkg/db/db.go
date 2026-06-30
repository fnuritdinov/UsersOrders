package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func New(o Options) (*pgxpool.Pool, error) {
	/*
		pool, err := pgxpool.New(context.Background(), fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			o.Host, o.Port, o.User, o.Password, o.DBName))
		if err != nil {
			return nil, err
		}

		err = pool.Ping(context.Background())
		if err != nil {
			return nil, err
		}

		fmt.Println("DSN =", pool)

		return pool, err

	*/

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		o.Host, o.Port, o.User, o.Password, o.DBName,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

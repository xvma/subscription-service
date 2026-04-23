package database

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
    log "github.com/sirupsen/logrus"
)

type PostgresDB struct {
    DB *sql.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
    log.WithField("dsn", dsn).Info("Connecting to PostgreSQL database")
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.WithError(err).Error("Failed to open database connection")
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    if err := db.Ping(); err != nil {
        log.WithError(err).Error("Failed to ping database")
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    log.Info("Successfully connected to PostgreSQL database")
    return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
    log.Info("Closing database connection")
    return p.DB.Close()
}
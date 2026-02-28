package db

import (
	"context"
	"fmt"

	"my-voice-billing/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Manager struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func (m *Manager) Reader() *sqlx.DB {
	if m.reader != nil {
		return m.reader
	}
	return m.writer
}

func (m *Manager) Writer() *sqlx.DB {
	return m.writer
}

func (m *Manager) Connect(ctx context.Context, cfg *config.Config) error {
	writeConn, ok := cfg.Connector("write")
	if !ok {
		return fmt.Errorf("database connector 'write' not found")
	}
	writeDSN := buildDSN(writeConn)
	writer, err := sqlx.Connect("pgx", writeDSN)
	if err != nil {
		return fmt.Errorf("connect write: %w", err)
	}
	if err := writer.PingContext(ctx); err != nil {
		_ = writer.Close()
		return fmt.Errorf("ping write: %w", err)
	}
	m.writer = writer

	readConn, ok := cfg.Connector("read")
	if ok && (readConn.Host != writeConn.Host || readConn.Port != writeConn.Port || readConn.Database != writeConn.Database) {
		readDSN := buildDSN(readConn)
		reader, err := sqlx.Connect("pgx", readDSN)
		if err != nil {
			_ = m.writer.Close()
			return fmt.Errorf("connect read: %w", err)
		}
		if err := reader.PingContext(ctx); err != nil {
			_ = reader.Close()
			_ = m.writer.Close()
			return fmt.Errorf("ping read: %w", err)
		}
		m.reader = reader
	}
	return nil
}

func (m *Manager) Close() {
	if m.reader != nil {
		_ = m.reader.Close()
	}
	if m.writer != nil {
		_ = m.writer.Close()
	}
}

func buildDSN(c config.DBConnector) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

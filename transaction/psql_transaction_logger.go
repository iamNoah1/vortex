package transaction

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}

type PsqlTransactionLogger struct {
	db *sql.DB
}

func NewPsqlTransactionLogger(config PostgresDBParams, tableName string) (TransactionLogger, error) {

	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", config.Host, config.DbName, config.User, config.Password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping() // Test the database connection
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PsqlTransactionLogger{db: db}

	exists, err := logger.verifyTableExists(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = logger.createTable(tableName); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}

func (p *PsqlTransactionLogger) verifyTableExists(tableName string) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_name = $1
            AND table_schema = 'public'
        );
    `

	var exists bool
	err := p.db.QueryRow(query, tableName).Scan(&exists)
	return exists, err
}

func (p *PsqlTransactionLogger) createTable(tableName string) error {
	return nil
}

func (p *PsqlTransactionLogger) Put(key string, value string) error {
	return nil
}

func (p *PsqlTransactionLogger) Delete(key string) error {
	return nil
}

func (p *PsqlTransactionLogger) ReplayEvents() error {
	return nil
}

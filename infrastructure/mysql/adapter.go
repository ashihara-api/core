package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/usk81/aveo"
)

type (
	// DBType is the data type to define database type
	DBType string
)

const (
	connectOptions = "charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo"

	// DriverName ...
	DriverName = "mysql"

	// TypeReader means to access reader DB
	TypeReader DBType = "READER"
	// TypeWriter means to access writer DB
	TypeWriter DBType = "WRITER"
)

var (
	readerConnection *sql.DB
	writerConnection *sql.DB
)

// ConnectReader creates database connection for reader DB
func ConnectReader(ctx context.Context, env aveo.Env) (*sql.DB, error) {
	conn, err := newDB(ctx, env, TypeReader)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("fail to connect reader DB. %w", err)
	}
	readerConnection = conn
	return readerConnection, nil
}

// ConnectWriter creates database connection for writer DB
func ConnectWriter(ctx context.Context, env aveo.Env) (*sql.DB, error) {
	conn, err := newDB(ctx, env, TypeWriter)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("fail to connect writer DB. %w", err)
	}
	writerConnection = conn
	return writerConnection, nil
}

// PingDB ...
func PingDB(dt DBType) error {
	switch dt {
	case TypeReader:
		return readerConnection.Ping()
	case TypeWriter:
		return writerConnection.Ping()
	}
	return fmt.Errorf("invalid database type: %s", dt)
}

// Close a database connection
func Close(dt DBType) error {
	switch dt {
	case TypeReader:
		return readerConnection.Close()
	case TypeWriter:
		return writerConnection.Close()
	}
	return fmt.Errorf("invalid database type: %s", dt)
}

// CloseAll closes all of database connections
func CloseAll() error {
	er := Close(TypeReader)
	ew := Close(TypeWriter)
	if er != nil || ew != nil {
		if ew == nil {
			// fail reader only
			return fmt.Errorf("fail to close connection. reader : %w", er)
		}
		if er == nil {
			// fail writer only
			return fmt.Errorf("fail to close connection.  writer : %w", ew)
		}
		// fail both
		return fmt.Errorf("fail to close connection. reader : %s / writer : %s", er.Error(), ew.Error())
	}
	return nil
}

func newDB(ctx context.Context, env aveo.Env, dt DBType) (*sql.DB, error) {
	cf := NewConfig(ctx, env, string(dt))
	conn, err := sql.Open("mysql", getConnectionString(&cf))
	if err != nil {
		return nil, err
	}
	appendConnSetting(conn, &cf)
	return conn, nil
}

func appendConnSetting(db *sql.DB, cf *Config) {
	db.SetConnMaxLifetime(time.Hour * 24)
	db.SetMaxIdleConns(cf.DBMaxIdleConnections)
	db.SetMaxOpenConns(cf.DBMaxOpenConnections)
}

func getConnectionString(cf *Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cf.DBUser,
		cf.DBPassword,
		cf.DBHost,
		cf.DBPort,
		cf.DBName,
		connectOptions,
	)
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

const defaultPingTimeout = 5 * time.Second

var (
	db *sql.DB
	mu sync.RWMutex
)

type Config struct {
	DSN             string
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	PingTimeout     time.Duration
}

func Init(ctx context.Context) error {
	cfg := ConfigFromApp()
	database, err := Open(ctx, cfg)
	if err != nil {
		return err
	}

	mu.Lock()
	if db != nil {
		_ = db.Close()
	}
	db = database
	mu.Unlock()

	return nil
}

func Open(ctx context.Context, cfg Config) (*sql.DB, error) {
	dsn, err := cfg.ConnectionString()
	if err != nil {
		return nil, err
	}

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		database.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		database.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		database.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	pingTimeout := cfg.PingTimeout
	if pingTimeout <= 0 {
		pingTimeout = defaultPingTimeout
	}

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := database.PingContext(pingCtx); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return database, nil
}

func DB() (*sql.DB, error) {
	mu.RLock()
	defer mu.RUnlock()

	if db == nil {
		return nil, errors.New("postgres database is not initialized")
	}

	return db, nil
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if db == nil {
		return nil
	}

	err := db.Close()
	db = nil
	return err
}

func ConfigFromApp() Config {
	return Config{
		DSN:             beego.AppConfig.DefaultString("ordering_postgres_dsn", ""),
		Host:            beego.AppConfig.DefaultString("ordering_postgres_host", "localhost"),
		Port:            beego.AppConfig.DefaultInt("ordering_postgres_port", 5432),
		User:            beego.AppConfig.DefaultString("ordering_postgres_user", "postgres"),
		Password:        beego.AppConfig.DefaultString("ordering_postgres_password", "postgres"),
		Database:        beego.AppConfig.DefaultString("ordering_postgres_database", "firstbeegoapi"),
		SSLMode:         beego.AppConfig.DefaultString("ordering_postgres_sslmode", "disable"),
		MaxOpenConns:    beego.AppConfig.DefaultInt("ordering_postgres_max_open_conns", 10),
		MaxIdleConns:    beego.AppConfig.DefaultInt("ordering_postgres_max_idle_conns", 5),
		ConnMaxLifetime: time.Duration(beego.AppConfig.DefaultInt("ordering_postgres_conn_max_lifetime_seconds", 300)) * time.Second,
		PingTimeout:     time.Duration(beego.AppConfig.DefaultInt("ordering_postgres_ping_timeout_seconds", 5)) * time.Second,
	}
}

func (c Config) ConnectionString() (string, error) {
	if c.DSN != "" {
		return c.DSN, nil
	}

	if c.Host == "" {
		return "", errors.New("ordering_postgres_host is required")
	}
	if c.User == "" {
		return "", errors.New("ordering_postgres_user is required")
	}
	if c.Database == "" {
		return "", errors.New("ordering_postgres_database is required")
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Database,
	}

	q := u.Query()
	q.Set("sslmode", c.SSLMode)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

package apiserver

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/gorilla/sessions"
	pgstore "github.com/honyshyota/constanta-rest-api/internal/app/store/pg"
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL, config)
	if err != nil {
		return err
	}

	defer db.Close()
	store := pgstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string, config *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if db != nil {
		logrus.Println("Running PostgreSQL main migrations")
		if err := runPgMigrations(config, config.PgURL); err != nil {
			return nil, fmt.Errorf("runPgMigrations main failed: %w", err)
		}
		logrus.Println("Main migrations done")

		logrus.Println("Running PostgreSQL test migrations")
		if err := runPgMigrations(config, config.PgTest); err != nil {
			return nil, fmt.Errorf("runPgMigrations test failed: %w", err)
		}
		logrus.Println("Test migrations done")
	}

	return db, nil
}

func runPgMigrations(config *Config, pgURL string) error {
	if config.PgMigrationsPath == "" {
		return nil
	}

	if pgURL == "" {
		return errors.New("no cfg.PgURL provided")
	}

	m, err := migrate.New(
		config.PgMigrationsPath,
		pgURL,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

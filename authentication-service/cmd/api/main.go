package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/getsentry/sentry-go"

	"authentication/data"
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int64

type Config struct {
	Repo        data.UserRepository
	Port        string `env:"WEB_PORT"`
	PostgresDSN string `env:"POSTGRES_DSN"`
	SentryDSN   string `env:"SENTRY_DSN"`
}

func main() {
	app := Config{}
	if err := env.Parse(&app); err != nil {
		log.Fatal(err.Error())
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              app.SentryDSN,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	conn := connectToDB(app.PostgresDSN)
	if conn == nil {
		log.Panic("Can't connect to postgres")
	}
	app.setupRepo(conn)

	log.Println("Starting service at port ", app.Port)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", app.Port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB(dsn string) *sql.DB {
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready..")
			counts++
		} else {
			log.Println("Connected to postgres")
			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds..")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}

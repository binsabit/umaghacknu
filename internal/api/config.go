package api

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/binsabit/umaghacknu/internal/repository/data"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		host     string
		port     int
		password string
		user     string
		dbname   string
	}
}

type application struct {
	config config
	models data.Models
}

func configure() config {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "indicate the port that server will be running on")
	flag.StringVar(&cfg.db.host, "dbhost", "localhost", "indicate host of psql database")
	flag.IntVar(&cfg.db.port, "dbport", 5432, "indicate port of psql database")
	flag.StringVar(&cfg.db.host, "dbuser", "hiidoskd", "indicate user of psql database")
	flag.StringVar(&cfg.db.host, "dbname", "umag", "indicate database name of psql database")
	flag.StringVar(&cfg.db.host, "password", "", "indicate database password of psql database")

	//flag.StringVar(&cfg.db.dsn, "dsn", "postgres://hiidoskd:@localhost:5432/umag", "")
	flag.Parse()

	return cfg
}

func StartAndConfigure() {
	cfg := configure()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	app := application{
		config: cfg,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  1 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		ReadTimeout:  10 * time.Minute,
	}
	log.Println("server has started")

	log.Fatal(srv.ListenAndServe())

}

// "host=localhost port=5432 user=hiidoskd dbname=umag sslmode=disable"
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s  sslmode=disable", cfg.db.host, cfg.db.port, cfg.db.user, cfg.db.password, cfg.db.dbname))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

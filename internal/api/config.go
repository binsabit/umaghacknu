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
		dsn string
	}
}

type application struct {
	config config
	models data.Models
}

func configure() config {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "indicate the port that server will be running on")
	//flag.StringVar(&cfg.db.dsn, "dsn", "postgres://hiidoskd:@localhost:5432/umag", "")
	flag.Parse()
	
	cfg.db.dsn =  "host=localhost port=5432 user=hiidoskd dbname=umag sslmode=disable"
	return cfg
}

func StartAndConfigure() {
	cfg := configure()

	
	db, err := openDB(cfg)
	if err != nil{
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

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
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

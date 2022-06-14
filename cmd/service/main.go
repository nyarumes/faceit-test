package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Nyarum/faceit-test/cmd/service/user"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const (
	defaultDBName          = "service.db"
	defaultRPCPort         = ":8082"
	defaultHTTPPort        = ":8083"
	defaultHealthCheckPort = ":80"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Debug().Msg("Starting service...")

	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("can't init bolt database")
	}
	defer db.Close()

	userRepository, err := user.NewRepository(db)
	if err != nil {
		log.Fatal().Err(err).Msg("can't init user repository")
	}

	userHandler := user.NewHandler(userRepository)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return user.RpcServer(defaultRPCPort, userHandler)
	})

	g.Go(func() error {
		return user.HttpServer(defaultHTTPPort, userHandler)
	})

	g.Go(func() error {
		return httpHealthCheckServer(defaultHealthCheckPort)
	})

	err = g.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("can't run one of handlers (http / rpc)")
	}
}

// httpHealthCheckServer declares http server for health checker
func httpHealthCheckServer(port string) error {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	log.Debug().Msg("Health check is starting")

	err := http.ListenAndServe(port, r)
	if err != nil {
		return err
	}

	return nil
}

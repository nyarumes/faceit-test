package main

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/Nyarum/faceit-test/cmd/service/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const (
	defaultDBName   = "service.db"
	defaultRPCPort  = ":8082"
	defaultHTTPPort = ":8083"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Debug().Msg("Starting service...")

	userRepository, err := user.NewRepository()
	if err != nil {
		log.Fatal().Err(err).Msg("can't init user repository")
	}

	userHandler := user.NewHandler(userRepository)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return rpcServer(defaultRPCPort, userHandler)
	})

	g.Go(func() error {
		return httpServer(defaultHTTPPort, userHandler)
	})

	g.Go(func() error {
		return healthCheck(":80")
	})

	err = g.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("can't run one of handlers (http / rpc)")
	}
}

// rpcServer delcares RPC server
func rpcServer(port string, handler *user.Handler) error {
	err := rpc.Register(handler)
	if err != nil {
		return err
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	err = http.Serve(listener, nil)
	if err != nil {
		return err
	}

	return nil
}

// httpServer declares http server
func httpServer(port string, handler *user.Handler) error {
	r := mux.NewRouter()
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		handler.FindAllWithPaginationAndFilter()
	}).Methods("GET")
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		handler.Add()
	}).Methods("POST")
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		handler.Update()
	}).Methods("PUT")
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		handler.Remove()
	}).Methods("DELETE")

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		return err
	}

	return nil
}

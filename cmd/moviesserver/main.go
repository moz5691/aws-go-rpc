package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/moz5691/twirp-lambda/internal/moviesserver"
	"github.com/moz5691/twirp-lambda/rpc/movies"

	"github.com/apex/gateway"
	"github.com/rs/cors"
)

func main() {
	//Dump all env vars
	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}

	mux := http.NewServeMux()

	server := &moviesserver.Server{}

	moviesHandler := movies.NewMoviesServer(server, nil)
	mux.Handle(movies.MoviesPathPrefix, moviesHandler)

	listenPort, exists := os.LookupEnv("LISTEN_PORT")
	if !exists {
		listenPort = "8080"
	}
	fmt.Println(listenPort)

	handler := cors.AllowAll().Handler(mux)

	// Check if it is in Lambda env.
	awsExecutionEnv, exists := os.LookupEnv("AWS_EXECUTION_ENV")
	if !exists {
		fmt.Println("listen port", listenPort)
		http.ListenAndServe(":"+listenPort, handler)
	} else {
		fmt.Println(awsExecutionEnv)
		log.Fatal(gateway.ListenAndServe("", mux))
	}
}

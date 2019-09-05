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

	// You can use any mux you like
	mux := http.NewServeMux()

	server := &moviesserver.Server{}

	moviesHandler := movies.NewMoviesServer(server, nil)
	mux.Handle(movies.MoviesPathPrefix, moviesHandler)

	// Un-comment below to test locally
	listenPort, exists := os.LookupEnv("LISTEN_PORT")
	if !exists {
		listenPort = "8080"
	}
	fmt.Println(listenPort)

	// log.Print("Listening on " + listenPort + " in stage " + appStage + " docker image: --CodeImage--")
	//Comment out below if serving over lambda

	handler := cors.AllowAll().Handler(mux)

	// http.ListenAndServe(":"+listenPort, mux)
	// http.ListenAndServe(":8080", mux)

	// Un-comment below before deploying to Lambda

	awsExecutionEnv, exists := os.LookupEnv("AWS_EXECUTION_ENV")
	if !exists {
		fmt.Println("listen port", listenPort)
		http.ListenAndServe(":8080", handler)
	} else {
		fmt.Println(awsExecutionEnv)
		log.Fatal(gateway.ListenAndServe("", mux))
	}
}

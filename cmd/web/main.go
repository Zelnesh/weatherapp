package main

import (

	"net/http"
	"flag"
	"log"
	"os"
)



type application struct {
	logInfo *log.Logger
	logError *log.Logger
}


func main(){

	port := flag.String("port", ":4040", "HTTP server address")
	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = ":" + envPort
	}

	logInfo := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	logError := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application {
		logInfo: logInfo,
		logError: logError,
	}


	srv := &http.Server {
		Addr: *port,
		ErrorLog: logError,
		Handler: app.routes(),
	}


	logInfo.Printf("Starting app on %s port", *port)

	err := srv.ListenAndServe()
	logError.Fatal(err)

}

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"simple-web-server/handlers"
	"simple-web-server/utils"
)

func main() {
	port := flag.String("p", "8080", "port to serve on")
	directory := flag.String("d", "files", "the directory of static file to host")
	logLevel := flag.String("l", "DEBUG", "log level: debug, info, warning, error")
	flag.Parse()

	logHandler, err := utils.NewLogHandler("json", os.Stdout, &slog.HandlerOptions{Level: utils.LogLevelMap[strings.ToUpper(*logLevel)]})
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	logger := utils.NewLogger(logHandler)

	// Create a new mux
	mux := http.NewServeMux()

	// Register a simple handler
	indexHandler := handlers.NewIndexHandler(logger)
	mux.HandleFunc("/", indexHandler.ServeHTTP)
	fs := handlers.NewCustomFileServer(*directory, logger)
	mux.HandleFunc("/files", fs.ServeHTTP)
	mux.HandleFunc("/files/", fs.ServeHTTP)

	// Start the server
	logger.Info(fmt.Sprintf("Starting server on %s", *port))
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		logger.Error("Could not start server: %s", err)
	}
}

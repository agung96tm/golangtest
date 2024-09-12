package main

import (
	"fmt"
	"net/http"
)

func (app *Application) serve() error {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.Port),
		Handler: app.routes(),
	}

	app.logger.Printf("Starting server on port :%d", app.config.Port)

	return srv.ListenAndServe()
}

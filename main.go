package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	digest "github.com/FeNoMeNa/goha"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
)

func init() {
	log.SetFormatter(new(log.JSONFormatter))
}

func main() {
	app := cli.App("genres-transformer", "A RESTful API for transforming TME genres to UP json")
	username := app.String(cli.StringOpt{
		Name:   "structure-service-username",
		Value:  "",
		Desc:   "Structure service username used for http digest authentication",
		EnvVar: "STRUCTURE_SERVICE_USERNAME",
	})
	password := app.String(cli.StringOpt{
		Name:   "structure-service-password",
		Value:  "",
		Desc:   "Structure service password used for http digest authentication",
		EnvVar: "STRUCTURE_SERVICE_PASSWORD",
	})
	principalHeader := app.String(cli.StringOpt{
		Name:   "principal-header",
		Value:  "",
		Desc:   "Structure service principal header used for authentication",
		EnvVar: "PRINCIPAL_HEADER",
	})
	baseURL := app.String(cli.StringOpt{
		Name:   "base-url",
		Value:  "http://localhost:8080/transformers/genres/",
		Desc:   "Base url",
		EnvVar: "BASE_URL",
	})
	structureServiceBaseURL := app.String(cli.StringOpt{
		Name:   "structure-service-base-url",
		Value:  "http://metadata.internal.ft.com:83",
		Desc:   "Structure service base url",
		EnvVar: "STRUCTURE_SERVICE_BASE_URL",
	})
	port := app.Int(cli.IntOpt{
		Name:   "port",
		Value:  8080,
		Desc:   "Port to listen on",
		EnvVar: "PORT",
	})

	app.Action = func() {
		c := digest.NewClient(*username, *password)
		c.Timeout(10 * time.Second)
		s, err := newGenreService(newTmeRepository(c, *structureServiceBaseURL, *principalHeader), *baseURL)
		if err != nil {
			log.Errorf("Error while creating GenreService: [%v]", err.Error())
		}
		h := newGenresHandler(s)
		m := mux.NewRouter()
		m.HandleFunc("/transformers/genres", h.getGenres).Methods("GET")
		m.HandleFunc("/transformers/genres/{uuid}", h.getGenreByUUID).Methods("GET")
		http.Handle("/", m)

		log.Printf("listening on %d", *port)
		http.ListenAndServe(fmt.Sprintf(":%d", *port),
			httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry,
				httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), m)))
	}
	app.Run(os.Args)
}
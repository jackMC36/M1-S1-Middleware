package main

import (
	agendas "middleware/example/internal/controllers/agendas"
	events "middleware/example/internal/controllers/events"
	"middleware/example/internal/helpers"
	_ "middleware/example/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()

	r.Route("/agendas", func(r chi.Router) { // route /agendas
		r.Get("/", agendas.GetAgendas)        // GET /agendas
		r.Route("/{id}", func(r chi.Router) { // route /agendas/{id}
			r.Use(agendas.Context)        // Use Context method to get agenda ID
			r.Get("/", agendas.GetAgenda) // GET /agendas/{id}

			// Events for the agenda
			r.Route("/events", func(r chi.Router) {
				r.Get("/", events.GetEvents) // GET /agendas/{id}/events
			})
		})
	})

		 
	r.Route("/events", func(r chi.Router) {// route /events
		 r.Route("/{id}", func(r chi.Router) { // route /events/{id}
			 r.Use(events.Context) 	  // Use Context method to get event ID
			 r.Get("/", events.Getevent) // GET /events/{id}
		 })
	 })

	logrus.Info("[INFO] Web server started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS agendas (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			ucaid VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS alertes (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL,
			agendaid VARCHAR(255),
			FOREIGN KEY (agendaid) REFERENCES agendas(id)
		);`,

		`CREATE TABLE IF NOT EXISTS events (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			uid VARCHAR(255) NOT NULL UNIQUE,
			description VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			start DATETIME NOT NULL,
			"end" DATETIME NOT NULL,
			location VARCHAR(255) NOT NULL,
			last_update DATETIME
		);`,

		`CREATE TABLE IF NOT EXISTS events_agendas (
			event_id VARCHAR(255) NOT NULL,
			agenda_id VARCHAR(255) NOT NULL,
			PRIMARY KEY (event_id,agenda_id),
			FOREIGN KEY (event_id) REFERENCES events(id),
			FOREIGN KEY (agenda_id) REFERENCES agendas(id)
		);`,


	}

	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}

package main

import (
	events "middleware/example/internal/controllers/events"
	agendas "middleware/example/internal/controllers/agendas"
	"middleware/example/internal/helpers"
	_ "middleware/example/internal/models"
	"net/http"

	events_consumers "middleware/example/internal/consumers/events"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()

	r.Route("/agendas", func(r chi.Router) { // route /agendas
  		r.Route("/{id}/events", func(r chi.Router) { // route /agendas/{id}/events
    		r.Use(agendas.Context)
    		r.Get("/", events.GetEventsByAgendaID)
		})
	})
		 
	r.Route("/events", func(r chi.Router) {// route /events
		 r.Route("/{id}", func(r chi.Router) { // route /events/{id }
			 r.Use(events.Context) 	  // Use Context method to get event ID
			 r.Get("/", events.Getevent) // GET /events/{id}
		 })
		 r.Get("/", events.GetEvents) // GET /events
	 })

	if err := helpers.InitNats(); err != nil {
        logrus.Fatal(err)
    }

	if err := helpers.EnsureStreams(); err != nil {
		logrus.Fatal(err)
	}

	go func() {
        consumer, err := events_consumers.EventConsumer()
        if err != nil {
            logrus.Warnf("error during nats consumer creation : %v", err)
            return
        }
        if err := events_consumers.Consume(*consumer); err != nil {
            logrus.Warnf("error during nats consume : %v", err)
        }
    }() 

	logrus.Info("[INFO] Web server started. Now listening on *:8081")
	logrus.Fatalln(http.ListenAndServe(":8081", r))
}

func init() {

	// create events tables in a separate events.db
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening events database : %s", err.Error())
	}
	schemes := []string{
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
			FOREIGN KEY (event_id) REFERENCES events(id)
		);`,
	}

	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate events table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}

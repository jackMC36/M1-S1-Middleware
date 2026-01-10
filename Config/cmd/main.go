package main

import (
	agendas "middleware/example/internal/controllers/agendas"
	"middleware/example/internal/controllers/alertes"
	"middleware/example/internal/helpers"
	_ "middleware/example/internal/models"
	"net/http"

	alertes_consumers "middleware/example/internal/consumers"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()

	r.Route("/agendas", func(r chi.Router) { // route /agendas
		r.Get("/", agendas.GetAgendas)        // GET /agendas
		r.Post("/", agendas.PostNewAgenda)    // POST /agendas
		r.Route("/{id}", func(r chi.Router) { // route /agendas/{id}
			r.Use(agendas.Context)              // Use Context method to get agenda ID
			r.Get("/", agendas.GetAgenda)       // GET /agendas/{id}
			r.Put("/", agendas.PutAgenda)       // PUT /agendas/{id}
			r.Delete("/", agendas.DeleteAgenda) // DELETE /agendas/{id}
		})
	})

	r.Route("/alertes", func(r chi.Router) { // route /alertes
		r.Get("/", alertes.GetAlertes)        // GET /agendas
		r.Post("/", alertes.PostNewAlerte)    // POST /alertes
		r.Route("/{id}", func(r chi.Router) { // route /alertes/{id}
			r.Use(alertes.Context)              // Use Context method to get alertes ID
			r.Get("/", alertes.GetAlerte)       // GET /alertes/{id}
			r.Put("/", alertes.PutAlerte)       // PUT /alertes/{id}
			r.Delete("/", alertes.DeleteAlerte) // DELETE /alertes/{id}
		})
	})

	logrus.Info("[INFO] Web server started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))

	go func() {
		consumer, err := alertes_consumers.AlertConsumer()
		if err != nil {
			logrus.Warnf("error during nats consumer creation : %v", err)
			return
		}
		if err := alertes_consumers.Consume(*consumer); err != nil {
			logrus.Warnf("error during nats consume : %v", err)
		}
	}()
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
	}

	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}

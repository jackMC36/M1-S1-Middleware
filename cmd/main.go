package main

import (
	agendas "middleware/example/internal/controllers/agendas"
	"middleware/example/internal/helpers"
	_ "middleware/example/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()

	r.Route("/agendas", func(r chi.Router) {
		r.Get("/", agendas.GetAgendas)
		r.Post("/", agendas.PostNewAgenda)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(agendas.Context)
			r.Get("/", agendas.GetAgenda)
		})
	})

	// Add a simple HTML form for testing
	r.Get("/form", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Create Agenda</title>
</head>
<body>
    <h1>Create New Agenda</h1>
    <form id="agendaForm">
        <label for="id">ID (UUID):</label><br>
        <input type="text" id="id" name="id" placeholder="550e8400-e29b-41d4-a716-446655440000"><br><br>
        
        <label for="name">Name:</label><br>
        <input type="text" id="name" name="name" placeholder="My Agenda"><br><br>
        
        <label for="ucaid">User ID (UUID):</label><br>
        <input type="text" id="ucaid" name="ucaid" placeholder="123e4567-e89b-12d3-a456-426614174000"><br><br>
        
        <button type="submit">Create Agenda</button>
    </form>

    <script>
        document.getElementById('agendaForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = {
                id: document.getElementById('id').value,
                name: document.getElementById('name').value,
                ucaid: document.getElementById('ucaid').value
            };

            try {
                const response = await fetch('/agendas', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });

                if (response.ok) {
                    const result = await response.json();
                    alert('Agenda created successfully: ' + JSON.stringify(result));
                } else {
                    alert('Error: ' + response.statusText);
                }
            } catch (error) {
                alert('Error: ' + error.message);
            }
        });
    </script>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
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
	}

	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}

	helpers.CloseDB(db)
}

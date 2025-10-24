package alertes

import (
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"

	"github.com/gofrs/uuid"
)

func GetAllAlertes() ([]models.Alerte, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM alertes")
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	// parsing datas in object slice
	alertes := []models.Alerte{}
	for rows.Next() {
		var data models.Alerte
		err = rows.Scan(&data.Id, &data.Email, &data.AgendaId)
		if err != nil {
			return nil, err
		}
		alertes = append(alertes, data)
	}
	// don't forget to close rows
	_ = rows.Close()

	return alertes, err
}

func GetAlerteById(id uuid.UUID) (*models.Alerte, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT * FROM alertes WHERE id=?", id.String())
	helpers.CloseDB(db)

	var alerte models.Alerte
	err = row.Scan(&alerte.Id, &alerte.Email, &alerte.AgendaId)
	if err != nil {
		return nil, err
	}
	return &alerte, err
}

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

	alertes := []models.Alerte{}
	for rows.Next() {
		var data models.Alerte
		err = rows.Scan(&data.Id, &data.Email, &data.AgendaId)
		if err != nil {
			return nil, err
		}
		alertes = append(alertes, data)
	}
	_ = rows.Close()

	return alertes, err
}

func GetAlertesByAgendaId(agendaId string) ([]models.Alerte, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM alertes WHERE agendaid = ?", agendaId)
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	alertes := []models.Alerte{}
	for rows.Next() {
		var data models.Alerte
		err = rows.Scan(&data.Id, &data.Email, &data.AgendaId)
		if err != nil {
			return nil, err
		}
		alertes = append(alertes, data)
	}
	_ = rows.Close()

	return alertes, nil
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

func DeleteAlerteById(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM alertes WHERE id=?", id.String())
	helpers.CloseDB(db)

	if err != nil {
		return err
	}
	return err
}
func PostNewAlerte(agendaId string, email string) (*models.Alerte, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	alerteId := uuid.Must(uuid.NewV4())

	_, err = db.Exec("INSERT INTO alertes (id, email, agendaid) VALUES (?, ?, ?)",
		alerteId.String(), email, agendaId)
	helpers.CloseDB(db)

	if err != nil {
		return nil, err
	}

	alerte := &models.Alerte{
		Id:       &alerteId,
		Email:    email,
		AgendaId: agendaId,
	}

	return alerte, nil
}

func UpdateAlerteById(id uuid.UUID, email string, agendaId string) (*models.Alerte, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("UPDATE alertes SET email = ?, agendaid = ? WHERE id = ?", email, agendaId, id.String())
	helpers.CloseDB(db)

	if err != nil {
		return nil, err
	}

	updatedAlerte := &models.Alerte{
		Id:       &id,
		Email:    email,
		AgendaId: agendaId,
	}

	return updatedAlerte, nil
}

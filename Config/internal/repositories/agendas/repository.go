package agendas

import (
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"

	"github.com/gofrs/uuid"
)

func GetAllAgendas() ([]models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM agendas")
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	agendas := []models.Agenda{}
	for rows.Next() {
		var data models.Agenda
		err = rows.Scan(&data.Id, &data.UcaId, &data.Name)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, data)
	}
	_ = rows.Close()

	return agendas, err
}

func GetAgendaById(id uuid.UUID) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT * FROM agendas WHERE id=?", id.String())
	helpers.CloseDB(db)

	var agenda models.Agenda
	err = row.Scan(&agenda.Id, &agenda.UcaId, &agenda.Name)
	if err != nil {
		return nil, err
	}
	return &agenda, err
}

func DeleteAgendaById(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}

	row := db.QueryRow("SELECT * FROM agendas WHERE id=?", id.String())
	var agenda models.Agenda
	err = row.Scan(&agenda.Id, &agenda.UcaId, &agenda.Name)
	if err != nil {
		helpers.CloseDB(db)
		return err
	}

	_, err = db.Exec("DELETE FROM agendas WHERE id=?", id.String())
	helpers.CloseDB(db)

	if err != nil {
		return err
	}
	return nil
}

func PostNewAgenda(name string, ucaid uuid.UUID) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err

	}
	id := uuid.Must(uuid.NewV4())

	_, err = db.Exec("INSERT INTO agendas (id, ucaid, name) VALUES (?, ?, ?)", id.String(), ucaid.String(), name)
	helpers.CloseDB(db)

	if err != nil {
		return nil, err
	}

	agenda := &models.Agenda{
		Id:    &id,
		Name:  name,
		UcaId: &ucaid,
	}
	return agenda, nil
}

func UpdateAgendaById(name string, ucaid uuid.UUID) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}

	id := uuid.Must(uuid.NewV4())

	_, err = db.Exec("UPDATE agendas SET name = ?, ucaid = ? WHERE id = ?", name, ucaid.String(), id.String())
	helpers.CloseDB(db)

	if err != nil {
		return nil, err
	}

	updatedAgenda := &models.Agenda{
		Id:    &id,
		Name:  name,
		UcaId: &ucaid,
	}

	return updatedAgenda, nil
}

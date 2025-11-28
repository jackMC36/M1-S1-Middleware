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
		err = rows.Scan(&data.Id, &data.Name)
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
	err = row.Scan(&agenda.Id, &agenda.Name, &agenda.UcaId)
	if err != nil {
		return nil, err
	}
	return &agenda, err
}

func DeleteAgendaById(id uuid.UUID) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("DELETE * FROM agendas WHERE id=?", id.String())
	helpers.CloseDB(db)

	var agenda models.Agenda
	err = row.Scan(&agenda.Id, &agenda.Name, &agenda.UcaId)
	if err != nil {
		return nil, err
	}
	return &agenda, err
}

func PostNewAgenda(id uuid.UUID, name string, ucaid uuid.UUID) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("INSERT into agendas * VALUES (?,?,?)", id.String(), name, ucaid.String())
	helpers.CloseDB(db)

	var agenda models.Agenda
	err = row.Scan(&agenda.Id, &agenda.Name, &agenda.UcaId)
	if err != nil {
		return nil, err
	}
	return &agenda, err
}

func ReplaceAgendaById(old_id uuid.UUID, new_id uuid.UUID, name string, ucaid uuid.UUID) (*models.Agenda, error) {
	DeleteAgendaById(old_id)
	agenda, err := PostNewAgenda(new_id, name, ucaid)
	if err != nil {
		return nil, err
	}
	return agenda, nil
}

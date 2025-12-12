package agendas

import (
	"database/sql"
	"fmt"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/agendas"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GetAllagendas() ([]models.Agenda, error) {
	var err error
	// calling repository
	agendas, err := repository.GetAllAgendas()
	if err != nil {
		logrus.Errorf("error retrieving agendas : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving agendas",
		}
	}

	return agendas, nil
}

func GetAgendaById(id uuid.UUID) (*models.Agenda, error) {
	agenda, err := repository.GetAgendaById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "agenda not found",
			}
		}
		logrus.Errorf("error retrieving agenda %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving agenda %s", id.String()),
		}
	}

	return agenda, err
}

func DeleteAgendaById(id uuid.UUID) error {
	err := repository.DeleteAgendaById(id)
	if err != nil {
		logrus.Errorf("error deleting agenda %s : %s", id.String(), err.Error())
		return &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while deleting agenda %s", id.String()),
		}
	}
	return nil
}

func PostNewAgenda(id uuid.UUID, name string, ucaid uuid.UUID) (*models.Agenda, error) {
	agenda, err := repository.PostNewAgenda(id, name, ucaid)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "agenda not found (sql.errNoRows)",
			}
		}
		logrus.Errorf("services/agendas: error posting agenda %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while posting agenda %s", id.String()),
		}
	}

	return agenda, err
}

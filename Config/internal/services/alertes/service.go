package alertes

import (
	"database/sql"
	"fmt"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/alertes"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GetAllAlertes() ([]models.Alerte, error) {
	var err error
	// calling repository
	alertes, err := repository.GetAllAlertes()
	if err != nil {
		logrus.Errorf("error retrieving alertes : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving alertes",
		}
	}

	return alertes, nil
}

func GetalerteById(id uuid.UUID) (*models.Alerte, error) {
	alerte, err := repository.GetAlerteById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "alerte not found",
			}
		}
		logrus.Errorf("error retrieving alerte %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving alerte %s", id.String()),
		}
	}

	return alerte, err
}

func DeleteAlerteById(id uuid.UUID) error {
	err := repository.DeleteAlerteById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &models.ErrorNotFound{
				Message: "alerte not found",
			}
		}
		logrus.Errorf("error deleting alerte %s : %s", id.String(), err.Error())
		return &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while deleting alerte %s", id.String()),
		}
	}

	return err
}

func PostNewAlerte(agendaId string, email string) (*models.Alerte, error) {
	if email == "" {
		return nil, &models.ErrorGeneric{
			Message: "Email is required",
		}
	}

	alerte, err := repository.PostNewAlerte(agendaId, email)
	if err != nil {
		logrus.Errorf("error creating alerte for agenda %s : %s", agendaId, err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while creating alerte",
		}
	}

	return alerte, nil
}

func UpdateAlerteById(id uuid.UUID, email string, agendaId string) (*models.Alerte, error) {
	if email == "" {
		return nil, &models.ErrorGeneric{
			Message: "Email is required",
		}
	}

	alerte, err := repository.UpdateAlerteById(id, email, agendaId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "alerte not found",
			}
		}
		logrus.Errorf("error updating alerte %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while updating alerte %s", id.String()),
		}
	}

	return alerte, nil
}

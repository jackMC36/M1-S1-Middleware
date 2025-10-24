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

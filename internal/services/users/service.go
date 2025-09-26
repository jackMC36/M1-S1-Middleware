package users

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/users"
)

func GetAllUsers() ([]models.User, error) {
	var err error
	// calling repository
	users, err := repository.GetAllUsers()
	// managing errors
	if err != nil {
		logrus.Errorf("error retrieving users : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving users",
		}
	}

	return users, nil
}

func GetUserById(id uuid.UUID) (*models.User, error) {
	user, err := repository.GetUserById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "user not found",
			}
		}
		logrus.Errorf("error retrieving user %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving user %s", id.String()),
		}
	}

	return user, err
}

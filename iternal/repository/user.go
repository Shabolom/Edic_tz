package repository

import (
	"Edos_Docer/config"
	"Edos_Docer/iternal/domain"
	"Edos_Docer/iternal/tools"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) Register(user domain.User) tools.BaseError {
	_, err, _ := ur.Find(user.Login)
	if err == nil {
		return tools.BaseError{
			Code:   http.StatusConflict,
			String: "такой пользователь уже существует",
		}
	}

	err = config.DB.
		Create(&user).
		Error
	if err != nil {
		return tools.BaseError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}
	return tools.BaseError{
		Code:   http.StatusCreated,
		String: "успешно зарегестрировались",
		Result: user,
	}
}

func (ur *UserRepo) Find(login string) (int, error, domain.User) {
	var user domain.User
	err := config.DB.Where("login =?", login).
		First(&user).
		Error

	if err == nil {
		return http.StatusOK, nil, user
	}
	return http.StatusUnauthorized, err, domain.User{}
}

func (ur *UserRepo) PostCSV(csvFile []domain.VideoInfo, settings domain.PostCSV) tools.BaseError {
	tx := config.DB.Begin()

	for _, csv := range csvFile {
		err := tx.
			Create(&csv).
			Error
		if err != nil {
			tx.Rollback()
			return tools.BaseError{
				Code: 500,
				Err:  err,
			}
		}
	}
	tx.Commit()

	err := config.DB.
		Create(&settings).
		Error
	if err != nil {
		return tools.BaseError{
			Code: 500,
			Err:  err,
		}
	}

	return tools.BaseError{
		Code:   http.StatusCreated,
		String: "успешно сохранено",
		Err:    nil,
	}
}

func (ur *UserRepo) SelfDelete(userID uuid.UUID) tools.BaseError {
	var user domain.User

	err := config.DB.
		Where("id =?", userID).
		Delete(&user).
		Error

	if err != nil {
		return tools.BaseError{
			Code: 500,
			Err:  err,
		}
	}

	return tools.BaseError{
		Code:   http.StatusOK,
		String: "вы удалены",
	}
}

func (ur *UserRepo) Swap(userID uuid.UUID, newLog domain.User) tools.BaseError {
	var user domain.User

	err := config.DB.
		Model(&user).
		Where("id =?", userID).
		Select("updated_at", "login", "password").
		Updates(domain.User{
			Base: domain.Base{
				UpdatedAt: time.Now(),
			},
			Login:    newLog.Login,
			Password: newLog.Password,
		}).
		Error

	if err != nil {
		return tools.BaseError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}

	return tools.BaseError{
		Code:   http.StatusOK,
		String: "логин и пароль изменены",
	}
}

func (ur *UserRepo) GetBatch(limit, skip int) (error, []domain.User) {
	var domainUsers []domain.User

	err := config.DB.
		Order("id asc").
		Limit(limit).
		Offset(skip).
		Find(&domainUsers).
		Error

	if err != nil {
		return err, []domain.User{}
	}

	return nil, domainUsers
}

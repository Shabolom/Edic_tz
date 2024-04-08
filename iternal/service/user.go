package service

import (
	"Edos_Docer/iternal/domain"
	"Edos_Docer/iternal/models"
	"Edos_Docer/iternal/repository"
	"Edos_Docer/iternal/tools"
	"github.com/gofrs/uuid"
	"mime/multipart"
	"net/http"
	"time"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

var userRepo = repository.NewUserRepo()

func (us *UserService) Register(user models.Register) tools.BaseError {
	id, _ := uuid.NewV4()
	perm := 0

	if user.Login == "admin" && user.Password == "admin" {
		perm = 3
	}

	password, err := tools.HashPassword(user.Password)
	if err != nil {
		return tools.BaseError{
			Code: http.StatusConflict,
			Err:  err,
		}
	}

	userEntity := domain.User{
		Login:    user.Login,
		Password: password,
		PermLVL:  perm,
	}
	userEntity.Base.ID = id

	result := userRepo.Register(userEntity)
	if result.Err != nil {
		return tools.BaseError{
			Code: result.Code,
			Err:  result.Err,
		}
	}
	return tools.BaseError{
		Code:   result.Code,
		String: result.String,
		Result: result.Result,
	}
}

func (us *UserService) Login(user models.Register) tools.BaseError {

	code, err, result := userRepo.Find(user.Login)
	if err != nil {
		return tools.BaseError{
			Code:   code,
			Err:    nil,
			String: "не верный логин или пароль",
			Result: result,
		}
	}
	if !tools.CheckPasswordHash(user.Password, result.Password) {
		return tools.BaseError{
			Code:   http.StatusUnauthorized,
			String: "не верный логин или пароль",
		}
	}
	return tools.BaseError{
		Code:   code,
		String: "успешно авторизировались",
		Err:    nil,
		Result: result,
	}
}

func (us *UserService) PostCSV(data []models.InfVidCsv, settings *multipart.FileHeader, adminID uuid.UUID) tools.BaseError {
	var massCsv []domain.VideoInfo

	for _, inf := range data {
		csvPart := domain.VideoInfo{
			VideoID:             inf.VideoID,
			TrendingDate:        inf.TrendingDate,
			Title:               inf.Title,
			ChannelTitle:        inf.ChannelTitle,
			CategoryId:          inf.CategoryId,
			PublishTime:         inf.PublishTime,
			Tags:                inf.Tags,
			Likes:               inf.Likes,
			Dislikes:            inf.Dislikes,
			CommentCount:        inf.CommentCount,
			ThumbnailLink:       inf.ThumbnailLink,
			CommentsDisabled:    inf.CommentsDisabled,
			RatingsDisabled:     inf.RatingsDisabled,
			VideoErrorOrRemoved: inf.VideoErrorOrRemoved,
			Description:         inf.Description,
		}
		massCsv = append(massCsv, csvPart)
	}

	csvSettings := domain.PostCSV{
		AdminID:   adminID,
		CreatedAt: time.Now(),
		Size:      settings.Size,
		FileName:  settings.Filename,
	}

	result := userRepo.PostCSV(massCsv, csvSettings)
	if result.Err != nil {
		return tools.BaseError{
			Code: result.Code,
			Err:  result.Err,
		}
	}
	return result
}

func (us *UserService) SelfDelete(userID uuid.UUID) tools.BaseError {
	result := userRepo.SelfDelete(userID)

	return result
}

func (us *UserService) Swap(userID uuid.UUID, newLogs models.Register) tools.BaseError {
	id, _ := uuid.NewV4()

	password, err := tools.HashPassword(newLogs.Password)
	if err != nil {
		return tools.BaseError{
			Code: 500,
			Err:  err,
		}
	}

	newLogEntity := domain.User{
		Login:    newLogs.Login,
		Password: password,
	}
	newLogEntity.Base.ID = id

	_, err, _ = userRepo.Find(newLogEntity.Login)
	if err == nil {
		return tools.BaseError{
			Code:   http.StatusBadRequest,
			String: "нужен другой логин",
		}
	}

	result := userRepo.Swap(userID, newLogEntity)

	return result
}

func (us *UserService) GetBatch(page, limit int) (error, []domain.User) {
	skip := page*limit - limit

	err, result := userRepo.GetBatch(limit, skip)
	if err != nil {
		return err, []domain.User{}
	}

	return nil, result
}

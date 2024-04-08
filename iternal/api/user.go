package api

import (
	"Edos_Docer/iternal/models"
	"Edos_Docer/iternal/service"
	"Edos_Docer/iternal/tools"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type UserAPI struct {
}

func NewUserApi() *UserAPI {
	return &UserAPI{}
}

var userService = service.NewUserService()

// @Summary	регистрация пользователя с выдачей токена
// @Produce	json
// @Accept	json
// @Tags	Authorization
// @Param	ввод	body		models.Register	true	"ввести логин и пароль"
// @Success	201		{string}	string	"вы зарегестрировались"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router		/api/register [post]
func (ua *UserAPI) Register(c *gin.Context) {
	var userRegister models.Register
	var token tools.BaseError

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err = json.Unmarshal(data, &userRegister)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.Register(userRegister)
	if result.Err != nil {
		tools.CreateError(result.Code, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	if result.Code != http.StatusConflict {
		token = tools.CreateToken(&result.Result)
		if token.Err != nil {
			tools.CreateError(token.Code, token.Err, c)
			log.WithField("component", "rest").Warn(err)
			return
		}
		c.Writer.Header().Set("Authorization", token.String)
	}

	c.String(result.Code, result.String)
}

// @Summary	авторизация с выдачей токена в хэдерсе
// @Produce	json
// @Accept	json
// @Tags	Authorization
// @Param	ввод	body		models.Register	true	"авторизация"
// @Success	200		{string}	string	"успешно авторизировались"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	401		{object}	models.Error
// @Router		/api/login [post]
func (ua *UserAPI) Login(c *gin.Context) {
	var userRegister models.Register
	var token tools.BaseError

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err = json.Unmarshal(data, &userRegister)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.Login(userRegister)
	if result.Err != nil {
		tools.CreateError(result.Code, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	if result.Code != http.StatusUnauthorized {
		token = tools.CreateToken(&result.Result)
		if token.Err != nil {
			tools.CreateError(token.Code, token.Err, c)
			log.WithField("component", "rest").Warn(token.Err)
			return
		}
		c.Writer.Header().Set("Authorization", token.String)
	}

	c.String(result.Code, result.String)
}

// @Summary	заносим данные из csv файла
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	ввод	formData		file	true	"вставьте файл"
// @Success	200		{string}	string	"успешно внесено"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	401		{object}	models.Error
// @Router		/api/post_csv [post]
func (ua *UserAPI) PostCSV(c *gin.Context) {
	var infoVideo []models.InfVidCsv

	PermLVl, adminID := tools.ParseToken(c)

	if PermLVl != 3 {
		tools.CreateError(http.StatusForbidden, errors.New("не достаточно прав"), c)
		log.WithField("component", "rest").Warn(errors.New("не достаточно прав"))
		return
	}

	file, settings, err := c.Request.FormFile("csv")
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	fmt.Println(settings.Header)

	err = gocsv.UnmarshalMultipartFile(&file, &infoVideo)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.PostCSV(infoVideo, settings, adminID)
	if result.Err != nil {
		tools.CreateError(result.Code, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.String(result.Code, result.String)
}

// @Summary	удаление аккаунта
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Success	200		{string}	string	"успешно удалено"
// @Failure	500		{object}	models.Error
// @Router		/api/delete_acc [delete]
func (ua *UserAPI) SelfDelete(c *gin.Context) {
	_, userID := tools.ParseToken(c)

	result := userService.SelfDelete(userID)
	if result.Err != nil {
		tools.CreateError(result.Code, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.String(result.Code, result.String)
}

// @Summary	изменяем логин и пароль
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	ввод	body		models.Register	true	"вставьте новый логин и пароль"
// @Success	200		{string}	string	"успешно"
// @Failure	500		{object}	models.Error
// @Router		/api/swap_login_password [post]
func (ua *UserAPI) Swap(c *gin.Context) {
	var newLogs models.Register

	_, userID := tools.ParseToken(c)

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err = json.Unmarshal(data, &newLogs)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.Swap(userID, newLogs)
	if result.Err != nil {
		tools.CreateError(result.Code, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.String(result.Code, result.String)
}

// @Summary	пагинация получение определенного количества элементов пропуская их в соответствии с страницей
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	page	query     	string		true  "введите страницу"
// @Param	limit	query     	string		true  "введите количество элементов"
// @Success	200		{string}	string	"успешно удалено"
// @Failure	500		{object}	models.Error
// @Router		/api/get_size_of_elements/:value [get]
func (ua *UserAPI) GetBatch(c *gin.Context) {
	perm, _ := tools.ParseToken(c)
	if perm < 3 {
		tools.CreateError(http.StatusForbidden, errors.New("у вас нет прав для данного действия"), c)
		log.WithField("component", "rest").Warn(errors.New("у вас нет прав для данного действия"))
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, errors.New("не валидные параметры"), c)
		log.WithField("component", "rest").Warn(errors.New("не валидные параметры"))
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, errors.New("не валидные параметры"), c)
		log.WithField("component", "rest").Warn(errors.New("не валидные параметры"))
		return
	}

	err, result := userService.GetBatch(page, limit)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

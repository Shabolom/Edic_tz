package tools

import (
	"Edos_Docer/config"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID  uuid.UUID
	PermLVl int
}

func CreateToken(user GetFields) BaseError {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Тимур",                                           //Эмитент: Организация или лицо, выдающее утверждение.
			Subject:   "user",                                            //Субъект: Лицо или объект, о котором делается утверждение.
			Audience:  nil,                                               //Аудитория: Целевая аудитория, для которой предназначено утверждение.
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)), //Не раньше: Дата и время, начиная с которых утверждение действительно.
			//NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)), //Истекает в: Дата и время истечения срока действия утверждения
			IssuedAt: jwt.NewNumericDate(time.Now()), //Выдано в: Дата и время выдачи утверждения
			ID:       "timurJWT",                     //Идентификатор JWT: Уникальный идентификатор утверждения.
		},
		UserID:  user.ID(),
		PermLVl: user.Permissions(),
	})

	strToken, err := token.SignedString([]byte(config.Env.SecretKey))

	if err != nil {
		return BaseError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return BaseError{
		Code:   http.StatusOK,
		String: "Bearer " + strToken,
		Err:    nil,
	}
}

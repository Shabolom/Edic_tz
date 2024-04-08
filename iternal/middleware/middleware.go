package middleware

import (
	"Edos_Docer/config"
	"Edos_Docer/iternal/tools"
	"compress/gzip"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func GZIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gxip" {
			reader, err := gzip.NewReader(c.Request.Body)
			defer reader.Close()
			if err != nil {
				log.WithField("component", "Gzip").Warn(err)
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				return
			}
			c.Request.Body = reader
		}
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// продолжаем работать с хэндлером который идет после мидлвейра  (который был вызван изначально))
		c.Next()

		latency := time.Since(t)
		log.WithField("component", "latency").Info(latency)
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") == "" {
			tools.CreateError(http.StatusUnauthorized, errors.New("You're Unauthorized"), c)
			return
		}

		strToken := c.Request.Header.Get("Authorization")

		token, err := jwt.Parse(strToken,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					tools.CreateError(http.StatusUnauthorized, errors.New("You're Unauthorized"), c)
				}
				return []byte(config.Env.SecretKey), nil
			})

		if err != nil {
			tools.CreateError(http.StatusUnauthorized, errors.New("You're Unauthorized"), c)
			return
		}

		if !token.Valid {
			tools.CreateError(http.StatusUnauthorized, errors.New("You're Unauthorized"), c)
			return
		}
		c.Next()
	}
}

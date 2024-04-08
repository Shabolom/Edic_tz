package main

import (
	"Edos_Docer/config"
	_ "Edos_Docer/docs"
	"Edos_Docer/iternal/routes"
	"Edos_Docer/iternal/tools"
	"Edos_Docer/migrate"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func main() {
	//	@title		User API
	//	@version	1.0.0

	// 	@description 	Это выпускной проэкт с использованием свагера и докуера
	// 	@termsOfService  сдесь были бы условия использования еслиб я их мог обозначить
	// 	@contact.url    тут моя контактная информация (https://vk.com/id192672036)
	// 	@contact.email  tima.gorenskiy@mail.ru

	// 	@securityDefinitions.apikey  ApiKeyAuth
	//  @in header
	//  @name Authorization

	//	@host		localhost:8008

	config.CheckFlagEnv()

	err := tools.InitLogger()
	if err != nil {
		fmt.Println("ошибка при инициализаии логера")
	}

	err = config.InitPgSQL()
	if err != nil {
		log.WithField("component", "initialization").Panic(err)
	}

	migrate.Migrate()

	r := routes.SetupRouter()

	// запуск сервера
	if err = r.Run(config.Env.Host + ":" + config.Env.Port); err != nil {
		log.WithField("component", "run").Panic(err)
	}
}

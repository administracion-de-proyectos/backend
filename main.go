package main

import (
	"backend-admin-proyect/src/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	dbUrl = "DB_URL"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	PanicOnError(godotenv.Load(".env"))
	url := os.Getenv(dbUrl)
	if url == "" {
		panic(fmt.Sprintf("missing url from db with format %s", dbUrl))
	}

	log.Infof("url is: %s", url)
	r := routes.Routes{
		Router: gin.Default(),
	}
	PanicOnError(r.AddUserRoutes(url))
	log.Infof("starting to run")
	PanicOnError(r.Router.Run(":8001"))
}

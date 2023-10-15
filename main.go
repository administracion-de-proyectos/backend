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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Errorf("env not found, not panicking but something to check out if not using okteto, %s", err.Error())
	}

	url := os.Getenv(dbUrl)
	if url == "" {
		panic(fmt.Sprintf("missing url from db with format %s", dbUrl))
	}

	log.Infof("url is: %s", url)
	r := routes.Routes{
		Router: gin.Default(),
	}
	r.Router.Use(CORSMiddleware())
	PanicOnError(r.AddUserRoutes(url))
	PanicOnError(r.AddSwaggerRoutes())
	log.Infof("starting to run")
	PanicOnError(r.Router.Run(":8001"))
}

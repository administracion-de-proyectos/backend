package routes

import (
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	"backend-admin-proyect/docs"
)

func (r Routes) AddSwaggerRoutes() error {
	docs.SwaggerInfo.Title = "Swagger Example API"
	group := r.Router.Group("/swagger")
	group.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return nil
}
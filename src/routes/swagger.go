package routes

import (
	"backend-admin-proyect/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func (r Routes) AddSwaggerRoutes() error {
	docs.SwaggerInfo.Title = "Swagger Example API"
	group := r.Router.Group("/swagger")
	group.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return nil
}

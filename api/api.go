package api

import (
	_ "app/api/docs"
	"app/api/handler"
	"app/config"
	"app/pkg/logger"
	"app/storage"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewApi(r *gin.Engine, cfg *config.Config, storage storage.StorageI, logger logger.LoggerI) {

	handler := handler.NewHandler(cfg, storage, logger)

	r.POST("/category", handler.CreateCategory)
	r.GET("/category/:id", handler.GetByIdCategory)
	r.GET("/category", handler.GetListCategory)
	r.PUT("/category",handler.UpdateCategory)
	r.DELETE("/category/:id",handler.DeleteCategory)

	r.POST("/product", handler.CreateProduct)
	r.GET("/product/:id", handler.GetProductByID)
	r.GET("/product", handler.GetListProduct)
	r.PUT("/product", handler.UpdateProduct)
	r.DELETE("/product/:id",handler.DeleteProduct)


	r.POST("/user", handler.CreateUser)
	r.GET("/user/:id", handler.GetUserByID)
	r.GET("/user", handler.GetListUser)
	r.PUT("/user", handler.UpdateUser)
	r.DELETE("/user/:id",handler.DeleteUser)



	url:=ginSwagger.URL("swagger/doc.json") //// The url pointing to API definition
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,url))
}

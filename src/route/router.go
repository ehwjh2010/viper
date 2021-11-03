package route

import (
	demoHandlers "ginLearn/src/controller/demo"
	projectHandlers "ginLearn/src/controller/project"
	userHandlers "ginLearn/src/controller/user"
	"github.com/gin-gonic/gin"
)

func BindRoute(handler *gin.Engine) {

	api := handler.Group("/api")

	user := api.Group("/user")

	{
		user.GET("", userHandlers.BasicUserInfo)
	}

	config := api.Group("/project")

	{
		config.GET("/config", projectHandlers.GetProjectConfig)
	}

	demo := api.Group("/demo")

	{
		demo.GET("/get", demoHandlers.MethodGetDemo)
		demo.GET("/path/:name/*action", demoHandlers.MethodPathDemo)
		demo.POST("/upload", demoHandlers.MethodUploadDemo)
		demo.POST("/json", demoHandlers.MethodJson)
		demo.POST("/login", demoHandlers.BindJson)
	}
}
package server

import (
	"github.com/alexanderbkl/vidre-back/internal/api"
	"github.com/alexanderbkl/vidre-back/internal/config"
	"github.com/alexanderbkl/vidre-back/internal/middlewares"
	"github.com/alexanderbkl/vidre-back/pkg/token"
	"github.com/gin-gonic/gin"
)

var APIv1 *gin.RouterGroup
var AuthAPIv1 *gin.RouterGroup

func registerRoutes(router *gin.Engine) {
	// Enables automatic redirection if the current route cannot be matched but a
	// handler for the path with (without) the trailing slash exists.
	// router.RedirectTrailingSlash = true

	// Create API router group.
	APIv1 = router.Group("/api")
	// Create AuthAPI router group.
	tokenMaker, err := token.NewPasetoMaker(config.Env().TokenSymmetricKey)
	if err != nil {
		log.Errorf("cannot create token maker: %s", err)
		panic(err)
	}



	AuthAPIv1 := router.Group("/api")
	AuthAPIv1.Use(middlewares.AuthMiddleware(tokenMaker))
	// routes
	api.Ping(APIv1)

	api.GetWorkers(APIv1)
	api.CreateWorker(APIv1)
	api.ModifyWorker(APIv1)
	api.DeleteWorker(APIv1)
	api.GetExtraHours(APIv1)
	api.ToggleExtraHours(APIv1)
	api.GetWorkDay(APIv1)
	api.PostWorkDay(APIv1)

}

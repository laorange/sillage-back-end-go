package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pilinux/gorest/config"
	"github.com/pilinux/gorest/database"
	"github.com/pilinux/gorest/lib/middleware"
	"sillage-back-end-go/controller"
)

var configure = config.Config()

func main() {
	if configure.Database.MongoDB.Activate == "yes" {
		// Initialize MONGO client
		if _, err := database.InitMongo(); err != nil {
			fmt.Println(err)
			return
		}
	}

	router, err := SetupRouter()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = router.Run(":" + configure.Server.ServerPort)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// SetupRouter ...
func SetupRouter() (*gin.Engine, error) {
	router := gin.Default()

	// Which proxy to trust
	if configure.Security.TrustedIP == "nil" {
		err := router.SetTrustedProxies(nil)
		if err != nil {
			return router, err
		}
	} else {
		if configure.Security.TrustedIP != "" {
			err := router.SetTrustedProxies([]string{configure.Security.TrustedIP})
			if err != nil {
				return router, err
			}
		}
	}

	router.Use(middleware.CORS())
	router.Use(middleware.SentryCapture(configure.Logger.SentryDsn))
	router.Use(middleware.Firewall(
		configure.Security.Firewall.ListType,
		configure.Security.Firewall.IP,
	))

	// Render HTML
	router.Use(middleware.Pongo2())

	// API:v1
	v1 := router.Group("/api/v1/")
	{
		v1.GET("course/", controller.CourseRetrieveList)
		v1.POST("course/", controller.CourseCreate)
	}

	return router, nil
}

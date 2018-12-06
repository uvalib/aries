package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const version = "1.2.1"

// favHandler is a dummy handler to silence browser API requests that look for /favicon.ico
func favHandler(c *gin.Context) {
}

// versionHandler reports the version of the serivce
func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "Aries version %s", version)
}

// healthCheckHandler reports the health of the serivce
func healthCheckHandler(c *gin.Context) {
	hcMap := make(map[string]string)
	hcMap["Aries"] = "true"
	for _, svc := range services {
		if pingService(svc, false) {
			hcMap[svc.Name] = "true"
		} else {
			hcMap[svc.Name] = "false"
		}
	}
	c.JSON(http.StatusOK, hcMap)
}

/**
 * MAIN
 */
func main() {
	log.Printf("===> Aries staring up <===")

	// Get config params; service port and redis info
	log.Printf("Read configuration...")
	var port, redisPort int
	var redisHost, redisPass string
	flag.IntVar(&port, "port", 8080, "Aries port (default 8080)")
	flag.StringVar(&redisHost, "redis_host", "localhost", "Redis host (default localhost)")
	flag.IntVar(&redisPort, "redis_port", 6379, "Redis port (default 6379)")
	flag.StringVar(&redisPass, "redis_pass", "", "Redis password")

	// NOTE: redisKeyPrefix is a string define in services.go
	flag.StringVar(&redisKeyPrefix, "redis_prefix", "aries", "Redis key prefix")

	flag.Parse()

	// Populate the service array with services known to Aries
	log.Printf("Init Services...")
	err := initServices(redisHost, redisPort, redisPass)
	if err != nil {
		log.Fatal("Unable to load services info")
	}

	log.Printf("Setup routes...")
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()
	router.GET("/favicon.ico", favHandler)
	router.GET("/version", versionHandler)
	router.GET("/healthcheck", healthCheckHandler)
	api := router.Group("/api")
	{
		api.GET("/resources/:id", resourcesHandler)
		api.GET("/services", servicesHandler)
		api.POST("/services", serviceAddHandler)
		api.PUT("/services", serviceUpdateHandler)
	}
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	// add a catchall route that renders the index page.
	// based on no-history config setup info here:
	//    https://router.vuejs.org/guide/essentials/history-mode.html#example-server-configurations
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	portStr := fmt.Sprintf(":%d", port)
	log.Printf("Start Aries v%s on port %s", version, portStr)
	log.Fatal(router.Run(portStr))
}

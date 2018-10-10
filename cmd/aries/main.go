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
const Version = "1.0.0"

// versionHandler reports the version of the serivce
func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "Aries version %s", Version)
}

// healthCheckHandler reports the health of the serivce
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alive": "true"})
}

// servicesHandler reports name and URL of all services seached by aries
func servicesHandler(c *gin.Context) {
	c.String(http.StatusOK, "None")
}

// resourcesHandler polls all services for info about the specified identifier
func resourcesHandler(c *gin.Context) {
	id := c.Param("id")
	c.String(http.StatusNotFound, "%s not found", id)
}

/**
 * MAIN
 */
func main() {
	log.Printf("===> Aries staring up <===")

	// Get config params
	var port, https int
	var key, crt string
	flag.IntVar(&port, "port", 8080, "Aries port (default 8080)")
	flag.IntVar(&https, "https", 0, "Use HTTPS? (default 0)")
	flag.StringVar(&key, "key", "", "Key for https connection")
	flag.StringVar(&crt, "crt", "", "Crt for https connection")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/version", versionHandler)
	router.GET("/healthcheck", healthCheckHandler)
	api := router.Group("/api")
	{
		api.GET("/resources/:id", resourcesHandler)
		api.GET("/services", servicesHandler)
	}
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	// add a catchall route that renders the index page.
	// based on no-history config setup info here:
	//    https://router.vuejs.org/guide/essentials/history-mode.html#example-server-configurations
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	portStr := fmt.Sprintf(":%d", port)
	if https == 1 {
		log.Printf("Start HTTPS Aries service on port %s", portStr)
		log.Fatal(router.RunTLS(portStr, crt, key))
	} else {
		log.Printf("Start HTTP Aries service on port %s", portStr)
		log.Fatal(router.Run(portStr))
	}
}

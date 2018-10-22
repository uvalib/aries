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
const version = "1.0.0"

// serviceInfo holds name and URL information for a service known to Aries
type serviceInfo struct {
	Name string
	URL  string
	OK   bool
}

// resoursesResponse is the respoonse format for a resources request
type resoursesResponse struct {
	SystemsSearched     int           `json:"systems_searched"`
	Hits                int           `json:"hits"`
	TotalResponseTimeMS int64         `json:"total_response_time_ms"`
	Responses           []interface{} `json:"responses"`
}

// services is a list of services known to Aries
var services []*serviceInfo

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
	hcMap["alive"] = "true"
	for _, svc := range services {
		if pingService(svc) {
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

	// Get config params
	var port int
	flag.IntVar(&port, "port", 8080, "Aries port (default 8080)")
	flag.Parse()

	// Populate the service array with services known to Aries
	err := initServices()
	if err != nil {
		log.Fatal("Unable to load services info")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/favison.ico", favHandler)
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
	log.Printf("Start Aries service on port %s", portStr)
	log.Fatal(router.Run(portStr))
}

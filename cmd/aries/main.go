package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const version = "1.0.0"

// serviceInfo holds name and URL information for a service known to Aries
type serviceInfo struct {
	Name string
	URL  string
}

// services is a list of services known to Aries
var services []serviceInfo

// versionHandler reports the version of the serivce
func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "Aries version %s", version)
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

func initServices() error {
	file, err := os.Open("services.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Got service: %s", line)
		bits := strings.Split(line, ",")
		services = append(services, serviceInfo{Name: bits[0], URL: bits[1]})
	}

	return nil
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

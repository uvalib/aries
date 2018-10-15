package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

// resoursesResponse is the respoonse format for a resources request
type resoursesResponse struct {
	SystemsSearched     int           `json:"systems_searched"`
	Hits                int           `json:"hits"`
	TotalResponseTimeMS int64         `json:"total_response_time_ms"`
	Responses           []interface{} `json:"responses"`
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
	c.JSON(http.StatusOK, services)
}

// resourcesHandler polls all services for info about the specified identifier
func resourcesHandler(c *gin.Context) {
	id := c.Param("id")
	channel := make(chan string)
	out := resoursesResponse{SystemsSearched: len(services)}
	start := time.Now()

	// kick off all requests in parallel
	for _, svc := range services {
		url := fmt.Sprintf("%s/%s", svc.URL, id)
		log.Printf("Check %s : %s for identifier %s", svc.Name, svc.URL, id)
		go getAriesResponse(svc.Name, url, channel)
	}

	// wait for all to be done and get respnses as they come in
	for range services {
		var jsonMap map[string]interface{}
		jsonRespStr := <-channel
		json.Unmarshal([]byte(jsonRespStr), &jsonMap)

		if int(jsonMap["status"].(float64)) == 200 {
			out.Hits++
		}
		out.Responses = append(out.Responses, jsonMap)
	}

	elapsedNanoSec := time.Since(start)
	out.TotalResponseTimeMS = int64(elapsedNanoSec / time.Millisecond)

	c.JSON(http.StatusOK, out)
}

func getAriesResponse(system string, url string, channel chan string) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	start := time.Now()
	resp, err := client.Get(url)
	elapsedNanoSec := time.Since(start)
	elapsedMS := int64(elapsedNanoSec / time.Millisecond)

	if err != nil {
		log.Printf("ERROR: GET %s failed : %s", url, err.Error())
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "Timeout") {
			status = http.StatusRequestTimeout
		}
		channel <- fmt.Sprintf(`{"system": "%s", "status": %d, "response": "%s", "response_time_ms": %d}`,
			system, status, err.Error(), elapsedMS)
		return
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	respString := string(bodyBytes)
	if resp.StatusCode != 200 {
		channel <- fmt.Sprintf(`{"system": "%s", "status": %d, "response": "%s", "response_time_ms": %d}`,
			system, resp.StatusCode, respString, elapsedMS)
		return
	}
	channel <- fmt.Sprintf(`{"system": "%s", "status": %d, "response": %s, "response_time_ms": %d}`,
		system, resp.StatusCode, respString, elapsedMS)
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

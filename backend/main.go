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

// versionHandler reports the version of the serivce
func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "Aries version %s", version)
}

// healthCheckHandler reports the health of the serivce
func healthCheckHandler(c *gin.Context) {
	hcMap := make(map[string]string)
	hcMap["alive"] = "true"
	for _, svc := range services {
		log.Printf("svs %s curr status %t", svc.Name, svc.OK)
		if pingService(svc) {
			hcMap[svc.Name] = "true"
		} else {
			hcMap[svc.Name] = "false"
		}
	}
	c.JSON(http.StatusOK, hcMap)
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
	outstandingRequests := 0
	for _, svc := range services {
		if svc.OK == false {
			log.Printf("Service %s is currently not active. See if it is now available...", svc.Name)
			if !pingService(svc) {
				log.Printf("   ...Skipping")
				out.Responses = append(out.Responses, gin.H{
					"system": svc.Name, "status": http.StatusServiceUnavailable,
					"response": "system is offline", "response_time_ms": 0,
				})
				continue
			}
		}
		log.Printf("Check %s : %s for identifier %s", svc.Name, svc.URL, id)
		outstandingRequests++
		go getAriesResponse(svc, id, channel)
	}

	// wait for all to be done and get respnses as they come in
	for outstandingRequests > 0 {
		var jsonMap map[string]interface{}
		jsonRespStr := <-channel
		json.Unmarshal([]byte(jsonRespStr), &jsonMap)

		if int(jsonMap["status"].(float64)) == 200 {
			out.Hits++
		}
		out.Responses = append(out.Responses, jsonMap)
		outstandingRequests--
	}

	elapsedNanoSec := time.Since(start)
	out.TotalResponseTimeMS = int64(elapsedNanoSec / time.Millisecond)

	c.JSON(http.StatusOK, out)
}

func getAriesResponse(svc *serviceInfo, id string, channel chan string) {
	url := fmt.Sprintf("%s/%s", svc.URL, id)
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
		errMsg := err.Error()
		if strings.Contains(err.Error(), "Timeout") {
			status = http.StatusRequestTimeout
			errMsg = "request timed out"
		} else if strings.Contains(err.Error(), "connection refused") {
			status = http.StatusServiceUnavailable
			errMsg = "system is offline"
		}
		svc.OK = false
		channel <- fmt.Sprintf(`{"system": "%s", "status": %d, "response": "%s", "response_time_ms": %d}`,
			svc.Name, status, errMsg, elapsedMS)
		return
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	respString := string(bodyBytes)
	if resp.StatusCode != 200 {
		channel <- fmt.Sprintf(`{"system": "%s", "status": %d, "response": "%s", "response_time_ms": %d}`,
			svc.Name, resp.StatusCode, respString, elapsedMS)
		return
	}

	channel <- fmt.Sprintf(`{"system": "%s", "status": %d, "response": %s, "response_time_ms": %d}`,
		svc.Name, resp.StatusCode, respString, elapsedMS)
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

		// create a and track a service; assume it is not alive by default
		// ping  will test and update this alive status
		svc := &serviceInfo{Name: bits[0], URL: bits[1], OK: false}
		services = append(services, svc)
		if !pingService(svc) {
			log.Printf("ERROR: Service %s is not available", svc.Name)
		}
	}

	return nil
}

// pingService will ping the service URL and return true if the serice responds
// The status of the service object will be updated based on this test
func pingService(svc *serviceInfo) bool {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	// Flag service as dead. If the test below is OK, update to alive
	svc.OK = false
	log.Printf("Ping sevice %s - %s...", svc.Name, svc.URL)
	resp, err := client.Get(svc.URL)
	if err != nil {
		log.Printf("Ping FAILED : %s", err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Ping returned bad status code : %d: ", resp.StatusCode)
		return false
	}
	log.Printf("%s is alive", svc.Name)
	svc.OK = true
	return true
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

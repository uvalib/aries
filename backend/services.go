package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// servicesHandler is a gin GET request handler. It
// reports name URL and status of all services seached by aries
func servicesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, services)
}

// serviceAddHandler will add/update a service contained in the post params. It will then
// ping that service and ensure the response is the expected format
// Call example: curl -d '{"name":"NAME", "url":"URL"}' -H "Content-Type: application/json" -X POST https://aries.lib.virginia.edu/api/services
func serviceAddHandler(c *gin.Context) {
	// Pull ServiceName and ServiceURL from query params
	var newSvc serviceInfo
	err := c.Bind(&newSvc)
	if err != nil {
		log.Printf("Bad request to update service: %s", err.Error())
		c.String(http.StatusBadRequest, "invalid request")
		return
	}
	log.Printf("Request to add/update service: %s - %s", newSvc.Name, newSvc.URL)

	// Before anything gets updated, ping the service and verify the response matches
	// expected format: [ServiceName] Aries API
	if !pingService(&newSvc, true) {
		c.String(http.StatusBadRequest, "invalid request")
		return
	}

	updated := false
	for _, svc := range services {
		if svc.Name == newSvc.Name {
			svc.URL = newSvc.URL
			log.Printf("%s URL updated to %s", svc.Name, svc.URL)
		}
	}
	if updated == false {
		services = append(services, &newSvc)
		log.Printf("Added new service")
	}

	c.String(http.StatusOK, "%s added", newSvc.Name)
}

// initServices will parse the services CSV file and create a list of available services.
// Each one will be will be checked to see if it is alive and this status is tracked.
func initServices() error {
	file, err := os.Open("services.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		bits := strings.Split(line, ",")

		// create a and track a service; assume it is not alive by default
		// ping  will test and update this alive status
		svc := &serviceInfo{Name: bits[0], URL: bits[1], OK: false}
		services = append(services, svc)
		log.Printf("Init new sevice %s - %s...", svc.Name, svc.URL)
		if !pingService(svc, false) {
			log.Printf("   * %s is not available", svc.Name)
		} else {
			log.Printf("   * %s is alive", svc.Name)
		}
	}

	// Start a ticker to periodically poll servivces and mark them
	// active or inactive. The weird syntax puts the polling of
	// the ticker channel an a goroutine so it doesn't block
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			log.Printf("Service check heartbeat")
			pingAllServices()
		}
	}()

	return nil
}

func pingAllServices() {
	errors := false
	for _, svc := range services {
		if !pingService(svc, false) {
			errors = true
		}
	}
	if errors == false {
		log.Printf("   * All services online")
	}
}

// getServiceURL will look up a service by name and return its URL
func getServiceURL(name string) string {
	for _, svc := range services {
		if svc.Name == name {
			return svc.URL
		}
	}
	return ""
}

// pingService will ping the service URL and return true if the service responds.
// The status of the service object will be updated based on this test.
func pingService(svc *serviceInfo, nameCheck bool) bool {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	ariesURL := fmt.Sprintf("%s/aries", svc.URL)
	resp, err := client.Get(ariesURL)
	if err != nil {
		log.Printf("   * FAIL: Service %s error : %s", svc.Name, err.Error())
		svc.OK = false
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("   * FAIL: Service %s returned bad status code : %d: ", svc.Name, resp.StatusCode)
		svc.OK = false
		return false
	}

	// read response and make sure it contains the name of the service
	respTxt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("   * FAIL: Service %s returned unreadable response : %s: ", svc.Name, err.Error())
		svc.OK = false
		return false
	}

	if nameCheck {
		expected := fmt.Sprintf("%s Aries API", svc.Name)
		if string(respTxt) != expected {
			log.Printf("   * FAIL: Service %s returned unexpected response [%s]", svc.Name, respTxt)
			svc.OK = false
			return false
		}
	}

	svc.OK = true
	return true
}

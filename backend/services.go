package main

import (
	"bufio"
	"fmt"
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
		if !pingService(svc) {
			log.Printf("ERROR: Service %s is not available", svc.Name)
		}
	}

	return nil
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
func pingService(svc *serviceInfo) bool {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	// Flag service as dead. If the test below is OK, update to alive
	svc.OK = false
	log.Printf("Ping sevice %s - %s...", svc.Name, svc.URL)
	ariesURL := fmt.Sprintf("%s/aries", svc.URL)
	resp, err := client.Get(ariesURL)
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

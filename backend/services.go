package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// serviceInfo holds name and URL information for a service known to Aries
type serviceInfo struct {
	ID int 		`json:"id"`
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
	OK   bool   `json:"alive"`
}
// services is a list of services known to Aries
var services []*serviceInfo

// The redis connection
var redisClient *redis.Client

// initServices will read services from the specified redis instance and create a list of services.
// Each one will be will be checked to see if it is alive and this status is tracked.
func initServices(host string, port int, pass string) error {
	redisHost := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Connect to redis instance at %s", redisHost)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: pass,
		DB:       0, // use default DB
	})

	// See if the connection is good...
	_, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}

	// Notes on redis data:
	//   aries:services contains a list of IDs for each service present
	//   aries:service:[id] contains a hash with service details; name and url
	//   aries:next_service_id is the next available ID for a new service

	// Get all of the service IDs, iterate them to get details and
	// establish connection / status
	svcIDs := redisClient.SMembers("aries:services").Val()
	for _, svcID := range svcIDs {
		redisID := fmt.Sprintf("aries:service:%s", svcID)
		svcInfo, svcErr := redisClient.HGetAll(redisID).Result()
		if svcErr != nil {
			log.Printf("Unable to get info for service %s", redisID)
			continue
		}

		// create a and track a service; assume it is not alive by default
		// ping  will test and update this alive status
		svc := &serviceInfo{Name: svcInfo["name"], URL: svcInfo["url"], OK: false}
		services = append(services, svc)
		log.Printf("Init %s - %s...", svc.Name, svc.URL)
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

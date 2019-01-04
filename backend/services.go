package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// serviceInfo holds name and URL information for a service known to Aries
type serviceInfo struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
	OK   bool   `json:"alive"`
}

// services is a list of services known to Aries
var services []*serviceInfo
var redisKeyPrefix string

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
	log.Printf("Connected; reading services...")
	servicesKey := fmt.Sprintf("%s:services", redisKeyPrefix)
	svcIDs := redisClient.SMembers(servicesKey).Val()
	for _, svcIDStr := range svcIDs {
		svcID, _ := strconv.ParseInt(svcIDStr, 10, 64)
		redisID := fmt.Sprintf("%s:service:%d", redisKeyPrefix, svcID)
		svcInfo, svcErr := redisClient.HGetAll(redisID).Result()
		if svcErr != nil {
			log.Printf("Unable to get info for service %s", redisID)
			continue
		}

		// create a and track a service; assume it is not alive by default
		// ping  will test and update this alive status
		svc := &serviceInfo{ID: svcID, Name: svcInfo["name"], URL: svcInfo["url"], OK: false}
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
	ticker := time.NewTicker(60 * time.Minute)
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

// serviceAddHandler will add a service contained in the post params. It will then
// ping that service and ensure the response is the expected format
// Call example: curl -d '{"id": "ID", "name":"NAME", "url":"URL"}' -H "Content-Type: application/json" -X POST https://aries.lib.virginia.edu/api/services
func serviceAddHandler(c *gin.Context) {
	var postedSvc serviceInfo
	err := c.BindJSON(&postedSvc)
	if err != nil {
		log.Printf("Bad request to update service: %s", err.Error())
		c.String(http.StatusBadRequest, "invalid request")
		return
	}
	log.Printf("Request to add service: %s - %s. Ping URL...", postedSvc.Name, postedSvc.URL)

	// Before anything gets updated, ping the service and verify the response matches
	// expected format: [ServiceName] Aries API
	if !pingService(&postedSvc, true) {
		c.String(http.StatusBadRequest, "Unable to reach target service")
		return
	}

	log.Printf("Get ID for %s,", postedSvc.Name)
	serviceIDKey := fmt.Sprintf("%s:next_service_id", redisKeyPrefix)
	newID, err := redisClient.Incr(serviceIDKey).Result()
	if err != nil {
		log.Printf("Unable to get ID new service")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	postedSvc.ID = newID
	redisErr := updateRedis(&postedSvc, true)
	if redisErr != nil {
		log.Printf("Unable to get update redis %s", redisErr.Error())
		c.String(http.StatusInternalServerError, redisErr.Error())
		return
	}
	services = append(services, &postedSvc)

	c.String(http.StatusOK, "%d", newID)
}

// serviceUpdateHandler will update the service contained in the post params
func serviceUpdateHandler(c *gin.Context) {
	var postedSvc serviceInfo
	err := c.BindJSON(&postedSvc)
	if err != nil {
		log.Printf("Bad request to update service: %s", err.Error())
		c.String(http.StatusBadRequest, "invalid request")
		return
	}
	log.Printf("Request to update service: %d: %s - %s", postedSvc.ID, postedSvc.Name, postedSvc.URL)

	// Find existing service...
	var existingService *serviceInfo
	for _, svc := range services {
		if svc.ID == postedSvc.ID {
			existingService = svc
		}
	}

	if existingService == nil {
		log.Printf("Service %d not found", postedSvc.ID)
		c.String(http.StatusBadRequest, "Service not found")
		return
	}

	// Only need to ping service if the URL is to be changed by the update
	if existingService.URL != postedSvc.URL {
		if !pingService(&postedSvc, true) {
			c.String(http.StatusBadRequest, "New URL is not valid")
			return
		}
	}

	redisErr := updateRedis(&postedSvc, false)
	if redisErr != nil {
		log.Printf("Unable to get update redis %s", redisErr.Error())
		c.String(http.StatusInternalServerError, redisErr.Error())
		return
	}
	existingService.Name = postedSvc.Name
	existingService.URL = postedSvc.URL

	c.String(http.StatusOK, "updated")
}

func updateRedis(svcInfo *serviceInfo, newService bool) error {
	// hmset aries:service:[ID] name [NAME] url [URL]
	redisID := fmt.Sprintf("%s:service:%d", redisKeyPrefix, svcInfo.ID)
	_, err := redisClient.HMSet(redisID, map[string]interface{}{
		"id":   svcInfo.ID,
		"name": svcInfo.Name,
		"url":  svcInfo.URL,
	}).Result()
	if err != nil {
		return err
	}

	// This is a new service.. add the ID to aries:services
	if newService {
		servicesKey := fmt.Sprintf("%s:services", redisKeyPrefix)
		_, err = redisClient.SAdd(servicesKey, svcInfo.ID).Result()
	}
	return err
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
	timeout := time.Duration(5 * time.Second)
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
		if strings.Contains(string(respTxt), " Aries API") == false {
			log.Printf("   * FAIL: Service %s returned unexpected response [%s]", svc.Name, respTxt)
			svc.OK = false
			return false
		}
	}

	svc.OK = true
	return true
}

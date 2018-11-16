package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// resourcesHandler is a gin GET handler for /api/respices/:ID.
// it polls all services for info about the specified identifier
func resourcesHandler(c *gin.Context) {
	id := c.Param("id")
	channel := make(chan string)
	out := resoursesResponse{SystemsSearched: len(services)}
	start := time.Now()

	// kick off all requests in parallel
	outstandingRequests := 0
	for _, svc := range services {
		if svc.OK == false {
			log.Printf("Service %s is currently not active", svc.Name)
			out.Responses = append(out.Responses, gin.H{
				"system": svc.Name, "status": http.StatusServiceUnavailable,
				"response": "system is offline", "response_time_ms": 0,
			})
			continue
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
	url := fmt.Sprintf("%s/aries/%s", svc.URL, id)
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

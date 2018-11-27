# Aries

As opposed to a monolithic repository, the UVA digital repository solution involves a great number of independent systems.  Experienced and knowledgeable staff make finding, accessing and integrating the pieces possible, but it's time to bring the information together in a more machine-actionable way.

The goal of Aries is to provide a method for seeing the full picture in a coherent way.

### System Requirements
* GO version 1.11 or greater
* Node version 8.11.1 or higher (https://nodejs.org/en/)
* Yarn version 1.9.3 or greater
* vue-cli 3 version 3.0.3 or greater
* Vue 2.5 or greater
* Redis 3.2.6 or greater

### Build Instructions

1. After clone, `cd frontend` and execute `yarn install` to prepare the front end
2. Run the Makefile target `build all` to generate binaries for linux, darwin and the front end.  All results will be in the bin directory.

### Current API

* GET /version : return service version info
* GET /healthcheck : test health of system components; results returned as JSON.
* GET /resources/:ID : Get a block of JSON data containing details about the specified resource.
* GET /services : Get a JSON list services that are a part of the aries constellation.
* POST /services : Add a new service or update an existing service. 
   * JSON Payload: {"name":"NAME", "url":"URL"}. 
   * Example: `curl -d '{"name":"NAME", "url":"URL"}' -H "Content-Type: application/json" -X POST https://aries.lib.virginia.edu/api/services`
   * Note: before add or update, the service will be pinged and the response scanned for expected content. If not present, the POT will fail. 

### Notes

To run in a development mode for the frontend only, build and launch aries, and supply a port param: `./aries.linux -port 8085`. Set an ENV variable:

* `ARIES_API` - set to localhost at the port used to launch the server (8085 from the example)

Run the frontend in development mode with `yarn serve` from the frontend directory. All API requests from the frontend will be redirected to the local instance of the backend services.

### Redis Notes

This repo includes a file containg some initial services. It is found in `redis_seed.txt`. To run it:

* `cat redis_seed.txt |  host] [-p port] -[a auth]`
  
The services Aries knows about are manintained in a redis instance. To see what keys it current uses, execute:

* `redis-cli [-h host] [-p port] -[a auth] --scan --pattern aries:*`

To clean up all of these keys, execute:

* `redis-cli [-h host] [-p port] -[a auth] --scan --pattern aries:* | xargs redis-cli [-h host] [-p port] -[a auth] del` 
